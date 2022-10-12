package app

import (
	"Currency/internal/config"
	"Currency/internal/infrastructure/model"
	"Currency/internal/infrastructure/service"
	"Currency/internal/infrastructure/use_case/get_exchanges"
	"Currency/internal/infrastructure/use_case/update_exchanges"
	"context"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net"
	"net/http"
	"time"
)

type App struct {
	config     *config.Config
	router     *httprouter.Router
	httpServer *http.Server
	db         *gorm.DB
}

func NewKernel(config *config.Config) App {
	return App{
		config: config,
	}
}

func (a *App) Run() {
	log.Print("Start HTTP server")

	a.startHttp()
}

func (a *App) ConfigureDatabase() *App {
	appDbConf := a.config.AppConfig.Database
	dsn := "host=" + appDbConf.Host + " user=" + appDbConf.User + " password=" + appDbConf.Password + " dbname=currency port=" + appDbConf.Port + " sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&model.CurrencyRate{})
	if err != nil {
		log.Fatal(err)
	}

	a.db = db

	return a
}

func (a *App) ConfigureRoutes() *App {
	log.Print("Configure routes")
	r := httprouter.New()

	exchangeRateService := service.NewExchangeRateService(a.db)
	getExchangeHandler := get_exchanges.NewHandler(exchangeRateService)
	updateExchangeHandler := update_exchanges.NewHandler(exchangeRateService)

	r.HandlerFunc(http.MethodGet, "/exchange", getExchangeHandler.ExchangeRate)
	r.HandlerFunc(http.MethodGet, "/convert", getExchangeHandler.Convert)

	//todo: move to scheduler call
	r.HandlerFunc(http.MethodGet, "/rates", updateExchangeHandler.GetCbrExchangeRates)

	updateExchangeHandler.SyncRatesOnStartup()

	a.router = r

	return a
}

func (a *App) startHttp() {
	log.Printf("bind application to host: %s and port: %s", a.config.Listen.BindIP, a.config.Listen.Port)

	var listener net.Listener
	var err error
	listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", a.config.Listen.BindIP, a.config.Listen.Port))
	if err != nil {
		log.Fatal(err)
	}

	c := cors.New(cors.Options{
		AllowedMethods:     []string{http.MethodGet, http.MethodPost},
		AllowedOrigins:     []string{"http://localhost:3000", "http://localhost:8080"},
		AllowCredentials:   true,
		AllowedHeaders:     []string{"Location", "Charset", "Access-Control-Allow-Origin", "Content-Type", "content-type", "Origin", "Access"},
		OptionsPassthrough: true,
		ExposedHeaders:     []string{"Location", "Authorization", "Content-Disposition"},
		Debug:              a.config.AppDebug,
	})

	handler := c.Handler(a.router)

	a.httpServer = &http.Server{
		Handler:      handler,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	if err := a.httpServer.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			log.Fatal("server shutdown")
		default:
			log.Fatal(err)
		}
	}

	err = a.httpServer.Shutdown(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
