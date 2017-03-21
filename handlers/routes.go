package handlers

// DrawRoutes configure all routes
func (r *Router) DrawRoutes() {
	r.GET("/", IndexHandler)
	r.POST("/v1/collect", CollectHandler)
}
