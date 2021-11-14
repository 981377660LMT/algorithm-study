// 冒泡排序有两种优化方式。

// 一种是外层循环的优化，我们可以记录当前循环中是否发生了交换，如果没有发生交换，则说明该序列已经为有序序列了。
// 因此我们不需要再执行之后的外层循环，此时可以直接结束。

// 一种是内层循环的优化，我们可以记录当前循环中最后一次元素交换的位置，该位置以后的序列都是已排好的序列，因此下
// 一轮循环中无需再去比较。

function bubbleSort(nums: number[]): number[] {
  if (nums.length <= 1) return nums

  return nums
}

const arr = [1, 4, 2, 5, 3, 6, 7]
bubbleSort(arr)
console.log(arr)
export {}

// 这个方法有问题
