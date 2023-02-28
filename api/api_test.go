package api_test

import (
	"fmt"
	"go-boilerplate/api"
	"go-boilerplate/repository"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	err := repository.Setup()
	if err != nil {
		fmt.Printf("error starting api tests %s \n", err)
		os.Exit(-1)
	}
	go api.Setup()
	time.Sleep(1 * time.Second)
	os.Exit(m.Run())
}
