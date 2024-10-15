package handler

import (
    "context"
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/TechBowl-japan/go-stations/model"
    "github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
    svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
    return &TODOHandler{
        svc: svc,
    }
}

// ServeHTTP implements http.Handler interface.
func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodPost:
        var req model.CreateTODORequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        ctx := r.Context()
        if req.Subject == "" {
            http.Error(w, "Subject cannot be empty", http.StatusBadRequest)
            return
        }
        resp, err := h.Create(ctx, &req)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        json.NewEncoder(w).Encode(resp)
    case http.MethodPut:
        var req model.UpdateTODORequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        ctx := r.Context()
        if req.ID == 0 {
            http.Error(w, "ID cannot be empty", http.StatusBadRequest)
            return
        }
        if req.Subject == "" {
            http.Error(w, "Subject cannot be empty", http.StatusBadRequest)
            return
        }
        resp, err := h.Update(ctx, &req)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        json.NewEncoder(w).Encode(resp)
    case http.MethodGet:
        var req model.ReadTODORequest
        QueryParams := r.URL.Query()
        prevID, err := strconv.ParseInt(QueryParams.Get("prev_id"), 10, 64)
        if err != nil {
            prevID = 0
        }
        size, err := strconv.ParseInt(QueryParams.Get("size"), 10, 64)
        if err != nil {
            size = -1
        }
        req.PrevID = prevID
        req.Size = size
        ctx := r.Context()
        resp, err := h.Read(ctx, &req)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        json.NewEncoder(w).Encode(resp)
    case http.MethodDelete:
        var req model.DeleteTODORequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        if len(req.IDs) == 0 {
            http.Error(w, "IDs cannot be empty", http.StatusBadRequest)
            return
        }
        ctx := r.Context()
        resp, err := h.Delete(ctx, &req)
        if err != nil {
            if _, ok := err.(*model.ErrNotFound); ok {
                http.Error(w, err.Error(), http.StatusNotFound)
                return
            }
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        json.NewEncoder(w).Encode(resp)
    default:
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
    }
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	resp, err := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}
	return &model.CreateTODOResponse{
		model.TODO{
			ID:          resp.ID,
			Subject:     resp.Subject,
			Description: resp.Description,
			CreatedAt:   resp.CreatedAt,
			UpdatedAt:   resp.UpdatedAt,
		},
	}, nil

}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
    resp, err := h.svc.ReadTODO(ctx, req.PrevID, req.Size)
    if err != nil {
        return nil, err
    }
    todos := make([]model.TODO, len(resp))
    for i, todo := range resp {
        todos[i] = model.TODO{
            ID:          todo.ID,
            Subject:     todo.Subject,
            Description: todo.Description,
            CreatedAt:   todo.CreatedAt,
            UpdatedAt:   todo.UpdatedAt,
        }
    }
    return &model.ReadTODOResponse{TODOs: todos}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	resp, err := h.svc.UpdateTODO(ctx, req.ID, req.Subject, req.Description)
    if err != nil {
        return nil, err
    }
    return &model.UpdateTODOResponse{
        model.TODO{
            ID:          resp.ID,
            Subject:     resp.Subject,
            Description: resp.Description,
            CreatedAt:   resp.CreatedAt,
            UpdatedAt:   resp.UpdatedAt,
        },
    }, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
    err := h.svc.DeleteTODO(ctx, req.IDs)
    if err != nil {
        return nil, err
    }
	return &model.DeleteTODOResponse{}, nil
}
