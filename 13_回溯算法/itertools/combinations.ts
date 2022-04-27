function* combinations<T extends ArrayLike<any>>(sequence: T, select?: number): Generator<T[]> {
  select = select ?? sequence.length
  yield* bt(0, [])

  function* bt(pos: number, path: T[]): Generator<T[]> {
    if (path.length === select) {
      yield path.slice()
      return
    }

    for (let i = pos; i < sequence.length; i++) {
      path.push(sequence[i])
      yield* bt(i + 1, path) // 唯一的区别在此：是否可取重复的元素
      path.pop()
    }
  }
}

if (require.main === module) {
  console.log(...combinations([1, 1, 3, 4], 2))
}

export { combinations }
