/**
 * 带有默认值的字典.
 * 模拟python的`collections.defaultdict`.
 */
class DefaultDict<K, V> extends Map<K, V> {
  private readonly _defaultFactory: () => V

  constructor(defaultFactory: () => V, iterable?: Iterable<readonly [K, V]> | null) {
    super(iterable)
    this._defaultFactory = defaultFactory
  }

  setDefault(key: K, value: V): V {
    const old = super.get(key)
    if (old !== void 0) return old // has key
    super.set(key, value)
    return value
  }

  override get(key: K): V {
    const old = super.get(key)
    if (old !== void 0) return old // has key
    const value = this._defaultFactory()
    super.set(key, value)
    return value
  }
}

export { DefaultDict }

if (require.main === module) {
  const mp = new DefaultDict<string, number>(() => 0)
  console.log(mp.get('12'))
}
