package main

import (
	"log"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/clients"
	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/clients/kuber"
	domainLogic "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/domain_logic"
	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/global"
	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/http"
	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/jobs"
	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/pkg/exitmanager"
)

// @title Uffizzi Pro Controller
// @version 1.0
// @description.markdown
// @query.collection.format multi
// @contact.name Uffizzi Pro Support
// @contact.url https://support.uffizzi.com/
// @contact.email admin@uffizzi.cloud
// @securityDefinitions.basic BasicAuth
func main() {
	log.Println("Starting uffizzi controller")

	// Initialize exit manager
	exitMgr := exitmanager.Init()

	// Setup controller
	go setup(exitMgr)

	// Block until gracefully exited
	exitMgr.Wait()

	log.Println("Exiting uffizzi controller")
}

func setup(exitMgr *exitmanager.ExitMgr) {
	defer func() {
		err := recover()
		if err != nil {
			var secondsToSleep time.Duration = 5

			sentry.CurrentHub().Recover(err)
			sentry.Flush(time.Second * secondsToSleep)
		}

		log.Print(err)
	}()

	// Set global variables
	if err := global.Init(getEnv()); err != nil {
		panic(err)
	}

	log.Printf("started env=%v settings=%+v", global.Env, global.Settings)

	// Initialize sentry SDK
	if err := sentry.Init(sentry.ClientOptions{AttachStacktrace: true, Environment: global.Env}); err != nil {
		panic(err)
	}

	log.Println("Sentry: initialized")

	config, err := clients.InitializeKubeConfig()
	if err != nil {
		panic(err)
	}

	kuberClient, err := kuber.NewClient(config)
	if err != nil {
		panic(err)
	}

	log.Println("Kubernetes client: initialized")

	workerСancellationСhan := jobs.Init()

	exitMgr.AddTeardownCallback(func() {
		close(workerСancellationСhan)
	})

	log.Println("jobs: initialized")

	logic := domainLogic.NewLogic(kuberClient)

	log.Println("HTTP server: initializing")

	// Initialize HTTP server
	srv := http.Init(logic, exitMgr)
	defer exitMgr.AddTeardownCallback(srv.Shutdown)
}

func getEnv() string {
	env := os.Getenv("ENV")
	return env
}
