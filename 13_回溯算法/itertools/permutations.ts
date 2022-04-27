/**
 * @description 返回值是按顺序从小到大的排列
 */
function* permutations<T extends ArrayLike<any>>(sequnce: T, select?: number): Generator<T[]> {
  select = select ?? sequnce.length
  yield* bt([], Array(sequnce.length).fill(false))

  function* bt(path: T[], visited: boolean[]): Generator<T[]> {
    if (path.length === select) {
      yield path.slice()
      return
    }

    for (let i = 0; i < sequnce.length; i++) {
      if (visited[i]) continue
      visited[i] = true
      path.push(sequnce[i])
      yield* bt(path, visited)
      path.pop()
      visited[i] = false
    }
  }
}

if (require.main === module) {
  console.log(...permutations([3, 2, 1], 2))
}

export { permutations }
