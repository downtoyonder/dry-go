package utils

import (
	"fmt"
	"sort"
)

type User struct {
	ID   int
	Name string
	Age  int
}

// ExampleSelectAll demonstrates how to use SelectAll to create a selector that always includes the field.
func ExampleSelectAll() {
	users := []User{
		{ID: 1, Name: "Alice", Age: 30},
		{ID: 2, Name: "Bob", Age: 25},
	}

	// SelectAll wraps a simple field extractor
	selector := SelectAll(func(u User) string { return u.Name })

	names := PluckFn(users, selector)
	fmt.Println(names)
	// Output: [Alice Bob]
}

// ExamplePluckFn demonstrates how to extract field values from a slice of structs.
func ExamplePluckFn() {
	users := []User{
		{ID: 1, Name: "Alice", Age: 30},
		{ID: 2, Name: "Bob", Age: 25},
		{ID: 3, Name: "Charlie", Age: 35},
	}

	// Extract names
	names := PluckFn(users, SelectAll(func(u User) string { return u.Name }))
	fmt.Println(names)
	// Output: [Alice Bob Charlie]
}

// ExamplePluckFn_conditional demonstrates conditional field extraction.
func ExamplePluckFn_conditional() {
	users := []User{
		{ID: 1, Name: "Alice", Age: 30},
		{ID: 2, Name: "Bob", Age: 25},
		{ID: 3, Name: "Charlie", Age: 35},
	}

	// Extract names only for users older than 28
	names := PluckFn(users, func(u User) (string, bool) {
		if u.Age > 28 {
			return u.Name, true
		}
		return "", false
	})

	// Filter out empty strings that were skipped
	var filtered []string
	for _, name := range names {
		if name != "" {
			filtered = append(filtered, name)
		}
	}
	fmt.Println(filtered)
	// Output: [Alice Charlie]
}

// ExamplePluckFn_ids demonstrates extracting integer IDs.
func ExamplePluckFn_ids() {
	users := []User{
		{ID: 1, Name: "Alice", Age: 30},
		{ID: 2, Name: "Bob", Age: 25},
		{ID: 3, Name: "Charlie", Age: 35},
	}

	ids := PluckFn(users, SelectAll(func(u User) int { return u.ID }))
	fmt.Println(ids)
	// Output: [1 2 3]
}

// ExamplePluckUniqFn demonstrates extracting unique field values.
func ExamplePluckUniqFn() {
	users := []User{
		{ID: 1, Name: "Alice", Age: 30},
		{ID: 2, Name: "Bob", Age: 25},
		{ID: 3, Name: "Alice", Age: 28},
		{ID: 4, Name: "Bob", Age: 35},
	}

	// Extract unique names
	uniqueNames := PluckUniqFn(users, SelectAll(func(u User) string { return u.Name }))

	// Sort for consistent output
	sort.Strings(uniqueNames)
	fmt.Println(uniqueNames)
	// Output: [Alice Bob]
}

// ExampleFieldMapStructFn demonstrates creating a map from field to struct.
func ExampleFieldMapStructFn() {
	users := []User{
		{ID: 1, Name: "Alice", Age: 30},
		{ID: 2, Name: "Bob", Age: 25},
		{ID: 3, Name: "Charlie", Age: 35},
	}

	// Create a map from ID to User
	userByID := FieldMapStructFn(users, SelectAll(func(u User) int { return u.ID }))

	fmt.Println(userByID[2].Name)
	// Output: Bob
}

// ExampleFieldMapStructFn_byName demonstrates creating a map using name as key.
func ExampleFieldMapStructFn_byName() {
	users := []User{
		{ID: 1, Name: "Alice", Age: 30},
		{ID: 2, Name: "Bob", Age: 25},
	}

	// Create a map from Name to User
	userByName := FieldMapStructFn(users, SelectAll(func(u User) string { return u.Name }))

	fmt.Println(userByName["Alice"].ID)
	// Output: 1
}

// ExampleFieldMapFieldFn demonstrates creating a map from one field to another.
func ExampleFieldMapFieldFn() {
	users := []User{
		{ID: 1, Name: "Alice", Age: 30},
		{ID: 2, Name: "Bob", Age: 25},
		{ID: 3, Name: "Charlie", Age: 35},
	}

	// Create a map from ID to Name
	nameByID := FieldMapFieldFn(
		users,
		SelectAll(func(u User) int { return u.ID }),
		SelectAll(func(u User) string { return u.Name }),
	)

	fmt.Println(nameByID[2])
	// Output: Bob
}

// ExampleFieldMapFieldFn_conditional demonstrates conditional mapping.
func ExampleFieldMapFieldFn_conditional() {
	users := []User{
		{ID: 1, Name: "Alice", Age: 30},
		{ID: 2, Name: "Bob", Age: 25},
		{ID: 3, Name: "Charlie", Age: 35},
	}

	// Map Name to Age, but only for users older than 28
	ageByName := FieldMapFieldFn(
		users,
		func(u User) (string, bool) {
			if u.Age > 28 {
				return u.Name, true
			}
			return "", false
		},
		SelectAll(func(u User) int { return u.Age }),
	)

	// Get sorted names for consistent output
	names := make([]string, 0, len(ageByName))
	for name := range ageByName {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		fmt.Printf("%s: %d\n", name, ageByName[name])
	}
	// Output:
	// Alice: 30
	// Charlie: 35
}

// ExampleSetCmp demonstrates comparing two slices to find new, overlapped, and deleted elements.
func ExampleSetCmp() {
	current := []int{1, 2, 3, 4}
	target := []int{3, 4, 5, 6}

	new, overlapped, deleted := SetCmp(current, target)

	// Sort for consistent output
	sort.Ints(new)
	sort.Ints(overlapped)
	sort.Ints(deleted)

	fmt.Println("New:", new)
	fmt.Println("Overlapped:", overlapped)
	fmt.Println("Deleted:", deleted)
	// Output:
	// New: [5 6]
	// Overlapped: [3 4]
	// Deleted: [1 2]
}

// ExampleSetCmp_strings demonstrates SetCmp with string slices.
func ExampleSetCmp_strings() {
	current := []string{"apple", "banana", "cherry"}
	target := []string{"banana", "cherry", "date"}

	new, overlapped, deleted := SetCmp(current, target)

	// Sort for consistent output
	sort.Strings(new)
	sort.Strings(overlapped)
	sort.Strings(deleted)

	fmt.Println("New:", new)
	fmt.Println("Overlapped:", overlapped)
	fmt.Println("Deleted:", deleted)
	// Output:
	// New: [date]
	// Overlapped: [banana cherry]
	// Deleted: [apple]
}

// ExampleUniq demonstrates removing duplicate elements from a slice.
func ExampleUniq() {
	numbers := []int{1, 2, 2, 3, 3, 3, 4, 5, 5}
	unique := Uniq(numbers)

	// Sort for consistent output (Uniq doesn't guarantee order)
	sort.Ints(unique)
	fmt.Println(unique)
	// Output: [1 2 3 4 5]
}

// ExampleUniq_strings demonstrates Uniq with strings.
func ExampleUniq_strings() {
	words := []string{"hello", "world", "hello", "go", "world"}
	unique := Uniq(words)

	// Sort for consistent output
	sort.Strings(unique)
	fmt.Println(unique)
	// Output: [go hello world]
}

// ExampleNewSet demonstrates creating a new set and checking membership.
func ExampleNewSet() {
	set := NewSet(1, 2, 3, 4, 5)

	fmt.Println(set.Contain(3))
	fmt.Println(set.Contain(10))
	// Output:
	// true
	// false
}

// ExampleSet_Add demonstrates adding elements to a set.
func ExampleSet_Add() {
	set := NewSet[string]()

	set.Add("apple")
	set.Add("banana")
	set.Add("apple") // Duplicate, won't be added again

	slice := set.ToSlice()
	sort.Strings(slice) // Sort for consistent output
	fmt.Println(slice)
	// Output: [apple banana]
}

// ExampleSet_Remove demonstrates removing elements from a set.
func ExampleSet_Remove() {
	set := NewSet("apple", "banana", "cherry")

	set.Remove("banana")

	fmt.Println(set.Contain("banana"))
	fmt.Println(set.Contain("apple"))
	// Output:
	// false
	// true
}

// ExampleSet_Contain demonstrates checking if an element exists in a set.
func ExampleSet_Contain() {
	set := NewSet(10, 20, 30, 40)

	fmt.Println(set.Contain(20))
	fmt.Println(set.Contain(25))
	fmt.Println(set.Contain(40))
	// Output:
	// true
	// false
	// true
}

// ExampleSet_ToSlice demonstrates converting a set to a slice.
func ExampleSet_ToSlice() {
	set := NewSet("dog", "cat", "bird")

	slice := set.ToSlice()
	sort.Strings(slice) // Sort for consistent output since sets are unordered
	fmt.Println(slice)
	// Output: [bird cat dog]
}

// ExampleSet_workflow demonstrates a complete workflow with a set.
func ExampleSet_workflow() {
	// Create a set to track unique visitor IDs
	visitors := NewSet[int]()

	// Simulate visitors arriving
	visitors.Add(101)
	visitors.Add(102)
	visitors.Add(101) // Duplicate visitor
	visitors.Add(103)

	fmt.Println("Total unique visitors:", len(visitors))

	// Check if a specific visitor has arrived
	fmt.Println("Visitor 102 arrived:", visitors.Contain(102))
	fmt.Println("Visitor 999 arrived:", visitors.Contain(999))

	// Remove a visitor
	visitors.Remove(102)
	fmt.Println("After 102 left:", visitors.Contain(102))

	// Get all current visitors
	currentVisitors := visitors.ToSlice()
	sort.Ints(currentVisitors)
	fmt.Println("Current visitors:", currentVisitors)
	// Output:
	// Total unique visitors: 3
	// Visitor 102 arrived: true
	// Visitor 999 arrived: false
	// After 102 left: false
	// Current visitors: [101 103]
}
