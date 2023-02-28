package main

import (
	"go-boilerplate/cmd"
	"go-boilerplate/common"
	"go-boilerplate/repository"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func main() {
	if common.Config.Get("datadogEnabled") == "true" {
		tracer.Start(tracer.WithRuntimeMetrics())
		defer tracer.Stop()
	}

	if common.Config.Get("environment") == "prod" {
		defer sentry.Flush(time.Second * 5)
		sentry.Init(sentry.ClientOptions{
			Dsn:   common.Config.Get("sentryDsn"),
			Debug: false,
		})
	}

	err := repository.Setup()
	if err != nil {
		common.HandleError("error on repository setup", err)
		sentry.Flush(time.Second * 2)
		os.Exit(-1)
	}

	cmd.Execute()
}
