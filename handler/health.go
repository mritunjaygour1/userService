package handler

import "net/http"

type HealthHandlerImpl struct{}

func NewHealthHandlerService() *HealthHandlerImpl {
	return &HealthHandlerImpl{}
}

func (h *HealthHandlerImpl) HealthCheck(w http.ResponseWriter, r *http.Request) {
	ReturnResponse(w, "Healthy", http.StatusOK)
}
