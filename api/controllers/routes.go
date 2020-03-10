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
	s.Router.GET("/sites/latest", s.GetLatestSites)
	site := s.Router.Group("/site")
	{
		site.GET("/:id", s.GetSite)
		site.GET("/:id/history", s.GetSiteHistory)
		site.POST("/", s.CreateSite)
		site.PUT("/:id", s.UpdateSite)
		site.DELETE("/:id", s.DeactivateSite)
		site.PATCH("/:id", s.ReactivateSite)
		/*site.DELETE("/:id/later", s.TerminateSiteLater)
		site.DELETE("/:id", s.TerminateSiteNow)*/

		//Retailer Site Relation
		site.GET("/:id/terminal-relations", s.GetSiteTerminalRelationBySiteID)
		site.POST("/:id/terminal-relation", s.CreateSiteTerminalRelation)
		site.PUT("/:id/terminal-relation/:relation_id", s.UpdateSiteTerminalRelation)
		site.DELETE("/:id/terminal-relation/:relation_id", s.UnlinkSiteTerminalRelation)
	}

	//Retailer
	s.Router.GET("/retailers", s.GetRetailers)
	s.Router.GET("/retailers/latest", s.GetLatestRetailers)
	s.Router.GET("/retailer-payment-terms", s.GetPaymentTerms)
	s.Router.GET("/retailer-reimbursement-cycles", s.GetReimbursementCycles)
	retailer := s.Router.Group("/retailer")
	{
		retailer.GET("/:id", s.GetRetailer)
		retailer.GET("/:id/history", s.GetRetailerHistory)
		retailer.POST("/", s.CreateRetailer)
		retailer.PUT("/:id", s.UpdateRetailer)
		retailer.POST(":id/deactivate", s.DeactivateRetailer)
		retailer.DELETE("/:id", s.DeactivateRetailer)
		retailer.PATCH("/:id", s.ReactivateRetailer)
		/*retailer.DELETE("/:id/later", s.TerminateRetailerLater)
		retailer.DELETE("/:id", s.TerminateRetailerNow)*/

		//Retailer Site Relation
		retailer.GET("/:id/site-relations", s.GetRetailerSiteRelationByRetailerID)
		retailer.POST("/:id/site-relation", s.CreateRetailerSiteRelation)
		retailer.PUT("/:id/site-relation/:relation_id", s.UpdateRetailerSiteRelation)
		retailer.DELETE("/:id/site-relation/:relation_id", s.UnlinkRetailerSiteRelation)
	}

	//Terminal
	s.Router.GET("/terminals", s.GetTerminals)
	s.Router.GET("/terminals/latest", s.GetLatestTerminals)
	terminal := s.Router.Group("/terminal")
	{
		terminal.GET("/:id", s.GetTerminal)
		terminal.GET("/:id/history", s.GetTerminalHistory)
		terminal.POST("/", s.CreateTerminal)
		terminal.PUT("/:id", s.UpdateTerminal)
		terminal.DELETE("/:id", s.DeactivateTerminal)
		terminal.PATCH("/:id", s.ReactivateTerminal)
	}
}
