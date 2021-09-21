class TrieNode {
  children: [TrieNode | undefined, TrieNode | undefined]
  count: number
  constructor() {
    this.children = [undefined, undefined]
    this.count = 0
  }
}

class Trie {
  private root: TrieNode
  private N: number

  constructor(N: number) {
    this.root = new TrieNode()
    this.N = N
  }

  insert(num: number) {
    let root = this.root
    for (let i = this.N; ~i; i--) {
      const bit = (num >> i) & 1
      if (!root.children[bit]) {
        root.children[bit] = new TrieNode()
      }
      root.children[bit]!.count++
      root = root.children[bit]!
    }
  }

  /**
   *
   * @param num
   * @summary
   * 求num与树中异或值严格小于high的个数
   * @description
   * 传入当前要比较的数字和limit数字，同时右移动
   * 如果limit在某位为0，那么校验数字不能出现异或为1的情况,校验数字结果会大于limit，不符合我们的要求
   * 如果limit为1，这种情况下当前位所有异或为0的数字，都是结果，需要进行累加,
   * 同时将字典树偏移到另一个分支。
   */
  search(num: number, high: number) {
    let root: TrieNode | undefined = this.root
    let res = 0

    for (let i = this.N; i >= 0; i--) {
      if (!root) break
      const bit = (num >> i) & 1
      const bitLimit = (high >> i) & 1

      if (bitLimit === 1) {
        if (root.children[bit]) {
          res += root.children[bit]!.count
        }
        // 切换分支
        root = root.children[1 - bit]
      } else {
        root = root.children[bit]
      }
    }

    return res
  }
}

/**
 *
 * @param nums 1 <= nums.length <= 2 * 104
 * @param low
 * @param high
 * 求low <= (nums[i] XOR nums[j]) <= high 的个数
 */
function countPairs(nums: number[], low: number, high: number): number {
  let res = 0
  const trie = new Trie(14)

  for (const num of nums) {
    trie.insert(num)
    res += trie.search(num, high + 1) - trie.search(num, low)
  }

  return res
}

console.log(countPairs([1, 4, 2, 7], 2, 6))

export {}
// console.log(1 ^ 4)
