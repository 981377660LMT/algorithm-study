// Container is base interface that all data structures implement.
interface Container<V> {
  empty(): boolean
  size(): number
  clear(): void
  values(): V[]
  toString(): string
}

function getSortedValues<V>(container: Container<V>, comparator: (a: V, b: V) => number): V[] {
  const values = container.values()
  if (values.length <= 1) return values
  values.sort(comparator)
  return values
}

getSortedValues
export {}
