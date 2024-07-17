package handler

import (
	"api-gateway/config"
	"api-gateway/genproto/dish"
	"api-gateway/genproto/kitchen"
	"api-gateway/genproto/order"
	"api-gateway/genproto/payment"
	"api-gateway/genproto/review"
	"api-gateway/genproto/user"
	"api-gateway/pkg"
	"api-gateway/pkg/logger"
	"log/slog"
)

type Handler struct {
	UserClient    user.UserClient
	KitchenClient kitchen.KitchenClient
	DishClient    dish.DishClient
	OrderClient   order.OrderClient
	ReviewClient  review.ReviewClient
	PaymentClient payment.PaymentClient
	Logger        *slog.Logger
}

func NewHandler(cfg *config.Config) *Handler {
	return &Handler{
		UserClient:    pkg.NewUserClient(cfg),
		KitchenClient: pkg.NewKitchenClient(cfg),
		DishClient:    pkg.NewDishClient(cfg),
		OrderClient:   pkg.NewOrderClient(cfg),
		ReviewClient:  pkg.NewReviewClient(cfg),
		PaymentClient: pkg.NewPaymentClient(cfg),
		Logger:        logger.NewLogger(),
	}
}
