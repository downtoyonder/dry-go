package utils

import (
	mapset "github.com/deckarep/golang-set/v2"
)

type SelectFn[StructT any, FieldT any] func(StructT) (FieldT, bool)

func SelectAll[FieldT, StructT any](f func(structure StructT) FieldT) SelectFn[StructT, FieldT] {
	return func(structure StructT) (FieldT, bool) {
		return f(structure), true
	}
}

// PluckFn extracts values using a field selector function.
// This provides compile-time safety - if the field changes, the code won't compile.
// This is the RECOMMENDED approach for extracting fields from structs.
//
// Example:
//
//	type User struct { ID int; Name string }
//	users := []User{{ID: 1, Name: "Alice"}, {ID: 2, Name: "Bob"}}
//	names := PluckFn(users, func(u User) string { return u.Name })
//	ids := PluckFn(users, func(u User) int { return u.ID })
func PluckFn[StructT any, FieldT any](list []StructT, fn SelectFn[StructT, FieldT]) []FieldT {
	result := make([]FieldT, len(list))
	for i, item := range list {
		field, add := fn(item)
		if !add {
			continue
		}
		result[i] = field
	}
	return result
}

func PluckUniqFn[StructT any, FieldT comparable](list []StructT, fn SelectFn[StructT, FieldT]) []FieldT {
	return Uniq(PluckFn(list, fn))
}

// FieldStructMapFn creates a map from a slice of structs, using the value of a specified field as the key.
func FieldMapStructFn[FieldT comparable, StructT any](list []StructT, fn SelectFn[StructT, FieldT]) map[FieldT]StructT {
	result := make(map[FieldT]StructT)
	for _, item := range list {
		field, add := fn(item)
		if !add {
			continue
		}
		result[field] = item
	}
	return result
}

func FieldMapFieldFn[KeyT comparable, ValueT, StructT any](slice []StructT, keySel SelectFn[StructT, KeyT], valueSel SelectFn[StructT, ValueT]) map[KeyT]ValueT {
	result := make(map[KeyT]ValueT)
	for _, item := range slice {
		key, addKey := keySel(item)
		if !addKey {
			continue
		}
		value, addValue := valueSel(item)
		if !addValue {
			continue
		}
		result[key] = value
	}
	return result
}

// SetCmp compares two slices and returns the elements that are new, overlapped, and deleted.
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
//	new, overlapped, deleted := SetCmp(current, target)
//	// new: [4]
//	// overlapped: [2 3]
//	// deleted: [1]
func SetCmp[E comparable](current, target []E) (added, overlapped, deleted []E) {
	currentSet := mapset.NewSet(current...)
	targetSet := mapset.NewSet(target...)

	added = targetSet.Difference(currentSet).ToSlice()
	deleted = currentSet.Difference(targetSet).ToSlice()
	overlapped = currentSet.Intersect(targetSet).ToSlice()

	return added, overlapped, deleted
}

// Uniq returns a slice containing only the unique elements from the input slice 'list'.
// The order of elements in the returned slice is not guaranteed.
// Elements must be of a comparable type.
//
// Example:
//
//	names := []string{"Alice", "Bob", "Alice", "Eve"}
//	uniqueNames := Uniq(names)
//	// uniqueNames might be: []string{"Alice", "Bob", "Eve"}
func Uniq[T comparable](list []T) []T {
	m := make(map[T]struct{}, len(list))
	result := make([]T, 0, len(list))
	for _, item := range list {
		if _, ok := m[item]; !ok {
			m[item] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}

type Set[T comparable] map[T]struct{}

func NewSet[T comparable](items ...T) Set[T] {
	s := make(Set[T])
	for _, item := range items {
		s[item] = struct{}{}
	}
	return s
}

func (s Set[T]) Add(item T) {
	s[item] = struct{}{}
}

func (s Set[T]) Remove(item T) {
	delete(s, item)
}

func (s Set[T]) Contain(item T) bool {
	_, exists := s[item]
	return exists
}

func (s Set[T]) ToSlice() []T {
	result := make([]T, 0, len(s))
	for item := range s {
		result = append(result, item)
	}
	return result
}
