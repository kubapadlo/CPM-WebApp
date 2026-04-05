package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"sort"
)

type TaskEventInput struct {
	ID       string `json:"id"`
	FromNode int    `json:"from_node"`
	ToNode   int    `json:"to_node"`
	Duration int    `json:"duration"`
}

// Wynik dla pojedynczego zdarzenia (bloczka)
type EventResult struct {
	Node       int  `json:"node"`
	ET         int  `json:"et"`
	LT         int  `json:"lt"`
	Slack      int  `json:"slack"`
	IsCritical bool `json:"is_critical"`
}

// Wynik dla pojedynczej czynności (strzałki)
type ActivityResult struct {
	ID         string `json:"id"`
	FromNode   int    `json:"from_node"`
	ToNode     int    `json:"to_node"`
	Duration   int    `json:"duration"`
	ES         int    `json:"es"`
	EF         int    `json:"ef"`
	LS         int    `json:"ls"`
	LF         int    `json:"lf"`
	Slack      int    `json:"slack"`
	IsCritical bool   `json:"is_critical"`
}

// Główny obiekt zwracany przez API
type AoAResponse struct {
	Events     []EventResult    `json:"events"`
	Activities []ActivityResult `json:"activities"`
}

func HandleEvent(w http.ResponseWriter, r *http.Request) {
	var inputs []TaskEventInput
	if err := json.NewDecoder(r.Body).Decode(&inputs); err != nil {
		sendError(w, "Błędny format JSON", http.StatusBadRequest)
		return
	}

	results, err := calculateCPMEvent(inputs)
	if err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func calculateCPMEvent(inputs []TaskEventInput) (*AoAResponse, error) {
	adj := make(map[int][]TaskEventInput)
	inDegree := make(map[int]int)
	eventsMap := make(map[int]*EventResult)
	allNodeIDs := make(map[int]bool)

	for _, i := range inputs {
		allNodeIDs[i.FromNode] = true
		allNodeIDs[i.ToNode] = true
		adj[i.FromNode] = append(adj[i.FromNode], i)

		inDegree[i.ToNode]++
		if _, exists := inDegree[i.FromNode]; !exists {
			inDegree[i.FromNode] = 0
		}
	}

	for id := range allNodeIDs {
		eventsMap[id] = &EventResult{Node: id, ET: 0, LT: 0}
	}

	queue := []int{}
	for id, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, id)
		}
	}

	var topoOrder []int
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		topoOrder = append(topoOrder, u)

		for _, activity := range adj[u] {
			v := activity.ToNode
			if eventsMap[u].ET+activity.Duration > eventsMap[v].ET {
				eventsMap[v].ET = eventsMap[u].ET + activity.Duration
			}
			inDegree[v]--
			if inDegree[v] == 0 {
				queue = append(queue, v)
			}
		}
	}

	if len(topoOrder) != len(allNodeIDs) {
		return nil, errors.New("wykryto cykl w grafie zdarzeń")
	}

	maxET := 0
	for _, e := range eventsMap {
		if e.ET > maxET {
			maxET = e.ET
		}
	}

	for _, e := range eventsMap {
		e.LT = maxET
	}

	for i := len(topoOrder) - 1; i >= 0; i-- {
		u := topoOrder[i]
		for _, activity := range adj[u] {
			v := activity.ToNode
			if eventsMap[v].LT-activity.Duration < eventsMap[u].LT {
				eventsMap[u].LT = eventsMap[v].LT - activity.Duration
			}
		}
	}

	response := &AoAResponse{
		Events:     make([]EventResult, 0, len(allNodeIDs)),
		Activities: make([]ActivityResult, 0, len(inputs)),
	}

	for id := range allNodeIDs {
		e := eventsMap[id]
		e.Slack = e.LT - e.ET
		e.IsCritical = e.Slack == 0
		response.Events = append(response.Events, *e)
	}

	for _, i := range inputs {
		es := eventsMap[i.FromNode].ET
		ef := es + i.Duration
		lf := eventsMap[i.ToNode].LT
		ls := lf - i.Duration
		slack := lf - ef

		response.Activities = append(response.Activities, ActivityResult{
			ID:         i.ID,
			FromNode:   i.FromNode,
			ToNode:     i.ToNode,
			Duration:   i.Duration,
			ES:         es,
			EF:         ef,
			LS:         ls,
			LF:         lf,
			Slack:      slack,
			IsCritical: slack == 0,
		})
	}

	sort.Slice(response.Events, func(i, j int) bool {
		return response.Events[i].Node < response.Events[j].Node
	})

	sort.Slice(response.Activities, func(i, j int) bool {
		return response.Activities[i].ID < response.Activities[j].ID
	})

	return response, nil
}
