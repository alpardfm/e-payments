package domain

import (
	"github.com/alpardfm/e-payment/src/utils/config"
	"github.com/alpardfm/go-toolkit/log"
	"github.com/alpardfm/go-toolkit/parser"
	"github.com/alpardfm/go-toolkit/sql"
)

type Domains struct{}

func Init(log log.Interface, db sql.Interface, parser parser.JSONInterface, cfg config.Application) *Domains {
	return &Domains{}
}
