package utils

import (
	"fmt"
	"maps"
	"reflect"
	"slices"

	mapset "github.com/deckarep/golang-set/v2"
)

// Pluck extracts the values of a specified field from a slice of structs.
//
// The field name is provided as a string, its type is specified by the type parameter FieldT, and the result is a slice of type FieldT.
// Returns an error if any item is not a struct, the field does not exist, or the field type does not match FieldT.
// If the field is not exported (i.e., starts with a lowercase letter), it cannot be accessed via reflection.
//
// If the field is not unique across the structs, the resulting slice will contain all values in the order they appear in the input slice.
//
// This function is useful for extracting a specific field from a slice of structs, such as IDs
// Example:
//
//	type User struct { ID int; Name string }
//	users := []User{{ID: 1, Name: "Alice"}, {ID: 1, Name: "Alice"}, {ID: 2, Name: "Bob"}}
//	ids, err := Pluck[int](users, "ID") // ids = [1, 1, 2]
func Pluck[FieldT, StructT any](list []StructT, fieldName string) ([]FieldT, error) {
	result := make([]FieldT, 0, len(list))

	for _, item := range list {
		v := reflect.ValueOf(item)
		// If the item is a pointer, dereference it to get the underlying value
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		// Check if the item is a struct. If not, return an error
		if v.Kind() != reflect.Struct {
			return nil, fmt.Errorf("Pluck: expected struct but got %T", item)
		}

		field := v.FieldByName(fieldName)
		if !field.IsValid() {
			return nil, fmt.Errorf("Pluck: no such field %q in %T", fieldName, item)
		}
		if !field.CanInterface() {
			return nil, fmt.Errorf("Pluck: field %q in %T cannot be accessed", fieldName, item)
		}

		val, ok := field.Interface().(FieldT)
		if !ok {
			return nil, fmt.Errorf("Pluck: field %q in %T is not of type %T", fieldName, item, *new(FieldT))
		}
		result = append(result, val)
	}

	return result, nil
}

// FieldMap creates a map from a slice of structs, using the value of a specified field as the key.
//
// The field name is provided as a string, its type is specified by the type parameter FieldT, and the result is a map[FieldT]StructT.
// Returns an error if any item is not a struct, the field does not exist, or the field type does not match FieldT.
// If the field is not exported (i.e., starts with a lowercase letter), it cannot be accessed via reflection.
//
// If the field is not unique across the structs, later values will overwrite earlier ones in the resulting map.
// Example:
//
//	type User struct { ID int; Name string }
//	users := []User{{ID: 1, Name: "Alice"}, {ID: 2, Name: "Bob"}}
//	m, err := FieldMap[int](users, "ID") // m = map[1:User{...}, 2:User{...}]
func FieldMap[FieldT comparable, StructT any](list []StructT, fieldName string) (map[FieldT]StructT, error) {
	result := make(map[FieldT]StructT, len(list))

	for _, item := range list {
		v := reflect.ValueOf(item)
		// If the item is a pointer, dereference it to get the underlying value
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		// Check if the item is a struct. If not, return an error
		if v.Kind() != reflect.Struct {
			return nil, fmt.Errorf("FieldMap: expected struct but got %T", item)
		}

		field := v.FieldByName(fieldName)
		if !field.IsValid() {
			return nil, fmt.Errorf("FieldMap: no such field %q in %T", fieldName, item)
		}
		if !field.CanInterface() {
			return nil, fmt.Errorf("FieldMap: field %q in %T cannot be accessed", fieldName, item)
		}

		val, ok := field.Interface().(FieldT)
		if !ok {
			return nil, fmt.Errorf("FieldMap: field %q in %T is not of type %T", fieldName, item, *new(FieldT))
		}
		result[val] = item
	}

	return result, nil
}

// SetCompare compares two slices and returns the elements that are new, overlapped, and deleted.
// It takes two slices of comparable elements: `current` and `target`.
// Returns three slices:
//   - new: elements present in `target` but not in `current`
//   - overlapped: elements present in both `current` and `target`
//   - deleted: elements present in `current` but not in `target`
//
// Example:
//
//	current := []int{1, 2, 3}
//	target := []int{2, 3, 4}
//	new, overlapped, deleted := SetCompare(current, target)
//	// new: [4]
//	// overlapped: [2 3]
//	// deleted: [1]
func SetCompare[E comparable](current, target []E) ([]E, []E, []E) {
	currentSet := mapset.NewSet(current...)
	targetSet := mapset.NewSet(target...)

	new := targetSet.Difference(currentSet).ToSlice()
	deleted := currentSet.Difference(targetSet).ToSlice()
	overlapped := currentSet.Intersect(targetSet).ToSlice()

	return new, overlapped, deleted
}

// Unique returns a slice containing only the unique elements from the input slice 'list'.
// The order of elements in the returned slice is not guaranteed.
// Elements must be of a comparable type.
//
// Example:
//
//	names := []string{"Alice", "Bob", "Alice", "Eve"}
//	uniqueNames := Unique(names)
//	// uniqueNames might be: []string{"Alice", "Bob", "Eve"}
func Unique[T comparable](list []T) []T {
	m := make(map[T]struct{}, len(list))
	for _, item := range list {
		m[item] = struct{}{}
	}

	return slices.Collect(maps.Keys(m))
}
