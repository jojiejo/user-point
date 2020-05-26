package controllers

func (s *Server) initializeRoutes() {
	//Global Variable
	s.Router.GET("/global-variables", s.GetGlobalVariables)
	gv := s.Router.Group("/global-variable")
	{
		gv.GET("/:id/detail", s.GetGlobalVariableDetailByGlobalVariableID)
	}
	gvd := s.Router.Group("/global-variable-detail")
	{
		gvd.GET("/:id", s.GetGlobalVariableDetail)
		gvd.POST("/", s.CreateGlobalVariableDetail)
		gvd.PUT("/:id", s.UpdateGlobalVariableDetail)
	}

	//Card Status
	s.Router.GET("/card-status", s.GetAllCardStatus)

	//Unit
	s.Router.GET("/units", s.GetAllUnits)

	//Province
	s.Router.GET("/provinces", s.GetProvinces)
	s.Router.GET("/province/:id", s.GetProvince)
	s.Router.GET("/province/:id/cities", s.GetCitiesByProvinceID)

	//City
	s.Router.GET("/cities", s.GetCities)
	s.Router.GET("/city/:id", s.GetCity)

	//Site
	s.Router.GET("/sites", s.GetSites)
	s.Router.GET("/site-types", s.GetSiteTypes)
	s.Router.GET("/sites/latest", s.GetLatestSites)
	s.Router.GET("/sites/active", s.GetActiveSites)
	site := s.Router.Group("/site")
	{
		site.GET("/:id", s.GetSite)
		site.GET("/:id/history", s.GetSiteHistory)
		site.POST("/", s.CreateSite)
		site.PUT("/:id", s.UpdateSite)
		site.DELETE("/:id/now", s.DeactivateSiteNow)
		site.DELETE("/:id/later", s.DeactivateSiteLater)
		site.PATCH("/:id", s.ReactivateSite)
		/*site.DELETE("/:id/later", s.TerminateSiteLater)
		site.DELETE("/:id", s.TerminateSiteNow)*/

		//Site Terminal Relation
		site.GET("/:id/terminal-relations", s.GetTerminalBySiteID)
		/*site.POST("/:id/terminal-relation", s.CreateSiteTerminalRelation)
		site.PUT("/:id/terminal-relation/:relation_id", s.UpdateSiteTerminalRelation)
		site.DELETE("/:id/terminal-relation/:relation_id", s.UnlinkSiteTerminalRelation)*/
	}

	//Retailer
	s.Router.GET("/retailers", s.GetRetailers)
	s.Router.GET("/retailers/latest", s.GetLatestRetailers)
	s.Router.GET("/retailers/active", s.GetActiveRetailers)
	s.Router.GET("/retailer-payment-terms", s.GetPaymentTerms)
	s.Router.GET("/retailer-reimbursement-cycles", s.GetReimbursementCycles)
	retailer := s.Router.Group("/retailer")
	{
		retailer.GET("/:id", s.GetRetailer)
		retailer.GET("/:id/history", s.GetRetailerHistory)
		retailer.POST("/", s.CreateRetailer)
		retailer.PUT("/:id", s.UpdateRetailer)
		retailer.DELETE("/:id/now", s.DeactivateRetailerNow)
		retailer.DELETE("/:id/later", s.DeactivateRetailerLater)
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
	s.Router.GET("/retailer/:id/site/:site_id/terminal-overview", s.GetTerminalOverview)
	s.Router.GET("/terminals/latest", s.GetLatestTerminals)
	terminal := s.Router.Group("/terminal")
	{
		terminal.GET("/:id", s.GetTerminal)
		terminal.GET("/:id/history", s.GetTerminalHistory)
		terminal.POST("/", s.CreateTerminal)
		terminal.PUT("/:id", s.UpdateTerminal)
		terminal.DELETE("/:id/now", s.DeactivateTerminalNow)
		terminal.DELETE("/:id/later", s.DeactivateTerminalLater)
		terminal.PATCH("/:id", s.ReactivateTerminal)
	}

	//Payer Association
	s.Router.GET("/payer-associations", s.GetPayerAssociations)
	s.Router.GET("/payer-association/:id/payers", s.GetPayerByPayerAssociationID)

	//Payer
	s.Router.GET("/payers", s.GetPayers)
	payer := s.Router.Group("/payer")
	{
		payer.GET("/:id", s.GetPayer)
		payer.GET("/:id/branches", s.GetBranchByCCID)
		payer.PATCH("/:id/configuration", s.UpdateConfiguration)
		payer.GET("/:id/charged-fees/automated", s.GetChargedAutomatedFeesOnSelectedAccount)
	}

	//Branch
	branch := s.Router.Group("/branch")
	{
		branch.GET("/:id", s.GetBranch)
		branch.GET("/:id/card-groups", s.GetCardGroupsByBranchID)
		branch.PUT("/:id", s.UpdateCardGroupFlagInSelectedBranch)
	}

	//Card Group
	cardGroup := s.Router.Group("/card-group")
	{
		cardGroup.GET("/:id", s.GetCardGroupByID)
		cardGroup.POST("/", s.CreateCardGroup)
		cardGroup.PUT("/:id", s.UpdateCardGroup)
		cardGroup.DELETE("/:id", s.DeactivateCardGroup)
	}

	//GSAP Master Data
	gsap_master_data := s.Router.Group("/gsap-customer-master-data")
	{
		gsap_master_data.GET("/:id", s.GetCustomerMasterData)
	}

	//Fee
	s.Router.GET("/initial-fees", s.GetInitialFees)
	s.Router.GET("/initial-fees/type/:id", s.GetInitialFeeByFeeType)
	initialFee := s.Router.Group("/initial-fee")
	{
		initialFee.GET("/:id", s.GetInitialFee)
		initialFee.POST("/", s.CreateFee)
		initialFee.PUT("/:id", s.UpdateFee)
		initialFee.DELETE("/:id", s.DeactivateFee)
	}

	//Charge Ad Hoc Fee
	s.Router.GET("/charged-fees/ad-hoc", s.GetChargedAdHocFees)
	adHocFee := s.Router.Group("/charged-fee/ad-hoc")
	{
		adHocFee.GET("/:id", s.GetChargedAdHocFee)
		adHocFee.POST("/", s.ChargeAdHocFee)
		adHocFee.POST("/bulk-check", s.CheckBulkChargeAdHocFee)
		adHocFee.POST("/bulk-charge", s.BulkChargeAdHocFee)
	}

	//Charge Automated Fee
	automatedFee := s.Router.Group("/charged-fee/automated")
	{
		automatedFee.GET("/:id", s.GetChargedAutomatedFee)
		automatedFee.PUT("/:id", s.UpdateAutomatedFee)
	}

	//Product
	s.Router.GET("/products", s.GetProducts)

	//Rebate
	s.Router.GET("/rebate/calculation-types", s.GetRebateCalculationTypes)
	s.Router.GET("/rebate/types", s.GetRebateTypes)
	s.Router.GET("/rebate/periods", s.GetRebatePeriods)
	rebatePrograms := s.Router.Group("/rebate/programs")
	{
		rebatePrograms.GET("/", s.GetRebatePrograms)
		rebatePrograms.GET("/type/:id", s.GetRebateProgramsByTypeID)
	}
	rebateProgram := s.Router.Group("/rebate/program")
	{
		rebateProgram.GET("/:id", s.GetRebateProgram)
		rebateProgram.POST("/", s.CreateRebateProgram)
		rebateProgram.PUT("/:id", s.UpdateRebateProgram)
		//rebateProgram.DELETE("/:id", s.GetRebateProgram)
	}

	//Rebate to Account
	s.Router.GET("/rebate/payer-relations", s.GetRebatePayerRelations)
	s.Router.POST("/rebate/main/payer-relations", s.CreateMainRebatePayer)
	s.Router.POST("/rebate/promotional/payer-relations", s.CreatePromotionalRebatePayer)

	//Posting Matrix by Product
	s.Router.GET("/posting-matrix/products", s.GetPostingMatrixProducts)
	postingMatrixByProduct := s.Router.Group("/posting-matrix/product")
	{
		postingMatrixByProduct.GET("/:id", s.GetPostingMatrixProduct)
		postingMatrixByProduct.POST("/", s.CreatePostingMatrixProduct)
		postingMatrixByProduct.PUT("/:id", s.UpdatePostingMatrixProduct)
	}

	//Posting Matrix by Fee
	s.Router.GET("/posting-matrix/fees", s.GetPostingMatrixFees)
	postingMatrixByFee := s.Router.Group("/posting-matrix/fee")
	{
		postingMatrixByFee.GET("/:id", s.GetPostingMatrixFee)
		postingMatrixByFee.POST("/", s.CreatePostingMatrixFee)
		postingMatrixByFee.PUT("/:id", s.UpdatePostingMatrixFee)
	}

	//Posting Matrix by Tax
	s.Router.GET("/posting-matrix/taxes", s.GetPostingMatrixTaxes)
	postingMatrixByTax := s.Router.Group("/posting-matrix/tax")
	{
		postingMatrixByTax.GET("/:id", s.GetPostingMatrixTax)
		postingMatrixByTax.POST("/", s.CreatePostingMatrixTax)
		postingMatrixByTax.PUT("/:id", s.UpdatePostingMatrixTax)
	}

	//Tax Type
	s.Router.GET("/tax-types", s.GetAllTaxTypes)

	//Tax
	s.Router.GET("/taxes", s.GetAllTaxes)
	tax := s.Router.Group("/tax")
	{
		tax.GET("/:id", s.GetTax)
	}
}
