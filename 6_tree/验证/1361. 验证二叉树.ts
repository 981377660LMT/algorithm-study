// 二叉树的充分条件：
// 1.根节点唯一
// 2.从根结点出发，其他节点都被访问到且只被访问一次

// 如果节点 i 没有左子节点，那么 leftChild[i] 就等于 -1。右子节点也符合该规则。
// bfs更好 因为可以直接return
function validateBinaryTreeNodes(n: number, leftChild: number[], rightChild: number[]): boolean {
  let root = Infinity
  const children = new Set<number>([...leftChild, ...rightChild])
  for (let i = 0; i < n; i++) {
    if (!children.has(i)) {
      if (root === Infinity) root = i
      else return false
    }
  }

  const queue = [root]
  const visited = Array<boolean>(n).fill(false)

  while (queue.length > 0) {
    const cur = queue.shift()!
    if (visited[cur]) return false
    visited[cur] = true
    if (leftChild[cur] !== -1) queue.push(leftChild[cur])
    if (rightChild[cur] !== -1) queue.push(rightChild[cur])
  }

  return visited.every(bool => bool)
}

export {}
console.log(validateBinaryTreeNodes(4, [1, -1, 3, -1], [2, -1, -1, -1]))
