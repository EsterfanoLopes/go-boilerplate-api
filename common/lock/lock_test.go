package lock_test

import (
	"fmt"
	"go-boilerplate/common/lock"
	"go-boilerplate/repository"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	err := repository.Setup()
	if err != nil {
		fmt.Printf("error starting lock tests %s \n", err)
		os.Exit(-1)
	}

	os.Exit(m.Run())
}

func TestAcquireRelease(t *testing.T) {
	err := lock.Acquire(34)
	if err != nil {
		t.Errorf("error acquiring lock %s", err)
		return
	}

	err = lock.Release(34)
	if err != nil {
		t.Errorf("error releasing lock %s", err)
	}
}
