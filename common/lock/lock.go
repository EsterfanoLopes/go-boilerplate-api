// Package lock holds distributed lock utility
package lock

import "go-boilerplate/repository"

// Key to lock
type Key int

const (
	// AdvertiserPortalsConfigUpdate key to lock when advertisers config file is updated in S3
	AdvertiserPortalsConfigUpdate Key = iota
)

// Acquire lock of given key
func Acquire(key Key) error {
	_, err := repository.DB.Exec(`SELECT pg_advisory_lock($1)`, key)
	if err != nil {
		return err
	}

	return nil
}

// Release lock of given key
func Release(key Key) error {
	_, err := repository.DB.Exec(`SELECT pg_advisory_unlock($1)`, key)
	if err != nil {
		return err
	}

	return nil
}
