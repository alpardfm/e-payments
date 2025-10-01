package usecase

import (
	"github.com/alpardfm/e-payment/src/business/domain"
	"github.com/alpardfm/e-payment/src/utils/config"
	"github.com/alpardfm/go-toolkit/log"
	"github.com/alpardfm/go-toolkit/parser"
)

type Usecases struct{}

func Init(log log.Interface, d *domain.Domains, jsonParser parser.JSONInterface, cfg config.Application) *Usecases {
	return &Usecases{}

}
