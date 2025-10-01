package rest

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/alpardfm/e-payment/docs/swagger"
	"github.com/alpardfm/e-payment/src/business/usecase"
	"github.com/alpardfm/e-payment/src/utils/config"
	"github.com/alpardfm/go-toolkit/configreader"
	"github.com/alpardfm/go-toolkit/log"
	"github.com/alpardfm/go-toolkit/parser"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gopkg.in/yaml.v2"
)

const (
	infoRequest  string = `httpclient Sent Request: uri=%v method=%v`
	infoResponse string = `httpclient Received Response: uri=%v method=%v resp_code=%v`
)

var once = &sync.Once{}

type REST interface {
	Run()
}

type rest struct {
	http         *gin.Engine
	conf         config.Application
	configreader configreader.Interface
	json         parser.JSONInterface
	log          log.Interface
	uc           *usecase.Usecases
}

func Init(conf config.Application, configreader configreader.Interface, log log.Interface, json parser.JSONInterface, uc *usecase.Usecases) REST {
	r := &rest{}
	once.Do(func() {

		switch conf.Gin.Mode {
		case gin.DebugMode:
			gin.SetMode(gin.DebugMode)
		case gin.ReleaseMode:
			gin.SetMode(gin.ReleaseMode)
		case gin.TestMode:
			gin.SetMode(gin.TestMode)
		default:
			gin.SetMode("")
		}

		httpServer := gin.New()

		r = &rest{
			conf:         conf,
			configreader: configreader,
			log:          log,
			json:         json,
			http:         httpServer,
			uc:           uc,
		}

		// Set CORS
		switch r.conf.Gin.CORS.Mode {
		case "allowall":
			r.http.Use(cors.New(cors.Config{
				AllowAllOrigins: true,
				AllowHeaders:    []string{"*"},
				AllowMethods: []string{
					http.MethodHead,
					http.MethodGet,
					http.MethodPost,
					http.MethodPut,
					http.MethodPatch,
					http.MethodDelete,
				},
			}))
		default:
			r.http.Use(cors.New(cors.DefaultConfig()))
		}

		// Set Recovery
		r.http.Use(gin.Recovery())

		// Set Timeout
		r.http.Use(r.SetTimeout)

		r.Register()
	})

	return r
}

func (r *rest) Run() {
	if r.conf.Gin.Port != "" {
		r.http.Run(fmt.Sprintf(":%s", r.conf.Gin.Port))
	} else {
		r.http.Run(":3001")
	}
}

func (r *rest) registerSwaggerRoutes() {
	if r.conf.Gin.Swagger.Enabled {
		swagger.SwaggerInfo.Title = r.conf.Meta.Title
		swagger.SwaggerInfo.Description = r.conf.Meta.Description
		swagger.SwaggerInfo.Version = r.conf.Meta.Version
		swagger.SwaggerInfo.Host = r.conf.Meta.Host
		swagger.SwaggerInfo.BasePath = r.conf.Meta.BasePath

		swaggerAuth := gin.Accounts{
			r.conf.Gin.Swagger.BasicAuth.Username: r.conf.Gin.Swagger.BasicAuth.Password,
		}

		r.http.GET(fmt.Sprintf("%s/*any", r.conf.Gin.Swagger.Path),
			gin.BasicAuthForRealm(swaggerAuth, "Restricted"),
			ginSwagger.WrapHandler(swaggerfiles.Handler))
	}
}

func (r *rest) registerPlatformRoutes() {
	if r.conf.Gin.Platform.Enabled {
		platformAuth := gin.Accounts{
			r.conf.Gin.Platform.BasicAuth.Username: r.conf.Gin.Platform.BasicAuth.Password,
		}

		r.http.GET(r.conf.Gin.Platform.Path,
			gin.BasicAuthForRealm(platformAuth, "Restricted"),
			r.platformConfig)
	}
}

func (r *rest) platformConfig(ctx *gin.Context) {
	switch ctx.Query("output") {
	case "json":
		ctx.IndentedJSON(http.StatusOK, r.configreader.AllSettings())
	default:
		c, err := yaml.Marshal(r.configreader.AllSettings())
		if err != nil {
			r.httpRespError(ctx, err)
			return
		}
		ctx.String(http.StatusOK, string(c))
	}
}
