package api

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"applicationDesignTest/internal/model"
)

type booker interface {
	Book(context.Context, model.Order) error
}

type API struct {
	booker booker
	server *http.Server
}

func New(s booker) *API {
	return &API{
		booker: s,
	}
}

func (a *API) Start(listen string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/orders", a.createOrder)

	a.server = &http.Server{Addr: listen, Handler: mux}
	slog.Info("Starting api-server on " + listen)
	return a.server.ListenAndServe()
}

func (a *API) Stop(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}

func (a *API) createOrder(w http.ResponseWriter, r *http.Request) {
	var order model.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, "got invalid json-request", http.StatusBadRequest)
		return
	}

	err = order.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.booker.Book(r.Context(), order)
	if err != nil {
		if errors.Is(err, model.ErrNoRoomsAvailable) {
			http.Error(w, err.Error(), http.StatusExpectationFailed)
			slog.Info("Hotel room is not available for selected dates", "from", order.From, "to", order.To)
			return
		}

		slog.Error("got internal error: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	slog.Info("Order successfully created", "order", order)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(order)
	if err != nil {
		slog.Error("failed to encode response: " + err.Error())
	}
}
