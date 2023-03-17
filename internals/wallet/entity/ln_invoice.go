package entity

type LNInvoice struct {
	UserID         string `json:"-"`
	RequestHash    string `json:"requestHash"`
	PaymentRequest string `json:"paymentRequest"`
	AddIndex       uint64 `json:"addIndex"`
	PaymentAddress string `json:"paymentAddress"`
}
