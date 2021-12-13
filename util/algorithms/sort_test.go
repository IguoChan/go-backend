package algorithms

import (
	"fmt"
	"math/rand"
	"testing"
)

var a []int32

func TestInsertionSortAscend(t *testing.T) {
	a := []int32{2, 1, 4, 7, 3, 5}
	InsertionSortAscend(a)
	fmt.Println(a)
}

func TestInsertionSortDescend(t *testing.T) {
	a := []int32{2, 1, 4, 7, 3, 5}
	InsertionSortDescend(a)
	fmt.Println(a)
}

func TestSelectionSortAscend(t *testing.T) {
	a := []int32{2, 1, 2, 4, 7, 3, 5}
	SelectionSortAscend(a)
	fmt.Println(a)
}

func TestSelectionSortDescend(t *testing.T) {
	a := []int32{2, 1, 2, 4, 7, 3, 5}
	SelectionSortDescend(a)
	fmt.Println(a)
}

func TestMergeSortAscend(t *testing.T) {
	a := []int32{2, 1, 2, 4, 7, 3, 5, 2, 9, 11}
	MergeSortAscend(a)
	fmt.Println(a)
}

func TestMergeSortDescend(t *testing.T) {
	a := []int32{2, 1, 2, 4, 7, 3, 5}
	MergeSortDescend(a)
	fmt.Println(a)
}

func TestBubbleSortAscend(t *testing.T) {
	a := []int32{2, 1, 2, 4, 7, 3, 5, 2, 9, 11}
	BubbleSortAscend(a)
	fmt.Println(a)
}

func TestBubbleSortDescend(t *testing.T) {
	a := []int32{2, 1, 2, 4, 7, 3, 5, 2, 9, 11}
	BubbleSortDescend(a)
	fmt.Println(a)
}

func BenchmarkInsertionSortAscend(b *testing.B) {
	b.ResetTimer()
	InsertionSortAscend(a)
	fmt.Println("insert")
	//fmt.Println(a)
}

func BenchmarkMergeSortAscend(b *testing.B) {
	b.ResetTimer()
	MergeSortAscend(a)
	fmt.Println("merge")
	//fmt.Println(a)
}

func init() {
	a = make([]int32, 0, 100000)
	for i := 0; i < 100000; i++ {
		a = append(a, int32(rand.Int()))
	}
	fmt.Println(a)
}
