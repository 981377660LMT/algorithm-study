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
  private readonly N: number

  constructor(N: number) {
    this.root = new TrieNode()
    this.N = N
  }

  insert(num: number): void {
    let root = this.root
    for (let i = this.N; ~i; i--) {
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

    for (let i = this.N; ~i; i--) {
      if (!root) break
      const bit = (num >> i) & 1
      const needBit = 1 ^ bit

      if (root.children[needBit] != undefined && root.children[needBit]!.count > 0) {
        res += 1 << i
        root = root.children[needBit]!
      } else if (root.children[bit] != undefined) {
        root = root.children[bit]!
      }
    }

    return res
  }

  delete(num: number): void {
    let root = this.root
    for (let i = this.N; ~i; i--) {
      const bit = (num >> i) & 1
      root.children[bit]!.count--
      root = root.children[bit]!
    }
  }
}

export { XORTrie }
