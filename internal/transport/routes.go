package transport

func (a *Api) registerRoutes() {
	registerHealthRoutes(a.api)
	registerAuthRoutes(a.api, a.auth)
}
