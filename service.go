package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/agungdwiprasetyo/agungkiki-backend/config"
	"github.com/agungdwiprasetyo/agungkiki-backend/helper"
	"github.com/agungdwiprasetyo/agungkiki-backend/middleware"
	"github.com/agungdwiprasetyo/agungkiki-backend/src/presenter"
	"github.com/agungdwiprasetyo/agungkiki-backend/src/repository"
	"github.com/agungdwiprasetyo/agungkiki-backend/src/usecase"
	tokenModule "github.com/agungdwiprasetyo/agungkiki-backend/token"
	"github.com/labstack/echo"
)

// Service main model
type Service struct {
	conf  config.Config
	token *tokenModule.Token
}

// NewService create new service
func NewService(conf config.Config) *Service {
	token := tokenModule.New(conf.LoadPrivateKey(), conf.LoadPublicKey(), 12*time.Hour)

	service := new(Service)
	service.conf = conf
	service.token = token
	return service
}

// ServeHTTP service
func (serv *Service) ServeHTTP(port int) {
	repositoryDecorator := repository.NewRepository(serv.conf.LoadDB())
	bearerMiddleware := middleware.Bearer(serv.token)
	uc := usecase.NewInvitationUsecase(serv.token, repositoryDecorator)
	invitationPresenter := presenter.NewInvitationPresenter(uc, bearerMiddleware)

	app := echo.New()
	app.Use(middleware.Recover(), middleware.SetCORS(), middleware.Logger())

	app.GET("/", func(c echo.Context) error {
		response := helper.NewHTTPResponse(http.StatusOK, "ok")
		return response.SetResponse(c)
	}, bearerMiddleware)

	invitationGroup := app.Group("/invitation")
	invitationPresenter.Mount(invitationGroup)

	appPort := fmt.Sprintf(":%d", port)
	if err := app.Start(appPort); err != nil {
		log.Fatal(err)
	}
}
