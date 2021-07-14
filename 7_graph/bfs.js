const graph = require('./graph')

const visited = new Set()
const q = [2]
while (q.length) {
  const n = q.shift()
  console.log(n)
  visited.add(n)
  graph[n].forEach(c => {
    if (!visited.has(c)) {
      q.push(c)
    }
  })
}
