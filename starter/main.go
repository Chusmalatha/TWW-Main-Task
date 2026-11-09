
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"
)

var (
	mu     sync.RWMutex
	events = make(map[string]Event)
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /events", handlePostEvents)
	mux.HandleFunc("GET /campaigns/{campaignID}/stats", handleGetStats)

	addr := ":8080"
	fmt.Printf("listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func handlePostEvents(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 5<<20)

	var rawEvents []json.RawMessage

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rawEvents); err != nil || rawEvents == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "request must contain a JSON array",
		})
		return
	}

	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "only one JSON array is allowed",
		})
		return
	}

	accepted := 0
	duplicates := 0
	rejected := 0
	var errors []string

	for i, raw := range rawEvents {
		var event Event

		if err := json.Unmarshal(raw, &event); err != nil {
			rejected++
			errors = append(errors, fmt.Sprintf("item %d: invalid JSON", i))
			continue
		}

		event.EventID = strings.TrimSpace(event.EventID)
		event.CampaignID = strings.TrimSpace(event.CampaignID)
		event.ContactID = strings.TrimSpace(event.ContactID)
		event.Type = strings.ToLower(strings.TrimSpace(event.Type))

		if err := validateEvent(event); err != nil {
			rejected++
			errors = append(errors, fmt.Sprintf("item %d: %v", i, err))
			continue
		}

		mu.Lock()

		existing, exists := events[event.EventID]
		if exists {
			mu.Unlock()

			if reflect.DeepEqual(existing, event) {
				duplicates++
			} else {
				rejected++
				errors = append(errors, fmt.Sprintf(
					"item %d: event ID already exists with different data", i,
				))
			}
			continue
		}

		events[event.EventID] = event
		mu.Unlock()

		accepted++
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"received":   len(rawEvents),
		"accepted":   accepted,
		"duplicates": duplicates,
		"rejected":   rejected,
		"errors":     errors,
	})
}

func validateEvent(event Event) error {
	if event.EventID == "" {
		return fmt.Errorf("event_id is required")
	}
	if event.CampaignID == "" {
		return fmt.Errorf("campaign_id is required")
	}
	if event.ContactID == "" {
		return fmt.Errorf("contact_id is required")
	}

	switch event.Type {
	case "sent", "delivered", "opened", "clicked":
	default:
		return fmt.Errorf("unsupported event type")
	}

	if _, err := time.Parse(time.RFC3339, event.Timestamp); err != nil {
		return fmt.Errorf("invalid timestamp")
	}

	return nil
}

func handleGetStats(w http.ResponseWriter, r *http.Request) {
	campaignID := r.PathValue("campaignID")

	stats := map[string]int{
		"sent":      0,
		"delivered": 0,
		"opened":    0,
		"clicked":   0,
	}

	mu.RLock()
	for _, event := range events {
		if event.CampaignID == campaignID {
			stats[event.Type]++
		}
	}
	mu.RUnlock()

	writeJSON(w, http.StatusOK, map[string]any{
		"campaign_id": campaignID,
		"sent":        stats["sent"],
		"delivered":   stats["delivered"],
		"opened":      stats["opened"],
		"clicked":     stats["clicked"],
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("writeJSON: %v", err)
	}
}
