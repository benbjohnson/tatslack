package tatslack

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Handler struct {
	DB      *DB
	Channel string
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/assets") {
		http.ServeFile(w, r, r.URL.Path[1:])
		return
	}

	switch r.URL.Path {
	case "/":
		h.serveRoot(w, r)
	case "/messages.json":
		h.serveMessagesJSON(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *Handler) serveRoot(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "assets/index.html")
}

func (h *Handler) serveMessagesJSON(w http.ResponseWriter, r *http.Request) {
	// Retrieve messages for the channel.
	a, err := h.DB.Messages(h.Channel)
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	// Write out messages.
	json.NewEncoder(w).Encode(a)
}
