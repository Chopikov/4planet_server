package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewsType_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		newsType NewsType
		expected bool
	}{
		{"valid achievement", NewsTypeAchievement, true},
		{"valid invite", NewsTypeInvite, true},
		{"valid update", NewsTypeUpdate, true},
		{"invalid empty", "", false},
		{"invalid random", "random", false},
		{"invalid mixed case", "Achievement", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.newsType.IsValid()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseNewsType(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		expected   NewsType
		expectedOk bool
	}{
		{"valid achievement", "achievement", NewsTypeAchievement, true},
		{"valid invite", "invite", NewsTypeInvite, true},
		{"valid update", "update", NewsTypeUpdate, true},
		{"invalid empty", "", "", false},
		{"invalid random", "random", "random", false},
		{"invalid mixed case", "Achievement", "Achievement", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := ParseNewsType(tt.input)
			assert.Equal(t, tt.expected, result)
			assert.Equal(t, tt.expectedOk, ok)
		})
	}
}
