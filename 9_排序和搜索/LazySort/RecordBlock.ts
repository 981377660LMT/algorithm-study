interface IRecordBlock<E> {
  add(value: E): void
  remove(value: E): void
  get(index: number): E
  reset(compareFn: (a: E, b: E) => number): void
}

class RecordBlock<E> implements IRecordBlock<E> {
  private _data: E[] = []
  private _compareFn: (a: E, b: E) => number

  constructor(compareFn: (a: E, b: E) => number) {
    this._compareFn = compareFn
  }

  add(value: E): void {
    this._data.push(value)
  }

  remove(value: E): void {
    const index = this._data.findIndex(v => v === value)
    if (index !== -1) {
      this._data.splice(index, 1)
    }
  }

  get(index: number): E {
    return this._data[index]
  }

  reset(compareFn: (a: E, b: E) => number): void {
    this._compareFn = compareFn
  }
}
