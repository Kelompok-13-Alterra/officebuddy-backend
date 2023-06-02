package midtrans

const (
	VaBNI = "va-bni"
)

type CreateOrderParam struct {
	PaymentID       string
	OrderID         uint
	GrossAmount     int64
	ItemsDetails    ItemsDetails
	CustomerDetails CustomerDetails
}

type ItemsDetails struct {
	ID    string
	Price int64
	Qty   int
	Name  string
}

type CustomerDetails struct {
	Name  string
	Email string
}
