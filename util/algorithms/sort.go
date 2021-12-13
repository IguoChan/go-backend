package algorithms

import "math"

// 插入排序升序
func InsertionSortAscend(arr []int32) {
	l := len(arr)
	if l <= 1 {
		return
	}
	for i := 1; i < l; i++ {
		key := arr[i]
		j := i - 1
		for j >= 0 && arr[j] > key {
			arr[j+1] = arr[j]
			j--
		}
		arr[j+1] = key
	}
}

// 插入排序降序
func InsertionSortDescend(arr []int32) {
	l := len(arr)
	if l <= 1 {
		return
	}
	for i := 1; i < l; i++ {
		key := arr[i]
		j := i - 1
		for j >= 0 && arr[j] < key {
			arr[j+1] = arr[j]
			j--
		}
		arr[j+1] = key
	}
}

// 选择排序升序
func SelectionSortAscend(arr []int32) {
	l := len(arr)
	if l <= 1 {
		return
	}
	for i := 0; i < l; i++ {
		minIndex := i
		for j := i + 1; j < l; j++ {
			if arr[j] < arr[minIndex] {
				minIndex = j
			}
		}
		arr[i], arr[minIndex] = arr[minIndex], arr[i]
	}
}

// 选择排序升序
func SelectionSortDescend(arr []int32) {
	l := len(arr)
	if l <= 1 {
		return
	}
	for i := 0; i < l; i++ {
		minIndex := i
		for j := i + 1; j < l; j++ {
			if arr[j] > arr[minIndex] {
				minIndex = j
			}
		}
		arr[i], arr[minIndex] = arr[minIndex], arr[i]
	}
}

func BubbleSortAscend(arr []int32) {
	l := len(arr)
	if l <= 1 {
		return
	}
	for i := 0; i < l-1; i++ {
		for j := l - 1; j > i; j-- {
			if arr[j] < arr[j-1] {
				arr[j], arr[j-1] = arr[j-1], arr[j]
			}
		}
	}
}

func BubbleSortDescend(arr []int32) {
	l := len(arr)
	if l <= 1 {
		return
	}
	for i := 0; i < l-1; i++ {
		for j := l - 1; j > i; j-- {
			if arr[j] > arr[j-1] {
				arr[j], arr[j-1] = arr[j-1], arr[j]
			}
		}
	}
}

func MergeSortAscend(arr []int32) {
	l := len(arr)
	if l <= 1 {
		return
	}

	merge := func(a []int32, p, q, r int) {
		lArr := make([]int32, 0, q-p+2)
		rArr := make([]int32, 0, r-q+1)
		for i := 0; i < q-p+1; i++ {
			lArr = append(lArr, a[p+i])
		}
		lArr = append(lArr, math.MaxInt32)
		for i := 0; i < r-q; i++ {
			rArr = append(rArr, a[q+i+1])
		}
		rArr = append(rArr, math.MaxInt32)
		i := 0
		j := 0
		for k := p; k < r+1; k++ {
			if lArr[i] <= rArr[j] {
				a[k] = lArr[i]
				i++
			} else {
				a[k] = rArr[j]
				j++
			}
		}

	}

	var mergeSort func(a []int32, p, r int)
	mergeSort = func(a []int32, p, r int) {
		if p < r {
			q := (p + r) / 2
			mergeSort(a, p, q)
			mergeSort(a, q+1, r)
			merge(a, p, q, r)
		}
	}
	mergeSort(arr, 0, l-1)
}

func MergeSortDescend(arr []int32) {
	l := len(arr)
	if l <= 1 {
		return
	}

	merge := func(a []int32, p, q, r int) {
		lArr := make([]int32, 0, q-p+2)
		rArr := make([]int32, 0, r-q+1)
		for i := 0; i < q-p+1; i++ {
			lArr = append(lArr, a[p+i])
		}
		lArr = append(lArr, math.MinInt32)
		for i := 0; i < r-q; i++ {
			rArr = append(rArr, a[q+i+1])
		}
		rArr = append(rArr, math.MinInt32)
		i := 0
		j := 0
		for k := p; k < r+1; k++ {
			if lArr[i] >= rArr[j] {
				a[k] = lArr[i]
				i++
			} else {
				a[k] = rArr[j]
				j++
			}
		}

	}

	var mergeSort func(a []int32, p, r int)
	mergeSort = func(a []int32, p, r int) {
		if p < r {
			q := (p + r) / 2
			mergeSort(a, p, q)
			mergeSort(a, q+1, r)
			merge(a, p, q, r)
		}
	}
	mergeSort(arr, 0, l-1)
}
