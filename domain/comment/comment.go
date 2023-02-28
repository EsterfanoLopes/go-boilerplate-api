// Package comment holds all comment related stuff, a comment can be done by an owner in any advertiser applications (eg CRM)
package comment

import (
	"encoding/json"
	"fmt"
	"go-boilerplate/common"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// Type of comment
type Type int

const (
	// TypeNone zero value for this enum
	TypeNone Type = iota
	// Lead comment
	Lead
	// Schedule comment
	Schedule
	// Negotiation comment
	Negotiation
	// Credit comment
	Credit
	// Transaction comment
	Transaction
)

var typeValues = [...]string{
	"",
	"LEAD",
	"SCHEDULE",
	"NEGOTIATION",
	"CREDIT",
	"TRANSACTION",
}

func (t Type) String() string {
	return typeValues[t]
}

// MarshalJSON marshals the enum as a quoted json string
func (t Type) MarshalJSON() ([]byte, error) {
	return common.QuotedStringBytes(t.String()), nil
}

// UnmarshalJSON unmarshals a quoted json string to the enum value
func (t *Type) UnmarshalJSON(b []byte) error {
	x := ""
	err := json.Unmarshal(b, &x)
	if err != nil {
		return err
	}
	value, err := TypeValueOf(x)
	if err != nil {
		return err
	}
	*t = value
	return nil
}

// TypeValueOf converts a comment type value into a comment type
func TypeValueOf(v string) (Type, error) {
	for i, value := range typeValues {
		if value == v {
			return Type(i), nil
		}
	}
	return 0, fmt.Errorf("unknown comment type value %s", v)
}

// Comment done by a user about some entity
type Comment struct {
	ID           int       `json:"id"`
	Type         Type      `json:"type"`
	Description  string    `json:"description"`
	AdvertiserID string    `json:"advertiserId"`
	AccountID    string    `json:"accountId"`
	ListingID    string    `json:"listingId"`
	Updated      bool      `json:"updated"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// Validate the given comment
func (c Comment) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Type, validation.Required),
		validation.Field(&c.Description, validation.Required),
		validation.Field(&c.AdvertiserID, validation.Required, is.UUID),
		validation.Field(&c.AccountID, validation.Required, is.UUID),
		validation.Field(&c.ListingID, validation.Required, is.Digit),
	)
}
