package internal

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/MarySmirnova/news_reader/internal/api"
	"github.com/MarySmirnova/news_reader/internal/config"
	"github.com/MarySmirnova/news_reader/internal/database"
	"github.com/MarySmirnova/news_reader/internal/rss"
	"github.com/chatex-com/process-manager"
	log "github.com/sirupsen/logrus"
)

type Application struct {
	sigChan <-chan os.Signal
	cfg     config.Application
	db      *database.Store
	manager *process.Manager
}

func NewApplication(cfg config.Application) (*Application, error) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	app := &Application{
		sigChan: sigChan,
		cfg:     cfg,
		manager: process.NewManager(),
	}

	if err := app.bootstrap(); err != nil {
		return nil, err
	}

	return app, nil
}

func (a *Application) bootstrap() error {
	// init dependencies
	if err := a.initDatabase(); err != nil {
		return err
	}

	// init workers
	a.bootRSSWorker()
	a.bootServerWorker()

	return nil
}

func (a *Application) initDatabase() error {
	db, err := database.NewPostgresDB(a.cfg.Postgres)
	if err != nil {
		log.WithError(err).Error("database connection error")
		return err
	}

	log.Info("database connection established")
	a.db = db
	return nil
}

func (a *Application) bootRSSWorker() {
	rssParser := rss.NewNewsParser(a.cfg.RSS, a.db)
	rssWorker := process.NewCallbackWorker("rss", rssParser.Start)
	a.manager.AddWorker(rssWorker)
}

func (a *Application) bootServerWorker() {
	server := api.New(a.cfg.API, a.db)
	serverWorker := process.NewServerWorker("api", server.GetHTTPServer())
	a.manager.AddWorker(serverWorker)
}

func (a *Application) Run() {
	a.manager.StartAll()
	a.registerShutdown()
}

func (a *Application) registerShutdown() {
	go func(manager *process.Manager) {
		<-a.sigChan

		manager.StopAll()
	}(a.manager)

	defer a.db.GetPGXPool().Close()

	a.manager.AwaitAll()
}
