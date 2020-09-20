package controllers

import "github.com/jojiejo/user-point/api/middlewares"

func (s *Server) initializeRoutes() {
	api := s.Router.Group("/api")
	{
		//JWT => Create testing token
		api.GET("/jwt", s.GetJWT)

		//User
		api.GET("/users", s.GetUsers)
		user := api.Group("/user")
		{
			user.GET("/:id", s.GetUserByID)
			user.GET("/:id/point-history", s.GetUserPointByUserID)
			user.POST("/", middlewares.TokenAuthMiddleware(), s.CreateUser)
			user.PATCH("/:id/point", middlewares.TokenAuthMiddleware(), s.UpdateUserPoint)
			user.DELETE("/:id", middlewares.TokenAuthMiddleware(), s.DeleteUser)
		}
	}
}
