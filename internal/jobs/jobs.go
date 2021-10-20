package jobs

import (
	"log"

	"github.com/getsentry/sentry-go"
	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/global"
)

var (
	deploymentJobs chan func() error
)

func Init() chan struct{} {
	cancelChan := make(chan struct{})
	deploymentJobs = make(chan func() error, global.Settings.MaxCountProcessDeploymentJobs)

	go Worker(deploymentJobs, cancelChan)

	return cancelChan
}

func Worker(jobs chan func() error, cancelChan chan struct{}) {
	for {
		select {
		case job := <-jobs:
			err := job()
			if err != nil {
				SendError(err)
			}
		case <-cancelChan:
			return
		}
	}
}

func AddDeploymentJob(job func() error) {
	deploymentJobs <- job
}

func SendError(err error) {
	localHub := sentry.CurrentHub().Clone()

	log.Printf("job error: %v\n", err)

	localHub.CaptureException(err)
}
