package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	headersMap := make(map[string][]string)

	for name, values := range r.Header {
		headersMap[name] = values
		fmt.Printf("%s: %s\n", name, values)
	}

	jsonHeaders, err := json.Marshal(headersMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonHeaders)
}
