package lib

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

const problemTypeBaseURL = "https://api.<your-domain>/problems"

const (
	ProblemTypeValidationError = problemTypeBaseURL + "/validation-error"
	ProblemTypeDuplicateRes    = problemTypeBaseURL + "/duplicate-resource"
	ProblemTypeNotFound        = problemTypeBaseURL + "/not-found"
	ProblemTypeInternalError   = problemTypeBaseURL + "/internal-error"
)

type ProblemDetails struct {
	Type      string            `json:"type"`
	Title     string            `json:"title"`
	Status    int               `json:"status"`
	Detail    string            `json:"detail"`
	Instance  string            `json:"instance"`
	RequestID string            `json:"request_id"`
	Errors    map[string]string `json:"errors,omitempty"`
}

func WriteJSON(w http.ResponseWriter, r *http.Request, status int, payload any) {
	body, err := json.Marshal(payload)
	if err != nil {
		WriteProblem(w, r, ProblemDetails{
			Type:   ProblemTypeInternalError,
			Title:  "Internal Server Error",
			Status: http.StatusInternalServerError,
			Detail: "failed to serialize response",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(body)
}

func WriteProblem(w http.ResponseWriter, r *http.Request, problem ProblemDetails) {
	if problem.Title == "" {
		problem.Title = http.StatusText(problem.Status)
	}
	if problem.Detail == "" {
		problem.Detail = http.StatusText(problem.Status)
	}
	if problem.Instance == "" && r != nil && r.URL != nil {
		problem.Instance = r.URL.Path
	}
	if r != nil {
		problem.RequestID = middleware.GetReqID(r.Context())
	}

	body, err := json.Marshal(problem)
	if err != nil {
		fallback := `{"type":"` + ProblemTypeInternalError + `","title":"Internal Server Error","status":500,"detail":"failed to serialize problem response","instance":"","request_id":""}`
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fallback))
		return
	}

	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(problem.Status)
	_, _ = w.Write(body)
}
