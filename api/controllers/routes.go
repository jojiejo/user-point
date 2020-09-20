package controllers

func (s *Server) initializeRoutes() {
	api := s.Router.Group("/api")
	{
		//User
		api.GET("/users", s.GetUsers)
		user := api.Group("/user")
		{
			user.GET("/:id", s.GetUserByID)
			user.GET("/:id/point-history", s.GetUserPointByUserID)
			user.POST("/", s.CreateUser)
			user.PATCH("/:id/point", s.UpdateUserPoint)
			user.DELETE("/:id", s.DeleteUser)
		}
	}
}
