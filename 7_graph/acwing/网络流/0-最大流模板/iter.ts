export {}

function* makeIter<T>(iteable: Iterable<T>): Generator<T, string, undefined> {
  yield* iteable
  return 'stop'
}

const iter = makeIter([1, 2, 3])

console.log(iter.next().value)
console.log(iter.next().value)
console.log(iter.next().value)
console.log(iter.next().value)
