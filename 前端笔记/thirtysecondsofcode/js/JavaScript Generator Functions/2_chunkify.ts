const x = new Set([1, 2, 1, 3, 4, 1, 2, 5])
console.log([...chunkify(x, 2)]) // [[1, 2], [3, 4], [5]]

function* chunkify<T>(iter: Iterable<T>, size: number): Generator<T[], void, unknown> {
  let chunk: T[] = []

  for (const val of iter) {
    chunk.push(val)
    if (chunk.length === size) {
      yield chunk
      chunk = []
    }
  }

  if (chunk.length) yield chunk
}

export {}
