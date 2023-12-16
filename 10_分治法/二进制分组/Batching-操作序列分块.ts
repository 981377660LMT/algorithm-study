// Batching-操作序列分块

interface IPreprocessor<V> {
  add(value: V): void
  build(): void
  clear(): void
}

class Batching<V> {}

export { Batching }
