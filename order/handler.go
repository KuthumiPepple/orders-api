package order

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

var ErrNotExist = errors.New("order does not exist")

type Handler struct {
	Repo Repo
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		CustomerID uuid.UUID  `json:"customer_id"`
		LineItems  []LineItem `json:"line_items"`
	}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Println("failed to unmarshal:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()

	order := Order{
		OrderID:    rand.Uint64(),
		CustomerID: body.CustomerID,
		LineItems:  body.LineItems,
		CreatedAt:  &now,
	}

	err = h.Repo.Insert(r.Context(), order)
	if err != nil {
		log.Println("failed to insert:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(order)
	if err != nil {
		log.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64

	orderID, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		log.Println("failed to parse id:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	foundOrder, err := h.Repo.FindByID(
		r.Context(),
		orderID,
	)
	if errors.Is(err, ErrNotExist) {
		log.Println("order not found:", err)
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		log.Println("failed to find by id:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(foundOrder); err != nil {
		log.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	cursorStr := r.URL.Query().Get("cursor")
	if cursorStr == "" {
		cursorStr = "0"
	}

	const base = 10
	const bitSize = 64

	cursor, err := strconv.ParseUint(cursorStr, base, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	const count = 50

	res, err := h.Repo.FindAll(
		r.Context(),
		PaginationOptions{
			Count:  count,
			Cursor: cursor,
		},
	)
	if err != nil {
		log.Println("failed to find all:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var response struct {
		Items []Order `json:"items"`
		Next  uint64  `json:"next,omitempty"`
	}

	response.Items = res.Orders
	response.Next = res.Cursor

	data, err := json.Marshal(response)
	if err != nil {
		log.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

func (h *Handler) UpdateByID(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Println("failed to unmarshal:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	idParam := chi.URLParam(r, "id")
	const base = 10
	const bitSize = 64

	orderID, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		log.Println("failed to parse:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	foundOrder, err := h.Repo.FindByID(
		r.Context(),
		orderID,
	)
	if errors.Is(err, ErrNotExist) {
		log.Println("order not found:", err)
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		log.Println("failed to find by id:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	const shippedStatus = "shipped"
	const completedStatus = "completed"
	now := time.Now().UTC()

	switch body.Status {
	case shippedStatus:
		if foundOrder.ShippedAt != nil {
			log.Println("`ShippedAt` field already set")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		foundOrder.ShippedAt = &now
	case completedStatus:
		if foundOrder.CompletedAt != nil || foundOrder.ShippedAt == nil {
			log.Println("`CompletedAt` field already set or `ShippedAt` field not set")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		foundOrder.CompletedAt = &now
	default:
		log.Println("`Status` field not set to shipped or completed")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.Repo.Update(r.Context(), foundOrder); err != nil {
		log.Println("failed to update:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(foundOrder); err != nil {
		log.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64

	orderID, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		log.Println("failed to parse id:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.Repo.DeleteByID(
		r.Context(),
		orderID,
	)
	if errors.Is(err, ErrNotExist) {
		log.Println("order not found:", err)
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		log.Println("failed to delete by id:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
