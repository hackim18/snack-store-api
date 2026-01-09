package test

import (
	"testing"

	"snack-store-api/internal/entity"
)

func TestPointsEarned(t *testing.T) {
	testCases := []struct {
		name       string
		totalPrice int
		expected   int
	}{
		{name: "zero", totalPrice: 0, expected: 0},
		{name: "below_threshold", totalPrice: 999, expected: 0},
		{name: "exact_threshold", totalPrice: 1000, expected: 1},
		{name: "multiple", totalPrice: 2000, expected: 2},
		{name: "with_remainder", totalPrice: 10500, expected: 10},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := entity.PointsEarned(tc.totalPrice)
			if got != tc.expected {
				t.Fatalf("expected %d, got %d", tc.expected, got)
			}
		})
	}
}

func TestPointsCost(t *testing.T) {
	testCases := []struct {
		name     string
		size     string
		expected int
	}{
		{name: "small", size: entity.SizeSmall, expected: 200},
		{name: "medium", size: entity.SizeMedium, expected: 300},
		{name: "large", size: entity.SizeLarge, expected: 500},
		{name: "with_whitespace", size: " Small ", expected: 200},
		{name: "unknown", size: "Extra", expected: 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := entity.PointsCost(tc.size)
			if got != tc.expected {
				t.Fatalf("expected %d, got %d", tc.expected, got)
			}
		})
	}
}
