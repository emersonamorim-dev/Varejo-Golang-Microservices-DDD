package routes

import (
	"Varejo-Golang-Microservices/auth"
	"Varejo-Golang-Microservices/middleware"

	customerHandler "Varejo-Golang-Microservices/services/customer-service/api/handler"
	customerRepository "Varejo-Golang-Microservices/services/customer-service/domain/repository"
	customerService "Varejo-Golang-Microservices/services/customer-service/domain/service"
	integrationHandler "Varejo-Golang-Microservices/services/integration-service/api/handler"
	integrationRepository "Varejo-Golang-Microservices/services/integration-service/domain/repository"
	integrationService "Varejo-Golang-Microservices/services/integration-service/domain/service"
	locationHandler "Varejo-Golang-Microservices/services/location-service/api/handler"
	locationRepository "Varejo-Golang-Microservices/services/location-service/domain/repository"
	locationService "Varejo-Golang-Microservices/services/location-service/domain/service"
	orderHandler "Varejo-Golang-Microservices/services/order-service/api/handler"
	orderRepository "Varejo-Golang-Microservices/services/order-service/domain/repository"
	orderService "Varejo-Golang-Microservices/services/order-service/domain/service"
	paymentHandler "Varejo-Golang-Microservices/services/payment-service/api/handler"
	paymentRepository "Varejo-Golang-Microservices/services/payment-service/domain/repository"
	paymentService "Varejo-Golang-Microservices/services/payment-service/domain/service"
	productHandler "Varejo-Golang-Microservices/services/product-service/api/handler"
	productRepository "Varejo-Golang-Microservices/services/product-service/domain/repository"
	productService "Varejo-Golang-Microservices/services/product-service/domain/service"
	promotionHandler "Varejo-Golang-Microservices/services/promotion-service/api/handler"
	promotionRepository "Varejo-Golang-Microservices/services/promotion-service/domain/repository"
	promotionService "Varejo-Golang-Microservices/services/promotion-service/domain/service"
	reportHandler "Varejo-Golang-Microservices/services/report-service/api/handler"
	reportRepository "Varejo-Golang-Microservices/services/report-service/domain/repository"
	reportService "Varejo-Golang-Microservices/services/report-service/domain/service"
	supportHandler "Varejo-Golang-Microservices/services/support-service/api/handler"
	supportRepository "Varejo-Golang-Microservices/services/support-service/domain/repository"
	supportService "Varejo-Golang-Microservices/services/support-service/domain/service"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, mongoURI string, kafkaBroker string) {
	// Define a rota de autenticação
	r.POST("/login", auth.Authenticate)

	// Grupo para endpoints protegidos por autenticação
	authorized := r.Group("/")
	authorized.Use(middleware.AuthMiddleware())

	// Inicialize conexões, repositórios e serviços do cliente.
	customerRepo := customerRepository.NewMongoCustomerRepository(mongoURI, kafkaBroker)
	custService := customerService.NewCustomerService(customerRepo)
	custHandler := customerHandler.NewCustomerHandler(custService)

	// Inicialização do integration-service
	integrationRepo := integrationRepository.NewMongoIntegrationRepository(mongoURI, kafkaBroker)
	integrationServ := integrationService.NewIntegrationService(integrationRepo)
	integrationHand := integrationHandler.NewIntegrationHandler(integrationServ)

	// Inicialização do Location Service
	locationRepo := locationRepository.NewMongoLocationRepository(mongoURI, kafkaBroker)
	locService := locationService.NewLocationService(locationRepo)
	locHandler := locationHandler.NewLocationHandler(locService)

	// Inicialize conexões, repositórios e serviços do cliente.
	orderRepo := orderRepository.NewMongoOrderRepository(mongoURI, kafkaBroker)
	ordService := orderService.NewOrderService(orderRepo)
	ordHandler := orderHandler.NewOrderHandler(ordService)

	// Inicialize conexões, repositórios e serviços do cliente
	payRepo := paymentRepository.NewMongoPaymentRepository(mongoURI, kafkaBroker)
	payService := paymentService.NewPaymentService(payRepo)
	payHandler := paymentHandler.NewPaymentHandler(payService)

	// Initialize product connections, repositories, services, and handlers.
	prodRepo := productRepository.NewMongoProductRepository(mongoURI, kafkaBroker)
	prodServ := productService.NewProductService(prodRepo)
	prodHand := productHandler.NewProductHandler(prodServ)

	// Inicialização do Promotion Service
	promotionRepo := promotionRepository.NewMongoPromotionRepository(mongoURI, kafkaBroker)
	promService := promotionService.NewPromotionService(promotionRepo)
	promHandler := promotionHandler.NewPromotionHandler(promService)

	// Inicialize conexões, repositórios e serviços do report-service.
	reportRepo := reportRepository.NewMongoReportRepository(mongoURI, kafkaBroker)
	rptService := reportService.NewReportService(reportRepo)
	rptHandler := reportHandler.NewReportHandler(rptService)

	// Initialize product connections, repositories, services, and handlers. do support-service:
	supportRepo := supportRepository.NewMongoSupportRepository(mongoURI, kafkaBroker)
	supService := supportService.NewSupportService(supportRepo)
	supHandler := supportHandler.NewSupportHandler(supService)

	// Configura routes para o custumer-service
	r.GET("/customers", custHandler.GetAllCustomers)
	r.GET("/customers/:id", custHandler.GetCustomerByID)
	r.POST("/customers", custHandler.AddCustomer)
	r.PUT("/customers/:id", custHandler.UpdateCustomer)
	r.DELETE("/customers/:id", custHandler.DeleteCustomer)

	// Configura routes para o integration-service
	r.GET("/integrations", integrationHand.ListIntegrationData)
	r.GET("/integrations/:id", integrationHand.GetIntegrationDataByID)
	r.POST("/integrations", integrationHand.AddIntegrationData)
	r.PUT("/integrations/:id", integrationHand.UpdateIntegrationData)
	r.DELETE("/integrations/:id", integrationHand.DeleteIntegrationData)

	// Configura routes para o location-service
	r.GET("/locations", locHandler.GetLocation)
	r.GET("/locations/:id", locHandler.GetLocationByID)
	r.POST("/locations", locHandler.AddLocation)
	r.PUT("/locations/:id", locHandler.UpdateLocation)
	r.DELETE("/locations/:id", locHandler.DeleteLocation)

	// Configura routes para o order-service
	r.GET("/orders", ordHandler.GetAllOrders)
	r.GET("/orders/:id", ordHandler.GetOrderByID)
	r.POST("/orders", ordHandler.AddOrder)
	r.PUT("/orders/:id", ordHandler.UpdateOrderStatus)
	r.DELETE("/orders/:id", ordHandler.DeleteOrder)

	// Configura routes para o payment-service
	r.GET("/payments", payHandler.GetAllPayments)
	r.GET("/payments/:id", payHandler.GetPaymentByID)
	r.POST("/payments", payHandler.AddPayment)
	r.PUT("/payments/:id", payHandler.UpdatePayment)
	r.DELETE("/payments/:id", payHandler.DeletePayment)

	// Configura routes para o product-service
	r.GET("/products", prodHand.ListProducts)
	r.GET("/products/:id", prodHand.GetProductByID)
	r.POST("/products", prodHand.AddProduct)
	r.PUT("/products/:id", prodHand.UpdateProduct)
	r.DELETE("/products/:id", prodHand.DeleteProduct)

	// Configura routes para o promotion-service
	r.GET("/promotions", promHandler.ListPromotions)
	r.GET("/promotions/:id", promHandler.GetPromotionByID)
	r.POST("/promotions", promHandler.AddPromotion)
	r.PUT("/promotions/:id", promHandler.UpdatePromotion)
	r.DELETE("/promotions/:id", promHandler.DeletePromotion)

	// Configura routes para o report-service
	r.GET("/reports", rptHandler.ListReports)
	r.GET("/reports/:id", rptHandler.GetReportByID)
	r.POST("/reports", rptHandler.AddReport)
	r.PUT("/reports/:id", rptHandler.UpdateReport)
	r.DELETE("/reports/:id", rptHandler.DeleteReport)

	// Configura routes para rotas o support-service:
	r.GET("/supports", supHandler.ListSupports)
	r.GET("/supports/:id", supHandler.GetSupportByID)
	r.POST("/supports", supHandler.AddSupport)
	r.PUT("/supports/:id", supHandler.UpdateSupport)
	r.DELETE("/supports/:id", supHandler.DeleteSupport)
}
