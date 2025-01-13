// Package model defines the core data structures used in the application
// for storing and managing sms information. The package also provides
// structures for handling query parameters and editing operations.
package model

import "time"

// Sms represents the core data entity, storing sms data,
// including the OrderItemID, PhoneNumber, Status and associated timestamps.
type Sms struct {
	ID          string
	OrderItemID string
	PhoneNumber string
	Status      string
	SentAt      time.Time
	CreatedAt   time.Time
}

// GetPars defines parameters for querying specific records,
// allowing filtering by ID, OrderItemID, PhoneNumber or Status.
type GetPars struct {
	ID          string
	OrderItemID string
	PhoneNumber string
	Status      string
}

// IsValid checks if at least one field in GetPars is populated.
func (m *GetPars) IsValid() bool {
	return m.ID != "" || m.OrderItemID != "" || m.PhoneNumber != "" || m.Status != ""
}

// ListPars defines parameters for listing records with optional filters,
// supporting filtering by IDs, OrderItemIDs, PhoneNumber and timestamps.
type ListPars struct {
	ID            *string
	IDs           *[]string
	OrderItemID   *string
	OrderItemIDs  *[]string
	PhoneNumber   *string
	PhoneNumbers  *[]string
	SentBefore    *time.Time
	SentAfter     *time.Time
	CreatedBefore *time.Time
	CreatedAfter  *time.Time
}

// Edit represents the editable fields for updating an existing record,
// allowing partial updates to fields like Status and timestamps.
type Edit struct {
	ID     string
	Status *string
	SentAt *time.Time
}
