package pipeline

import "fmt"

// multiply and add are considered BATCH PROCESSING functions
// because they modify a batch (=collection) of data at once

// multiply is a stage that multiplies each value of values by the value of multiplier
func multiply(values []int, multiplier int) []int {
	m := make([]int, len(values))

	for i, v := range values {
		m[i] = v * multiplier
	}

	return m
}

// add is a stage that adds additive to each value of values
func add(values []int, additive int) []int {
	a := make([]int, len(values))

	for i, v := range values {
		a[i] = v + additive
	}

	return a
}

func batchProcessingExec() {
	ints := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	for _, v := range add(multiply(ints, 2), 1) {
		fmt.Println(v)
	}
}

// multiplyStream is a stage that multiplies one value by the value of multiplier
func multiplyStream(value int, multiplier int) int {
	return value * multiplier
}

// addStream is a stage that adds additive to one value
func addStream(value int, additive int) int {
	return value + additive
}

func streamProcessingExec() {
	ints := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	for _, v := range ints {
		fmt.Println(addStream(multiplyStream(v, 2), 1))
	}
}
