package utils

import (
	"context"
	"log"

	"cloud.google.com/go/errorreporting"
)

// ErrorService structure of a custom error that contains code and an error
type ErrorService struct {
	Code int
	Err  error
}

// ErrorReporter structure of the ErrorReporter instance
type ErrorReporter struct {
	Reporter *errorreporting.Client
}

// ErrorReport error reporter client
var ErrorReport *ErrorReporter

// InitErrorReporting initialize instance of ErrorReport
func InitErrorReporting(projectID string) {
	ctx := context.Background()

	errorReporting, err := errorreporting.NewClient(ctx, projectID, errorreporting.Config{
		ServiceName: "bitcoin-functions",
		OnError: func(err error) {
			log.Printf("Could not log error: %v", err)
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	ErrorReport = &ErrorReporter{
		Reporter: errorReporting,
	}
}

// LogAndPrintError log and print error into GCP journals
func (e *ErrorReporter) LogAndPrintError(err error) {
	e.Reporter.Report(errorreporting.Entry{
		Error: err,
	})

	log.Print(err)
}
