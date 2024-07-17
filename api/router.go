package api

import (
	"api-gateway/api/handler"
	"api-gateway/api/middleware"
	"api-gateway/config"

	_ "api-gateway/api/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Local Eats
// @version 1.0
// @description API Gateway of Local Eats
// @host localhost:8080
// @BasePath /local-eats
// @schemes http
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func NewRouter(cfg *config.Config) *gin.Engine {
	h := handler.NewHandler(cfg)

	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/local-eats")
	api.Use(middleware.Check)

	u := api.Group("/users")
	{
		u.GET(":id", h.GetUser)
		u.PUT(":id", h.UpdateUser)
		u.DELETE(":id", h.DeleteUser)
		u.GET(":id/activity", h.TrackActivity)
	}

	k := api.Group("/kitchens")
	{
		k.POST("", h.CreateKitchen)
		k.GET(":id", h.GetKitchen)
		k.PUT(":id", h.UpdateKitchen)
		k.DELETE(":id", h.DeleteKitchen)
		k.GET("", h.FetchKitchens)
		k.GET("/search", h.SearchKitchens)
		k.GET(":id/dishes", h.FetchDishes)
		k.GET(":id/orders", h.FetchOrdersForKitchen)
		k.GET(":id/reviews", h.GetReviews)
		k.GET(":id/statistics", h.GetStatistics)
		k.POST(":id/working-hours", h.SetWorkingHours)
	}

	d := api.Group("/dishes")
	{
		d.POST("", h.CreateDish)
		d.GET(":id", h.GetDish)
		d.PUT(":id", h.UpdateDish)
		d.DELETE(":id", h.DeleteDish)
		d.GET(":id/nutrition", h.GetNutrition)
	}

	o := api.Group("/orders")
	{
		o.POST("", h.CreateOrder)
		o.GET(":id", h.GetOrderByID)
		o.PUT(":id/status", h.ChangeStatus)
		o.GET("", h.FetchOrdersForCustomer)
	}

	r := api.Group("/reviews")
	{
		r.POST("", h.CreateReview)
	}

	p := api.Group("/payments")
	{
		p.POST("", h.CreatePayment)
		p.GET(":id", h.GetPayment)
	}

	return router
}
