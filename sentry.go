package qoingohelper

import (
	"log"
	"time"

	"github.com/getsentry/sentry-go"
)

type SentryData struct {
	Dsn             string
	Environment     string
	Release         string
	Debug           bool
	TraceSampleRate float64
	Event           *sentry.Event
}

var Sentries *SentryData

func NewSentryData() *SentryData {
	return new(SentryData)
}

func InitSentry(Dsn, Environment, Release string, Debug bool) {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              Dsn,
		Environment:      Environment,
		Release:          Release,
		Debug:            Debug,
		AttachStacktrace: true,
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	defer sentry.Flush(10 * time.Second)
}

func SendSentryMessage(message string, service, module, function string) {
	sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetLevel(sentry.LevelError)
		scope.AddBreadcrumb(&sentry.Breadcrumb{
			Type:     "Info",
			Category: "Information",
			Message:  "Details of message stack",
			Data: map[string]interface{}{
				"Service":  service,
				"Module":   module,
				"Function": function,
			},
		}, 5)
		sentry.CaptureMessage(message)
	})
}

func SendSentryError(err error, service, module, function string) {
	sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetLevel(sentry.LevelError)
		scope.AddBreadcrumb(&sentry.Breadcrumb{
			Type:     "Info",
			Category: "Information",
			Message:  "Details of error stack",
			Data: map[string]interface{}{
				"Service":  service,
				"Module":   module,
				"Function": function,
			},
		}, 5)
		sentry.CaptureException(err)
	})
}

func SendSentryEvent(event *sentry.Event, service, module, function string) {
	sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetLevel(sentry.LevelInfo)
		scope.AddBreadcrumb(&sentry.Breadcrumb{
			Type:     "Info",
			Category: "Information",
			Message:  "Details of event stack",
			Data: map[string]interface{}{
				"Service":  service,
				"Module":   module,
				"Function": function,
			},
		}, 5)
		sentry.CaptureEvent(event)
	})
}
