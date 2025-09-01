// Package main demonstrates value vs pointer behavior
// This file contains functions only - no main function
// Run all examples with: go run examples/*.go
package main

import "fmt"

// =============================================================================
// VALUE OPERATIONS (WORKING WITH COPIES)
// =============================================================================

// modifyStudentValue works with a copy - original unchanged
func modifyStudentValue(student Student) {
	student.Grade = 100.0
	student.Name = "Modified " + student.Name
	fmt.Printf("Inside modifyStudentValue: %+v\n", student)
	// Changes only affect the copy!
}

// calculateGPAValue takes copies - safe but potentially slower for large structs
func calculateGPAValue(student1, student2, student3 Student) float64 {
	total := student1.Grade + student2.Grade + student3.Grade
	return total / 3.0
}

// createStudentCopy demonstrates making copies
func createStudentCopy(template Student, newName string, newID int64) Student {
	// Create a new student based on template
	newStudent := template // This makes a complete copy
	newStudent.Name = newName
	newStudent.ID = newID
	return newStudent
}

// =============================================================================
// POINTER OPERATIONS (WORKING WITH ORIGINALS)
// =============================================================================

// modifyStudentPointer works with the original - changes persist
func modifyStudentPointer(student *Student) {
	if student == nil { // Always check for nil!
		return
	}

	student.Grade = 100.0
	student.Name = "Modified " + student.Name
	fmt.Printf("Inside modifyStudentPointer: %+v\n", *student)
	// Changes affect the original!
}

// calculateGPAPointer takes pointers - faster, but need to be careful
func calculateGPAPointer(student1, student2, student3 *Student) float64 {
	// Check for nil pointers
	if student1 == nil || student2 == nil || student3 == nil {
		return 0.0
	}

	total := student1.Grade + student2.Grade + student3.Grade
	return total / 3.0
}

// updateGradeInPlace modifies the original student
func updateGradeInPlace(student *Student, newGrade float64) {
	if student == nil {
		return
	}
	student.Grade = newGrade
}

// =============================================================================
// SLICE BEHAVIOR (REFERENCE TYPE)
// =============================================================================

// modifyGradesSlice - slices are reference types!
func modifyGradesSlice(grades []float64) {
	// Even though we pass the slice by "value",
	// the slice points to the underlying array
	for i := range grades {
		grades[i] += 5.0 // Bonus points - affects original array!
	}
	fmt.Printf("Inside modifyGradesSlice: %v\n", grades)
}

// addStudentToClass - modifying slice structure
func addStudentToClass(students *[]Student, newStudent Student) {
	// Need pointer to slice to modify the slice itself
	*students = append(*students, newStudent)
}

// =============================================================================
// PRACTICAL EXAMPLES
// =============================================================================

func demonstrateValueBehavior() {
	fmt.Println("=== VALUE BEHAVIOR ===")

	original := Student{ID: 1, Name: "Alice", Grade: 85.0, Age: 20}
	fmt.Printf("Original student: %+v\n", original)

	// Pass by value - original unchanged
	modifyStudentValue(original)
	fmt.Printf("After modifyStudentValue: %+v\n", original) // UNCHANGED!

	// Create copy
	copy := createStudentCopy(original, "Alice Clone", 2)
	fmt.Printf("Created copy: %+v\n", copy)
	fmt.Printf("Original still: %+v\n", original) // UNCHANGED!
}

func demonstratePointerBehavior() {
	fmt.Println("\n=== POINTER BEHAVIOR ===")

	original := Student{ID: 3, Name: "Bob", Grade: 78.0, Age: 21}
	fmt.Printf("Original student: %+v\n", original)

	// Pass by pointer - original WILL change
	modifyStudentPointer(&original)
	fmt.Printf("After modifyStudentPointer: %+v\n", original) // CHANGED!

	// Multiple pointers to same student
	ptr1 := &original
	ptr2 := &original

	fmt.Printf("Pointer 1 points to: %+v\n", *ptr1)
	fmt.Printf("Pointer 2 points to: %+v\n", *ptr2)

	// Change through one pointer
	ptr1.Grade = 95.0

	fmt.Printf("After changing via ptr1:\n")
	fmt.Printf("  original: %+v\n", original) // CHANGED!
	fmt.Printf("  ptr1: %+v\n", *ptr1)        // CHANGED!
	fmt.Printf("  ptr2: %+v\n", *ptr2)        // CHANGED! (same object)
}

func demonstrateSliceBehavior() {
	fmt.Println("\n=== SLICE BEHAVIOR (REFERENCE TYPE) ===")

	grades := []float64{85.0, 92.0, 78.0, 88.0}
	fmt.Printf("Original grades: %v\n", grades)

	// Slices are reference types - changes affect original
	modifyGradesSlice(grades)
	fmt.Printf("After modifyGradesSlice: %v\n", grades) // CHANGED!

	// Working with slice of students
	var class []Student
	fmt.Printf("Empty class: %v\n", class)

	// Add students to class
	addStudentToClass(&class, Student{ID: 1, Name: "Charlie", Grade: 82.0, Age: 19})
	addStudentToClass(&class, Student{ID: 2, Name: "Diana", Grade: 91.0, Age: 20})

	fmt.Printf("Class after adding students: %d students\n", len(class))
	for _, student := range class {
		fmt.Printf("  %s: %.1f\n", student.Name, student.Grade)
	}
}

// =============================================================================
// WHEN TO USE EACH
// =============================================================================

func demonstrateWhenToUseEach() {
	fmt.Println("\n=== WHEN TO USE VALUES vs POINTERS ===")

	students := []Student{
		{ID: 1, Name: "Alice", Grade: 85.0, Age: 20},
		{ID: 2, Name: "Bob", Grade: 92.0, Age: 21},
		{ID: 3, Name: "Charlie", Grade: 78.0, Age: 19},
	}

	fmt.Println("\n1. For READ-ONLY operations - use VALUES:")
	gpa := calculateGPAValue(students[0], students[1], students[2])
	fmt.Printf("   GPA (by value): %.2f\n", gpa)
	fmt.Printf("   Students unchanged: %v\n", students[0].Grade == 85.0)

	fmt.Println("\n2. For MODIFICATIONS - use POINTERS:")
	fmt.Printf("   Before: %s has grade %.1f\n", students[0].Name, students[0].Grade)
	updateGradeInPlace(&students[0], 95.0)
	fmt.Printf("   After: %s has grade %.1f\n", students[0].Name, students[0].Grade)

	fmt.Println("\n3. For LARGE structs - use POINTERS (efficiency):")
	gpa2 := calculateGPAPointer(&students[0], &students[1], &students[2])
	fmt.Printf("   GPA (by pointer): %.2f\n", gpa2)
}

// =============================================================================
// MAIN DEMONSTRATION
// =============================================================================

// runValueVsPointerDemo demonstrates value vs pointer behavior
func runValueVsPointerDemo() {
	fmt.Println("Value vs Pointer Demonstration")
	fmt.Println("==============================")

	// Basic behavior examples
	demonstrateValueBehavior()
	demonstratePointerBehavior()
	demonstrateSliceBehavior()

	// Practical usage examples
	demonstrateWhenToUseEach()

	fmt.Println("\n=== SUMMARY ===")
	fmt.Println("VALUES (copies):")
	fmt.Println("  • Safe - can't accidentally modify original")
	fmt.Println("  • Good for small structs and read-only operations")
	fmt.Println("  • Each function gets its own copy to work with")

	fmt.Println("\nPOINTERS (shared access):")
	fmt.Println("  • Efficient - no copying")
	fmt.Println("  • Required when you want to modify the original")
	fmt.Println("  • Multiple pointers can refer to same object")
	fmt.Println("  • Always check for nil!")

	fmt.Println("\nSLICES & MAPS:")
	fmt.Println("  • Reference types - contain pointers internally")
	fmt.Println("  • Modifications affect the underlying data")
	fmt.Println("  • Be aware of this behavior!")
}
