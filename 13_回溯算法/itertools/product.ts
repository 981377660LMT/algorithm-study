/**
 *
 * @param arr
 * @returns 笛卡尔积
 */
function* product<T extends Iterable<any>>(...iterables: T[]): Generator<T[]> {
  yield* bt(0, [])

  function* bt(i: number, path: T[]): Generator<T[]> {
    if (path.length === iterables.length) {
      yield path.slice()
      return
    }

    for (const choose of iterables[i]) {
      path.push(choose)
      yield* bt(i + 1, path)
      path.pop()
    }
  }
}

if (require.main === module) {
  console.log(...product(['A', 'a'], ['1'], ['B', 'b'], ['2']))
}

export { product }
