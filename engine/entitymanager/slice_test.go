package utils

import (
	"reflect"
	"testing"
)

func TestGetTopKValues(t *testing.T) {
	t.Run("整数切片", func(t *testing.T) {
		values := []int{1, 2, 2, 3, 3, 3, 4}
		got := GetTopKValues(values, 2)
		want := []int{3, 2} // 3出现3次，2出现2次
		if !reflect.DeepEqual(got, want) {
			t.Errorf("GetTopKValues() = %v, want %v", got, want)
		}
	})

	t.Run("字符串切片", func(t *testing.T) {
		values := []string{"apple", "banana", "banana", "cherry", "cherry", "cherry"}
		got := GetTopKValues(values, 2)
		want := []string{"cherry", "banana"} // cherry出现3次，banana出现2次
		if !reflect.DeepEqual(got, want) {
			t.Errorf("GetTopKValues() = %v, want %v", got, want)
		}
	})

	t.Run("空切片", func(t *testing.T) {
		var values []int
		got := GetTopKValues(values, 3)
		want := make([]int, 0)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("GetTopKValues() = %v, want %v", got, want)
		}
	})

	t.Run("k大于切片长度", func(t *testing.T) {
		values := []int{1, 2, 2, 3}
		got := GetTopKValues(values, 5)
		want := []int{2, 1, 3} // 2出现2次排第一，1和3各出现1次
		if !reflect.DeepEqual(got, want) {
			t.Errorf("GetTopKValues() = %v, want %v", got, want)
		}
	})

	t.Run("相同出现次数", func(t *testing.T) {
		values := []int{1, 1, 2, 2, 3, 3}
		got := GetTopKValues(values, 2)
		want := []int{1, 2} // 都出现2次，保持原有顺序
		if !reflect.DeepEqual(got, want) {
			t.Errorf("GetTopKValues() = %v, want %v", got, want)
		}
	})
}
