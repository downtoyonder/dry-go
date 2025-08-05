package utils

import (
	"reflect"
	"testing"
)

func TestPluck(t *testing.T) {
	// Define the struct type
	type User struct {
		ID   int
		Name string
		Age  int
	}

	t.Run("ValidField", func(t *testing.T) {
		users := []User{
			{ID: 1, Name: "Alice", Age: 25},
			{ID: 2, Name: "Bob", Age: 30},
			{ID: 3, Name: "Charlie", Age: 35},
		}

		// Test plucking the "Name" field
		names, err := Pluck[string](users, "Name")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		expectedNames := []string{"Alice", "Bob", "Charlie"}
		if !reflect.DeepEqual(names, expectedNames) {
			t.Errorf("expected %v, got %v", expectedNames, names)
		}

		// Test plucking the "Age" field
		ages, err := Pluck[int](users, "Age")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		expectedAges := []int{25, 30, 35}
		if !reflect.DeepEqual(ages, expectedAges) {
			t.Errorf("expected %v, got %v", expectedAges, ages)
		}
	})

	t.Run("InvalidField", func(t *testing.T) {
		users := []User{
			{ID: 1, Name: "Alice", Age: 25},
			{ID: 2, Name: "Bob", Age: 30},
			{ID: 3, Name: "Charlie", Age: 35},
		}

		// Test plucking a non-existent field
		_, err := Pluck[string](users, "NonExistentField")
		if err == nil {
			t.Errorf("expected error but got none")
		}
	})

	t.Run("InvalidType", func(t *testing.T) {
		// Define the struct type
		type User struct {
			ID   int
			Name string
			Age  int
		}

		users := []User{
			{ID: 1, Name: "Alice", Age: 25},
			{ID: 2, Name: "Bob", Age: 30},
			{ID: 3, Name: "Charlie", Age: 35},
		}

		// Test plucking the "Name" field but with incorrect type (e.g., expecting int instead of string)
		_, err := Pluck[int](users, "Name")
		if err == nil {
			t.Errorf("expected error due to type mismatch, but got none")
		}
	})

	t.Run("EmptyList", func(t *testing.T) {
		// Define the struct type
		type User struct {
			ID   int
			Name string
			Age  int
		}

		// Test plucking from an empty list
		users := []User{}
		names, err := Pluck[string](users, "Name")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		// Expected result is an empty slice
		if len(names) != 0 {
			t.Errorf("expected empty result, got %v", names)
		}
	})
}

func TestFieldMap(t *testing.T) {
	// Define the struct type
	type User struct {
		ID   int
		Name string
		Age  int
	}

	t.Run("ValidKey", func(t *testing.T) {
		users := []User{
			{ID: 1, Name: "Alice", Age: 25},
			{ID: 2, Name: "Bob", Age: 30},
			{ID: 3, Name: "Charlie", Age: 35},
		}

		// Test mapping by "ID" field
		idMap, _ := FieldMap[int](users, "ID")
		expectedIDMap := map[int]User{
			1: {ID: 1, Name: "Alice", Age: 25},
			2: {ID: 2, Name: "Bob", Age: 30},
			3: {ID: 3, Name: "Charlie", Age: 35},
		}
		if !reflect.DeepEqual(idMap, expectedIDMap) {
			t.Errorf("expected %v, got %v", expectedIDMap, idMap)
		}

		// Test mapping by "Name" field
		nameMap, _ := FieldMap[string](users, "Name")
		expectedNameMap := map[string]User{
			"Alice":   {ID: 1, Name: "Alice", Age: 25},
			"Bob":     {ID: 2, Name: "Bob", Age: 30},
			"Charlie": {ID: 3, Name: "Charlie", Age: 35},
		}
		if !reflect.DeepEqual(nameMap, expectedNameMap) {
			t.Errorf("expected %v, got %v", expectedNameMap, nameMap)
		}
	})

	t.Run("InvalidKey", func(t *testing.T) {
		users := []User{
			{ID: 1, Name: "Alice", Age: 25},
			{ID: 2, Name: "Bob", Age: 30},
			{ID: 3, Name: "Charlie", Age: 35},
		}

		// Test mapping by a non-existent field
		invalidMap, _ := FieldMap[string](users, "NonExistentField")
		if len(invalidMap) != 0 {
			t.Errorf("expected empty map, got %v", invalidMap)
		}
	})

	t.Run("EmptyList", func(t *testing.T) {
		// Define the struct type
		type User struct {
			ID   int
			Name string
			Age  int
		}

		// Test mapping from an empty list
		users := []User{}
		idMap, _ := FieldMap[int](users, "ID")
		if len(idMap) != 0 {
			t.Errorf("expected empty map, got %v", idMap)
		}
	})
}
