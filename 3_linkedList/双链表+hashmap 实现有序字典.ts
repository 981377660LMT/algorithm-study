// 这个就是Map
interface OrderedDict<K, V> {
  clear(): void
  delete(key: K): boolean
  forEach(callbackfn: (value: V, key: K, map: OrderedDict<K, V>) => void, thisArg?: any): void
  get(key: K): V | undefined
  has(key: K): boolean
  set(key: K, value: V): this
  readonly size: number
}

class OrderedDict<K, V> implements OrderedDict<K, V> {}

export default 1
