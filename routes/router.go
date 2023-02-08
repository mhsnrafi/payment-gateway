package routes

import (
	"checkout-task/docs"
	"checkout-task/middlewares"
	"checkout-task/models"
	"checkout-task/services"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

func New() *gin.Engine {
	r := gin.New()
	initRoute(r)

	r.Use(gin.LoggerWithWriter(middlewares.LogWriter()))
	r.Use(gin.CustomRecovery(middlewares.AppRecovery()))
	r.Use(middlewares.CORSMiddleware())
	r.Use(middlewares.ResponseTimeMiddleware())

	v1 := r.Group("/v1")
	{
		AuthRoute(v1)
		Payments(v1)

	}

	docs.SwaggerInfo.BasePath = v1.BasePath()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	// Add Prometheus endpoint for metrics collection
	r.GET("/metrics", func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	})

	// Add pprof endpoint for profiling
	pprof.Register(r)

	return r
}

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests processed by the service.",
		},
		[]string{"code", "method"},
	)
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Time taken to complete HTTP requests.",
			Buckets: []float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0},
		},
		[]string{"code", "method"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal, httpRequestDuration)
}

func initRoute(r *gin.Engine) {
	_ = r.SetTrustedProxies(nil)
	r.RedirectTrailingSlash = false
	r.HandleMethodNotAllowed = true

	r.NoRoute(func(c *gin.Context) {
		models.SendErrorResponse(c, http.StatusNotFound, c.Request.RequestURI+" not found")
	})

	r.NoMethod(func(c *gin.Context) {
		models.SendErrorResponse(c, http.StatusMethodNotAllowed, c.Request.Method+" is not allowed here")
	})
}

func InitGin() {
	gin.DisableConsoleColor()
	gin.SetMode(services.Config.Mode)
}
