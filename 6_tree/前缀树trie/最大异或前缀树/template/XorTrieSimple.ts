/* eslint-disable no-inner-declarations */
/* eslint-disable eqeqeq */
/* eslint-disable @typescript-eslint/no-non-null-assertion */

class TrieNode {
  count = 0
  zero: TrieNode | undefined = undefined
  one: TrieNode | undefined = undefined
}

class XORTrieSimple {
  private readonly _root: TrieNode

  /**
   * 位数，即树的最大高度，离根节点越近表示位数越高(大).
   */
  private readonly _bitLength: number

  /**
   * @param upperInt32 最大值，用于确定树的最大高度.
   * @warning 需要保证答案不超过32位(js中位运算会强制转为int32).
   */
  constructor(upperInt32: number) {
    this._bitLength = 32 - Math.clz32(upperInt32)
    this._root = new TrieNode()
  }

  insert(int32: number): TrieNode {
    let root = this._root
    for (let i = this._bitLength; ~i; i--) {
      const bit = (int32 >>> i) & 1
      if (!bit) {
        if (!root.zero) root.zero = new TrieNode()
        root = root.zero
      } else {
        if (!root.one) root.one = new TrieNode()
        root = root.one
      }
      root.count++
    }
    return root
  }

  /**
   * 从树中删除元素.
   * 必须保证树中存在该元素.
   */
  remove(int32: number): TrieNode {
    let root = this._root
    for (let i = this._bitLength; ~i; i--) {
      const bit = (int32 >>> i) & 1
      root = bit ? root.one! : root.zero!
      root.count--
    }
    return root
  }

  /**
   * 求int32与树中异或最大值.
   */
  query(int32: number): number {
    let root: TrieNode | undefined = this._root
    let res = 0
    for (let i = this._bitLength; ~i; i--) {
      if (!root) break
      let bit = (int32 >>> i) & 1
      if (!bit) {
        if (root.one && root.one.count > 0) {
          res |= 1 << i
          root = root.one
        } else {
          root = root.zero
        }
      } else if (root.zero && root.zero.count > 0) {
        res |= 1 << i
        root = root.zero
      } else {
        root = root.one
      }
    }
    return res
  }
}

export { XORTrieSimple }

if (require.main === module) {
  // 2935. 找出强数对的最大异或值 II
  // https://leetcode.cn/problems/maximum-strong-pair-xor-ii/description/
  function maximumStrongPairXor(nums: number[]): number {
    nums.sort((a, b) => a - b)
    let res = 0
    let left = 0
    const n = nums.length
    const trie = new XORTrieSimple(Math.max(...nums))
    for (let right = 0; right < n; right++) {
      trie.insert(nums[right])
      while (left <= right && nums[right] > 2 * nums[left]) {
        trie.remove(nums[left])
        left++
      }
      res = Math.max(res, trie.query(nums[right]))
    }
    return res
  }

  // https://leetcode.cn/problems/maximum-xor-of-two-numbers-in-an-array/description/
  function findMaximumXOR(nums: number[]): number {
    const trie = new XORTrieSimple(Math.max(...nums))
    let res = 0
    for (let i = 0; i < nums.length; i++) {
      trie.insert(nums[i])
      res = Math.max(res, trie.query(nums[i]))
    }
    return res
  }
}
