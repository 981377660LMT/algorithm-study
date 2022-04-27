/**
 *
 * @param arr
 * @param k
 * @returns 组合可取重复元素
 */
function* combinationsWithReplacement<T extends ArrayLike<any>>(
  sequence: T,
  select?: number
): Generator<T[]> {
  select = select ?? sequence.length
  yield* bt(0, [])

  function* bt(pos: number, path: T[]): Generator<T[]> {
    if (path.length === select) {
      yield path.slice()
      return
    }

    for (let i = pos; i < sequence.length; i++) {
      path.push(sequence[i])
      yield* bt(i, path) // 唯一的区别在此：可取重复的元素
      path.pop()
    }
  }
}

if (require.main === module) {
  console.log(...combinationsWithReplacement([1, 1, 3, 4], 2))
}

export { combinationsWithReplacement }
