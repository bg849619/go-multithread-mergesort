# Go Multithreaded Mergesort

This is by no means a perfect program, but it's more me learning Go by playing with some algorithms. My favorite part of Go is the ease of multithreading using Goroutines. I've found starting new threads obnoxious in other programming languages (although, I was much less experienced at the time), but in Go, it's ridiculously simple.

## The algorithm

Most implementations of MergeSort will use some sort of recursive call: MergeSort one half, MergeSort the other, then merge the resultant sorted lists. However, it's also possible to start this algorithm from the individual items in the list. As long as the merge algorithm has two sorted lists, it will work. A list of a single item is sorted.

Although it would work if we started the merge at the individual items and worked our way up, the overhead of so many recursion calls, or in this case job queues, would greatly increase the complexity of the algorithm. A solution is to decide on the smallest chunk of a list you'd want to work with, then run some type of non-recursive sort on each of those chunks in the list. Since these are seperate sublists, this can also be multithreaded. Now that we have a bunch of sorted lists, we can start merging the lists. Once a list is merged, we can merge it with another sorted list, until we return to a list with the same size as the original.

## Multithreading

For each thread wanted, we can run a worker function as a goroutine. The workers are just functions with an input and output channel (Basically pipes/queues), and a loop. The input channel will send a custom struct, which can contain two lists. When the worker receives this struct, it will run the merge function on the two lists in its own thread. (Or, if only one of the lists is defined, it will run a bubblesort on that list). Finally, the sorted list will be sent back through the output channel.

The dispatching function, `multithreadMergeSort` is what starts the goroutines and sends the jobs to the workers. Initially, it fills the job queue with the chunks of the list as described earlier to be bubblesorted. It will then wait for 2 sorted lists to be returned through the output channel, then push those two lists as another job on the queue. If it receives a output list that is the same length as the original, it kills the workers, and the function returns the sorted list.

## Performance

With the chunk size set to 100, and for a list of size 99999999, my computer (8c/16t) completed the sort in the following times:
|Threads|Completion Time|
|:-------|---------------:|
|1|20.0126s|
|2|9.8935s|
|4|5.7323s|
|8|4.3197s|
|16|3.9385s|
|32|4.0180s|

## Improvements

Things I might be looking to improve:

### Space complexity and slice copying

Currently, the algorithm will create new slices on every merge. There might be a way utilizing slices to do the merge in-place. Slices in Go are references to a section of an underlying array. It should be noted that because of the multithreading of the algorithm, any two lists which are being merged might not be adjacent in the underlying array.

### Chunk size optimization

The chunk size being used for the initial bubble-sort has so far been trial and error. Either doing some math or running more complex tests might be useful to find an optimized chunk.
