package mb_broker

type FetchCompletedOrdersRepSt struct {
	Page       int     `json:"page"`
	PageSize   int     `json:"page_size"`
	TotalCount int     `json:"total_count"`
	Results    []OrdSt `json:"results"`
}

type OrdSt struct {
	PrvCode  string        `json:"prv_code"`
	Customer OrdCustomerSt `json:"customer"`
}

type OrdCustomerSt struct {
	CellPhone string `json:"cell_phone"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type FetchProductCodesReqSt struct {
	PrvCode string `json:"prv_code"`
}
