class TrieNode {
  children: [TrieNode | undefined, TrieNode | undefined]
  constructor() {
    this.children = [undefined, undefined]
  }
}

class Trie {
  private root: TrieNode

  constructor() {
    this.root = new TrieNode()
  }

  insert(num: number) {
    let root = this.root
    for (let i = 31; ~i; i--) {
      const bit = (num >> i) & 1
      if (!root.children[bit]) {
        root.children[bit] = new TrieNode()
      }
      root = root.children[bit]!
    }
  }

  /**
   *
   * @param num
   * @returns
   * 求num与树中异或最大值
   */
  search(num: number) {
    let root = this.root
    let res = 0

    for (let i = 31; ~i; i--) {
      if (!root) break
      const bit = (num >> i) & 1
      const needBit = 1 - bit
      // 贪心
      if (root.children[needBit]) {
        res = (res << 1) + 1
        root = root.children[needBit]!
      } else if (root.children[bit]) {
        res = res << 1
        root = root.children[bit]!
      }
    }

    return res
  }
}

/**
 *
 * @param nums
 * @param queries  查询数组 queries ，其中 queries[i] = [xi, mi] 。
 * 第 i 个查询的答案是 xi 和任何 nums 数组中不超过 mi 的元素按位异或（XOR）得到的最大值
 * 如果 nums 中的所有元素都大于 mi，最终答案就是 -1 。
 */
function maximizeXor(nums: number[], queries: number[][]): number[] {
  let idx = 0
  nums.sort((a, b) => a - b)
  // 对query离线排序
  // queries = queries.map((item, index) => [...item, index]).sort((a, b) => a[1] - b[1])
  for (let i = 0; i < queries.length; i++) {
    queries[i][2] = i
  }
  queries.sort((a, b) => a[1] - b[1])
  const trie = new Trie()
  const res = Array<number>(queries.length).fill(0)

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
