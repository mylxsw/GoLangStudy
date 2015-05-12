package algorithm

func SelectSort(items []int) []int {
	items_len := len(items)
	for i := 0; i < items_len-1; i++ {
		p := i
		for j := i + 1; j < items_len; j++ {
			if items[p] > items[j] {
				p = j
			}
		}

		if p != i {
			items[p], items[i] = items[i], items[p]
		}
	}
	return items
}
