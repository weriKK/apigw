package service

import (
	"apigw/internal/util"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func parseServiceRegistration(r *http.Request) (*ServiceRegistration, error) {

	var sr ServiceRegistration

	if err := json.NewDecoder(r.Body).Decode(&sr); err != nil {
		return nil, util.NewJSONErrorMessage("unable to parse JSON payload", err.Error())
	}

	if err := mandatoryFieldCheck(sr); err != nil {
		return nil, util.NewJSONErrorMessage("missing mandatory fields", err.Error())
	}

	if err := supportedHttpMethodsCheck(sr); err != nil {
		return nil, util.NewJSONErrorMessage("unsupported HTTP method", err.Error())
	}

	return &sr, nil
}

func mandatoryFieldCheck(sr ServiceRegistration) error {

	var missingFields []string

	if sr.Name == "" {
		missingFields = append(missingFields, "name")
	}

	if len(sr.Endpoints) < 1 {
		missingFields = append(missingFields, "endpoints")
	}

	for i, ep := range sr.Endpoints {
		if ep.ApiRoot == "" {
			missingFields = append(missingFields, fmt.Sprintf("endpoints[%d].apiRoot", i))
		}
		if len(ep.Methods) < 1 {
			missingFields = append(missingFields, fmt.Sprintf("endpoints[%d].methods", i))
		}
	}

	if 0 < len(missingFields) {
		return fmt.Errorf("missing fields: %s", strings.Join(missingFields, ", "))
	}

	return nil
}

func supportedHttpMethodsCheck(sr ServiceRegistration) error {

	for _, ep := range sr.Endpoints {
		for i, method := range ep.Methods {
			switch method {
			case
				"GET",
				"HEAD",
				"POST",
				"PUT",
				"DELETE",
				//"CONNECT", // not supported
				"OPTIONS",
				//"TRACE", // not supported
				"PATCH":
			default:
				return fmt.Errorf("%q in endpoints[%d].methods is not a supported HTTP method", method, i)
			}
		}
	}
	return nil
}
