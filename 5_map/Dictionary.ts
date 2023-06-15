/**
 * A dictionary that maps values to unique ids.
 */
class Dictionary<V> {
  private readonly _valueToId = new Map<V, number>()
  private readonly _idToValue: V[] = []

  id(value: V): number {
    const res = this._valueToId.get(value)
    if (res !== void 0) return res
    const id = this._idToValue.length
    this._idToValue.push(value)
    this._valueToId.set(value, id)
    return id
  }

  value(id: number): V | undefined {
    if (id < 0 || id >= this._idToValue.length) return void 0
    return this._idToValue[id]
  }

  get size(): number {
    return this._idToValue.length
  }
}

export { Dictionary }

if (require.main === module) {
  const dict = new Dictionary<string>()
  console.log(dict.id('a'))
  console.log(dict.id('b'))
  console.log(dict.id('a'))
  console.log(dict.value(0))
  console.log(dict.value(1))
  console.log(dict.value(2))
  console.log(dict.size)
}
