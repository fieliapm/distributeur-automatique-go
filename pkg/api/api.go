package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"gitlab.rayark.com/fieliapm/distributeur_automatique/pkg/core"
)

type ExactPurchaseInput struct {
	Prices []int `json:"prices"`
	Budget int   `json:"budget"`
}

type ExactPurchaseOutput struct {
	SolutionCount int     `json:"solution_count"`
	Solutions     [][]int `json:"solutions"`
}

func findExactPurchase(w http.ResponseWriter, r *http.Request) error {
	defer func() {
		fmt.Fprintln(os.Stderr, "[server] request complete")
	}()

	fmt.Fprintln(os.Stderr, "[server] receiving request")

	var exactPurchaseInput ExactPurchaseInput
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&exactPurchaseInput)
	if err != nil {
		return WrapApiError(ErrInvalidRequestBody, err)
	}

	if !core.ValidatePrices(exactPurchaseInput.Prices) {
		return ErrInvalidPrices
	}

	solutionCount, solutions := core.FindExactPurchaseCache(exactPurchaseInput.Prices, exactPurchaseInput.Budget)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	errorResponse := ExactPurchaseOutput{SolutionCount: solutionCount, Solutions: solutions}
	err = encoder.Encode(errorResponse)
	if err != nil {
		panic(err)
	}

	return nil
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "template/distributeur-automatique.htm")
}

func RunServer(port int) *http.Server {
	defer func() {
		fmt.Fprintln(os.Stderr, "[server] server shutted down")
	}()

	fmt.Fprintln(os.Stderr, "[server] server starting")

	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/api/find_exact_purchase", ErrAwareHandle(findExactPurchase))
	serveMux.HandleFunc("/", index)

	server := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: serveMux}

	fmt.Fprintln(os.Stderr, "[server] server started")
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}

	return server
}
