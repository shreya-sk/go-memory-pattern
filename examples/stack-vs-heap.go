// Package main demonstrates stack vs heap performance differences
// This file contains functions only - no main function
// Run all examples with: go run examples/*.go
package main

import (
	"fmt"
	"runtime"
	"time"
)

// =============================================================================
// STACK-FRIENDLY APPROACHES (FAST)
// =============================================================================

// calculateAverageStack processes data without heap allocations
func calculateAverageStack(students []Student) float64 {
	var total float64
	var count int

	// Process in place - no extra allocations
	for _, student := range students {
		total += student.Grade
		count++
	}

	if count == 0 {
		return 0
	}
	return total / float64(count)
}

// createStudentByValue returns a value, not pointer
func createStudentByValue(id int64, name string, grade float64) Student {
	return Student{
		ID:    id,
		Name:  name,
		Grade: grade,
		Age:   18 + int(id%10),
	}
}

// =============================================================================
// HEAP-HEAVY APPROACHES (SLOWER)
// =============================================================================

// calculateAverageHeap creates intermediate slice (heap allocations)
func calculateAverageHeap(students []Student) float64 {
	var grades []float64 // This slice grows on heap

	for _, student := range students {
		grades = append(grades, student.Grade) // Potential reallocations
	}

	var total float64
	for _, grade := range grades {
		total += grade
	}

	return total / float64(len(grades))
}

// createStudentByPointer returns pointer, forces heap allocation
func createStudentByPointer(id int64, name string, grade float64) *Student {
	student := Student{
		ID:    id,
		Name:  name,
		Grade: grade,
		Age:   18 + int(id%10),
	}
	return &student // Forces heap allocation
}

// =============================================================================
// OPTIMIZED APPROACH (BALANCED)
// =============================================================================

// calculateAverageOptimized pre-allocates to minimize heap usage
func calculateAverageOptimized(students []Student) float64 {
	// Pre-allocate with exact capacity - one allocation only
	grades := make([]float64, 0, len(students))

	for _, student := range students {
		grades = append(grades, student.Grade)
	}

	var total float64
	for _, grade := range grades {
		total += grade
	}

	return total / float64(len(grades))
}

// =============================================================================
// PERFORMANCE MEASUREMENT
// =============================================================================

func runPerformanceTest() {
	const iterations = 50000

	// Create test data
	students := make([]Student, 100)
	for i := range students {
		students[i] = Student{
			ID:    int64(i),
			Name:  fmt.Sprintf("Student%d", i),
			Grade: 60.0 + float64(i%40), // Grades 60-100
			Age:   18 + i%10,
		}
	}

	fmt.Printf("Testing with %d students, %d iterations each\n", len(students), iterations)
	fmt.Println("==================================================")

	// Test stack approach
	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	start := time.Now()
	for i := 0; i < iterations; i++ {
		_ = calculateAverageStack(students)
	}
	stackTime := time.Since(start)

	runtime.ReadMemStats(&m2)
	stackAllocs := m2.TotalAlloc - m1.TotalAlloc

	// Test heap approach
	runtime.GC()
	runtime.ReadMemStats(&m1)

	start = time.Now()
	for i := 0; i < iterations; i++ {
		_ = calculateAverageHeap(students)
	}
	heapTime := time.Since(start)

	runtime.ReadMemStats(&m2)
	heapAllocs := m2.TotalAlloc - m1.TotalAlloc

	// Test optimized approach
	runtime.GC()
	runtime.ReadMemStats(&m1)

	start = time.Now()
	for i := 0; i < iterations; i++ {
		_ = calculateAverageOptimized(students)
	}
	optimizedTime := time.Since(start)

	runtime.ReadMemStats(&m2)
	optimizedAllocs := m2.TotalAlloc - m1.TotalAlloc

	// Display results
	fmt.Printf("Stack approach:     %8v  %10d bytes\n", stackTime, stackAllocs)
	fmt.Printf("Heap approach:      %8v  %10d bytes\n", heapTime, heapAllocs)
	fmt.Printf("Optimized approach: %8v  %10d bytes\n", optimizedTime, optimizedAllocs)

	if stackTime > 0 && heapTime > 0 {
		fmt.Printf("\nPerformance improvements:\n")
		fmt.Printf("• Stack is %.1fx faster than heap\n", float64(heapTime)/float64(stackTime))
		if optimizedTime > 0 {
			fmt.Printf("• Optimized is %.1fx faster than heap\n", float64(heapTime)/float64(optimizedTime))
		}
	}
}

// =============================================================================
// MAIN FUNCTION
// =============================================================================

// runStackVsHeapDemo demonstrates the performance differences
func runStackVsHeapDemo() {
	fmt.Println("Stack vs Heap Performance Comparison")
	fmt.Println("====================================")

	// Simple examples
	fmt.Println("\n--- Creating Students ---")

	stackStudent := createStudentByValue(1, "Alice", 85.0)
	fmt.Printf("Stack student: %+v\n", stackStudent)

	heapStudent := createStudentByPointer(2, "Bob", 92.0)
	fmt.Printf("Heap student: %+v\n", heapStudent)

	// Performance test
	fmt.Println("\n--- Performance Test ---")
	runPerformanceTest()

	fmt.Println("\n=== Key Takeaways ===")
	fmt.Println("• Stack allocation is much faster than heap")
	fmt.Println("• Avoid creating unnecessary intermediate slices")
	fmt.Println("• Pre-allocate when you know the size")
	fmt.Println("• Return values when possible, pointers when necessary")
}
