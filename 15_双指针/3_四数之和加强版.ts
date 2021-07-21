// 灵活选择键值

// 给定四个包含整数的数组列表 A , B , C , D ,计算有多少个元组 (i, j, k, l) ，
// 使得 A[i] + B[j] + C[k] + D[l] = 0。
// 所有的 A, B, C, D 具有相同的长度 N

// 思路:查找表记录当前值与需要的值
// 可以记a+b 记在map 循环c+d 匹配map O(n^2)
const fourSumCount = (
  nums1: number[],
  nums2: number[],
  nums3: number[],
  nums4: number[]
): number => {
  if (nums1.length === 0) return 0
  const n = nums1.length
  const map = new Map<number, number>()
  let res = 0

  for (let i = 0; i < n; i++) {
    for (let j = 0; j < n; j++) {
      const key = 0 - nums1[i] - nums2[j]
      map.set(key, map.get(key)! + 1 || 1)
    }
  }

  for (let i = 0; i < n; i++) {
    for (let j = 0; j < n; j++) {
      const sum = nums3[i] + nums4[j]
      if (map.has(sum)) res += map.get(sum)!
    }
  }

  return res
}

console.log(fourSumCount([1, 2], [-2, -1], [-1, 2], [0, 2]))

export {}
