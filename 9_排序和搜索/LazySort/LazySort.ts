// 局部排序/TopK
// n不大的话，直接排序
// 否则懒排序.

interface ILazySort<V> {
  get(index: number): V | undefined
  resetData(data: V[]): void
  resetComparator(comparator: (a: V, b: V) => number): void
}

type NthElement<V> = (arr: V[], start: number, end: number, nth: number) => number

class LazySort<V> implements ILazySort<V> {
  private readonly _nthElement: NthElement<V>
  private _data: V[]
  private _comparator: (a: V, b: V) => number
  private _prefixSortedTo = 0
  private _suffixSortedTo = 0

  constructor(data: V[], comparator: (a: V, b: V) => number) {}

  get(index: number): V | undefined {}

  resetData(data: V[]): void {}

  resetComparator(comparator: (a: V, b: V) => number): void {}
}

export { LazySort }
