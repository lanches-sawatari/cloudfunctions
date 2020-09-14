package error_utility

import (
	"cloud.google.com/go/errorreporting"
	"context"
	"log"
)

var errClient *errorreporting.Client

type ErrMsg struct {
	msg string
	err error
}

func (erm ErrMsg) Log() {
	log.Printf("[ERROR] %s: %v", erm.msg, erm.err)
}

func New(gcpProjectId string, serviceName string) {
	// err is pre-declared to avoid shadowing client.
	var err error
	errClient, err = errorreporting.NewClient(context.Background(), gcpProjectId, errorreporting.Config{
		ServiceName: serviceName,
		OnError: func(err error) {
			ErrMsg{"Could not log to ErrorReporting", err}.Log()
		},
	})
	if err != nil {
		ErrMsg{"init() failed for ErrorClient", err}.Log()
	}
}

func (erm ErrMsg) Empty() bool {
	return erm.err == nil
}

func (erm ErrMsg) Report() {
	if erm.Empty() {
		return
	}
	erm.Log()
	errClient.Report(errorreporting.Entry{Error: erm.err})
	errClient.Flush()
}

func (erm ErrMsg) ReportOnly() {
	if erm.Empty() {
		return
	}
	errClient.Report(errorreporting.Entry{Error: erm.err})
	errClient.Flush()
}
