// Binary Lifting 树上倍增
// 先预处理每一个点走 2^k步可以到达的祖先
// 之后求第t个祖先时，不断二进制拆分即可

class TreeAncestor2 {
  private parent: number[]
  private n: number

  // 树以父节点数组的形式给出，其中 parent[i] 是节点 i 的父节点。树的根节点是编号为 0 的节点。
  constructor(n: number, parent: number[]) {
    this.parent = parent
    this.n = n
  }

  // 函数返回节点 node 的第 k 个祖先节点。如果不存在这样的祖先节点，返回 -1
  getKthAncestor(node: number, k: number): number {
    let res = node
    for (let i = 0; i < k; i++) {
      res = this.parent[res]
      if (res < 0) return -1
    }
    return res
  }
}

export {}
//////////////////////////////////////////////////////////////////////////////////////
// 1 <= k <= n <= 5*10^4
// 1.如果每次查找的步长是1，会超时
// 如果把一个数字转换成二进制，从左至右第一个bit是1的，就是“尽量长的步子”
// 比如 13 = 8+4+1 = 1101
// 那我们查找的时候，步长 依次是8，4， 1
// 2.设计一个fa数组
// fa[i][j]: 结点i 的， 距离为2^j 的祖先结点
class TreeAncestor {
  private fa: number[][]

  /**
   *
   * @param n  节点个数
   * @param parent  树以父节点数组的形式给出，其中 parent[i] 是节点 i 的父节点。树的根节点是编号为 0 的节点。
   */
  constructor(n: number, parent: number[]) {
    this.fa = Array.from({ length: n }, () => Array(16).fill(-1))

    for (let i = 0; i < n; i++) {
      this.fa[i][0] = parent[i]
    }

    // 注意先指数
    for (let j = 0; j < 15; j++) {
      for (let i = 0; i < n; i++) {
        if (this.fa[i][j] === -1) this.fa[i][j + 1] = -1
        else this.fa[i][j + 1] = this.fa[this.fa[i][j]][j] // 2^i*2^i === 2^(i+1)
      }
    }
  }

  /**
   *
   * @param node
   * @param k 1 <= k <= n <= 5*10^4
   * @returns 函数返回节点 node 的第 k 个祖先节点。如果不存在这样的祖先节点，返回 -1
   */
  getKthAncestor(node: number, k: number): number {
    if (node <= 0) return -1
    let bit = 0
    while (k) {
      if (k & 1) node = this.fa[node][bit]
      if (node === -1) break
      bit++
      k >>>= 1
    }
    return node
  }
}

if (require.main === module) {
  const lca = new TreeAncestor2(7, [-1, 0, 0, 1, 1, 2, 2])
  console.log(lca.getKthAncestor(3, 1))
  console.log(lca.getKthAncestor(5, 2))
  console.log(lca.getKthAncestor(6, 3))
}

// 倍增法（英语：binary lifting），顾名思义就是翻倍。
// 它能够使线性的处理转化为对数级的处理，大大地优化时间复杂度。
// 这个方法在很多算法中均有应用，其中最常用的是RMQ 区间最大（最小）值问题和求LCA（最近公共祖先）

// 思考 求出 list以后怎么求两个点的LCA呢?
// 1. 根据depth数组 下面的结点往上跳 二进制分解跳到他们一样高
// 2. 最左能力二分 两个一起上跳mid 一样就缩小步伐 知道跳到LCA的下一层(在跳一步就到LCA)

// 实际上题中给出的 parent数组是可以在dfs中更新求出来的
export { TreeAncestor }
