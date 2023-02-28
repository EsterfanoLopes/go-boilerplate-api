// Package common deals with all common functions and configurations from the APIs
package common

import (
	"fmt"
	"math/rand"
	"runtime/debug"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/golang-jwt/jwt"
	"github.com/olxbr/ligeiro/envcfg"
	"github.com/olxbr/ligeiro/logger"
)

// ADay time in hours
const ADay = 24 * time.Hour

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Config holds the app configuration
var Config = envcfg.Load(envcfg.Map{
	// General
	"ENVIRONMENT":                   "dev",
	"SENTRY_DSN":                    "",
	"DATADOG_ENABLED":               "false",
	"BASE_URL":                      "http://localhost:9000",
	"VIVA_REAL_PORTAL_HOST":         "www.vivareal.com.br",
	"ZAP_PORTAL_HOST":               "www.zapimoveis.com.br",
	"MY_ACCOUNT_API_VIVA_REAL_HOST": "my-account-api.vivareal.com.br",
	"MY_ACCOUNT_API_ZAP_HOST":       "my-account-api.zapimoveis.com.br",
	"MEDIA_UPLOAD_JWT_SECRET":       "local-secret",
	"MEDIA_UPLOAD_JWT_EXP_HOURS":    "1",
	"WORKDAY_START_HOUR":            "8",
	"WORKDAY_PERIOD_IN_HOURS":       "10",

	// DB Config
	"DB_HOST":            "localhost",
	"DB_PORT":            "5432",
	"DB_USER":            "user",
	"DB_PASSWORD":        "pass",
	"DB_NAME":            "database",
	"DB_MIN_CONNECTIONS": "10",
	"DB_MAX_CONNECTIONS": "10",
	"DB_TIMEOUT_SECONDS": "2",

	// Http Server Config
	"HTTP_SERVER_READ_TIMEOUT_SECONDS":  "600",
	"HTTP_SERVER_WRITE_TIMEOUT_SECONDS": "600",
	"HTTP_HEALTHCHECK_ENDPOINT":         "",

	// AWS Config
	"AWS_REGION": "us-east-1",

	// HTTP Default Configurations
	"HTTP_TIMEOUT_SECONDS": "600",
	"HTTP_MIN_CONNECTIONS": "10",
	"HTTP_MAX_CONNECTIONS": "10",
	"HTTP_RESPONSE_DEBUG":  "false",
	"HTTP_MAX_RETRIES":     "1",
})

// Logger is the default app logger
var Logger = logger.WithFields(logger.Fields{
	"application": "go-boilerplate",
	"environment": Config.Get("environment"),
})

// ToTime converts a string in ISO format to a time type
func ToTime(v string) (time.Time, error) {
	return time.Parse(time.RFC3339Nano, v)
}

// HandleError handles errors sending them to sentry and logging
func HandleError(message string, err error) {
	Logger.WithFields(logger.Fields{
		"error": err,
		"stack": string(debug.Stack()),
	}).Error(message)
	sentry.CaptureException(err)
}

// CreateJWTToken given a set o claims and a secret create a jwt token
func CreateJWTToken(claims jwt.MapClaims, secret string) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}

// ParseJWTToken given a jwt token and secret, parse the token
func ParseJWTToken(tokenStr, secret string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
}

// QuotedStringBytes given a string returns its bytes quoted
func QuotedStringBytes(v string) []byte {
	return []byte(fmt.Sprint(`"`, v, `"`))
}
