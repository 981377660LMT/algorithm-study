/* eslint-disable no-inner-declarations */
/* eslint-disable eqeqeq */
/* eslint-disable @typescript-eslint/no-non-null-assertion */

class TrieNode {
  count = 0
  children: [TrieNode | undefined, TrieNode | undefined] = [undefined, undefined]
}

class XORTrieSimple {
  private readonly _root: TrieNode

  /**
   * 位数，即树的最大高度，离根节点越近表示位数越高(大).
   */
  private readonly _bitLength: number

  /**
   * @param upper 最大值，用于确定树的最大高度.
   * @warning 需要保证答案不超过32位(js中位运算会强制转为int32).
   */
  constructor(upper: number) {
    this._bitLength = 32 - Math.clz32(upper)
    this._root = new TrieNode()
  }

  insert(int32: number): TrieNode {
    let root = this._root
    for (let i = this._bitLength; ~i; i--) {
      const bit = (int32 >>> i) & 1
      if (root.children[bit] == undefined) root.children[bit] = new TrieNode()
      root = root.children[bit]!
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
      root.children[bit]!.count--
      root = root.children[bit]!
    }
    return root
  }

  /**
   * 求int32与树中异或最大值.
   */
  query(int32: number): number {
    let root = this._root
    let res = 0
    for (let i = this._bitLength; ~i; i--) {
      let bit = (int32 >>> i) & 1
      if (root.children[bit ^ 1] && root.children[bit ^ 1]!.count > 0) {
        res |= 1 << i
        bit ^= 1
      }
      root = root.children[bit]!
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
}
