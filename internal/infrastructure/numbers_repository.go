package infrastructure

import (
	"fmt"
	"github.com/OmgAbear/gosolve/internal/http_interface/dto"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/OmgAbear/gosolve/internal/config"
)

type NumbersRepo struct {
	data   []int
	logger *slog.Logger
}

func NewNumbersRepo(cfg *config.Config, logger *slog.Logger) *NumbersRepo {
	return &NumbersRepo{
		data:   loadData(cfg),
		logger: logger,
	}
}

func loadData(cfg *config.Config) []int {
	path, err := filepath.Abs(cfg.InputFilePath)
	if err != nil {
		log.Fatalf("Invalid file path: %v", err)
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	lines := strings.Split(string(bytes), "\n")
	var data []int
	for _, line := range lines {
		if line == "" {
			continue
		}
		num, err := strconv.Atoi(strings.TrimSpace(line))
		if err != nil {
			log.Fatalf("Failed to parse number from line: %q, error: %v", line, err)
		}
		data = append(data, num)
	}

	return data
}

// FindNearestIndex does a binary search on the slice of data for the given target value
// It returns a NumbersResult dto from the http_interface pkg
// Normally, I would not have used the same dto and generally create separate DTOs and mappings as needed
// between different packages
func (r *NumbersRepo) FindNearestIndex(target int) dto.NumbersResult {
	dataLen := len(r.data)

	// Binary search for exact or nearest match
	index := sort.SearchInts(r.data, target)

	// If exact match found
	if index < dataLen && r.data[index] == target {
		r.logger.Info(fmt.Sprintf("Found exact match for target %d", target))
		return dto.NumbersResult{
			Index: index,
			Value: r.data[index],
		}
	}

	foundIdx := -1
	foundValue := -1
	message := fmt.Sprintf("target %d not found", target)
	maxDeviation := target / 10

	// Check if we might be at the start and if so, return
	if index == 0 {
		deviationToNext := r.data[index] - target
		if deviationToNext <= maxDeviation {
			r.logger.Info(fmt.Sprintf("Found next number within deviation %d for target %d", maxDeviation, target))
			foundIdx = index
			foundValue = r.data[index]
			message = fmt.Sprintf("Exact value not found. Found %d next value in accepted deviation.", foundValue)
		}
		return dto.NumbersResult{
			Index:   foundIdx,
			Value:   foundValue,
			Message: &message,
		}
	}

	// Check if we might be at the end and if so, return
	if index == dataLen {
		deviationToPrev := target - r.data[index-1]
		if deviationToPrev <= maxDeviation {
			r.logger.Info(fmt.Sprintf("Found previous number within deviation %d for target %d", maxDeviation, target))
			foundIdx = index - 1
			foundValue = r.data[index-1]
			message = fmt.Sprintf("Exact value not found. Found %d previous value in accepted deviation.", foundValue)
		}
		return dto.NumbersResult{
			Index:   foundIdx,
			Value:   foundValue,
			Message: &message,
		}
	}

	// Check within bounds
	prevIdx := index - 1
	nextIdx := index
	deviationToPrev := target - r.data[prevIdx]
	deviationToNext := r.data[nextIdx] - target

	// Check which value is closer within the 10% deviation if both are viable
	if deviationToPrev <= maxDeviation && deviationToNext <= maxDeviation {
		if deviationToPrev <= deviationToNext {
			r.logger.Info(fmt.Sprintf("Found previous number within deviation %d for target %d", maxDeviation, target))
			foundIdx = prevIdx
			foundValue = r.data[prevIdx]
			message = fmt.Sprintf("Exact value not found. Found %d previous value in accepted deviation.", foundValue)
		} else {
			r.logger.Info(fmt.Sprintf("Found next number within deviation %d for target %d", maxDeviation, target))
			foundIdx = nextIdx
			foundValue = r.data[nextIdx]
			message = fmt.Sprintf("Exact value not found. Found %d next value in accepted deviation.", foundValue)
		}
	} else if deviationToPrev <= maxDeviation {
		r.logger.Info(fmt.Sprintf("Found previous number within deviation %d for target %d", maxDeviation, target))

		// Else, if only the prev is viable within deviation, use that
		foundIdx = prevIdx
		foundValue = r.data[prevIdx]
		message = fmt.Sprintf("Exact value not found. Found %d previous value in accepted deviation.", foundValue)
	} else if deviationToNext <= maxDeviation {
		r.logger.Info(fmt.Sprintf("Found next number within deviation %d for target %d", maxDeviation, target))
		// Else, if only the next is viable within deviation, use that
		foundIdx = nextIdx
		foundValue = r.data[nextIdx]
		message = fmt.Sprintf("Exact value not found. Found %d next value in accepted deviation.", foundValue)
	}

	return dto.NumbersResult{
		Index:   foundIdx,
		Value:   foundValue,
		Message: &message,
	}
}
