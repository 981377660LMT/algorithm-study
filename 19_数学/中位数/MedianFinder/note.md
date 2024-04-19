# 区间距离和问题

区间内所有数到某个 x=k 的距离之和

- 有序数组维护区间距离和 => 前缀和
- 无序数组维护区间距离和 => WaveletMatrix
- 有序容器维护区间距离和 => Sortedlist
  [DistSum](DistSum.go)
  [MedianFinderWaveletMatrix](MedianFinderWaveletMatrix.go)
  [MedianFinderSortedList](MedianFinderSortedList.go)

```
Api:
 required:
   GetMedian(sortedNums []int, start, end int) int // 求有序数组区间中位数.
   DistSum(sortedNums []int) func(k int) int // 有序数组中所有点到`x=k`的距离之和.
   DistSumRange(sortedNums []int) func(k, start, end int) int // 有序数组区间所有点到`x=k`的距离之和.

 optional:
   DistSumOfAllPairs(sortedNums []int) int // 有序数组中所有点对两两距离之和.
   DistSumOfAllPairsRange(sortedNums []int, start, end int) int // 有序数组区间中所有点对两两距离之和.
```
