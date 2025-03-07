package models

type Message struct {
    Indonesia string `json:"Indonesia"`
    English   string `json:"English"`
}

type BackendPostResponse struct {
    Code    string  `json:"Code"`
    Message Message `json:"Message"`
}
