package main

import "fmt"

// merge sort，使用goroutinues 和 channel 提高并发度
func merge_sort_conc(sli []float64, res_chan chan float64) {
	fmt.Println("[debug] sub_merge, sli:", sli, len(sli))

	N := len(sli)

	if N == 1 {
		res_chan <- sli[0]
		fmt.Printf("[debug][size_1_slice] sli:%v, send:%f, and return.\n",
			sli, sli[0])

		close(res_chan)
		return // bug fix 需要跳出
	}

	split_p := int(N / 2)

	var left_chan, right_chan chan float64 // init for nil
	var left_chan_tmp, right_chan_tmp chan float64

	if len(sli[:split_p]) >= 1 {
		left_chan = make(chan float64)
		go merge_sort_conc(sli[:split_p], left_chan)
	}
	if len(sli[split_p:]) >= 1 {
		right_chan = make(chan float64)
		go merge_sort_conc(sli[split_p:], right_chan)
	}

	var left_val, right_val float64
	var left_tag, right_tag bool

	left_left := len(sli[:split_p])
	right_left := len(sli[split_p:])
	fmt.Printf("[debug][] sli:%v, init_left_n:(%d, %d).\n",
		sli, left_left, right_left)

	i := 0
	for {
		i += 1

		if !left_tag && left_left > 0 {
			left_chan_tmp = left_chan
		} else {
			left_chan_tmp = nil
		}

		if !right_tag && right_left > 0 {
			right_chan_tmp = right_chan
		} else {
			right_chan_tmp = nil
		}

		if left_chan_tmp != nil || right_chan_tmp != nil {
			select {
			case left_val = <-left_chan_tmp:
				left_tag = true
				left_left -= 1

			case right_val = <-right_chan_tmp:
				right_tag = true
				right_left -= 1
			}
		}

		fmt.Printf("[debug][each loop] sli:%v, loop: %d,val:(%f, %f), received_n(%d, %d),.\n",
			sli, i, left_val, right_val, left_left, right_left)

		if left_tag && right_tag { // 都是新值需要比较
			fmt.Printf("[debug][sort_two_chan][two new] sli:%v, loop:%d val(%f, %f).\n", sli, i, left_val, right_val)
			if left_val < right_val {
				res_chan <- left_val
				left_tag = false // 旧值作废，需要通过select重新取
			} else {
				res_chan <- right_val
				right_tag = false
			}
		} else if left_tag && right_left == 0 {
			res_chan <- left_val
			left_tag = false // 旧值作废，需要通过select重新取
			fmt.Printf("[debug][sort_two_chan][left new] sli:%v, loop:%d val(%f, %f).\n", sli, i, left_val, right_val)
		} else if right_tag && left_left == 0 {
			res_chan <- right_val
			right_tag = false
			fmt.Printf("[debug][sort_two_chan][right new] sli:%v, loop:%d val(%f, %f).\n", sli, i, left_val, right_val)
		}

		fmt.Printf("[debug][sort_two_chan][loop_final_state] sli:%v, loop:%d,tag(%t, %t),left(%d, %d).\n",
			sli, i, left_tag, right_tag, left_left, right_left)

		fmt.Printf("[debug][sort_two_chan][break_tag] sli:%v, loop:%d,break-tag:%t.\n",
			sli, i, (left_left+right_left) == 0 && !left_tag && !right_tag)

		if (left_left+right_left) == 0 && !left_tag && !right_tag {
			break
		}
	}

}

func merge_sort(sli []float64) []float64 {
	res_chan := make(chan float64)
	go merge_sort_conc(sli, res_chan)

	N := len(sli)

	var sorted_slice []float64

	for i := 0; i < N; i++ {
		sorted_slice = append(sorted_slice, <-res_chan)
	}
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
