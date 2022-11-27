/* eslint-disable @typescript-eslint/no-non-null-assertion */
/* eslint-disable eqeqeq */

class TrieNode {
  children: [TrieNode | undefined, TrieNode | undefined] = [undefined, undefined]
  count = 0
}

class XORTrie {
  private readonly _root: TrieNode

  /**
   * 位数，即树的最大高度，离根节点越近表示位数越高(大)
   */
  private readonly _bitLength: number

  /**
   * @param upper 最大值，用于确定树的最大高度
   * @warning 需要保证答案不超过32位(js中位运算会强制转为int32)
   */
  constructor(upper: number) {
    this._bitLength = 32 - Math.clz32(upper)
    if (this._bitLength < 0) throw new Error('upper must be no more than 2**32')
    this._root = new TrieNode()
  }

  insert(num: number): void {
    let root = this._root
    for (let i = this._bitLength; ~i; i--) {
      const bit = (num >> i) & 1
      if (root.children[bit] == undefined) root.children[bit] = new TrieNode()
      root.children[bit]!.count++
      root = root.children[bit]!
    }
  }

  /**
   * 求num与树中异或最大值
   */
  search(num: number): number {
    let root = this._root
    let res = 0

    for (let i = this._bitLength; ~i; i--) {
      if (!root) break // Trie中未插入任何数时
      const bit = (num >> i) & 1
      const needBit = 1 ^ bit

      if (root.children[needBit] && root.children[needBit]!.count > 0) {
        res |= 1 << i
        root = root.children[needBit]!
      } else if (root.children[bit] && root.children[bit]!.count > 0) {
        root = root.children[bit]!
      }
    }

    return res
  }

  /**
   * 从树中删除元素num，不存在此元素会抛出错误
   */
  remove(num: number): void {
    let root = this._root
    for (let i = this._bitLength; ~i; i--) {
      if (!root) throw new Error(`fail to remove :${num} not in trie`)
      const bit = (num >> i) & 1
      root.children[bit]!.count--
      root = root.children[bit]!
    }
  }
}

export { XORTrie }
