"use strict";
// // 深度与广度遍历
// interface Node {
//   val: number
//   neighbors: Node[]
// }
// class CNode {
//   val: number
//   neighbors: Node[]
//   constructor(val: number, neighbors: Node[]) {
//     this.val = val
//     this.neighbors = neighbors
//   }
// }
// // 拷贝所有节点
// // 拷贝所有边
// const cloneGraph = (node: Node) => {
//   if (!node) return
//   const visited = new Map<Node, CNode>()
//   const dfs = (node: Node) => {
//     const nCopy = new CNode(node.val)
//     console.log(node.val)
//     visited.set(node, nCopy)
//     node.neighbors.forEach(nei => {
//       if (!visited.has(nei)) {
//         dfs(nei)
//       }
//       nCopy.neighbors?.push(visited.get(nei))
//     })
//   }
// }
