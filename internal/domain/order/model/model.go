package model

import "time"

type Order struct {
	ID              string
	ExternalOrderID string
	UserPhone       string
	UserName        string
	CreatedAt       time.Time
}

type GetPars struct {
	ID              string
	ExternalOrderID string
	UserPhone       string
}

func (m *GetPars) IsValid() bool {
	return m.ID != "" || m.ExternalOrderID != "" || m.UserPhone != ""
}

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

type Edit struct {
	ExternalOrderID string
	UserPhone       *string
	UserName        *string
	CreatedAt       *time.Time
}
