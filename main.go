package main

import (
	"fmt"
	"math/rand"
)

type sortJob struct {
	a []int
	b []int
}

func sortWorker(id int, jobs <-chan sortJob, result chan<- []int) {
	for j := range jobs {
		if j.b == nil {
			fmt.Println("worker", id, "bubble sorting list size", len(j.a))
			// Bubble sort
			result <- bubbleSort(j.a)
		} else {
			// Merge
			fmt.Println("worker", id, "merging lists of size", len(j.a), "and", len(j.b))
			result <- merge(j)
		}
	}
}

func bubbleSort(list []int) []int {
	sorted := false
	for !sorted {
		sorted = true
		for i := 1; i < len(list); i++ {
			if list[i] < list[i-1] {
				sorted = false
				temp := list[i]
				list[i] = list[i-1]
				list[i-1] = temp
			}
		}
	}
	return list
}

func merge(lists sortJob) []int {
	result := make([]int, len(lists.a)+len(lists.b))
	ac := 0
	bc := 0
	for i := 0; i < len(result); i++ {
		if bc == len(lists.b) {
			result[i] = lists.a[ac]
			ac++
		} else if ac == len(lists.a) {
			result[i] = lists.b[bc]
			bc++
		} else if lists.a[ac] < lists.b[bc] {
			result[i] = lists.a[ac]
			ac++
		} else {
			result[i] = lists.b[bc]
			bc++
		}
	}
	return result
}

func multithreadMergeSort(list []int) []int {
	const numWorkers = 1
	const minChunkSize = 20
	originalLength := len(list)
	var jobCount int = (originalLength / minChunkSize) + 2
	jobs := make(chan sortJob, jobCount)      // We'll have the most jobs in the queue at the start.
	results := make(chan []int, numWorkers*2) // Results in queue shouldn't exceed double the workers.

	// Preload jobs for bubble sort.
	for i := 0; i < originalLength; i += minChunkSize {
		end := i + minChunkSize
		if end > originalLength {
			end = originalLength
		}
		temp := sortJob{list[i:end], nil}
		jobs <- temp
	}
	// Start workers.
	for w := 1; w <= numWorkers; w++ {
		go sortWorker(w, jobs, results)
	}

	// Result handler
	var (
		a []int
		b []int
	)
	for true {
		a = <-results
		if len(a) == originalLength {
			// Result is an entire sorted list.
			close(jobs)
			break
		} else {
			b = <-results
			jobs <- sortJob{a, b}
		}
	}

	// Had planned to do this inside the above loop, but compiler didn't like it.
	return a
}

func isSorted(list []int) bool {
	for i := 1; i < len(list); i++ {
		if list[i] > list[i-1] {
			return false
		}
	}
	return true
}

func main() {
	const listSize = 100

	// Make a list of random numbers.
	fmt.Print("Building random list... ")
	numbers := make([]int, listSize)
	for i := 0; i < listSize; i++ {
		numbers[i] = rand.Intn(10 * listSize)
	}
	fmt.Println("Done")

	// Run the sort.
	result := multithreadMergeSort(numbers)

	// Test the sort
	fmt.Print("Checking if sort is valid... ")
	fmt.Println(isSorted(result))
}
