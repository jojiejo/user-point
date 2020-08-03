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
		payer.GET("/:id/transaction-invoice/:month/:year", s.GetTransactionInvoiceByPayer)
		payer.GET("/:id/fee-invoice/:month/:year", s.GetFeeInvoiceByPayer)
	}

	//Branch
	branch := s.Router.Group("/branch")
	{
		branch.GET("/:id", s.GetBranch)
		branch.GET("/:id/card-groups", s.GetCardGroupsByBranchID)
		branch.PUT("/:id", s.UpdateCardGroupFlagInSelectedBranch)
		branch.GET("/:id/transaction-invoice/:month/:year", s.GetTransactionInvoiceByBranch)
		branch.GET("/:id/fee-invoice/:month/:year", s.GetFeeInvoiceByBranch)
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
	product := s.Router.Group("/product")
	{
		product.GET("/:id", s.GetProduct)
		product.POST("/", s.CreateProduct)
		product.PUT("/:id", s.UpdateProduct)
	}

	//Product Group
	s.Router.GET("/product-groups", s.GetProductGroups)
	productGroup := s.Router.Group("/product-group")
	{
		productGroup.GET("/:id", s.GetProductGroup)
		productGroup.POST("/", s.CreateProductGroup)
		productGroup.PUT("/:id", s.UpdateProductGroup)

		/*productGroup.GET("/:id/product", s.GetProductInProductGroup)
		productGroup.PUT("/:id")*/
	}

	//Sales Rep
	s.Router.GET("/sales-reps", s.GetSalesReps)
	salesRep := s.Router.Group("/sales-rep")
	{
		salesRep.GET("/:id", s.GetSalesRep)
		salesRep.POST("/", s.CreateSalesRep)
		salesRep.PUT("/:id", s.UpdateSalesRep)
	}

	//Card Type
	s.Router.GET("/card-types", s.GetCardTypes)
	cardType := s.Router.Group("/card-type")
	{
		cardType.GET("/:id", s.GetCardType)
		cardType.POST("/", s.CreateCardType)
		cardType.PUT("/:id", s.UpdateCardType)
	}

	//Membership
	s.Router.GET("/memberships", s.GetMemberships)
	membership := s.Router.Group("/membership")
	{
		membership.GET("/:id", s.GetMembership)
		membership.POST("/", s.CreateMembership)
		membership.PUT("/:id", s.UpdateMembership)
	}

	//Industry Classification
	s.Router.GET("/industry-classifications", s.GetIndustryClassifications)
	industryClassification := s.Router.Group("/industry-classification")
	{
		industryClassification.GET("/:id", s.GetIndustryClassification)
		industryClassification.POST("/", s.CreateIndustryClassification)
		industryClassification.PUT("/:id", s.UpdateIndustryClassification)
	}

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
	assignedRebatePayerRelation := s.Router.Group("/rebate/payer-relation")
	{
		assignedRebatePayerRelation.GET("/:id", s.GetRebatePayerRelationByID)
		assignedRebatePayerRelation.PUT("/:id", s.UpdateRebatePayerRelation)
	}

	s.Router.POST("/rebate/main/payer-relations", s.CreateMainRebatePayer)
	//s.Router.POST("/rebate/:rebate_id/payer-association/:pa_id", s.GetRelationByRebateAndPA)
	assignedPromotionalRebate := s.Router.Group("/rebate/promotional")
	{
		assignedPromotionalRebate.POST("/payer-relations", s.CreatePromotionalRebatePayer)
		assignedPromotionalRebate.POST("/bulk-check", s.CheckBulkAssignRebateToPayer)
		assignedPromotionalRebate.POST("/bulk-assign", s.BulkAssignRebateToPayer)
	}

	//Posting Matrix VAT
	s.Router.GET("/posting-matrix/vats", s.GetPostingMatrixVAT)

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
		tax.POST("/", s.CreateTax)
		tax.PUT("/:id", s.UpdateTax)
	}

	//Faktur Pajak
	s.Router.GET("/faktur-pajak-ranges", s.GetAllFakturPajakRange)
	s.Router.GET("/faktur-pajak-ranges/available-number", s.GetNextAvailableFakturPajakNumber)
	s.Router.GET("/faktur-pajak-ranges/available-range", s.GetAvailableFakturPajakRange)
	fakturPajakRange := s.Router.Group("/faktur-pajak-range")
	{
		fakturPajakRange.POST("/", s.CreateFakturPajakRange)
	}

	//Transactions
	transaction := s.Router.Group("/transactions")
	{
		transaction.GET("/manual-settlement/:dateFrom/:dateTo", s.GetAllTransactionForManualSettlement)
	}

	//Manual Settlement
	s.Router.POST("/manual-settlement", s.ManualSettle)
	s.Router.POST("/manual-settlement/all", s.ManualSettleAllTransaction)

	//Card (Telematic device purpose)
	card := s.Router.Group("/card")
	{
		card.GET("/:id/telematic-device", s.GetTelematicDeviceByCardID)
		card.PUT("/:id/telematic-device", s.UpdateTelematicDevice)
	}
}
