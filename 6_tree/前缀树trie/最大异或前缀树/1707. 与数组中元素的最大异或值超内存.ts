import { XORTrie } from './XORTrie'

/**
 *
 * @param nums
 * @param queries  查询数组 queries ，其中 queries[i] = [xi, mi] 。
 * 第 i 个查询的答案是 xi 和任何 nums 数组中不超过 mi 的元素按位异或（XOR）得到的最大值
 * 如果 nums 中的所有元素都大于 mi，最终答案就是 -1 。
 */
function maximizeXor(nums: number[], queries: number[][]): number[] {
  nums.sort((a, b) => a - b)
  // 对query离线排序
  // queries = queries.map((item, index) => [...item, index]).sort((a, b) => a[1] - b[1])
  for (let i = 0; i < queries.length; i++) queries[i][2] = i

  queries.sort((a, b) => a[1] - b[1])
  const trie = new XORTrie(31)
  const res = Array<number>(queries.length).fill(0)

  let idx = 0
  for (let i = 0; i < queries.length; i++) {
    while (idx < nums.length && nums[idx] <= queries[i][1]) {
      trie.insert(nums[idx])
      idx++
    }

    if (nums[0] > queries[i][1]) {
      res[queries[i][2]] = -1
    } else {
      res[queries[i][2]] = trie.search(queries[i][0])
    }
  }

  return res
}

console.log(
  maximizeXor(
    [0, 1, 2, 3, 4],
    [
      [3, 1],
      [1, 3],
      [5, 6],
    ]
  )
)
// 输出：[3,3,7]
// 解释：
// 1) 0 和 1 是仅有的两个不超过 1 的整数。0 XOR 3 = 3 而 1 XOR 3 = 2 。二者中的更大值是 3 。
// 2) 1 XOR 2 = 3.
// 3) 5 XOR 2 = 7.
console.log(
  maximizeXor(
    [5, 2, 4, 6, 6, 3],
    [
      [12, 4],
      [8, 1],
      [6, 3],
    ]
  )
)
export {}

console.log(Math.log10(2))
