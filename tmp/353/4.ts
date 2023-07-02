import { CountingSort } from '../../9_排序和搜索/CountingSort'

export {}

// 一个长度为 n 下标从 0 开始的整数数组 arr 的 不平衡数字 定义为，在 sarr = sorted(arr) 数组中，满足以下条件的下标数目：

// 0 <= i < n - 1 ，和
// sarr[i+1] - sarr[i] > 1
// 这里，sorted(arr) 表示将数组 arr 排序后得到的数组。

// 给你一个下标从 0 开始的整数数组 nums ，请你返回它所有 子数组 的 不平衡数字 之和。

// 子数组指的是一个数组中连续一段 非空 的元素序列。

// !1.当数组元素个数比较少时，类型数组的性能并不比普通数组好.
// !2.js的sort并没有那么快, 桶排序的性能比sort快很多

function sumImbalanceNumbers(N: number[]): number {
  const nums = new Uint16Array(N)
  const C = new CountingSort(Math.max(...nums))

  let res = 0
  for (let i = 0; i < nums.length; i++) {
    for (let j = i; j < nums.length; j++) {
      const sorted = C.sort(nums.subarray(i, j + 1))
      for (let k = 0; k < sorted.length - 1; k++) res += +(sorted[k + 1] - sorted[k] > 1)
    }
  }

  return res
}

if (require.main === module) {
  const n = 1000
  const arr = new Array(n).fill(0).map((_, i) => n - i)
  console.time('sumImbalanceNumbers')
  console.log(sumImbalanceNumbers(arr))
  console.timeEnd('sumImbalanceNumbers')
}
