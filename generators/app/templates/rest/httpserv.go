package rest

//Response ...
type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (h *handler) handleRoot() error {
	response := Response{}
	response.Code = 200
	response.Status = "running"
	response.Message = "@build"

	return h.writeJSON(response)
}
