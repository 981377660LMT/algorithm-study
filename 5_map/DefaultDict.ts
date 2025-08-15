/* eslint-disable max-len */

/**
 * 带有默认值的字典.
 * 模拟python的`collections.defaultdict`.
 */
class DefaultDict<K, V> extends Map<K, V> {
  private readonly _defaultFactory: (self: DefaultDict<K, V>) => V

  constructor(
    defaultFactory: (self: DefaultDict<K, V>) => V,
    iterable?: Iterable<readonly [K, V]> | null
  ) {
    super(iterable)
    this._defaultFactory = defaultFactory
  }

  override get(key: K): V {
    if (super.has(key)) return super.get(key)!
    const value = this._defaultFactory(this)
    super.set(key, value)
    return value
  }

  setDefault(key: K, value: V): V {
    if (super.has(key)) return super.get(key)!
    super.set(key, value)
    return value
  }
}

export { DefaultDict }

if (require.main === module) {
  const mp = new DefaultDict<string, number>(d => d.size)
  console.log(mp.get('12'))
  console.log(mp.get('121'))
  console.log(mp.get('12'))

  const idPool = new DefaultDict<number, number>(d => d.size)
  console.log(idPool.get(1))
  console.log(idPool.get(2))
  console.log(idPool.get(1))
  console.log(idPool.get(3))
}
