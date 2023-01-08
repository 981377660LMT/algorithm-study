// 给你无向 连通 图中一个节点的引用，请你返回该图的 深拷贝（克隆）。

class Node {
  val: number
  neighbors: Node[]
  constructor(val = 0, neighbors = []) {
    this.val = val
    this.neighbors = neighbors
  }
}

/**
 *
 * @param node 给你无向 连通 图中一个节点的引用，请你返回该图的 深拷贝（克隆）。
 * 无向图是一个简单图，这意味着图中没有重复的边，也没有自环
 */
function cloneGraph(node: Node | null): Node | null {
  const visited = new WeakMap<Node, Node>()
  return dfs(node)

  /**
   * @param node 根据实际的node返回克隆的node
   * @returns
   */
  function dfs(node: Node | null): Node | null {
    if (!node) return null
    if (visited.has(node)) return visited.get(node)!

    const res = new Node(node.val, [])
    visited.set(node, res)
    for (const next of node.neighbors) {
      res.neighbors.push(dfs(next)!) // 将临边的克隆也加入克隆的临边
    }

    return res
  }
}

export {}
// 138. 复制带随机指针的链表
// 这道题的这种解决方法和一个克隆链表的题目很像，
// 其中 Node 中的一个指针随机指向链表中的任意一个节点。
// 解决方法就是，
// 先顺序遍历节点，然后记录实体和克隆 Node 的映射对。
// 然后再次遍历链表，通过映射关系再为克隆链表指向映射的随机节点。

// 哈希+dfs
