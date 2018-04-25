package main

type data struct {
	Msg         interface{} `json:"msg"`
	Event       string      `json:"event"`
	ReceiverIDs []string    `json:"receiver_ids"`
	Username    string      `json:"username"`
}

type payload struct {
	ID     *uint64 `json:"id"`
	Method string  `json:"method"`
	Data   data    `json:"params"`
}

type registerResponse struct {
	ID     *uint64                 `json:"id"`
	Result *map[string]interface{} `json:"result,omitempty"`
	Error  *eventError             `json:"error,omitempty"`
}

type registerMessage struct {
	Sender string                 `json:"sender"`
	Event  string                 `json:"event"`
	Data   map[string]interface{} `json:"data"`
}

type eventError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
