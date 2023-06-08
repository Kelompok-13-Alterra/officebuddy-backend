package rest

import (
	"context"
	"fmt"
	"go-clean/docs/swagger"
	"go-clean/src/business/usecase"
	"go-clean/src/lib/auth"
	"go-clean/src/lib/configreader"
	"go-clean/src/utils/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var once = &sync.Once{}

type REST interface {
	Run()
}

type rest struct {
	http         *gin.Engine
	conf         config.GinConfig
	configreader configreader.Interface
	uc           *usecase.Usecase
	auth         auth.Interface
}

func Init(conf config.GinConfig, confReader configreader.Interface, uc *usecase.Usecase, auth auth.Interface) REST {
	r := &rest{}
	once.Do(func() {
		switch conf.Mode {
		case gin.ReleaseMode:
			gin.SetMode(gin.ReleaseMode)
		case gin.DebugMode, gin.TestMode:
			gin.SetMode(gin.TestMode)
		default:
			gin.SetMode("")
		}

		httpServ := gin.Default()

		r = &rest{
			conf:         conf,
			configreader: confReader,
			http:         httpServ,
			uc:           uc,
			auth:         auth,
		}

		switch r.conf.CORS.Mode {
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

		r.Register()
	})

	return r
}

func (r *rest) Run() {
	port := ":8080"
	if r.conf.Port != "" {
		port = fmt.Sprintf(":%s", r.conf.Port)
	}

	server := &http.Server{
		Addr:    port,
		Handler: r.http,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println(fmt.Sprintf("Serving HTTP error: %s", err.Error()))
		}
	}()
	fmt.Println(fmt.Sprintf("Listening and Serving HTTP on %s", server.Addr))

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), r.conf.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(fmt.Sprintf("Server forced to shutdown: %v", err))
	}

	log.Println("Server exiting")
}

func (r *rest) Register() {
	r.registerSwaggerRoutes()
	publicApi := r.http.Group("/public")
	publicApi.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": "hello world",
		})
	})

	api := r.http.Group("/api")
	v1 := api.Group("/v1")

	v1.GET("/", r.VerifyUser, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": "hello mail",
		})
	})

	auth := v1.Group("/auth")
	auth.POST("/register", r.RegisterUser)
	auth.POST("/login", r.LoginUser)
	auth.POST("/admin-login", r.LoginAdmin)

	user := v1.Group("/user")
	user.PUT("/edit", r.VerifyUser, r.UpdateUser)
	user.GET("/me", r.VerifyUser, r.GetProfile)
	user.GET("", r.VerifyUser, r.VerifyAdmin, r.GetUserList)

	office := v1.Group("/office")
	office.POST("", r.VerifyUser, r.VerifyAdmin, r.CreateOffice)
	office.GET("", r.VerifyUser, r.GetOfficeList)
	office.GET("/:office_id", r.VerifyUser, r.GetOffice)
	office.PUT("/:office_id", r.VerifyUser, r.VerifyAdmin, r.UpdateOffice)
	office.DELETE("/:office_id", r.VerifyUser, r.VerifyAdmin, r.DeleteOffice)

	transaction := v1.Group("/transaction")
	transaction.POST("/office/:office_id/book", r.VerifyUser, r.CreateOrder)
	transaction.GET("/booked", r.VerifyUser, r.GetTransactionBookedList)
	transaction.GET("/history", r.VerifyUser, r.GetTransactionHistoryBookedList)
	transaction.GET("/:transaction_id/payment-detail", r.VerifyUser, r.VerifyTransaction, r.GetPaymentDetail)
	transaction.PUT("/:transaction_id/reschedule", r.VerifyUser, r.VerifyTransaction, r.RescheduleBooked)
	transaction.GET("", r.VerifyUser, r.VerifyAdmin, r.GetTransactionList)

	notification := v1.Group("/notification")
	notification.GET("", r.VerifyUser, r.GetNotificationList)

	midtransTransaction := v1.Group("/midtrans-transaction")
	midtransTransaction.POST("/handle", r.HandleNotification)

	rating := v1.Group("/rating")
	rating.GET("", r.VerifyUser, r.VerifyAdmin, r.GetRatingList)
	rating.GET("/:rating_id", r.VerifyUser, r.VerifyAdmin, r.GetRating)
	rating.POST("/:transaction_id", r.VerifyUser, r.CreateRating)

	widgetAnalytic := v1.Group("")
	widgetAnalytic.GET("/admin/dashboard-widget", r.VerifyUser, r.VerifyAdmin, r.GetDashboardWidget)
}

func (r *rest) registerSwaggerRoutes() {
	swagger.SwaggerInfo.Title = r.conf.Meta.Title
	swagger.SwaggerInfo.Description = r.conf.Meta.Description
	swagger.SwaggerInfo.Version = r.conf.Meta.Version
	swagger.SwaggerInfo.Host = r.conf.Meta.Host
	swagger.SwaggerInfo.BasePath = r.conf.Meta.BasePath

	r.http.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
