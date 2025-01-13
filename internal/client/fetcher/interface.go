package fetcher

type Fetcher interface {
	FetchCompletedOrders()
	FetchProductCodesByOrder(orderID string) ([]string, error)
}
