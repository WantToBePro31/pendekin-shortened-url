package utils

import (
	apih "pendekin/handler/api"
	webh "pendekin/handler/web"
	"pendekin/middleware"
	apim "pendekin/middleware/api"
	webm "pendekin/middleware/web"
	"pendekin/repository"
	"pendekin/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func SetUpRoutes(server *gin.Engine, db *gorm.DB, rdb *redis.Client) *gin.Engine {
	server.Use(gin.Logger())
	server.Use(gin.Recovery())
	server.Use(cors.New(cors.Config{
		AllowMethods:     []string{"POST", "PUT", "PATCH", "GET", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowAllOrigins:  true,
		AllowCredentials: true,
	}))
	server.LoadHTMLGlob("view/**/*")

	urlRepository := repository.InitUrlRepository(db)
	userRepository := repository.InitUserRepository(db)
	sessionRepository := repository.InitSessionRepository(db)

	urlService := service.InitUrlService(urlRepository)
	userService := service.InitUserService(userRepository)
	sessionService := service.InitSessionService(sessionRepository)
	redisService := service.InitRedisService(rdb)

	urlHandlerApi := apih.InitUrlHandlerApi(urlService, sessionService, redisService)
	userHandlerApi := apih.InitUserHandlerApi(userService, sessionService, redisService)

	api := server.Group("/api")
	{
		user := api.Group("/user")
		{
			user.POST("/register", middleware.PostMiddleware(), userHandlerApi.Register)
			user.POST("/login", middleware.PostMiddleware(), apim.AlreadyLoginMiddleware(), userHandlerApi.Login)
			user.POST("/logout", middleware.PostMiddleware(), apim.IsLoginMiddleware(), userHandlerApi.Logout)
		}
		url := api.Group("/url")
		{
			url.GET("/all", middleware.GetMiddleware(), apim.IsLoginMiddleware(), urlHandlerApi.GetShortenedUrls)
			url.POST("/create", middleware.PostMiddleware(), urlHandlerApi.CreateShortenedUrl)
			url.POST("/update/:url_id", middleware.PostMiddleware(), apim.IsLoginMiddleware(), urlHandlerApi.EditShortenedUrls)
			url.POST("/delete/:url_id", middleware.PostMiddleware(), apim.IsLoginMiddleware(), urlHandlerApi.RemoveShortenedUrls)
		}
	}

	urlHandlerWeb := webh.InitUrlHandlerWeb(urlService, sessionService, redisService)
	userHandlerWeb := webh.InitUserHandlerWeb(userService, sessionService, redisService)

	user := server.Group("/user")
	{
		user.GET("/register", middleware.GetMiddleware(), userHandlerWeb.Register)
		user.POST("/register", middleware.PostMiddleware(), userHandlerWeb.RegisterProcess)
		user.GET("/login", middleware.GetMiddleware(), webm.AlreadyLoginMiddleware(), userHandlerWeb.Login)
		user.POST("/login", middleware.PostMiddleware(), webm.AlreadyLoginMiddleware(), userHandlerWeb.LoginProcess)
		user.POST("/logout", middleware.PostMiddleware(), webm.IsLoginMiddleware(), userHandlerWeb.Logout)
	}

	base := server.Group("/")
	{
		base.GET("", middleware.GetMiddleware(), urlHandlerWeb.Index)
		base.GET("/:shortened", middleware.GetMiddleware(), urlHandlerWeb.Real)
		base.GET("/dashboard", middleware.GetMiddleware(), webm.IsLoginMiddleware(), urlHandlerWeb.Dashboard)
	}

	url := server.Group("/url")
	{
		url.GET("", middleware.GetMiddleware(), urlHandlerWeb.Home)
		url.POST("/create", middleware.PostMiddleware(), urlHandlerWeb.CreateShortenedUrl)
		url.GET("/:shortened", middleware.GetMiddleware(), urlHandlerWeb.Result)
		url.GET("/update/:url_id", middleware.GetMiddleware(), webm.IsLoginMiddleware(), urlHandlerWeb.Update)
		url.POST("/update/:url_id", middleware.PostMiddleware(), webm.IsLoginMiddleware(), urlHandlerWeb.EditShortenedUrls)
		url.POST("/delete/:url_id", middleware.PostMiddleware(), webm.IsLoginMiddleware(), urlHandlerWeb.RemoveShortenedUrls)
	}
	return server
}
