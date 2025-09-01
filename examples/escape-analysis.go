// Package main demonstrates escape analysis patterns in Go
// Run with: go run examples/escape-analysis.go
// See escape analysis: go build -gcflags='-m' examples/escape-analysis.go
package main

import "fmt"

// =============================================================================
// EXAMPLES THAT STAY ON STACK (FAST)
// =============================================================================

// staysOnStack1 - variable used only locally
func staysOnStack1() {
	student := Student{
		ID:    1,
		Name:  "Alice",
		Grade: 85.5,
		Age:   20,
	}

	// Just using the student locally - stays on stack
	fmt.Printf("Student %s has grade %.1f\n", student.Name, student.Grade)
} // student is automatically cleaned up here

// staysOnStack2 - returning value (not pointer)
func staysOnStack2() Student {
	student := Student{
		ID:    2,
		Name:  "Bob",
		Grade: 92.0,
		Age:   19,
	}

	return student // Returning value - original can stay on stack
}

// staysOnStack3 - passing to function temporarily
func staysOnStack3() {
	student := Student{ID: 3, Name: "Charlie", Grade: 78.0, Age: 21}

	// Pass by value - function gets a copy
	processStudent(student)

	// Original student still exists on stack
	fmt.Printf("Original student unchanged: %+v\n", student)
}

func processStudent(s Student) {
	// This works with a copy - safe and stack-friendly
	s.Grade += 10 // Only changes the copy
	fmt.Printf("Inside function (copy): %+v\n", s)
}

// =============================================================================
// EXAMPLES THAT ESCAPE TO HEAP (SLOWER)
// =============================================================================

// escapesToHeap1 - returning pointer
func escapesToHeap1() *Student {
	student := Student{
		ID:    4,
		Name:  "Diana",
		Grade: 95.0,
		Age:   20,
	}

	// Compiler thinks: "This pointer needs to survive after function ends!"
	return &student // ESCAPES TO HEAP
}

// escapesToHeap2 - storing in interface{}
func escapesToHeap2() {
	student := Student{
		ID:    5,
		Name:  "Eve",
		Grade: 88.0,
		Age:   22,
	}

	// fmt.Println takes interface{} - forces escape
	fmt.Println("Student via fmt.Println:", student) // ESCAPES TO HEAP
}

// escapesToHeap3 - large struct (simulated with array)
func escapesToHeap3() {
	// Large local variable might escape to heap
	largeData := [10000]int{} // If too big for stack
	largeData[0] = 42
	fmt.Printf("Large data first element: %d\n", largeData[0])
}

// escapesToHeap4 - stored in global or long-lived structure
var globalStudent *Student

func escapesToHeap4() {
	student := Student{
		ID:    6,
		Name:  "Frank",
		Grade: 82.0,
		Age:   23,
	}

	globalStudent = &student // Stored globally - must escape to heap
}

// =============================================================================
// COMPARISON FUNCTIONS
// =============================================================================

// calculateAverageStack - no heap allocations
func calculateAverageStacks(students []Student) float64 {
	var total float64

	// Process in place - nothing escapes
	for _, student := range students {
		total += student.Grade
	}

	return total / float64(len(students))
}

// calculateAverageHeap - creates intermediate slice
func calculateAverageHeaps(students []Student) float64 {
	var grades []float64 // This slice grows on heap

	for _, student := range students {
		grades = append(grades, student.Grade) // Reallocations possible
	}

	var total float64
	for _, grade := range grades {
		total += grade
	}

	return total / float64(len(grades))
}

// =============================================================================
// DEMONSTRATION
// =============================================================================

func main() {
	fmt.Println("Go Escape Analysis Examples")
	fmt.Println("===========================")

	fmt.Println("\n--- Variables That STAY ON STACK ---")
	staysOnStack1()

	stackStudent := staysOnStack2()
	fmt.Printf("Returned by value: %s (Grade: %.1f)\n", stackStudent.Name, stackStudent.Grade)

	staysOnStack3()

	fmt.Println("\n--- Variables That ESCAPE TO HEAP ---")
	heapStudent := escapesToHeap1()
	fmt.Printf("Returned by pointer: %s (Grade: %.1f)\n", heapStudent.Name, heapStudent.Grade)

	escapesToHeap2() // Uses fmt.Println

	escapesToHeap3() // Large array

	escapesToHeap4() // Global storage
	fmt.Printf("Global student: %+v\n", globalStudent)

	fmt.Println("\n--- Performance Impact ---")

	// Create test data
	testStudents := []Student{
		{ID: 1, Name: "Test1", Grade: 85.0, Age: 20},
		{ID: 2, Name: "Test2", Grade: 90.0, Age: 21},
		{ID: 3, Name: "Test3", Grade: 78.0, Age: 19},
		{ID: 4, Name: "Test4", Grade: 88.0, Age: 22},
		{ID: 5, Name: "Test5", Grade: 92.0, Age: 20},
	}

	// Compare approaches
	stackAvg := calculateAverageStacks(testStudents)
	heapAvg := calculateAverageHeaps(testStudents)

	fmt.Printf("Stack calculation: %.2f (no extra allocations)\n", stackAvg)
	fmt.Printf("Heap calculation: %.2f (creates intermediate slice)\n", heapAvg)

	fmt.Println("\n=== To See Escape Analysis ===")
	fmt.Println("Run: go build -gcflags='-m' examples/escape_analysis.go")
	fmt.Println("\nLook for output like:")
	fmt.Println("  'moved to heap: student'")
	fmt.Println("  '&student escapes to heap'")
	fmt.Println("\nThis tells you which variables Go puts on the heap!")
}
