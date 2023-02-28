// Package repository is responsible for all interations with data commit technologies
package repository

import (
	"errors"
	"go-boilerplate/common"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	awstrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/aws/aws-sdk-go/aws"
)

var (
	// ErrNotFound is when no resource was found
	ErrNotFound = errors.New("resource not found")

	AWSSession *session.Session
)

// Setup prepares the entire layer to be used
func Setup() error {
	setupHTTP()

	err := setupAWSSession()
	if err != nil {
		return err
	}

	err = setupDB()
	if err != nil {
		return err
	}

	return nil
}

// HealthcheckResponse response of healthcheck process
type HealthcheckResponse struct {
	DB   string
	HTTP string
}

const ok = "OK"

// Healthy checks if all dependencies are healthy
func (h HealthcheckResponse) Healthy() bool {
	return h.DB == ok && h.HTTP == ok
}

// Healthcheck checks if this layer is healthy
func Healthcheck() HealthcheckResponse {
	_, dbErr := dbHealthcheck()

	httpErr := httpHealthcheck()

	DB := ok
	if dbErr != nil {
		DB = dbErr.Error()
	}

	HTTP := ok
	if httpErr != nil {
		HTTP = httpErr.Error()
	}

	return HealthcheckResponse{
		DB:   DB,
		HTTP: HTTP,
	}
}

func setupAWSSession() error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(common.Config.Get("awsRegion")),
	})
	if err != nil {
		return err
	}

	AWSSession = awstrace.WrapSession(sess, awstrace.WithServiceName("go-boilerplate-aws"))

	return nil
}
