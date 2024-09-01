package handlers

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateRequestedFilters(t *testing.T) {
	tests := []struct {
		name         string
		inputFilters []string
		wantFilters  []string
		wantError    error
	}{
		{
			name:         "Valid filters",
			inputFilters: []string{"namespace", "phase", "labels"},
			wantFilters:  []string{"namespace", "phase", "labels"},
			wantError:    nil,
		},
		{
			name:         "Mixed valid and empty filters",
			inputFilters: []string{"namespace", "", "phase", " ", "labels"},
			wantFilters:  []string{"namespace", "phase", "labels"},
			wantError:    nil,
		},
		{
			name:         "All empty filters",
			inputFilters: []string{"", " ", "\t", "\n"},
			wantFilters:  nil,
			wantError:    errors.New("no valid filters provided"),
		},
	}

	for _, tt := range tests {
		ttp := tt
		t.Run(ttp.name, func(t *testing.T) {
			filters, err := validateRequestedFilters(ttp.inputFilters)
			if err != nil {
				assert.EqualError(t, ttp.wantError, err.Error(), "Unexpected error message")
			}

			assert.Equal(t, ttp.wantFilters, filters, "Unexpected filters")
		})
	}
}
