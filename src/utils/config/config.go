package config

import (
	"time"

	"github.com/alpardfm/go-toolkit/log"
	"github.com/alpardfm/go-toolkit/parser"
	"github.com/alpardfm/go-toolkit/sql"
)

type Application struct {
	Log    log.Config
	Meta   ApplicationMeta
	Gin    GinConfig
	SQL    sql.Config
	JWT    JWTConfig
	Parser parser.Options
}

type ApplicationMeta struct {
	Title       string
	Description string
	Host        string
	BasePath    string
	Version     string
}

type GinConfig struct {
	Port        string
	Mode        string
	LogRequest  bool
	LogResponse bool
	Timeout     time.Duration
	CORS        CORSConfig
	Swagger     SwaggerConfig
	Platform    PlatformConfig
	Dummy       DummyConfig
}

type CORSConfig struct {
	Mode string
}
type SwaggerConfig struct {
	Enabled   bool
	Path      string
	BasicAuth BasicAuthConf
}

type PlatformConfig struct {
	Enabled   bool
	Path      string
	BasicAuth BasicAuthConf
}

type DummyConfig struct {
	Enabled bool
	Path    string
}

type BasicAuthConf struct {
	Username string
	Password string
}

type JWTConfig struct {
	JWTTokenExpirationInMinute        int64
	DashboardJWTTokenExpirationMinute int64
	JWTTokenKey                       string
}

func Init() Application {
	return Application{}
}
