package repository

import (
	"database/sql"
	"fmt"
	"testing"
)

// GetNextID gets the next ID of a given table and register its cleanup in the current test
func GetNextID(t *testing.T, table string) int {
	ID := 0
	err := DB.QueryRow(fmt.Sprintf("SELECT NEXTVAL('%s_id_seq') + 1", table)).Scan(&ID)
	if err != nil {
		t.Errorf("error getting next %s id %s", table, err)
	}
	t.Cleanup(func() {
		DB.Exec(fmt.Sprintf("DELETE FROM %s_status WHERE %s_id = $1", table, table), ID)
		DB.Exec(fmt.Sprintf("DELETE FROM %s_owners WHERE %s_id = $1", table, table), ID)
		_, err := DB.Exec(fmt.Sprintf("DELETE FROM %s WHERE id = $1", table), ID)
		if err != nil {
			t.Errorf("error deleting id of %s %s", table, err)
		}
	})
	return ID
}

// Tx executes the given func in a tx
func Tx(t *testing.T, f func(tx *sql.Tx)) {
	tx, err := DB.Begin()
	if err != nil {
		t.Fatalf("error starting tx %s", err)
	}
	defer tx.Commit()
	f(tx)
}
