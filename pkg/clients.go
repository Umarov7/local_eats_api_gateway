package pkg

import (
	"api-gateway/config"
	pbd "api-gateway/genproto/dish"
	pbe "api-gateway/genproto/extra"
	pbk "api-gateway/genproto/kitchen"
	pbo "api-gateway/genproto/order"
	pbp "api-gateway/genproto/payment"
	pbr "api-gateway/genproto/review"
	pbu "api-gateway/genproto/user"
	"log"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewUserClient(cfg *config.Config) pbu.UserClient {
	conn, err := grpc.NewClient(cfg.AUTH_SERVICE_PORT,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Println(errors.Wrap(err, "failed to connect to the address"))
		return nil
	}

	return pbu.NewUserClient(conn)
}

func NewKitchenClient(cfg *config.Config) pbk.KitchenClient {
	conn, err := grpc.NewClient(cfg.AUTH_SERVICE_PORT,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Println(errors.Wrap(err, "failed to connect to the address"))
		return nil
	}

	return pbk.NewKitchenClient(conn)
}

func NewDishClient(cfg *config.Config) pbd.DishClient {
	conn, err := grpc.NewClient(cfg.ORDER_SERVICE_PORT,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Println(errors.Wrap(err, "failed to connect to the address"))
		return nil
	}

	return pbd.NewDishClient(conn)
}

func NewOrderClient(cfg *config.Config) pbo.OrderClient {
	conn, err := grpc.NewClient(cfg.ORDER_SERVICE_PORT,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Println(errors.Wrap(err, "failed to connect to the address"))
		return nil
	}

	return pbo.NewOrderClient(conn)
}

func NewReviewClient(cfg *config.Config) pbr.ReviewClient {
	conn, err := grpc.NewClient(cfg.ORDER_SERVICE_PORT,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Println(errors.Wrap(err, "failed to connect to the address"))
		return nil
	}

	return pbr.NewReviewClient(conn)
}

func NewPaymentClient(cfg *config.Config) pbp.PaymentClient {
	conn, err := grpc.NewClient(cfg.ORDER_SERVICE_PORT,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Println(errors.Wrap(err, "failed to connect to the address"))
		return nil
	}

	return pbp.NewPaymentClient(conn)
}

func NewExtraClient(cfg *config.Config) pbe.ExtraClient {
	conn, err := grpc.NewClient(cfg.ORDER_SERVICE_PORT,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Println(errors.Wrap(err, "failed to connect to the address"))
		return nil
	}

	return pbe.NewExtraClient(conn)
}
