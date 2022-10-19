function* combinations<E>(array: ArrayLike<E>, select: number): Generator<E[]> {
  yield* bt(0, [])

  function* bt(pos: number, path: E[]): Generator<E[]> {
    if (path.length === select) {
      yield path.slice()
      return
    }

    for (let i = pos; i < array.length; i++) {
      path.push(array[i])
      yield* bt(i + 1, path) // 唯一的区别在此：是否可取重复的元素
      path.pop()
    }
  }
}

if (require.main === module) {
  console.log(...combinations([1, 1, 3, 4], 2))
}

export { combinations }
