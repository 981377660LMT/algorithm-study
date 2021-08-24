const json = {
  a: { b: { c: 1 } },
  d: [1, 2],
}

const dfs = (n, path) => {
  if (typeof n !== 'object') {
    console.log(n, path)
  }
  Object.keys(n).forEach(k => {
    path.push(k)
    dfs(n[k], path)
    path.pop()
  })
}

dfs(json, [])
