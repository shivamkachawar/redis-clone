package server

type Client struct {
	FD          int
	InputBuffer []byte
}
