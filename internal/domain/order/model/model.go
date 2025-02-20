// Package model defines the core data structures used in the application
// for storing and managing order information. The package also provides
// structures for handling query parameters and editing operations.
package model

import "time"

// Order represents the core data entity, storing order data,
// including the ExternalOrderID, UserPhone and associated timestamps.
type Order struct {
	ID              string
	ExternalOrderID string
	UserPhone       string
	UserName        string
	CreatedAt       time.Time
}

// GetPars defines parameters for querying specific records,
// allowing filtering by ID, ExternalOrderID or UserPhone.
type GetPars struct {
	ID              string
	ExternalOrderID string
	UserPhone       string
}

// IsValid checks if at least one field in GetPars is populated.
func (m *GetPars) IsValid() bool {
	return m.ID != "" || m.ExternalOrderID != "" || m.UserPhone != ""
}

// ListPars defines parameters for listing records with optional filters,
// supporting filtering by IDs, ExternalOrderIDs, UserPhone and timestamps.
type ListPars struct {
	ID               *string
	IDs              *[]string
	ExternalOrderID  *string
	ExternalOrderIDs *[]string
	UserPhone        *string
	UserPhones       *[]string
	CreatedBefore    *time.Time
	CreatedAfter     *time.Time
}

// Edit represents the editable fields for updating an existing record,
// allowing partial updates to fields like UserPhone and timestamps.
type Edit struct {
	ExternalOrderID string
	UserPhone       *string
	UserName        *string
	CreatedAt       *time.Time
}
