package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/getsentry/sentry-go"
)

func sentryInit() {
	err := sentry.Init(sentry.ClientOptions{})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
}

func sentryFlush() {
	sentry.Flush(2 * time.Second)
}

func sentryRecover() {
	if err := recover(); err != nil {
		tracedErr := errors.New(fmt.Sprint(err))
		sentry.CaptureException(tracedErr)
		panic(tracedErr)
	}
}
