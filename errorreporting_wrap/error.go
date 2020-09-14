package errorreporting_wrap

import (
	"cloud.google.com/go/errorreporting"
	"context"
	"log"
)

type Client struct {
	errorReportingClient *errorreporting.Client
}

func NewClient(gcpProjectId string) (*Client, error) {
	erpClient, err := errorreporting.NewClient(context.Background(), gcpProjectId, errorreporting.Config{
		//ServiceName: serviceName,
		OnError: func(err error) {
			log.Printf("Could not log to ErrorReporting, %v", err)
		},
	})
	if err != nil {
		return nil, err
	}
	return &Client{errorReportingClient: erpClient}, nil
}

func (c Client) Log(msg string, err error) {
	log.Printf("[ERROR] %s: %v", msg, err)
}

func (c Client) Report(msg string, err error) {
	c.Log(msg, err)
	c.errorReportingClient.Report(errorreporting.Entry{Error: err})
	c.errorReportingClient.Flush()
}

func (c Client) ReportWithoutLog(err error) {
	c.errorReportingClient.Report(errorreporting.Entry{Error: err})
	c.errorReportingClient.Flush()
}
