package infrastructure

import (
	"fmt"
	"github.com/OmgAbear/gosolve/internal/http_interface/dto"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindNearestIndex(t *testing.T) {
	tests := []struct {
		name           string
		data           []int
		target         int
		expectedResult dto.NumbersResult
	}{
		{
			name:   "Target is found",
			data:   []int{11, 15},
			target: 11,
			expectedResult: dto.NumbersResult{
				Index: 0,
				Value: 11,
			},
		},
		{
			name:   "Target is lower than first idx but within deviation",
			data:   []int{11, 15},
			target: 10,
			expectedResult: dto.NumbersResult{
				Index: 0,
				Value: 11,
				Message: func() *string {
					ret := fmt.Sprintf("Exact value not found. Found %d next value in accepted deviation.", 11)
					return &ret
				}(),
			},
		},
		{
			name:   "Target is bigger than last value but within deviation",
			data:   []int{5, 10, 15, 19},
			target: 20,
			expectedResult: dto.NumbersResult{
				Index: 3,
				Value: 19,
				Message: func() *string {
					ret := fmt.Sprintf("Exact value not found. Found %d previous value in accepted deviation.", 19)
					return &ret
				}(),
			},
		},
		{
			name:   "Target is lower than first idx and outside deviation",
			data:   []int{11, 15},
			target: 1,
			expectedResult: dto.NumbersResult{
				Index: -1,
				Value: -1,
				Message: func() *string {
					ret := fmt.Sprintf("target %d not found", 1)
					return &ret
				}(),
			},
		},
		{
			name:   "Target is bigger than last idx and outside deviation",
			data:   []int{5, 10, 15, 19},
			target: 25,
			expectedResult: dto.NumbersResult{
				Index: -1,
				Value: -1,
				Message: func() *string {
					ret := fmt.Sprintf("target %d not found", 25)
					return &ret
				}(),
			},
		},
		{
			name:   "Target not found, but lower neighbour found that is within deviation",
			data:   []int{5, 10, 15, 19, 24},
			target: 20,
			expectedResult: dto.NumbersResult{
				Index: 3,
				Value: 19,
				Message: func() *string {
					ret := fmt.Sprintf("Exact value not found. Found %d previous value in accepted deviation.", 19)
					return &ret
				}(),
			},
		},
		{
			name:   "Target not found, but higher neighbour found that is within deviation",
			data:   []int{5, 10, 15, 19, 26},
			target: 25,
			expectedResult: dto.NumbersResult{
				Index: 4,
				Value: 26,
				Message: func() *string {
					ret := fmt.Sprintf("Exact value not found. Found %d next value in accepted deviation.", 26)
					return &ret
				}(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NumbersRepo{data: tt.data}
			result := repo.FindNearestIndex(tt.target)

			assert.Equal(t, tt.expectedResult.Index, result.Index)
			assert.Equal(t, tt.expectedResult.Value, result.Value)

			if tt.expectedResult.Message != nil && result.Message == nil ||
				tt.expectedResult.Message == nil && result.Message != nil {
				t.Error(fmt.Sprintf("expected message: %v, received message: %v", tt.expectedResult.Message != nil, result.Message != nil))
				t.Fail()
			}

			if tt.expectedResult.Message != nil && result.Message != nil {
				assert.Equal(t, *tt.expectedResult.Message, *result.Message)
			}
		})
	}
}
