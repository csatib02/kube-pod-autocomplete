package filter

import (
	"testing"

	"github.com/csatib02/kube-pod-autocomplete/internal/services/autocomplete/model"
	"github.com/stretchr/testify/assert"
)

var optionsTest = Options{}

func TestRemoveDuplicateValues(t *testing.T) {
	tests := []struct {
		name             string
		inputsuggestions []model.Suggestion
		wantSuggestions  []model.Suggestion
	}{
		{
			name: "No duplicates",
			inputsuggestions: []model.Suggestion{
				{Key: "namespace", Values: []string{"value1", "value2"}},
				{Key: "label", Values: []string{"value3", "value4"}},
			},
			wantSuggestions: []model.Suggestion{
				{Key: "namespace", Values: []string{"value1", "value2"}},
				{Key: "label", Values: []string{"value3", "value4"}},
			},
		},
		{
			name: "With duplicates",
			inputsuggestions: []model.Suggestion{
				{Key: "namespace", Values: []string{"value1", "value2", "value1"}},
				{Key: "label", Values: []string{"value3", "value4", "value3"}},
			},
			wantSuggestions: []model.Suggestion{
				{Key: "namespace", Values: []string{"value1", "value2"}},
				{Key: "label", Values: []string{"value3", "value4"}},
			},
		},
		{
			name: "All duplicates",
			inputsuggestions: []model.Suggestion{
				{Key: "namespace", Values: []string{"value1", "value1", "value1"}},
			},
			wantSuggestions: []model.Suggestion{
				{Key: "namespace", Values: []string{"value1"}},
			},
		},
	}

	for _, tt := range tests {
		ttp := tt
		t.Run(ttp.name, func(t *testing.T) {
			optionsTest.RemoveDuplicateValues(&ttp.inputsuggestions)
			assert.Equal(t, ttp.wantSuggestions, ttp.inputsuggestions)
		})
	}
}

func TestRemoveIgnoredKeys(t *testing.T) {
	tests := []struct {
		name             string
		inputsuggestions []model.Suggestion
		wantSuggestions  []model.Suggestion
	}{
		{
			name: "No ignored keys",
			inputsuggestions: []model.Suggestion{
				{Key: "namespace", Values: []string{"value1", "value2"}},
				{Key: "label", Values: []string{"value3", "value4"}},
			},
			wantSuggestions: []model.Suggestion{
				{Key: "namespace", Values: []string{"value1", "value2"}},
				{Key: "label", Values: []string{"value3", "value4"}},
			},
		},
		{
			name: "With ignored keys",
			inputsuggestions: []model.Suggestion{
				{Key: "annotations:kubectl.kubernetes.io/last-applied-configuration", Values: []string{"value1"}},
				{Key: "label", Values: []string{"value2"}},
			},
			wantSuggestions: []model.Suggestion{
				{Key: "label", Values: []string{"value2"}},
			},
		},
		{
			name: "All ignored keys",
			inputsuggestions: []model.Suggestion{
				{Key: "annotations:kubectl.kubernetes.io/last-applied-configuration", Values: []string{"value1"}},
			},
			wantSuggestions: []model.Suggestion{},
		},
	}

	for _, tt := range tests {
		ttp := tt
		t.Run(ttp.name, func(t *testing.T) {
			optionsTest.RemoveIgnoredKeys(&ttp.inputsuggestions)
			assert.Equal(t, ttp.wantSuggestions, ttp.inputsuggestions)
		})
	}
}
