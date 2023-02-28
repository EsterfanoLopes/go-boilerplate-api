// Package comment holds data access logic of comments
package comment

import (
	sql "database/sql"
	"go-boilerplate/common/pagination"
	"go-boilerplate/domain/comment"
	"go-boilerplate/repository"
	"time"

	sq "github.com/Masterminds/squirrel"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

const (
	columns = `
		id,
		description,
		type,
		updated,
		account_id,
		advertiser_id,
		listing_id,
		owner,
		created_at,
		updated_at
	`
)

var (
	instance = &repositoryImpl{}
)

// Repository to enable this repository to be mocked
type Repository interface {
	// Insert a comment
	Insert(tx *sql.Tx, cmt comment.Comment) (int, error)
	// Update a comment
	Update(tx *sql.Tx, cmt comment.Comment) error
	// FindByID a comment
	FindByID(tx *sql.Tx, ID int) (comment.Comment, error)
	// Find comments by a given query
	Find(tx *sql.Tx, q Query, p pagination.Pagination) ([]comment.Comment, error)
	// Count comments by a given query
	Count(tx *sql.Tx, q Query) (int, error)
	// Delete a comment
	Delete(tx *sql.Tx, ID int) error
}

type repositoryImpl struct{}

// Get this repository instance
func Get() Repository {
	return instance
}

// Query possible values to find comments
type Query struct {
	AccountID    string
	AdvertiserID string
	ListingID    string
}

// Validate validates negotiation query
func (q Query) Validate() error {
	return validation.ValidateStruct(&q,
		validation.Field(&q.AccountID, validation.Required, is.UUID),
		validation.Field(&q.AdvertiserID, validation.Required, is.UUID),
		validation.Field(&q.ListingID, is.Digit),
	)
}

func (r *repositoryImpl) Insert(tx *sql.Tx, cmt comment.Comment) (int, error) {
	insert, values, err := repository.Psq.Insert("comment").Columns(`
		description,
		type,
		updated,
		account_id,
		advertiser_id,
		listing_id,
		owner,
		created_at,
		updated_at
	`).Values(
		cmt.Description,
		cmt.Type.String(),
		false,
		cmt.AccountID,
		cmt.AdvertiserID,
		cmt.ListingID,
		time.Now(),
		time.Now(),
	).Suffix("RETURNING id").ToSql()
	if err != nil {
		return 0, err
	}

	ID := 0
	err = tx.QueryRow(insert, values...).Scan(&ID)
	if err != nil {
		return 0, err
	}

	return ID, nil
}

func (r *repositoryImpl) Update(tx *sql.Tx, cmt comment.Comment) error {
	update, values, err := repository.Psq.Update("comment").
		Set("updated_at", time.Now()).
		Set("updated", true).
		Set("description", cmt.Description).
		Where(sq.Eq{"id": cmt.ID}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(update, values...)
	if err != nil {
		return err
	}

	return nil
}

func (r *repositoryImpl) FindByID(tx *sql.Tx, ID int) (comment.Comment, error) {
	query, values, err := repository.Psq.Select(columns).From("comment").Where(sq.Eq{"id": ID}).ToSql()
	if err != nil {
		return comment.Comment{}, err
	}

	var rows *sql.Rows
	if tx == nil {
		rows, err = repository.DB.Query(query, values...)
	} else {
		rows, err = tx.Query(query, values...)
	}
	if err != nil {
		return comment.Comment{}, err
	}
	defer repository.CloseRows(rows)

	if rows.Next() {
		result, err := r.scanRow(rows)
		if err != nil {
			return comment.Comment{}, err
		}

		return result, nil
	}

	return comment.Comment{}, repository.ErrNotFound
}

func (r *repositoryImpl) Find(tx *sql.Tx, q Query, p pagination.Pagination) ([]comment.Comment, error) {
	query, values, err := r.commentSelect(columns, q).OrderBy("created_at DESC").ToSql()
	if err != nil {
		return nil, err
	}

	var rows *sql.Rows
	if tx == nil {
		rows, err = repository.DB.Query(p.PaginateQuery(query), values...)
	} else {
		rows, err = tx.Query(p.PaginateQuery(query), values...)
	}
	if err != nil {
		return nil, err
	}
	defer repository.CloseRows(rows)

	results := []comment.Comment{}
	for rows.Next() {
		result, err := r.scanRow(rows)
		if err != nil {
			return nil, err
		}

		results = append(results, result)
	}

	return results, nil
}

func (r *repositoryImpl) Count(tx *sql.Tx, q Query) (int, error) {
	countQ, values, err := r.commentSelect("count(1)", q).ToSql()
	if err != nil {
		return 0, err
	}

	count := 0
	if tx == nil {
		err = repository.DB.QueryRow(countQ, values...).Scan(&count)
	} else {
		err = tx.QueryRow(countQ, values...).Scan(&count)
	}
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *repositoryImpl) commentSelect(columns string, q Query) sq.SelectBuilder {
	sqq := repository.Psq.Select(columns).From("comment")
	if q.ListingID != "" {
		sqq = sqq.Where(sq.Eq{"listing_id": q.ListingID})
	}
	if q.AccountID != "" {
		sqq = sqq.Where(sq.Eq{"account_id": q.AccountID})
	}
	if q.AdvertiserID != "" {
		sqq = sqq.Where(sq.Eq{"advertiser_id": q.AdvertiserID})
	}

	return sqq
}

func (r *repositoryImpl) scanRow(rows *sql.Rows) (comment.Comment, error) {
	result := comment.Comment{}
	tpValue := ""
	onrBytes := []byte{}
	err := rows.Scan(
		&result.ID,
		&result.Description,
		&tpValue,
		&result.Updated,
		&result.AccountID,
		&result.AdvertiserID,
		&result.ListingID,
		&onrBytes,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		return result, err
	}

	tp, err := comment.TypeValueOf(tpValue)
	if err != nil {
		return result, err
	}
	result.Type = tp

	return result, nil
}

func (r *repositoryImpl) Delete(tx *sql.Tx, ID int) error {
	delete, values, err := repository.Psq.Delete("comment").Where(sq.Eq{"id": ID}).ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(delete, values...)
	if err != nil {
		return err
	}

	return nil
}
