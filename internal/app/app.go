package app

import (
	"Currency/internal/config"
	"Currency/internal/domain/rate/handlers/get_exchanges"
	"Currency/internal/domain/rate/handlers/update_exchanges"
	"Currency/internal/domain/rate/model"
	"Currency/internal/domain/rate/service"
	"github.com/julienschmidt/httprouter"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type App struct {
	config *config.Config
	router *httprouter.Router
	db     *gorm.DB
	DI     *ServiceLocator
	server *Server
}

func NewKernel(config *config.Config) App {
	return App{
		config: config,
		server: NewServer(),
		DI:     NewServiceLocator(),
	}
}

func (a *App) End() {
	a.server.ConfigureAndRun(a.config.Listen.BindIP, a.config.Listen.Port, a.config.AppDebug, a.router)
}

func (a *App) ConfigureDatabase() *App {
	log.Println("Start configure database connection")
	db, err := gorm.Open(postgres.Open(a.config.GetDsn()), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	a.db = db

	return a
}

func (a *App) ConfigureServiceLocator() *App {
	exchangeRateService := service.NewExchangeRateService(a.db)

	a.DI.Set(get_exchanges.ExchangeRateHandlerTag, get_exchanges.NewHandler(exchangeRateService))
	a.DI.Set(update_exchanges.UpdateRateHandlerTag, update_exchanges.NewHandler(exchangeRateService))

	return a
}

func (a *App) ConfigureRoutes() *App {
	log.Print("Configure routes")
	r := httprouter.New()

	r.GET("/exchange", a.DI.Get(get_exchanges.ExchangeRateHandlerTag).(*get_exchanges.GetExchangeRateHandler).ExchangeRate)
	r.GET("/convert", a.DI.Get(get_exchanges.ExchangeRateHandlerTag).(*get_exchanges.GetExchangeRateHandler).Convert)

	//todo: move to scheduler call
	r.HandlerFunc(http.MethodGet, "/rates", a.DI.Get(update_exchanges.UpdateRateHandlerTag).(*update_exchanges.UpdateExchangeHandler).GetCbrExchangeRates)
	//r.HandlerFunc(http.MethodGet, "/rush", a.DI.Get(update_exchanges.UpdateRateHandlerTag).(*update_exchanges.UpdateExchangeHandler).RushRates)

	a.router = r

	return a
}

func (a *App) AfterInitializationEvents() *App {
	log.Println("Run after initialization events")

	err := a.db.AutoMigrate(&model.CurrencyRate{})
	if err != nil {
		log.Fatal(err)
	}

	if a.config.AppConfig.SyncRatesAfterStartup == true {
		a.DI.Get(update_exchanges.UpdateRateHandlerTag).(*update_exchanges.UpdateExchangeHandler).SyncRatesOnStartup()
	}

	return a
}
