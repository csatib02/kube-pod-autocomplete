package autocomplete

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/csatib02/kube-pod-autocomplete/internal/services/autocomplete/model"
)

var serviceTest = Service{}

func TestProcessMapSuggestion(t *testing.T) {
	tests := []struct {
		name             string
		filterName       string
		mapData          map[string][]string
		inputSuggestions []model.Suggestion
		wantSuggestions  []model.Suggestion
	}{
		{
			name:             "Valid Map Data",
			filterName:       "labels",
			mapData:          map[string][]string{"app": {"nginx", "redis"}},
			inputSuggestions: []model.Suggestion{},
			wantSuggestions: []model.Suggestion{
				{Key: "labels:app", Values: []string{"nginx", "redis"}},
			},
		},
		{
			name:             "Empty Map Data",
			filterName:       "labels",
			mapData:          map[string][]string{},
			inputSuggestions: []model.Suggestion{},
			wantSuggestions:  []model.Suggestion{},
		},
	}

	for _, tt := range tests {
		ttp := tt
		t.Run(ttp.name, func(t *testing.T) {
			serviceTest.processMapSuggestion(&ttp.inputSuggestions, ttp.filterName, &ttp.mapData)
			assert.Equal(t, ttp.wantSuggestions, ttp.inputSuggestions)
		})
	}
}
