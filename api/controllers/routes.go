package controllers

func (s *Server) initializeRoutes() {
	api := s.Router.Group("/api")
	{
		//User
		s.Router.GET("/users", s.GetUsers)
		user := s.Router.Group("/api/user")
		{
			user.GET("/:id", s.GetUser)
			/*user.GET("/:id/points", s.GetUserPoint)
			user.PUT("/:id", s.UpdateUser)
			user.PUT("/:id/patch", s.UpdateUserPoint)
			user.DELETE(":/id", s.DeleteUser)*/
		}
	}
}
