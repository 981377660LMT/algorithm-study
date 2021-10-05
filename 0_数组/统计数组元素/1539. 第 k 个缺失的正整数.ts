// 给你一个 严格升序排列 的正整数数组 arr 和一个整数 k 。
// 请你找到这个数组里第 k 个缺失的正整数。

// 一个不缺失元素的序列，会有arr[i]=i+1这样的关系
// 我们要找第k个缺失的，换言之，只要arr[i]-i-1 == k,我们便找到了题目要找的数字。(每缺一个,arr[i]与(i+1)的差值都会扩大1)
function findKthPositive1(arr: number[], k: number): number {
  for (let i = 0; i < arr.length; i++) {
    if (arr[i] - i - 1 >= k) return k + i // 画图可知
  }
  return k + arr.length // 画图可知
}

// 上面的解法没用到排序性质
// ，已经找出了 arr[i]-i-1 > 0关系之后，我们可以利用上述的线性查找的方式改为二分查找的方式。
function findKthPositive2(arr: number[], k: number): number {
  let l = 0
  let r = arr.length - 1
  while (l <= r) {
    const mid = (l + r) >> 1
    const missing = arr[mid] - (mid + 1)
    if (missing >= k) r = mid - 1
    else l = mid + 1
  }
  return k + l
}

console.log(findKthPositive2([2, 3, 4, 7, 11], 5))
// 解释：缺失的正整数包括 [1,5,6,8,9,10,12,13,...] 。第 5 个缺失的正整数为 9 。
findKthPositive1([1, 2, 3, 4], 2)
