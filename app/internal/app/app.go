package app

import (
	"Currency/internal/config"
	"Currency/internal/domain/model"
	"Currency/internal/domain/service"
	"Currency/internal/domain/use_case/get_exchanges"
	"Currency/internal/domain/use_case/update_exchanges"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"gorm.io/driver/sqlite"
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

func NewKernel(config *config.Config) (App, error) {
	db := configureDatabase()
	return App{
		config: config,
		db:     db,
		router: configureRoutes(db),
	}, nil
}

func (a *App) Run() {
	log.Print("Start HTTP server")

	a.startHttp()
}

func configureDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("currency.db"), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&model.CurrencyRate{})
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func configureRoutes(db *gorm.DB) *httprouter.Router {
	log.Print("Configure routes")
	r := httprouter.New()

	srv := service.NewExchangeRateService(db)

	r.HandlerFunc(http.MethodGet, "/exchange", func(writer http.ResponseWriter, request *http.Request) {
		exchangeRate := get_exchanges.NewHandler(srv).ExchangeRate(request)

		marshalled, _ := json.Marshal(exchangeRate)
		_, err := writer.Write(marshalled)
		if err != nil {
			return
		}
	})

	r.HandlerFunc(http.MethodGet, "/rates", func(writer http.ResponseWriter, request *http.Request) {
		update_exchanges.
			NewHandler(srv).
			ExchangeCbrRates(request)
	})
	return r
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
