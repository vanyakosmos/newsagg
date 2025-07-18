package internal

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/getsentry/sentry-go"
)

func SentryInit() {
	err := sentry.Init(sentry.ClientOptions{})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
}

func SentryFlush() {
	sentry.Flush(2 * time.Second)
}

func SentryRecover() {
	if err := recover(); err != nil {
		tracedErr := errors.New(fmt.Sprint(err))
		sentry.CaptureException(tracedErr)
		panic(tracedErr)
	}
}
