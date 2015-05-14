package algorithm

func QuickSort(s *[]int, left, right int) {
	if left < right {
		i := left
		j := right
		x := (*s)[left]

		for i < j {
			for i < j && (*s)[j] >= x {
				j--
			}

			if i < j {
				(*s)[i] = (*s)[j]
				i++
			}

			for i < j && (*s)[i] < x {
				i++
			}

			if i < j {
				(*s)[j] = (*s)[i]
				j--
			}
		}

		(*s)[i] = x
		QuickSort(s, left, i-1)
		QuickSort(s, i+1, right)
	}
}
