/**
 * IteratorWithIndex is stateful iterator for ordered containers
 * whose values can be fetched by an index.
 */
interface IteratorWithIndex<V> {
  getIndex(): number
  getValue(): V

  /**
   * moves the iterator to the next element
   * and returns true if there was a next element in the container.
   */
  next(): boolean

  /**
   * `resets` the iterator to the beginning(one-before-first).
   */
  begin(): void

  /**
   * moves the iterator to the first element.
   */
  first(): boolean

  /**
   * moves the iterator to the next element from current position that `satisfies the condition
   * given by the passed function`, and returns true if there was a next element in the container.
   * @example
   * ```go
   * func (iterator *_Iterator) NextTo(f func(key interface{}, value interface{}) bool) bool {
      for iterator.Next() {
        key, value := iterator.Key(), iterator.Value()
        if f(key, value) {
          return true
        }
      }
      return false
    }
    ```
   */
  nextTo(f: (index: number, value: V) => boolean): boolean
}

interface IteratorWithKey<K, V> {
  getIndex(): number
  getValue(): V
  next(): boolean
  begin(): void
  first(): boolean
  nextTo(f: (key: K, value: V) => boolean): boolean
}

interface ReverseIteratorWithIndex<V> extends IteratorWithIndex<V> {
  prev(): boolean
  end(): void
  last(): boolean
  prevTo(f: (index: number, value: V) => boolean): boolean
}

interface ReverseIteratorWithKey<K, V> extends IteratorWithKey<K, V> {
  prev(): boolean
  end(): void
  last(): boolean
  prevTo(f: (key: K, value: V) => boolean): boolean
}

export {}
