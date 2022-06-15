package service

type ServiceDB struct {
	Service map[string]ServiceInformation
}

type ServiceInformation struct {
	Name         string
	ApiRoot      string
	ReadinessURI string
	Status       ServiceStatus
	Endpoints    []ServiceEndpoint
}

type ServiceStatus uint

const (
	ServiceStatusUnknown     ServiceStatus = 0
	ServiceStatusReady       ServiceStatus = 1
	ServiceStatusNotReady    ServiceStatus = 2
	ServiceStatusUnreachable ServiceStatus = 3
)

var serviceStatusText = map[ServiceStatus]string{
	ServiceStatusUnknown:     "Unknown",
	ServiceStatusReady:       "Ready",
	ServiceStatusNotReady:    "NotReady",
	ServiceStatusUnreachable: "Unreachable",
}

func (ss ServiceStatus) String() string {
	return serviceStatusText[ss]
}

var _DB *ServiceDB = &ServiceDB{
	Service: make(map[string]ServiceInformation),
}

func (db *ServiceDB) FindServices() map[string]ServiceInformation {
	return _DB.Service
}

func (db *ServiceDB) FindServiceByName(name string) *ServiceInformation {

	if v, found := db.Service[name]; found {
		return &v
	}

	return nil
}

func (db *ServiceDB) UpdateServiceStatus(name string, status ServiceStatus) {

	if tmp := db.FindServiceByName(name); tmp != nil {
		tmp.Status = status
		db.SaveService(tmp)
	}
}

func (db *ServiceDB) SaveService(s *ServiceInformation) {
	db.Service[s.Name] = *s
}
