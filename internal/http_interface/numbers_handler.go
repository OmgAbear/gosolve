package http_interface

import (
	"encoding/json"
	"github.com/OmgAbear/gosolve/internal/http_interface/dto"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"strconv"
)

// NumbersRepo is the interface required by the NumbersHandler
// defined here as the handler itself should declare what it needs
type NumbersRepo interface {
	FindNearestIndex(target int) dto.NumbersResult
}

// NumbersRepoFactory is the factory that provides an impl for NumbersRepo
type NumbersRepoFactory func() NumbersRepo

type NumbersHandler struct {
	dataRepo NumbersRepo
	logger   *slog.Logger
}

// NewNumbersHandler - creates a new NumbersHandler instance
func NewNumbersHandler(repoFactory NumbersRepoFactory, logger *slog.Logger) NumbersHandler {
	return NumbersHandler{
		dataRepo: repoFactory(),
		logger:   logger,
	}
}

// get -http handler that returns a given index and value for a searched value
func (handler NumbersHandler) get(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(request)
	valueStr := vars["value"]

	target, err := strconv.Atoi(valueStr)

	if err != nil {
		handler.logger.Error("Invalid number", "error", err)
		writer.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(writer).Encode(dto.NumbersResult{
			Index:   -1,
			Value:   -1,
			Message: func() *string { ret := "invalid number. Must be integer"; return &ret }(),
		})
		return
	}

	searchResult := handler.dataRepo.FindNearestIndex(target)

	resultStatus := http.StatusOK
	if searchResult.Index == -1 {
		resultStatus = http.StatusNotFound
	}

	writer.WriteHeader(resultStatus)

	_ = json.NewEncoder(writer).Encode(searchResult)
}
