package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/alpardfm/e-payment/src/entity"
	"github.com/alpardfm/go-toolkit/appcontext"
	"github.com/alpardfm/go-toolkit/codes"
	"github.com/alpardfm/go-toolkit/errors"
	"github.com/alpardfm/go-toolkit/header"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
)

func (r *rest) BodyLogger(ctx *gin.Context) {
	if r.conf.Gin.LogRequest {
		r.log.Info(ctx.Request.Context(),
			fmt.Sprintf(infoRequest, ctx.Request.RequestURI, ctx.Request.Method))
	}

	ctx.Next()
	if r.conf.Gin.LogResponse {
		if ctx.Writer.Status() < 300 {
			r.log.Info(ctx.Request.Context(),
				fmt.Sprintf(infoResponse, ctx.Request.RequestURI, ctx.Request.Method, ctx.Writer.Status()))
		} else {
			r.log.Error(ctx.Request.Context(),
				fmt.Sprintf(infoResponse, ctx.Request.RequestURI, ctx.Request.Method, ctx.Writer.Status()))
		}
	}
}

// timeout middleware wraps the request context with a timeout
func (r *rest) SetTimeout(ctx *gin.Context) {
	// wrap the request context with a timeout
	c, cancel := context.WithTimeout(ctx.Request.Context(), r.conf.Gin.Timeout)

	defer func() {
		// check if context timeout was reached
		if c.Err() == context.DeadlineExceeded {
			// write response and abort the request
			r.httpRespError(ctx, errors.NewWithCode(codes.CodeContextDeadlineExceeded, "Context Deadline Exceeded"))
		}

		//cancel to clear resources after finished
		cancel()
	}()

	// replace request with context wrapped request
	ctx.Request = ctx.Request.WithContext(c)
	ctx.Next()

}

func (r *rest) addFieldsToContext(ctx *gin.Context) {
	reqid := ctx.GetHeader(header.KeyRequestID)
	if reqid == "" {
		reqid = uuid.New().String()
	}

	c := ctx.Request.Context()
	c = appcontext.SetRequestId(c, reqid)
	c = appcontext.SetUserAgent(c, ctx.Request.Header.Get(header.KeyUserAgent))
	c = appcontext.SetAcceptLanguage(c, ctx.Request.Header.Get(header.KeyAcceptLanguage))
	c = appcontext.SetServiceVersion(c, r.conf.Meta.Version)
	ctx.Request = ctx.Request.WithContext(c)
	ctx.Next()
}

func (r *rest) httpRespError(ctx *gin.Context, err error) {
	httpStatus, displayError := errors.Compile(err, appcontext.GetAcceptLanguage(ctx))
	statusStr := http.StatusText(httpStatus)

	c := ctx.Request.Context()
	errResp := &entity.HTTPResp{
		Message: entity.HTTPMessage{
			Title: displayError.Title,
			Body:  displayError.Body,
		},
		Meta: entity.Meta{
			Path:       r.conf.Meta.Host + ctx.Request.URL.String(),
			StatusCode: httpStatus,
			Status:     statusStr,
			Message:    fmt.Sprintf("%s %s [%d] %s", ctx.Request.Method, ctx.Request.URL.RequestURI(), httpStatus, statusStr),
			Error: &entity.MetaError{
				Code:    int(displayError.Code),
				Message: err.Error(),
			},
		},
	}

	r.log.Error(c, err)
	ctx.Header(header.KeyRequestID, appcontext.GetRequestId(c))
	ctx.AbortWithStatusJSON(httpStatus, errResp)
}

func (r *rest) httpRespSuccess(ctx *gin.Context, code codes.Code, data interface{}, p *entity.Pagination) {
	successApp := codes.Compile(code, appcontext.GetAcceptLanguage(ctx))
	c := ctx.Request.Context()
	meta := entity.Meta{
		Path:       r.conf.Meta.Host + ctx.Request.URL.String(),
		StatusCode: successApp.StatusCode,
		Status:     http.StatusText(successApp.StatusCode),
		Message:    fmt.Sprintf("%s %s [%d] %s", ctx.Request.Method, ctx.Request.URL.RequestURI(), successApp.StatusCode, http.StatusText(successApp.StatusCode)),
	}

	resp := &entity.HTTPResp{
		Message: entity.HTTPMessage{
			Title: successApp.Title,
			Body:  successApp.Body,
		},
		Meta:       meta,
		Data:       data,
		Pagination: p,
	}

	raw, err := r.json.Marshal(&resp)
	if err != nil {
		r.httpRespError(ctx, errors.NewWithCode(codes.CodeInternalServerError, "MarshalHTTPResp"))
		return
	}

	ctx.Header(header.KeyRequestID, appcontext.GetRequestId(c))
	ctx.Data(successApp.StatusCode, header.ContentTypeJSON, raw)
}

func (r *rest) Bind(ctx *gin.Context, obj interface{}) error {
	return ctx.ShouldBindWith(obj, binding.Default(ctx.Request.Method, ctx.ContentType()))
}

func (r *rest) Ping(ctx *gin.Context) {
	r.httpRespSuccess(ctx, codes.CodeSuccess, "PONG!", nil)
}
