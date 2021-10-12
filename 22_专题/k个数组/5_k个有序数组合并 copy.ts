/**
 * @param {number[][]} arrList
 * non-descending integer array
 * @return {number[]}
 */
function merge(arrList: number[][]): number[] {
  // your code here
  if (arrList.length === 0) return []
  if (arrList.length === 1) return arrList[0]
  if (arrList.length === 2) return mergeTwo(arrList[0], arrList[1])
  const mid = arrList.length >> 1
  const left = arrList.slice(0, mid)
  const right = arrList.slice(mid, arrList.length)
  return mergeTwo(merge(left), merge(right))
}

function mergeTwo(nums1: number[], nums2: number[]) {
  const res: number[] = []
  let i = 0
  let j = 0

  while (i < nums1.length && j < nums2.length) {
    if (nums1[i] < nums2[j]) {
      res.push(nums1[i])
      i++
    } else {
      res.push(nums2[j])
      j++
    }
  }

  // 连接剩余的元素，防止没有把两个数组遍历完整
  return [...res, ...nums1.slice(i), ...nums2.slice(j)]
}

console.log(
  merge([
    [1, 1, 1, 100, 1000, 10000],
    [1, 2, 2, 2, 200, 200, 1000],
    [1000000, 10000001],
    [2, 3, 3],
  ])
)
