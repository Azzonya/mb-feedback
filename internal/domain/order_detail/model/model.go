package model

import "time"

type OrderDetail struct {
	ID          string
	OrderID     string
	ProductCode string
	CreatedAt   time.Time
}

type OrderDetailWithUserInfo struct {
	ID          string
	OrderID     string
	UserPhone   string
	UserName    string
	ProductCode string
}

type GetPars struct {
	ID          string
	OrderID     string
	ProductCode string
}

func (m *GetPars) IsValid() bool {
	return m.ID != "" || m.OrderID != "" || m.ProductCode != ""
}

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

type Edit struct {
	ID          string
	OrderID     string
	ProductCode *string
}
