package main

/*
精简merge sort的代码，降低滥用的并发度
*/

import (
	"fmt"
)

// merge sort，使用goroutinues 和 channel 提高并发度
func merge_sort(sli []float64) []float64 {
	fmt.Println("[debug] sub_merge, sli:", sli, len(sli))
	N := len(sli)

	if N <= 1 {
		return sli
	}

	// STEP1: sort sub_slice
	res_chan := make(chan []float64)
	split_p := int(N / 2)

	// sorted_left_part
	go func() {
		res_chan <- merge_sort(sli[:split_p])
	}()

	// sorted_right_part
	go func() {
		res_chan <- merge_sort(sli[split_p:])
	}()

	sorted_sli_part1 := <-res_chan
	sorted_sli_part2 := <-res_chan

	// STEP2: merge two sorted sub_slice
	fmt.Printf("[debug]sli-%v, sorted_sub-(%v, %v).\n",
		sli, sorted_sli_part1, sorted_sli_part2)

	var sorted_slice []float64

	i, j := 0, 0
	for i < len(sorted_sli_part1) && j < len(sorted_sli_part2) {
		if sorted_sli_part1[i] < sorted_sli_part2[j] {
			sorted_slice = append(sorted_slice, sorted_sli_part1[i])
			i += 1
		} else {
			sorted_slice = append(sorted_slice, sorted_sli_part2[j])
			j += 1
		}
	}
	for i < len(sorted_sli_part1) {
		sorted_slice = append(sorted_slice, sorted_sli_part1[i])
		i += 1
	}

	for j < len(sorted_sli_part2) {
		sorted_slice = append(sorted_slice, sorted_sli_part2[j])
		j += 1
	}

	fmt.Printf("[debug]sli%v, sorted_final:%v.\n",
		sli, sorted_slice)

	return sorted_slice

}

func harder9() {
	ori_slice := []float64{3.1, 3.2, 5, 2, 1, 6, 10}
	res_slice := merge_sort(ori_slice)
	fmt.Println(res_slice)
}

func main() {
	harder9()

	// panic(" new a painc")
}
