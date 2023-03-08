package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	domain "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/domain_logic"
	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/global"

	sentryHTTP "github.com/getsentry/sentry-go/http"
	"github.com/gorilla/mux"
	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/http/middlewares"
	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/pkg/exitmanager"
)

const (
	ShutdownTimeout = 5 * time.Second
)

var (
	domainLogic *domain.Logic
)

type Server struct {
	srv *http.Server
}

func Init(logic *domain.Logic, exitMgr *exitmanager.ExitMgr) *Server {
	// Set package variables
	domainLogic = logic

	// Create sentryRecoveryHandler
	sentryRecoveryHandler := sentryHTTP.New(sentryHTTP.Options{
		Repanic: true,
	})

	// Create router
	r := mux.NewRouter()

	// Create handlers
	h := &Handlers{}

	// Draw routes
	drawRoutes(r, h)

	// Setup middleware
	r.Use(middlewares.Logging)
	r.Use(middlewares.Authentication)

	// Set the server address
	addr := fmt.Sprintf(":%s", global.Settings.ControllerPort)

	// Create http.Server instance
	srv := &http.Server{
		Addr:              addr,
		Handler:           sentryRecoveryHandler.Handle(r),
		ReadHeaderTimeout: 6 * time.Second, //nolint: gomnd
	}

	// Listen and serve
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
			exitMgr.ServerError(err)

			return
		}
	}()

	return &Server{
		srv: srv,
	}
}

func (s *Server) Shutdown() {
	log.Println("Gracefully shutting down the http server")

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()

	// Gracefully shutdown the server
	err := s.srv.Shutdown(ctx)
	if err == nil {
		log.Println("Server gracefully shut down")
	} else {
		log.Printf("Server gracefully shut down with err: #+%v\n", err)
	}
}
