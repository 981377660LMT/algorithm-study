// 如果给定的数组已经排好序呢？你将如何优化你的算法？
// 如果 nums1 的大小比 nums2 小很多，哪种方法更优？
// 如果 nums2 的元素存储在磁盘上，内存是有限的，并且你不能一次加载所有的元素到内存中，你该怎么办？

// 双指针顺序查找相同的元素
const intersect = (nums1: number[], nums2: number[]) => {
  const res: number[] = []
  let left = 0
  let right = 0

  while (left < nums1.length && right < nums2.length) {
    if (nums1[left] < nums2[right]) left++
    else if (nums1[left] > nums2[right]) right++
    else {
      res.push(nums1[left])
      left++
      right++
    }
  }

  return res
}

// 2. 将较小的数组哈希计数，随后在另一个数组中根据哈希来寻找。
// 时间复杂度：O(max(n, m))O(max(n,m))
// 空间复杂度：O(min(n, m))O(min(n,m))

// 解法三：通过归并外排将两个数组排序后再使用排序双指针查找
// 归并排序是天然适合外部排序的算法，可以将分割后的子数组写到单个文件中，归并时将小文件合并为更大的文件。当两个数组均排序完成生成两个大文件后，即可使用双指针遍历两个文件，如此可以使空间复杂度最低。
