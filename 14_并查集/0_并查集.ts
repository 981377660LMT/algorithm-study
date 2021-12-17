// 并查集适用于合并集合，查找哪些元素属于同一组
// 如果a与b属于同一组，那么就uniona与b;以后找元素属于哪个组时只需要find(这个元素)找到属于哪个根元素
// 例如很多邮箱都属于同一个人，就union这些邮箱

// 并查集是一种数据结构，并(Union)代表合并，查(Find)代表查找，用来查找两根元素是否具有公共的根。
// 集代表这是一个以字典为基础的数据结构，基本功能是合并集合中的元素，
// 查找集合中的元素。
// 并查集的典型应用是有关连通分量的问题，
// 并查集解决单个问题（添加，合并，查找）的时间复杂度都是O(h)(因为都是用的map的set和get方法)。

interface IUnionFind {
  union: (key1: number, key2: number) => boolean
  find: (key: number) => number
  isConnected: (key1: number, key2: number) => boolean
}

class UnionFind implements IUnionFind {
  // 记录每个节点的父节点
  // 如果节点互相连通（从一个节点可以达到另一个节点），那么他们的祖先是相同的。
  private readonly parent: number[]
  // rank优化 union时连接到rank较大的根上 且rank表示每个根连了多少点
  private readonly rank: number[]
  // 记录无向图连通域数量
  public part: number

  constructor(size: number) {
    this.part = size
    this.parent = Array.from({ length: size }, (_, i) => i)
    this.rank = Array(size).fill(1)
  }

  union(key1: number, key2: number): boolean {
    let root1 = this.find(key1)
    let root2 = this.find(key2)
    if (root1 === root2) return false
    // 小树root1接到大树root2下面，平衡性优化
    if (this.rank[root1] > this.rank[root2]) {
      ;[root1, root2] = [root2, root1]
    }
    this.parent[root1] = root2
    this.rank[root2] += this.rank[root1]
    this.part--
    return true
  }

  find(key: number): number {
    while (this.parent[key] != undefined && this.parent[key] !== key) {
      this.parent[key] = this.parent[this.parent[key]] // 进行路径压缩
      key = this.parent[key]
    }
    return key
  }

  isConnected(key1: number, key2: number): boolean {
    return this.find(key1) === this.find(key2)
  }
}

if (require.main === module) {
  const uf = new UnionFind(5)
  uf.union(2, 3)
  uf.union(4, 3)
  console.dir(uf, { depth: null })
  console.log(uf.find(1))
  console.log(uf.isConnected(4, 2))
  console.log(uf.isConnected(4, 1))
}

export { UnionFind }
