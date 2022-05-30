import { KthTree } from './255. 第K小数-查询区间第k小数'

function medianSlidingWindow(nums: number[], k: number): number[] {
  const kthTree = new KthTree(nums)
  const res: number[] = []

  for (let left = 0; left + k - 1 < nums.length; left++) {
    if (k & 1) {
      res.push(kthTree.query(left, left + k - 1, Math.floor(k / 2) + 1))
    } else {
      res.push(
        (kthTree.query(left, left + k - 1, Math.floor(k / 2)) +
          kthTree.query(left, left + k - 1, Math.floor(k / 2) + 1)) /
          2
      )
    }
  }

  return res
}

console.log(medianSlidingWindow([1, 3, -1, -3, 5, 3, 6, 7], 3))

export {}
