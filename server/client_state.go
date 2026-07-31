package server

type Client struct {
	FD          int
	InputBuffer []byte

	InTransaction    bool
	TransactionQueue [][]string

	Subscriptions map[string]struct{}
}
