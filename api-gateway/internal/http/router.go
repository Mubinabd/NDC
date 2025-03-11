package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "posts/docs"
	"posts/internal/http/handlers"
)

// @title NDC Post Project API Documentation
// @version 1.0
// @description API for Instant Delivery resources
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func NewGin(h *handlers.Handler) *gin.Engine {
	router := gin.Default()

	router.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Adjust for your specific origins
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	//enforcer, err := casbin.NewEnforcer("./internal/http/casbin/model.conf", "./internal/http/casbin/policy.csv")
	//if err != nil {
	//	log.Println("k")
	//}
	//router.Use(middlerware.NewAuth(enforcer))

	posts := router.Group("/api/v1/posts")
	{
		posts.POST("/create", h.Create)
		posts.GET("/:id", h.GetDetail)
		posts.PATCH("/update", h.UpdatePost)
		posts.DELETE("/:id", h.Delete)
		posts.GET("/list", h.GetList)
	}

	logs := router.Group("/api/v1/logs")
	{
		logs.POST("/create", h.CreateLog)
		logs.GET("/:id", h.GetLogDetail)
		logs.PATCH("/update", h.UpdateLog)
		logs.DELETE("/:id", h.DeleteLog)
		logs.GET("/list", h.GetLogList)
	}

	users := router.Group("/api/v1/users")
	{
		users.POST("/create", h.CreateUser)
		users.GET("/:id", h.GetUserDetail)
		users.PATCH("/update", h.UpdateUser)
		users.DELETE("/:id", h.DeleteUser)
		users.GET("/list", h.GetUserList)
		users.PUT("/user-password", h.ChangeUserPassword)
	}
	router.POST("/api/v1/login", h.Login)

	return router
}
