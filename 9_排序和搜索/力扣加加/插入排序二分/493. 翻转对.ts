// 如果 i < j 且 nums[i] > 2*nums[j] 我们就将 (i, j) 称作一个重要翻转对。
function reversePairs(nums: number[]): number {
  let count = 0

  // 操作放到这里
  const merge = (nums1: number[], nums2: number[]): number[] => {
    // 统计阶段
    const ll = nums1.length
    const rl = nums2.length
    const res: number[] = []
    let l = 0
    let r = 0

    while (l < ll && r < rl) {
      if (nums1[l] > 2 * nums2[r]) {
        count += ll - l
        r++
      } else {
        l++
      }
    }

    // 排序阶段
    l = 0
    r = 0
    while (l < ll && r < rl) {
      if (nums1[l] < nums2[r]) {
        res.push(nums1[l])
        l++
      } else {
        res.push(nums2[r])
        r++
      }
    }

    return [...res, ...nums1.slice(l), ...nums2.slice(r)]
  }

  const mergeSort = (nums: number[]): number[] => {
    if (nums.length <= 1) return nums
    const mid = nums.length >> 1
    const left = nums.slice(0, mid)
    const right = nums.slice(mid)
    return merge(mergeSort(left), mergeSort(right))
  }

  mergeSort(nums)
  return count
}

console.log(reversePairs([2, 4, 3, 5, 1]))
