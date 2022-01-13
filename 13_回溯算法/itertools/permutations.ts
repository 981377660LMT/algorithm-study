function permutations<T>(arrs: T[]): T[][] {
  const res: T[][] = []

  const bt = (path: T[], visited: boolean[]) => {
    if (path.length === arrs.length) return res.push(path.slice())

    for (let i = 0; i < arrs.length; i++) {
      if (visited[i]) continue
      visited[i] = true
      path.push(arrs[i])
      bt(path, visited)
      path.pop()
      visited[i] = false
    }
  }
  bt([], [])

  return res
}

if (require.main === module) {
  console.log(permutations([1, 1, 0]))
}
export { permutations }
