package utils

import (
	"sort"

	"github.com/samber/lo"
)

// GetTopKValues 返回切片中出现频率最高的前k个值
func GetTopKValues[T comparable](values []T, k int) []T {
	// 统计每个值出现的次数
	countMap := lo.CountValues(values)

	// 将map转换为slice以便排序
	type kv struct {
		Key   interface{}
		Value int
	}
	var kvSlice []kv
	for key, value := range countMap {
		kvSlice = append(kvSlice, kv{key, value})
	}

	// 按照出现次数降序排序
	sort.Slice(kvSlice, func(i, j int) bool {
		return kvSlice[i].Value > kvSlice[j].Value
	})

	// 获取前k个值
	if k > len(kvSlice) {
		k = len(kvSlice)
	}
	topKValues := make([]T, k)
	for i := 0; i < k; i++ {
		topKValues[i] = kvSlice[i].Key.(T)
	}

	return topKValues
}
