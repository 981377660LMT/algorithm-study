// 二进制分组，模拟二进制加法。

interface IPreprocessor<V> {
  add(value: V): void
  build(): void
  clear(): void
}

class BinaryGrouping<V> {
  private _groups: V[][] = []
  private _preprocessors: IPreprocessor<V>[] = []
  private _createPreprocessor: () => IPreprocessor<V>

  constructor(createContainer: () => IPreprocessor<V>) {
    this._createPreprocessor = createContainer
  }

  add(value: V): void {
    let k = 0
    // 二进制加法进位，每次进位后都会将之前的所有元素都加入到新的分组中并清空之前的分组
    while (k < this._groups.length && this._groups[k].length) {
      k++
    }
    if (k === this._groups.length) {
      this._groups.push([])
      this._preprocessors.push(this._createPreprocessor())
    }
    this._groups[k].push(value)
    this._preprocessors[k].add(value)
    for (let i = 0; i < k; i++) {
      this._groups[i].forEach(v => this._preprocessors[k].add(v))
      this._groups[k].push(...this._groups[i])
      this._preprocessors[i].clear()
      this._groups[i].length = 0
    }
    this._preprocessors[k].build()
  }

  query(onQuery: (p: IPreprocessor<V>) => void, ignoreEmpty = true): void {
    for (let i = 0; i < this._preprocessors.length; i++) {
      if (ignoreEmpty && this._groups[i].length === 0) {
        continue
      }
      onQuery(this._preprocessors[i])
    }
  }
}

export { BinaryGrouping }
