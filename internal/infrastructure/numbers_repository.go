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

// NotFoundIndex represent the index that points to a "not found" equivalent
const NotFoundIndex = -1

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
//
// If a number is not found it also looks at the adjacent numbers and checks if they are within a 10% accept deviation
// It returns the one with the closest acceptable deviation
func (r *NumbersRepo) FindNearestIndex(target int) dto.NumbersResult {
	dataLen := len(r.data)

	// Binary search for exact or nearest match
	indexForTarget := sort.SearchInts(r.data, target)

	// If exact match found
	if indexForTarget < dataLen && r.data[indexForTarget] == target {
		r.logger.Info(fmt.Sprintf("Found exact match for target %d", target))
		return populateNumbersResult(
			indexForTarget,
			r.data[indexForTarget],
			nil,
		)
	}

	result := dto.NumbersResult{
		Index:   NotFoundIndex,
		Value:   -1,
		Message: func() *string { ret := fmt.Sprintf("target %d not found", target); return &ret }(),
	}
	maxDeviation := target / 10

	// Check if we might be at the start and if so return or if within deviation from start, return
	if indexForTarget == 0 {
		deviationToNext := r.data[indexForTarget] - target
		if deviationToNext <= maxDeviation {
			r.logger.Info(fmt.Sprintf("Found next number within deviation %d for target %d", maxDeviation, target))
			message := fmt.Sprintf("Exact value not found. Found %d next value in accepted deviation.", r.data[indexForTarget])

			result = populateNumbersResult(
				indexForTarget,
				r.data[indexForTarget],
				&message,
			)
		}
		return result
	}

	// Check if we might be at the end and return or if within deviation from end, return
	if indexForTarget == dataLen {
		deviationToPrev := target - r.data[indexForTarget-1]
		if deviationToPrev <= maxDeviation {
			r.logger.Info(fmt.Sprintf("Found previous number within deviation %d for target %d", maxDeviation, target))
			message := fmt.Sprintf("Exact value not found. Found %d previous value in accepted deviation.", r.data[indexForTarget-1])

			result = populateNumbersResult(
				indexForTarget-1,
				r.data[indexForTarget-1],
				&message,
			)
		}
		return result
	}

	// We did not find it relative to first or last value. Check within bounds
	prevIdx := indexForTarget - 1
	nextIdx := indexForTarget
	deviationToPrev := target - r.data[prevIdx]
	deviationToNext := r.data[nextIdx] - target

	// Check which value is closer within the 10% deviation if both are viable
	if deviationToPrev <= maxDeviation && deviationToNext <= maxDeviation {
		if deviationToPrev <= deviationToNext {
			r.logger.Info(fmt.Sprintf("Found previous number within deviation %d for target %d", maxDeviation, target))
			message := fmt.Sprintf("Exact value not found. Found %d previous value in accepted deviation.", r.data[prevIdx])

			result = populateNumbersResult(
				prevIdx,
				r.data[prevIdx],
				&message,
			)
		} else {
			r.logger.Info(fmt.Sprintf("Found next number within deviation %d for target %d", maxDeviation, target))
			message := fmt.Sprintf("Exact value not found. Found %d next value in accepted deviation.", r.data[nextIdx])
			result = populateNumbersResult(
				nextIdx,
				r.data[nextIdx],
				&message,
			)
		}
	} else if deviationToPrev <= maxDeviation {
		r.logger.Info(fmt.Sprintf("Found previous number within deviation %d for target %d", maxDeviation, target))

		// Else, if only the prev is viable within deviation, use that
		message := fmt.Sprintf("Exact value not found. Found %d previous value in accepted deviation.", r.data[prevIdx])
		result = populateNumbersResult(
			prevIdx,
			r.data[prevIdx],
			&message,
		)
	} else if deviationToNext <= maxDeviation {
		r.logger.Info(fmt.Sprintf("Found next number within deviation %d for target %d", maxDeviation, target))
		// Else, if only the next is viable within deviation, use that
		message := fmt.Sprintf("Exact value not found. Found %d next value in accepted deviation.", r.data[nextIdx])
		result = populateNumbersResult(
			nextIdx,
			r.data[nextIdx],
			&message,
		)
	}

	return result
}

func populateNumbersResult(foundIdx, foundValue int, message *string) dto.NumbersResult {
	return dto.NumbersResult{
		Index:   foundIdx,
		Value:   foundValue,
		Message: message,
	}
}
