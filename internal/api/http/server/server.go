package server

import (
	"net/http"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/secure"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	_ "encomers/license/docs"
)

type Server struct {
	router *gin.Engine
	logger *zap.Logger
}

// @title           License Plate Generator API
// @version         1.0.0
// @description     API description
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /
func New(logger *zap.Logger) *Server {

	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(secure.New(secure.Config{
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'",
	}))

	router.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "Not found")
	})

	return &Server{
		router: router,
		logger: logger,
	}
}

func (s *Server) Router() *gin.Engine {
	return s.router
}

func (s *Server) Run(addr string) error {
	s.logger.Info("starting server", zap.String("address", addr))
	return s.router.Run(addr)
}

func (s *Server) Get(relativePath string, handler gin.HandlerFunc) {
	s.router.GET(relativePath, handler)
}
