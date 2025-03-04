package model

import "time"

type Notification struct {
	ID          string
	OrderItemID string
	PhoneNumber string
	Status      string
	SentAt      time.Time
	CreatedAt   time.Time
}

type GetPars struct {
	ID          string
	OrderItemID string
	PhoneNumber string
	Status      string
}

func (m *GetPars) IsValid() bool {
	return m.ID != "" || m.OrderItemID != "" || m.PhoneNumber != "" || m.Status != ""
}

type ListPars struct {
	ID            *string
	IDs           *[]string
	OrderItemID   *string
	OrderItemIDs  *[]string
	PhoneNumber   *string
	PhoneNumbers  *[]string
	Status        *string
	Statuses      *[]string
	SentBefore    *time.Time
	SentAfter     *time.Time
	CreatedBefore *time.Time
	CreatedAfter  *time.Time
}

type Edit struct {
	ID          string
	OrderItemID *string
	PhoneNumber *string
	Status      *string
	SentAt      *time.Time
}
