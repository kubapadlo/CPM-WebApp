package handler

import (
	"encoding/json"
	"errors"
	"net/http"
)

// Pojedyncze zadanie - wariant czynnosciowy
type TaskInputActivity struct {
	ID           string   `json:"id"` // nazwa zadania
	Duration     int      `json:"duration"`
	Predecessors []string `json:"predecessors"`
}

type TaskResultActivity struct {
	ID           string   `json:"id"`
	Duration     int      `json:"duration"`
	Predecessors []string `json:"predecessors"`
	ES           int      `json:"es"`          // Early Start
	EF           int      `json:"ef"`          // Early Finish
	LS           int      `json:"ls"`          // Late Start
	LF           int      `json:"lf"`          // Late Finish
	Slack        int      `json:"slack"`       // Luz
	IsCritical   bool     `json:"is_critical"` // Ścieżka krytyczna 0/1
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func HandleActivity(w http.ResponseWriter, r *http.Request) {
	var inputs []TaskInputActivity
	if err := json.NewDecoder(r.Body).Decode(&inputs); err != nil {
		sendError(w, "Błędny format JSON", http.StatusBadRequest)
		return
	}

	results, err := calculateCPMActivity(inputs)
	if err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func sendError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

func calculateCPMActivity(inputs []TaskInputActivity) ([]TaskResultActivity, error) {
	n := len(inputs)
	adj := make(map[string][]string) // poprzednicy
	nPreds := make(map[string]int)   // liczba poprzednikow dla kazdego zadania
	nodes := make(map[string]*TaskResultActivity)
	preds := make(map[string][]string)

	for _, i := range inputs {
		nodes[i.ID] = &TaskResultActivity{ID: i.ID, Duration: i.Duration, Predecessors: i.Predecessors}
		preds[i.ID] = i.Predecessors
		if _, exists := nPreds[i.ID]; !exists {
			nPreds[i.ID] = 0
		}
		for _, p := range i.Predecessors {
			adj[p] = append(adj[p], i.ID)
			nPreds[i.ID]++
		}
	}

	// ----------------------------------------------
	// KROK DO PRZODU
	// ----------------------------------------------

	queue := []string{}
	for id, degree := range nPreds {
		if degree == 0 {
			queue = append(queue, id)
		}
	}

	var topoOrder []string
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		topoOrder = append(topoOrder, u)

		// EF = ES + Duration
		nodes[u].EF = nodes[u].ES + nodes[u].Duration

		for _, v := range adj[u] {
			// ES = max(EF poprzedników)
			if nodes[u].EF > nodes[v].ES {
				nodes[v].ES = nodes[u].EF
			}
			nPreds[v]--
			if nPreds[v] == 0 {
				queue = append(queue, v)
			}
		}
	}

	// Wykrywanie cyklu - idk rzekomo może sie zdarzyć
	if len(topoOrder) != n {
		return nil, errors.New("wykryto cykl w grafie zadań")
	}

	// Czas całkowity projektu
	maxEF := 0
	for _, t := range nodes {
		if t.EF > maxEF {
			maxEF = t.EF
		}
	}

	// ----------------------------------------------
	// KROK DO TYŁU
	// ----------------------------------------------
	for _, t := range nodes {
		t.LF = maxEF
	}

	for i := len(topoOrder) - 1; i >= 0; i-- {
		u := topoOrder[i]
		// LS = LF - Duration
		nodes[u].LS = nodes[u].LF - nodes[u].Duration
		// Slack = LS - ES
		nodes[u].Slack = nodes[u].LS - nodes[u].ES
		nodes[u].IsCritical = nodes[u].Slack == 0

		// Aktualizacja LF dla poprzedników: LF = min(LS następców)
		for _, pID := range preds[u] {
			if nodes[u].LS < nodes[pID].LF {
				nodes[pID].LF = nodes[u].LS
			}
		}
	}

	// Konwersja mapy na listę wynikową
	finalResults := make([]TaskResultActivity, 0, n)
	for _, id := range topoOrder {
		finalResults = append(finalResults, *nodes[id])
	}

	return finalResults, nil
}
