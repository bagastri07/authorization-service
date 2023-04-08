package grpc

import (
	"errors"

	"github.com/bagastri07/boilerplate-service/internal/model"
)

func (t *Server) InjectProductUsecase(productUC model.ProductUsecase) error {
	if productUC == nil {
		return errors.New("productUC can't be nil")
	}
	t.productUC = productUC
	return nil
}
