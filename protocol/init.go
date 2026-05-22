package protocol

import (
	"github.com/labstack/echo/v4"
	"gitlab.company.com/projectname/configs"
	"gitlab.company.com/projectname/infrastructure"
	"gitlab.company.com/projectname/internal/core/service"
	"gitlab.company.com/projectname/internal/handler"
	"gitlab.company.com/projectname/internal/handler/common"
	"gitlab.company.com/projectname/internal/repository"
	"gitlab.company.com/projectname/pkg/logx"
	"gorm.io/gorm"
)

// App is the root application struct. All deps flow through here.
type App struct {
	conf    configs.Config
	echo    *echo.Echo
	db      *gorm.DB
	repo    *repository.Repository
	handler *handler.Handler
	service *service.Service
}

// New bootstraps the entire application from config.
func New(conf configs.Config) (*App, error) {
	logx.Init(conf.App.LogLevel)
	log := logx.GetLog()

	db, err := infrastructure.NewPostgresDB(conf.Postgres)
	if err != nil {
		log.WithError(err).Error("failed to connect to postgres")
		return nil, err
	}
	log.Info("postgres connected")

	repo := repository.New(repository.Dependencies{DB: db})

	svc := service.New(service.Dependencies{Config: conf})

	common.SetAppPrefix(conf.App.Prefix)

	h := handler.New(handler.Dependencies{
		HealthService: svc.Health,
	})

	app := &App{
		conf:    conf,
		db:      db,
		repo:    repo,
		handler: h,
		service: svc,
	}

	app.echo = app.newEcho()
	return app, nil
}

// Start begins listening for HTTP requests.
func (a *App) Start() error {
	return a.echo.Start(":" + a.conf.App.Port)
}
