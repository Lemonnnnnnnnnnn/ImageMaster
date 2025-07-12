package config

type API struct {
	manager *Manager
}

func NewAPI(appName string) *API {
	return &API{
		manager: NewManager(appName),
	}
}

func (a *API) GetActiveLibrary() string {
	return a.manager.GetActiveLibrary()
}

func (a *API) SetActiveLibrary(library string) bool {
	return a.manager.SetActiveLibrary(library)
}

func (a *API) GetOutputDir() string {
	return a.manager.GetOutputDir()
}

func (a *API) SetOutputDir() bool {
	return a.manager.SetOutputDir()
}

func (a *API) GetProxy() string {
	return a.manager.GetProxy()
}

func (a *API) SetProxy(proxy string) bool {
	return a.manager.SetProxy(proxy)
}

func (a *API) GetLibraries() []string {
	return a.manager.GetLibraries()
}

func (a *API) AddLibrary() bool {
	return a.manager.AddLibrary()
}
