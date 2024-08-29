package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ahnlabio/tsm-appserver/config"
	"github.com/ahnlabio/tsm-appserver/docs"
	"github.com/ahnlabio/tsm-appserver/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gitlab.com/Blockdaemon/go-tsm-sdkv2/v64/tsm"

	_ "github.com/ahnlabio/tsm-appserver/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var tsmStatic1 = tsm.Configuration{URL: "http://localhost:8501"}.WithAPIKeyAuthentication("apikey1")
var tsmStatic2 = tsm.Configuration{URL: "http://localhost:8502"}.WithAPIKeyAuthentication("apikey2")

func main() {
	godotenv.Load()
	swagInit()
	router := getRouter()
	addMiddlewares(router)
	runServerApplication(router)
}

func getRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", rootHandler)
	r.GET("/swagger/*any", func(c *gin.Context) {
		ginSwagger.WrapHandler(swaggerFiles.Handler)(c)
	})
	r.POST("/generateKey", handlers.GenerateKeyHandler)
	r.POST("/copyKey", handlers.CopyKeyHandler)
	r.POST("/preSign", handlers.PreSignHandler)
	r.POST("/finalizeSign", handlers.FinalizeSignHandler)
	return r
}

func addMiddlewares(r *gin.Engine) {
	r.Use(cors.Default())
}

func runServerApplication(router *gin.Engine) {
	srv := &http.Server{
		Addr:    ":3000",
		Handler: router.Handler(),
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 1 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	<-ctx.Done()
	log.Println("Server exiting")
}

func swagInit() {
	godotenv.Load()
	apiHostName := os.Getenv("HOST_NAME")
	docs.SwaggerInfo.Schemes = []string{"https", "http"}
	if apiHostName == "" {
		docs.SwaggerInfo.Schemes = []string{"http"}
		apiHostName = "localhost:3000"
	}
	docs.SwaggerInfo.Host = apiHostName
	docs.SwaggerInfo.Version = "0.1.0"
	docs.SwaggerInfo.Description = fmt.Sprintf("ABC Core BTC API v.%s", docs.SwaggerInfo.Version)

	for _, scheme := range docs.SwaggerInfo.Schemes {
		log.Printf("Swagger API URL: %s://%s/swagger/index.html\n", scheme, docs.SwaggerInfo.Host)
	}
}

// @title ABC Core BTC API
// @version 0.0.1
// @description ABC Core BTC API v0.0.1

// @host localhost:3000
// @BasePath /
type RootResponse struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	BuildType string `json:"build_type"`
	Time      string `json:"time"`
}

// rootHandler godoc
// @Summary Show the application info
// @Description get application info
// @Tags info
// @Accept  json
// @Produce  json
// @Success 200 {object} RootResponse
// @Router / [get]
func rootHandler(c *gin.Context) {
	config := config.GetConfig()
	current_time := time.Now().Format(time.RFC3339)

	response := RootResponse{
		Name:      config.AppName,
		Version:   config.AppVersion,
		BuildType: config.BuildType,
		Time:      current_time,
	}

	c.JSON(200, response)
}
