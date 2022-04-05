class Node {
  weight: number
  isWord: boolean
  val: string
  children: Map<string, Node>

  constructor(val: string) {
    this.val = val
    this.isWord = false
    this.weight = 0
    this.children = new Map()
  }
}

class MapSum {
  private root: Node

  constructor() {
    this.root = new Node('')
  }

  /**
   *
   * @param key 如果键 key 已经存在，那么原来的键值对将被替代成新的键值对。
   * @param val
   */
  insert(key: string, val: number): void {
    let root = this.root
    for (const letter of key) {
      if (!root.children.has(letter)) root.children.set(letter, new Node(letter))
      root = root.children.get(letter)!
    }
    root.isWord = true
    root.weight = val
  }

  /**
   *
   * @param prefix 返回所有以该前缀 prefix 开头的键 key 的值的总和。
   */
  sum(prefix: string): number {
    let root = this.root
    for (const letter of prefix) {
      if (!root.children.has(letter)) return 0
      root = root.children.get(letter)!
    }
    // 此时将这个节点以下的节点的weight相加即可
    return this._sum(root)
  }

  // 遍历子树 将weight加起来
  private _sum(node: Node): number {
    let res = node.weight
    for (const next of node.children.values()) {
      res += this._sum(next)
    }
    return res
  }
}

const mapSum = new MapSum()
mapSum.insert('apple', 3)
console.log(mapSum.sum('ap')) // 3
mapSum.insert('app', 2)
console.log(mapSum.sum('ap')) // 5

export {}
