interface EnumerableWithIndex<V> {
  forEach(callback: (value: V, index: number) => void): void
  some(callback: (value: V, index: number) => boolean): boolean
  every(callback: (value: V, index: number) => boolean): boolean
  // [-1, undefined] if not found
  find(callback: (value: V, index: number) => boolean): [i: number, v: V | undefined]
}

interface EnumerableWithKey<K, V> {
  forEach(callback: (value: V, key: K) => void): void
  some(callback: (value: V, key: K) => boolean): boolean
  every(callback: (value: V, key: K) => boolean): boolean
  // [undefined, undefined] if not found
  find(callback: (value: V, key: K) => boolean): [k: K | undefined, v: V | undefined]
}

export {}
