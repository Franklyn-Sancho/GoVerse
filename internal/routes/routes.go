package routers

/* var userRouter routers.SetupUserRoutes();

func SetupRoutes(router *gin.Engine) {
	// Defina a chave secreta do JWT
	secretKey := os.Getenv("JWT_SECRET_KEY")
	log.Printf("SetupRoutes Secret Key: %s", secretKey)

	// Rotas públicas (não requerem autenticação)
	router.POST("/login", handlers.Login)
	router.POST("/register", handlers.RegisterUser)

	// Rotas protegidas (requerem autenticação)
	auth := router.Group("/")
	auth.Use(middleware.AuthMiddleware(secretKey))

	// Configuração das rotas de usuários e postagens
	routers.SetupUserRoutes(auth)
	routers.SetupPostRoutes(auth)
} */
