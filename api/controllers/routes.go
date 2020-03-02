package controllers

func (s *Server) initializeRoutes() {
	//Province
	s.Router.GET("/provinces", s.GetProvinces)
	s.Router.GET("/province/:id", s.GetProvince)
	s.Router.GET("/province/:id/cities", s.GetCitiesByProvinceID)

	//City
	s.Router.GET("/cities", s.GetCities)
	s.Router.GET("/city/:id", s.GetCity)

	//Site
	s.Router.GET("/sites", s.GetSites)
	site := s.Router.Group("/site")
	{
		site.GET("/:id", s.GetSite)
		site.POST("/", s.CreateSite)
		site.PUT("/:id", s.UpdateSite)
		site.DELETE("/:id", s.DeleteSite)
	}

	//Retailer
	s.Router.GET("/retailers", s.GetRetailers)
	s.Router.GET("/retailer-payment-terms", s.GetPaymentTerms)
	s.Router.GET("/retailer-reimbursement-cycles", s.GetReimbursementCycles)
	retailer := s.Router.Group("/retailer")
	{
		retailer.GET("/:id", s.GetRetailer)
		retailer.POST("/", s.CreateRetailer)
		retailer.PUT("/:id", s.UpdateRetailer)
		retailer.DELETE("/:id", s.DeleteRetailer)

		//Retailer Site Relation
		retailer.GET("/:id/site-relations", s.GetRetailerSiteRelationByRetailerID)
		retailer.POST("/:id/site-relation", s.CreateRetailerSiteRelation)
		//retailer.PUT("/:id/site-relation/:site_relation_id", s.UpdateRetailerSiteRelation)
	}

	//Terminal
	s.Router.GET("/terminals", s.GetTerminals)
	terminal := s.Router.Group("/terminal")
	{
		terminal.GET("/:id", s.GetTerminal)
		terminal.POST("/", s.CreateTerminal)
		terminal.PUT("/:id", s.UpdateTerminal)
		terminal.DELETE("/:id", s.DeleteTerminal)
	}
}
