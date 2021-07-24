interface Node {
  id: number
  children: Node[]
}

// 给出每个节点的子节点关系
const input = {
  0: [1, 2],
  1: [3, 4],
  3: [5],
} as any

// 先要规范每个节点的格式，再从根节点遍历
const map = Array.from({ length: 5 }, (_, i) => ({ id: i + 1, children: [] }))

const dfs = (id: number, root: Node) => {
  id in input &&
    input[id].forEach((childNum: number) => {
      const childNode = map[childNum - 1]
      root.children.push(childNode)
      dfs(childNum, childNode)
    })
}
const roots = map.filter(node => node.id === 1 || node.id === 2)
roots.forEach(root => dfs(root.id, root))
console.dir(roots, { depth: null })
// 重建为
// const output = [
//   {
//     id: 1,
//     children: [
//       {
//         id: 3,
//         children: [
//           {
//             id: 5,
//             children: [],
//           },
//         ],
//       },
//       {
//         id: 4,
//         children: [],
//       },
//     ],
//   },
//   {
//     id: 2,
//     children: [],
//   },
// ]

export {}
