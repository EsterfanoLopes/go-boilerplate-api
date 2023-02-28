package repository

import (
	"database/sql"
	"fmt"
	"go-boilerplate/common"
	"sync/atomic"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
)

const pgUniqueViolationSQLState = "23505"

var (
	dbHost     = common.Config.Get("dbHost")
	dbPort     = common.Config.Get("dbPort")
	dbUser     = common.Config.Get("dbUser")
	dbPassword = common.Config.Get("dbPassword")
	dbName     = common.Config.Get("dbName")
	dbTimeout  = common.Config.GetInt("dbTimeoutSeconds")

	// Psq query builder instance
	Psq = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	// DB is a db instance
	DB *sql.DB

	dbReady = int32(0)
)

// OrderDirection of the given order type
type OrderDirection int

const (
	Asc OrderDirection = iota
	Desc
)

var orderDirectionValues = [...]string{
	"asc",
	"desc",
}

func (s OrderDirection) String() string {
	return orderDirectionValues[s]
}

// OrderDirectionValueOf converts an order direction string into an order direction enum type
func OrderDirectionValueOf(v string) (OrderDirection, error) {
	for i, value := range orderDirectionValues {
		if value == v {
			return OrderDirection(i), nil
		}
	}
	return 0, fmt.Errorf("unknown order direction value %s", v)
}

func dbIsReady() {
	atomic.StoreInt32(&dbReady, 1)
}

func isDBReady() bool {
	return atomic.LoadInt32(&dbReady) == 1
}

func setupDB() error {
	if isDBReady() {
		return nil
	}

	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable connect_timeout=%d statement_timeout=%ds", dbHost, dbPort, dbUser, dbPassword, dbName, dbTimeout, dbTimeout)

	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		return err
	}
	db.SetMaxIdleConns(common.Config.GetInt("dbMinConnections"))
	db.SetMaxOpenConns(common.Config.GetInt("dbMaxConnections"))
	db.SetConnMaxLifetime(time.Duration(30 * time.Minute))

	DB = db
	err = db.Ping()
	if err != nil {
		return err
	}
	dbIsReady()

	return nil
}

func dbHealthcheck() (int, error) {
	result := 0
	err := DB.QueryRow("SELECT 1").Scan(&result)
	return result, err
}

// CloseRows closes the given rows
func CloseRows(rows *sql.Rows) {
	if rows != nil {
		rows.Close()
	}
}

// IsUniqueConstraintViolation checks if the given error is a sql constraint violation
func IsUniqueConstraintViolation(err error) bool {
	if pgerr, ok := err.(*pgconn.PgError); ok {
		return pgerr.SQLState() == pgUniqueViolationSQLState
	}
	return false
}

// GetOrderByStatusClause given a table alias and slice of status strings generates an order by clause
func GetOrderByStatusClause(alias string, orderByStatusWeights [][]string) string {
	clause := fmt.Sprintf("case %s.last_status ", alias)
	for weight, statuses := range orderByStatusWeights {
		for _, status := range statuses {
			clause += fmt.Sprintf(`when '%s' then %d `, status, weight)
		}
	}
	clause += "else 99 end"
	return clause
}
