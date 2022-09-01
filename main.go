// Create by Blake Gall <bg849619@ohio.edu>
// August 31, 2022
package main

import (
	"fmt"
	"math/rand"
	"time"
)

type sortJob struct {
	a []int
	b []int
}

// Worker to be ran as a goroutine.
// Runs sortjobs using a bubblesort if only A has a list, or merge if both.
func sortWorker(jobs <-chan sortJob, result chan<- []int) {
	for j := range jobs {
		if j.b == nil {
			// Bubble sort
			result <- bubbleSort(j.a)
		} else {
			// Merge
			result <- merge(j)
		}
	}
}

// Sorts a list using a O(N^2) algorithm
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

// Takes two sorted lists, and merges them into a singular sorted list.
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

// The dispatching function for mt merge sort.
// Splits the list into chunks, starts workers in goroutines, then dispatches new jobs to the workers until a completely sorted list is received.
// The first chunks created are sorted using bubble sort. The overhead of dispatching merge workers back and forth for smaller lists is more than multithreading a bubblesort for pieces of the list.
func multithreadMergeSort(list []int, workers int) []int {
	var minChunkSize = 100

	originalLength := len(list)
	var jobCount int = (originalLength / minChunkSize) + 2
	jobs := make(chan sortJob, jobCount)   // We'll have the most jobs in the queue at the start.
	results := make(chan []int, workers*2) // Results in queue shouldn't exceed double the workers.

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
	for w := 1; w <= workers; w++ {
		go sortWorker(jobs, results)
	}

	// Result handler
	var (
		a []int
		b []int
	)
	for {
		a = <-results
		if len(a) == originalLength {
			// Result is an entire sorted list.
			close(jobs)
			break
		} else {
			// Wait for another sorted list before dispatching to worker.
			b = <-results
			jobs <- sortJob{a, b}
		}
	}

	// Had planned to do this inside the above loop, but compiler didn't like it.
	return a
}

func isSorted(list []int) bool {
	for i := 1; i < len(list); i++ {
		if list[i] < list[i-1] {
			return false
		}
	}
	return true
}

// Entrypoint for testing the algorithm.
func main() {
	// Size of testing list. This should be large enough that cached instructions don't overcome
	// the benefit of multithreading. Can keep small for testing.

	const listSize = 99999999

	// Make a list of random numbers.
	fmt.Print("Building random list... ")
	var numbers [listSize]int
	for i := 0; i < listSize; i++ {
		numbers[i] = rand.Intn(10 * listSize)
	}
	fmt.Println("Done")

	runTest := func(threads int) {
		var (
			result  []int
			start   time.Time
			elapsed time.Duration
		)

		// Run test with 1 thread.
		// So the sort doesn't affect the original array.
		disjoint := make([]int, len(numbers))
		copy(disjoint, numbers[:])

		fmt.Println("Running test with ", threads, " threads...")

		start = time.Now()
		result = multithreadMergeSort(disjoint, threads)
		elapsed = time.Since(start)

		fmt.Println("Completed in ", elapsed, ". PASS? ", isSorted(result), "\n")
	}

	runTest(1)
	runTest(2)
	runTest(4)
	runTest(8)
	runTest(16)
	runTest(32)

}
