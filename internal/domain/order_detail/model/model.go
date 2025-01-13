// Package model defines the core data structures used in the application
// for storing and managing order items information. The package also provides
// structures for handling query parameters and editing operations.
package model

import "time"

// OrderDetail represents the core data entity, storing order data,
// including the OrderID, ProductCode and associated timestamps.
type OrderDetail struct {
	ID          string
	OrderID     string
	ProductCode string
	CreatedAt   time.Time
}

// GetPars defines parameters for querying specific records,
// allowing filtering by ID, OrderID or ProductCode.
type GetPars struct {
	ID          string
	OrderID     string
	ProductCode string
}

// IsValid checks if at least one field in GetPars is populated.
func (m *GetPars) IsValid() bool {
	return m.ID != "" || m.OrderID != "" || m.ProductCode != ""
}

// ListPars defines parameters for listing records with optional filters,
// supporting filtering by IDs, OrderIDs, ProductCodes and timestamps.
type ListPars struct {
	ID            *string
	IDs           *[]string
	OrderID       *string
	OrderIDs      *[]string
	ProductCode   *string
	ProductCodes  *[]string
	CreatedBefore *time.Time
	CreatedAfter  *time.Time
}

// Edit represents the editable fields for updating an existing record,
// allowing partial updates to fields like ProductCode.
type Edit struct {
	ID          string
	OrderID     string
	ProductCode *string
}
