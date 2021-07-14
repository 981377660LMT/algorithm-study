interface Graph {
  [key: number]: number[]
}

const graph: Graph = {
  0: [1, 2],
  1: [2],
  2: [0, 3],
  3: [3],
}

const visited = new Set()

// const dfs = (n: number) => {
//   console.log(n)
//   // 多了一步
//   visited.add(n)
//   graph[n].forEach(c => !visited.has(c) && dfs(c))
// }

const bfs = (n: number) => {
  visited.add(n)
  const queue = [n]

  while (queue.length) {
    const head = queue.shift()
    console.log(head)
    visited.add(head)

    graph[head!].forEach(c => {
      if (!visited.has(c)) {
        queue.push(c)
      }
    })
  }

  // visited.add(n)
}

// dfs(2)
bfs(2)

export {}
