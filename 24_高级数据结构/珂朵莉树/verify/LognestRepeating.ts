import { ErasableHeap } from '../../../8_heap/ErasableHeap'
import { ODT } from '../ODT-fastset'

/**
 * 维护相同元素的最长连续长度.
 */
class LongestRepeating<T> {
  private readonly _arr: ArrayLike<T>
  private readonly _odt: ODT<T>
  private readonly _sortedLens: ErasableHeap<number>

  constructor(arr: ArrayLike<T>) {}

  queryAll(): number {}

  query(start: number, end: number): number {}

  update(start: number, end: number, value: T): void {}

  set(index: number, value: T): void {}
}

export { LongestRepeating }
