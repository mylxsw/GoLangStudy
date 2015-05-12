package algorithm

func BubbleSort(items []int) []int {
	item_len := len(items)
	for i := 1; i < item_len; i++ {
		for k := 0; k < item_len-i; k++ {
			if items[k] > items[k+1] {
				items[k], items[k+1] = items[k+1], items[k]
			}
		}

	}

	return items
}
