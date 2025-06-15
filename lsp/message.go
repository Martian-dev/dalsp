package lsp

type Request struct {
	RPC    string `json:"jsonrpc"`
	ID     int    `json:"id"`
	Method string `json:"method"`

	// params
}

type Response struct {
	RPC string `json:"jsonrpc"`
	ID  *int   `json:"id,omitempty"`

	// result
	// errors
}

type Notification struct {
	RPC    string `josn:"jsonrpc"`
	Method string `json:"method"`
}
