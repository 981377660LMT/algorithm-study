class TrieNode {
  children: [TrieNode | undefined, TrieNode | undefined]
  count: number
  constructor() {
    this.children = [undefined, undefined]
    this.count = 0
  }
}

class XORTrie {
  private readonly root: TrieNode
  private readonly bitLength: number

  /**
   * @param bitLength 位数，即树的最大高度
   * @description 离根节点越近表示位数越高(大)
   */
  constructor(bitLength: number) {
    this.root = new TrieNode()
    this.bitLength = bitLength
  }

  insert(num: number): void {
    let root = this.root
    for (let i = this.bitLength; ~i; i--) {
      const bit = (num >> i) & 1
      if (root.children[bit] == undefined) {
        root.children[bit] = new TrieNode()
      }
      root.children[bit]!.count++
      root = root.children[bit]!
    }
  }

  /**
   *
   * @param num
   * @returns
   * 求num与树中异或最大值
   */
  search(num: number): number {
    let root = this.root
    let res = 0

    for (let i = this.bitLength; ~i; i--) {
      // if (!root) break // Trie中未插入任何数时
      const bit = (num >> i) & 1
      const needBit = 1 ^ bit

      if (root.children[needBit]?.count ?? 0 > 0) {
        res = (res << 1) | 1
        root = root.children[needBit]!
      } else if (root.children[bit]?.count ?? 0 > 0) {
        res = res << 1
        root = root.children[bit]!
      }
    }

    return res
  }

  /**
   * @param num 不存在此元素可能抛出错误
   */
  remove(num: number): void {
    let root = this.root
    for (let i = this.bitLength; ~i; i--) {
      const bit = (num >> i) & 1
      root.children[bit]!.count--
      root = root.children[bit]!
    }
  }
}

// 自举写法
class SimpleXORTrie {
  private static bitLength = 31
  private readonly children: [SimpleXORTrie | undefined, SimpleXORTrie | undefined] = [
    undefined,
    undefined,
  ]

  /**
   * @param bitLength 位数，即树的最大高度
   * @description 离根节点越近表示位数越高(大)
   */
  static setBitLength(bitLength: number): void {
    SimpleXORTrie.bitLength = bitLength
  }

  insert(num: number): void {
    let root: SimpleXORTrie = this
    for (let i = SimpleXORTrie.bitLength; ~i; i--) {
      const bit = (num >> i) & 1
      if (root.children[bit] == undefined) root.children[bit] = new SimpleXORTrie()
      root = root.children[bit]!
    }
  }

  /**
   *
   * @param num
   * @returns
   * 求num与树中异或最大值
   */
  search(num: number): number {
    let root: SimpleXORTrie = this
    let res = 0

    for (let i = SimpleXORTrie.bitLength; ~i; i--) {
      if (!root) break // Trie中未插入任何数时
      const bit = (num >> i) & 1
      const needBit = 1 ^ bit

      if (root.children[needBit] != undefined) {
        res = (res << 1) | 1
        root = root.children[needBit]!
      } else {
        res = res << 1
        root = root.children[bit]!
      }
    }

    return res
  }
}

export { XORTrie, SimpleXORTrie }
