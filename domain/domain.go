// Package domain keeps generic structs who could be used by any other domain like package
package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-boilerplate/common"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/golang-jwt/jwt"
)

var (
	ErrInvalidDocumentToken = errors.New("invalid document token")
	fredoVivaRealURL        = common.Config.Get("fredoVivaRealUrl")
	fredoZapURL             = common.Config.Get("fredoZapUrl")
	vivaRealPortalHost      = common.Config.Get("vivaRealPortalHost")
	zapPortalHost           = common.Config.Get("zapPortalHost")
)

/* Transaction */

// TransactionType possible transaction types
type TransactionType int

const (
	// TransactionTypeNone zero value for this enum
	TransactionTypeNone TransactionType = iota
	// Rental when the transaction is for the rental of a listing
	Rental
	// Sale when the transaction if for the sale of a listing
	Sale
)

var transactionTypeValues = [...]string{
	"",
	"RENTAL",
	"SALE",
}

func (t TransactionType) String() string {
	return transactionTypeValues[t]
}

// MarshalJSON marshals the enum as a quoted json string
func (t TransactionType) MarshalJSON() ([]byte, error) {
	return common.QuotedStringBytes(t.String()), nil
}

// UnmarshalJSON unmarshals a quoted json string to the enum value
func (t *TransactionType) UnmarshalJSON(b []byte) error {
	x := ""
	err := json.Unmarshal(b, &x)
	if err != nil {
		return err
	}
	value, err := TransactionTypeValueOf(x)
	if err != nil {
		return err
	}
	*t = value
	return nil
}

// TransactionTypeValueOf converts a transaction type value into a transaction type
func TransactionTypeValueOf(v string) (TransactionType, error) {
	for i, value := range transactionTypeValues {
		if value == v {
			return TransactionType(i), nil
		}
	}
	return 0, fmt.Errorf("unknown transaction type value %s", v)
}

// UsageType possible usage types
type UsageType int

const (
	// UsageTypeNone zero value for this enum
	UsageTypeNone UsageType = iota
	// Residential usage type
	Residential
	// Commercial usage type
	Commercial
)

var usageTypeValues = [...]string{
	"",
	"RESIDENTIAL",
	"COMMERCIAL",
}

func (t UsageType) String() string {
	return usageTypeValues[t]
}

// MarshalJSON marshals the enum as a quoted json string
func (t UsageType) MarshalJSON() ([]byte, error) {
	return common.QuotedStringBytes(t.String()), nil
}

// UnmarshalJSON unmarshals a quoted json string to the enum value
func (t *UsageType) UnmarshalJSON(b []byte) error {
	x := ""
	err := json.Unmarshal(b, &x)
	if err != nil {
		return err
	}
	value, err := UsageTypeValueOf(x)
	if err != nil {
		return err
	}
	*t = value
	return nil
}

// UsageTypeValueOf converts a usage type value into a usage type
func UsageTypeValueOf(v string) (UsageType, error) {
	for i, value := range usageTypeValues {
		if value == v {
			return UsageType(i), nil
		}
	}
	return 0, fmt.Errorf("unknown usage type value %s", v)
}

/* Week Day */

// IsValidWeekday tells if a weekday is valid uppercase needed
func IsValidWeekday(weekday string) bool {
	for day := time.Sunday; day <= time.Saturday; day++ {
		if weekday == strings.ToUpper(day.String()) {
			return true
		}
	}
	return false
}

/* Address */

// Address is an address
type Address struct {
	State        string       `json:"state"`
	City         string       `json:"city"`
	Neighborhood string       `json:"neighborhood"`
	Street       string       `json:"street,omitempty"`
	StreetNumber string       `json:"streetNumber,omitempty"`
	Complement   string       `json:"complement,omitempty"`
	ZipCode      string       `json:"zipCode,omitempty"`
	AddressPoint AddressPoint `json:"point"`
	Precision    string       `json:"precision,omitempty"`
}

// Anonymize returns an anonymized Address instance
func (a Address) Anonymize() Address {
	return Address{
		State:        a.State,
		City:         a.City,
		Neighborhood: a.Neighborhood,
		Street:       a.Street,
		StreetNumber: strings.Repeat("*", len(a.StreetNumber)),
		Complement:   strings.Repeat("*", len(a.Complement)),
		ZipCode:      a.ZipCode,
	}
}

// Validate if an address is valid
func (a Address) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.ZipCode, validation.Required, is.Digit),
		validation.Field(&a.State, validation.Required),
		validation.Field(&a.City, validation.Required),
		validation.Field(&a.Neighborhood, validation.Required),
		validation.Field(&a.Street, validation.Required),
		validation.Field(&a.StreetNumber, validation.Required),
	)
}

// AddressPoint geo location
type AddressPoint struct {
	Lat float32 `json:"lat"`
	Lon float32 `json:"lon"`
}

/* Origin */

// Origin possible origins
type Origin int

const (
	// OriginNone zero value for this enum
	OriginNone Origin = iota
	// VivaReal indicates an origin from VivaReal Portal
	VivaReal
	// Zap indicates an origin from Zap Portal
	Zap
	// Sms indicates an origin from SMS
	Sms
	// EmailMarketing indicates an origin from email marketing
	EmailMarketing
	// ExternalAdvertising indicates an origin from external advertising
	ExternalAdvertising
	// ActiveOffer indicates an origin from active offer
	ActiveOffer
	// Telephone indicates an origin from telephone
	Telephone
	// Recommendation indicates an origin from recommendation
	Recommendation
	// AdvertiserSite indicates an origin from advertiser site
	AdvertiserSite
	//Other indicates an origin different from others
	Other
)

var originValues = [...]string{
	"",
	"VIVAREAL",
	"ZAP",
	"SMS",
	"EMAIL_MARKETING",
	"EXTERNAL_ADVERTISING",
	"ACTIVE_OFFER",
	"TELEPHONE",
	"RECOMMENDATION",
	"ADVERTISER_SITE",
	"OTHER",
}

func (o Origin) String() string {
	return originValues[o]
}

// MarshalJSON marshals the enum as a quoted json string
func (o Origin) MarshalJSON() ([]byte, error) {
	return common.QuotedStringBytes(o.String()), nil
}

// UnmarshalJSON unmarshals a quoted json string to the enum value
func (o *Origin) UnmarshalJSON(b []byte) error {
	x := ""
	err := json.Unmarshal(b, &x)
	if err != nil {
		return err
	}
	value, err := OriginValueOf(x)
	if err != nil {
		return err
	}
	*o = value
	return nil
}

// OriginValueOf converts a origin type value into a origin type
func OriginValueOf(v string) (Origin, error) {
	for i, value := range originValues {
		if value == v {
			return Origin(i), nil
		}
	}
	return 0, fmt.Errorf("unknown origin value %s", v)
}

var originFromPortals = []Origin{
	VivaReal,
	Zap,
}

// IsFromPortals tells when a origin is from ZAP or Vivareal
func (o Origin) IsFromPortals() bool {
	for _, portal := range originFromPortals {
		if o == portal {
			return true
		}
	}
	return false
}

// RoleType holds possible roles who can do or perform some action
type RoleType int

const (
	// RoleTypeNone zero value for this enum
	RoleTypeNone RoleType = iota
	// ContactRole means when the contact requested the action
	ContactRole
	// AdvertiserRole means when the advertiser requested the action
	AdvertiserRole
	// ProposerRole is when the processing is related to a credit analysis proposer
	ProposerRole
)

var roleValues = [...]string{
	"",
	"CONTACT",
	"ADVERTISER",
	"PROPOSER",
}

func (r RoleType) String() string {
	return roleValues[r]
}

// MarshalJSON marshals the enum as a quoted json string
func (r RoleType) MarshalJSON() ([]byte, error) {
	return common.QuotedStringBytes(r.String()), nil
}

// UnmarshalJSON unmarshals a quoted json string to the enum value
func (r *RoleType) UnmarshalJSON(b []byte) error {
	x := ""
	err := json.Unmarshal(b, &x)
	if err != nil {
		return err
	}
	value, err := RoleValueOf(x)
	if err != nil {
		return err
	}
	*r = value
	return nil
}

// RoleValueOf converts a role value into a role type
func RoleValueOf(v string) (RoleType, error) {
	for i, value := range roleValues {
		if value == v {
			return RoleType(i), nil
		}
	}
	return 0, fmt.Errorf("unknown role type value %s", v)
}

// ListingOrigin possible Listing origins
type ListingOrigin int

const (
	// ListingOriginNone zero value for this enum
	ListingOriginNone ListingOrigin = iota
	// PortalVivaReal indicates an origin from VivaReal Portal
	PortalVivaReal
	// PortalZap indicates an origin from Zap Portal
	PortalZap
)

var listingOriginValues = [...]string{
	"",
	"VIVAREAL",
	"ZAP",
}

// Host given a listing origin
func (lo ListingOrigin) Host() string {
	if lo == PortalVivaReal {
		return vivaRealPortalHost
	}
	return zapPortalHost
}

// PortalURL given a listing origin
func (lo ListingOrigin) PortalURL() string {
	if lo == PortalVivaReal {
		return fmt.Sprint("https://", vivaRealPortalHost)
	}
	return fmt.Sprint("https://", zapPortalHost)
}

// FredoURL given a listing origin
func (lo ListingOrigin) FredoURL() string {
	if lo == PortalVivaReal {
		return fredoVivaRealURL
	}

	return fredoZapURL
}

func (lo ListingOrigin) String() string {
	return listingOriginValues[lo]
}

// MarshalJSON marshals the enum as a quoted json string
func (lo ListingOrigin) MarshalJSON() ([]byte, error) {
	return common.QuotedStringBytes(lo.String()), nil
}

// UnmarshalJSON unmarshals a quoted json string to the enum value
func (lo *ListingOrigin) UnmarshalJSON(b []byte) error {
	x := ""
	err := json.Unmarshal(b, &x)
	if err != nil {
		return err
	}
	value, err := ListingOriginValueOf(x)
	if err != nil {
		return err
	}
	*lo = value
	return nil
}

// ListingOriginValueOf converts a listing origin type value into a listing origin type
func ListingOriginValueOf(v string) (ListingOrigin, error) {
	for i, value := range listingOriginValues {
		if value == v {
			return ListingOrigin(i), nil
		}
	}
	return 0, fmt.Errorf("unknown origin value %s", v)
}

// MaritalStatus possible marital statuses
type MaritalStatus int

const (
	// MaritalStatusNone zero value for this enum
	MaritalStatusNone MaritalStatus = iota
	// Single proposer
	Single
	// Married proposer
	Married
	// Divorced proposer
	Divorced
	// Widowed proposer
	Widowed
	// Separated proposer
	Separated
)

var maritalStatusValues = [...]string{
	"",
	"SINGLE",
	"MARRIED",
	"DIVORCED",
	"WIDOWED",
	"SEPARATED",
}

func (s MaritalStatus) String() string {
	return maritalStatusValues[s]
}

// MarshalJSON marshals the enum as a quoted json string
func (s MaritalStatus) MarshalJSON() ([]byte, error) {
	return common.QuotedStringBytes(s.String()), nil
}

// UnmarshalJSON unmarshals a quoted json string to the enum value
func (s *MaritalStatus) UnmarshalJSON(b []byte) error {
	x := ""
	err := json.Unmarshal(b, &x)
	if err != nil {
		return err
	}
	value, err := MaritalStatusValueOf(x)
	if err != nil {
		return err
	}
	*s = value
	return nil
}

// MaritalStatusValueOf converts a signer status value into a contract status type
func MaritalStatusValueOf(v string) (MaritalStatus, error) {
	for i, value := range maritalStatusValues {
		if value == v {
			return MaritalStatus(i), nil
		}
	}
	return 0, fmt.Errorf("unknown marital status value %s", v)
}

// CancelReasonType holds reasons of cancel of any entity
type CancelReasonType int

const (
	// CancelReasonTypeNone zero value for this enum
	CancelReasonTypeNone CancelReasonType = iota
	// ContactGaveUp is when the contact is not interested anymore
	ContactGaveUp
	// ContactRejectedByCreditAnalysis is when the contact was rejected by credit analysis
	ContactRejectedByCreditAnalysis
	// ContactProposalNotCompleted is when the contact didn't finished the offer / proposal
	ContactProposalNotCompleted
	// ContactDontKnowCreditAnalysis is when the contact didn't recognize this credit analysis
	ContactDontKnowCreditAnalysis
	// ContactRentedAnotherProperty is when the contact has already rented another property
	ContactRentedAnotherProperty
	// PropertyAlreadyRent is when the property is not avaiable anymore because it is already rented
	PropertyAlreadyRented
	// PropertyOwnerGaveUp is when the property owner does'nt want to rent the property anymore
	PropertyOwnerGaveUp
	// PropertyOwnerRejected is when the property owner rejected the negotiation
	PropertyOwnerRejectedOffer
	// BackofficeCanceled is when a backoffice user cancels
	BackofficeCanceled
)

var cancelReasonTypeValues = [...]string{
	"",
	"CONTACT_GAVE_UP",
	"CONTACT_REJECTED_BY_CREDIT_ANALYSIS",
	"CONTACT_PROPOSAL_NOT_COMPLETED",
	"CONTACT_DONT_KNOW_CREDIT_ANALYSIS",
	"CONTACT_RENTED_ANOTHER_PROPERTY",
	"PROPERTY_ALREADY_RENTED",
	"PROPERTY_OWNER_GAVE_UP",
	"PROPERTY_OWNER_REJECTED_OFFER",
	"BACKOFFICE_CANCELED",
}

func (crt CancelReasonType) String() string {
	return cancelReasonTypeValues[crt]
}

// MarshalJSON marshals the enum as a quoted json string
func (crt CancelReasonType) MarshalJSON() ([]byte, error) {
	return common.QuotedStringBytes(crt.String()), nil
}

// UnmarshalJSON unmarshals a quoted json string to the enum value
func (crt *CancelReasonType) UnmarshalJSON(b []byte) error {
	x := ""
	err := json.Unmarshal(b, &x)
	if err != nil {
		return err
	}
	value, err := CancelReasonTypeValueOf(x)
	if err != nil {
		return err
	}
	*crt = value
	return nil
}

// CancelReasonTypeValueOf converts a cancel reason type value into a cancel reason type
func CancelReasonTypeValueOf(v string) (CancelReasonType, error) {
	for i, value := range cancelReasonTypeValues {
		if value == v {
			return CancelReasonType(i), nil
		}
	}
	return 0, fmt.Errorf("unknown cancel reason type value %s", v)
}

// GenerateAccessToken generate a token to provide access to some private document to not logged users
func GenerateAccessToken(accountID, advertiserID, JWTSecret string, expHours int) (string, error) {
	return common.CreateJWTToken(jwt.MapClaims{
		"accountId":    accountID,
		"advertiserId": advertiserID,
		"exp":          time.Now().Add(time.Duration(expHours) * time.Hour).Unix(),
	}, JWTSecret)
}

// ParseAccessToken given a previosuly generated token, parse its contents
func ParseAccessToken(tokenStr, JWTSecret string) (string, string, error) {
	token, err := common.ParseJWTToken(tokenStr, JWTSecret)
	if err != nil {
		return "", "", err
	}

	if !token.Valid {
		return "", "", ErrInvalidDocumentToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", ErrInvalidDocumentToken
	}

	advertiserID, ok := claims["advertiserId"].(string)
	if !ok {
		return "", "", ErrInvalidDocumentToken
	}

	accountID, ok := claims["accountId"].(string)
	if !ok {
		return "", "", ErrInvalidDocumentToken
	}

	return accountID, advertiserID, nil
}

// SummaryItem composed by value and total
type SummaryItem struct {
	Value int `json:"value"`
	Total int `json:"total,omitempty"`
}

// LiveWithType with who a person can live with
type LiveWithType int

const (
	// LiveWithTypeNone zero value for this enum
	LiveWithTypeNone = iota
	// Alone when the tenant pretends to move alone
	Alone
	// Family when the tenant pretends to move with her/his family
	Family
	// Friends when the tenant pretends to move with her/his friends
	Friends
)

var liveWithValues = [...]string{
	"",
	"ALONE",
	"FAMILY",
	"FRIENDS",
}

func (lw LiveWithType) String() string {
	return liveWithValues[lw]
}

// MarshalJSON marshals the enum as a quoted json string
func (lw LiveWithType) MarshalJSON() ([]byte, error) {
	return common.QuotedStringBytes(lw.String()), nil
}

// UnmarshalJSON unmarshals a quoted json string to the enum value
func (lw *LiveWithType) UnmarshalJSON(b []byte) error {
	x := ""
	err := json.Unmarshal(b, &x)
	if err != nil {
		return err
	}
	value, err := LiveWithValueOf(x)
	if err != nil {
		return err
	}
	*lw = value
	return nil
}

// LiveWithValueOf converts a live with value into a live with type
func LiveWithValueOf(v string) (LiveWithType, error) {
	for i, value := range liveWithValues {
		if value == v {
			return LiveWithType(i), nil
		}
	}
	return 0, fmt.Errorf("unknown live with value %s", v)
}

// TenantInfo is the data about the tenant that want to rent a place
type TenantInfo struct {
	Adults          int          `json:"adults"`
	Children        int          `json:"children"`
	Pets            int          `json:"pets"`
	PetsDescription string       `json:"petsDescription"`
	LiveWith        LiveWithType `json:"liveWith"`
	SelfCommentary  string       `json:"selfCommentary"`
}

// Validate checks if a tenantInfo is valid
func (ti TenantInfo) Validate() error {
	return validation.ValidateStruct(&ti,
		validation.Field(&ti.LiveWith, validation.Required),
		validation.Field(&ti.Adults, validation.Required, validation.Min(1)),
		validation.Field(&ti.Children, validation.Min(0)),
		validation.Field(&ti.Pets, validation.Min(0), validation.When(ti.PetsDescription != "", validation.Required, validation.Min(1))),
		validation.Field(&ti.PetsDescription, validation.When(ti.Pets > 0, validation.Required)),
	)
}

/* Temperature */

const (
	// WarmScoreTemperatureEnd final value for a warm lead
	WarmScoreTemperatureEnd = 69
	// WarmScoreTemperatureStart start value for a warm lead
	WarmScoreTemperatureStart = 50
)

// Temperature of a score
type Temperature int

const (
	// TemperatureNone zero value for this enum
	TemperatureNone Temperature = iota
	// TemperatureHot best temperature
	TemperatureHot
	// TemperatureWarm medium temperature
	TemperatureWarm
	// TemperatureCold disposable temperature
	TemperatureCold
)

var temperatureValues = [...]string{
	"",
	"HOT",
	"WARM",
	"COLD",
}

// String converts a temperature value to string
func (t Temperature) String() string {
	return temperatureValues[t]
}

// MarshalJSON marshals the enum as a quoted json string
func (t Temperature) MarshalJSON() ([]byte, error) {
	return common.QuotedStringBytes(t.String()), nil
}

// UnmarshalJSON unmarshals a quoted json string to the enum value
func (t *Temperature) UnmarshalJSON(b []byte) error {
	x := ""
	err := json.Unmarshal(b, &x)
	if err != nil {
		return err
	}
	value, err := TemperatureValueOf(x)
	if err != nil {
		return err
	}
	*t = value
	return nil
}

// TemperatureValueOf converts a temperature value into a temperature type
func TemperatureValueOf(v string) (Temperature, error) {
	for i, value := range temperatureValues {
		if value == v {
			return Temperature(i), nil
		}
	}

	return 0, fmt.Errorf("unknown temperature value %s", v)
}

// ScoreRange of a given temperature, 0 means not present
func (t Temperature) ScoreRange() (int, int) {
	if t == TemperatureHot {
		return WarmScoreTemperatureEnd + 1, 0
	}
	if t == TemperatureWarm {
		return WarmScoreTemperatureStart, WarmScoreTemperatureEnd
	}

	return 0, WarmScoreTemperatureStart - 1
}

// TemperatureFromScore given a score value tells its temperature
func TemperatureFromScore(s int) Temperature {
	if s > WarmScoreTemperatureEnd {
		return TemperatureHot
	}

	if s >= WarmScoreTemperatureStart && s <= WarmScoreTemperatureEnd {
		return TemperatureWarm
	}

	return TemperatureCold
}
