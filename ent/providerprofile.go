// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/paycrest/paycrest-protocol/ent/apikey"
	"github.com/paycrest/paycrest-protocol/ent/fiatcurrency"
	"github.com/paycrest/paycrest-protocol/ent/provideravailability"
	"github.com/paycrest/paycrest-protocol/ent/providerprofile"
	"github.com/paycrest/paycrest-protocol/ent/providerrating"
	"github.com/paycrest/paycrest-protocol/ent/user"
)

// ProviderProfile is the model entity for the ProviderProfile schema.
type ProviderProfile struct {
	config `json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// TradingName holds the value of the "trading_name" field.
	TradingName string `json:"trading_name,omitempty"`
	// HostIdentifier holds the value of the "host_identifier" field.
	HostIdentifier string `json:"host_identifier,omitempty"`
	// ProvisionMode holds the value of the "provision_mode" field.
	ProvisionMode providerprofile.ProvisionMode `json:"provision_mode,omitempty"`
	// IsPartner holds the value of the "is_partner" field.
	IsPartner bool `json:"is_partner,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ProviderProfileQuery when eager-loading is set.
	Edges                  ProviderProfileEdges `json:"edges"`
	fiat_currency_provider *uuid.UUID
	user_provider_profile  *uuid.UUID
	selectValues           sql.SelectValues
}

// ProviderProfileEdges holds the relations/edges for other nodes in the graph.
type ProviderProfileEdges struct {
	// User holds the value of the user edge.
	User *User `json:"user,omitempty"`
	// APIKey holds the value of the api_key edge.
	APIKey *APIKey `json:"api_key,omitempty"`
	// Currency holds the value of the currency edge.
	Currency *FiatCurrency `json:"currency,omitempty"`
	// ProvisionBuckets holds the value of the provision_buckets edge.
	ProvisionBuckets []*ProvisionBucket `json:"provision_buckets,omitempty"`
	// OrderTokens holds the value of the order_tokens edge.
	OrderTokens []*ProviderOrderToken `json:"order_tokens,omitempty"`
	// Availability holds the value of the availability edge.
	Availability *ProviderAvailability `json:"availability,omitempty"`
	// ProviderRating holds the value of the provider_rating edge.
	ProviderRating *ProviderRating `json:"provider_rating,omitempty"`
	// AssignedOrders holds the value of the assigned_orders edge.
	AssignedOrders []*LockPaymentOrder `json:"assigned_orders,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [8]bool
}

// UserOrErr returns the User value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ProviderProfileEdges) UserOrErr() (*User, error) {
	if e.loadedTypes[0] {
		if e.User == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: user.Label}
		}
		return e.User, nil
	}
	return nil, &NotLoadedError{edge: "user"}
}

// APIKeyOrErr returns the APIKey value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ProviderProfileEdges) APIKeyOrErr() (*APIKey, error) {
	if e.loadedTypes[1] {
		if e.APIKey == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: apikey.Label}
		}
		return e.APIKey, nil
	}
	return nil, &NotLoadedError{edge: "api_key"}
}

// CurrencyOrErr returns the Currency value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ProviderProfileEdges) CurrencyOrErr() (*FiatCurrency, error) {
	if e.loadedTypes[2] {
		if e.Currency == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: fiatcurrency.Label}
		}
		return e.Currency, nil
	}
	return nil, &NotLoadedError{edge: "currency"}
}

// ProvisionBucketsOrErr returns the ProvisionBuckets value or an error if the edge
// was not loaded in eager-loading.
func (e ProviderProfileEdges) ProvisionBucketsOrErr() ([]*ProvisionBucket, error) {
	if e.loadedTypes[3] {
		return e.ProvisionBuckets, nil
	}
	return nil, &NotLoadedError{edge: "provision_buckets"}
}

// OrderTokensOrErr returns the OrderTokens value or an error if the edge
// was not loaded in eager-loading.
func (e ProviderProfileEdges) OrderTokensOrErr() ([]*ProviderOrderToken, error) {
	if e.loadedTypes[4] {
		return e.OrderTokens, nil
	}
	return nil, &NotLoadedError{edge: "order_tokens"}
}

// AvailabilityOrErr returns the Availability value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ProviderProfileEdges) AvailabilityOrErr() (*ProviderAvailability, error) {
	if e.loadedTypes[5] {
		if e.Availability == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: provideravailability.Label}
		}
		return e.Availability, nil
	}
	return nil, &NotLoadedError{edge: "availability"}
}

// ProviderRatingOrErr returns the ProviderRating value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ProviderProfileEdges) ProviderRatingOrErr() (*ProviderRating, error) {
	if e.loadedTypes[6] {
		if e.ProviderRating == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: providerrating.Label}
		}
		return e.ProviderRating, nil
	}
	return nil, &NotLoadedError{edge: "provider_rating"}
}

// AssignedOrdersOrErr returns the AssignedOrders value or an error if the edge
// was not loaded in eager-loading.
func (e ProviderProfileEdges) AssignedOrdersOrErr() ([]*LockPaymentOrder, error) {
	if e.loadedTypes[7] {
		return e.AssignedOrders, nil
	}
	return nil, &NotLoadedError{edge: "assigned_orders"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*ProviderProfile) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case providerprofile.FieldIsPartner:
			values[i] = new(sql.NullBool)
		case providerprofile.FieldID, providerprofile.FieldTradingName, providerprofile.FieldHostIdentifier, providerprofile.FieldProvisionMode:
			values[i] = new(sql.NullString)
		case providerprofile.FieldCreatedAt, providerprofile.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		case providerprofile.ForeignKeys[0]: // fiat_currency_provider
			values[i] = &sql.NullScanner{S: new(uuid.UUID)}
		case providerprofile.ForeignKeys[1]: // user_provider_profile
			values[i] = &sql.NullScanner{S: new(uuid.UUID)}
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the ProviderProfile fields.
func (pp *ProviderProfile) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case providerprofile.FieldID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value.Valid {
				pp.ID = value.String
			}
		case providerprofile.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				pp.CreatedAt = value.Time
			}
		case providerprofile.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				pp.UpdatedAt = value.Time
			}
		case providerprofile.FieldTradingName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field trading_name", values[i])
			} else if value.Valid {
				pp.TradingName = value.String
			}
		case providerprofile.FieldHostIdentifier:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field host_identifier", values[i])
			} else if value.Valid {
				pp.HostIdentifier = value.String
			}
		case providerprofile.FieldProvisionMode:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field provision_mode", values[i])
			} else if value.Valid {
				pp.ProvisionMode = providerprofile.ProvisionMode(value.String)
			}
		case providerprofile.FieldIsPartner:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field is_partner", values[i])
			} else if value.Valid {
				pp.IsPartner = value.Bool
			}
		case providerprofile.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullScanner); !ok {
				return fmt.Errorf("unexpected type %T for field fiat_currency_provider", values[i])
			} else if value.Valid {
				pp.fiat_currency_provider = new(uuid.UUID)
				*pp.fiat_currency_provider = *value.S.(*uuid.UUID)
			}
		case providerprofile.ForeignKeys[1]:
			if value, ok := values[i].(*sql.NullScanner); !ok {
				return fmt.Errorf("unexpected type %T for field user_provider_profile", values[i])
			} else if value.Valid {
				pp.user_provider_profile = new(uuid.UUID)
				*pp.user_provider_profile = *value.S.(*uuid.UUID)
			}
		default:
			pp.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the ProviderProfile.
// This includes values selected through modifiers, order, etc.
func (pp *ProviderProfile) Value(name string) (ent.Value, error) {
	return pp.selectValues.Get(name)
}

// QueryUser queries the "user" edge of the ProviderProfile entity.
func (pp *ProviderProfile) QueryUser() *UserQuery {
	return NewProviderProfileClient(pp.config).QueryUser(pp)
}

// QueryAPIKey queries the "api_key" edge of the ProviderProfile entity.
func (pp *ProviderProfile) QueryAPIKey() *APIKeyQuery {
	return NewProviderProfileClient(pp.config).QueryAPIKey(pp)
}

// QueryCurrency queries the "currency" edge of the ProviderProfile entity.
func (pp *ProviderProfile) QueryCurrency() *FiatCurrencyQuery {
	return NewProviderProfileClient(pp.config).QueryCurrency(pp)
}

// QueryProvisionBuckets queries the "provision_buckets" edge of the ProviderProfile entity.
func (pp *ProviderProfile) QueryProvisionBuckets() *ProvisionBucketQuery {
	return NewProviderProfileClient(pp.config).QueryProvisionBuckets(pp)
}

// QueryOrderTokens queries the "order_tokens" edge of the ProviderProfile entity.
func (pp *ProviderProfile) QueryOrderTokens() *ProviderOrderTokenQuery {
	return NewProviderProfileClient(pp.config).QueryOrderTokens(pp)
}

// QueryAvailability queries the "availability" edge of the ProviderProfile entity.
func (pp *ProviderProfile) QueryAvailability() *ProviderAvailabilityQuery {
	return NewProviderProfileClient(pp.config).QueryAvailability(pp)
}

// QueryProviderRating queries the "provider_rating" edge of the ProviderProfile entity.
func (pp *ProviderProfile) QueryProviderRating() *ProviderRatingQuery {
	return NewProviderProfileClient(pp.config).QueryProviderRating(pp)
}

// QueryAssignedOrders queries the "assigned_orders" edge of the ProviderProfile entity.
func (pp *ProviderProfile) QueryAssignedOrders() *LockPaymentOrderQuery {
	return NewProviderProfileClient(pp.config).QueryAssignedOrders(pp)
}

// Update returns a builder for updating this ProviderProfile.
// Note that you need to call ProviderProfile.Unwrap() before calling this method if this ProviderProfile
// was returned from a transaction, and the transaction was committed or rolled back.
func (pp *ProviderProfile) Update() *ProviderProfileUpdateOne {
	return NewProviderProfileClient(pp.config).UpdateOne(pp)
}

// Unwrap unwraps the ProviderProfile entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (pp *ProviderProfile) Unwrap() *ProviderProfile {
	_tx, ok := pp.config.driver.(*txDriver)
	if !ok {
		panic("ent: ProviderProfile is not a transactional entity")
	}
	pp.config.driver = _tx.drv
	return pp
}

// String implements the fmt.Stringer.
func (pp *ProviderProfile) String() string {
	var builder strings.Builder
	builder.WriteString("ProviderProfile(")
	builder.WriteString(fmt.Sprintf("id=%v, ", pp.ID))
	builder.WriteString("created_at=")
	builder.WriteString(pp.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(pp.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("trading_name=")
	builder.WriteString(pp.TradingName)
	builder.WriteString(", ")
	builder.WriteString("host_identifier=")
	builder.WriteString(pp.HostIdentifier)
	builder.WriteString(", ")
	builder.WriteString("provision_mode=")
	builder.WriteString(fmt.Sprintf("%v", pp.ProvisionMode))
	builder.WriteString(", ")
	builder.WriteString("is_partner=")
	builder.WriteString(fmt.Sprintf("%v", pp.IsPartner))
	builder.WriteByte(')')
	return builder.String()
}

// ProviderProfiles is a parsable slice of ProviderProfile.
type ProviderProfiles []*ProviderProfile
