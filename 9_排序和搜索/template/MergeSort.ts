function mergeSort(arr: number[]): number[] {
  if (arr.length <= 1) return arr
  const mid = arr.length >>> 1
  const left = arr.slice(0, mid)
  const right = arr.slice(mid)
  return merge(mergeSort(left), mergeSort(right))

  function merge(nums1: number[], nums2: number[]): number[] {
    const res: number[] = []
    let i = 0
    let j = 0
    while (i < nums1.length && j < nums2.length) {
      if (nums1[i] < nums2[j]) {
        res.push(nums1[i++])
      } else {
        res.push(nums2[j++])
      }
    }

    while (i < nums1.length) {
      res.push(nums1[i++])
    }

    while (j < nums2.length) {
      res.push(nums2[j++])
    }

    return res
  }
}

export {}

if (require.main === module) {
  const arr = [4, 2, 100, 99, 10000, -1, 99, 2]
  console.log(mergeSort(arr))
}
