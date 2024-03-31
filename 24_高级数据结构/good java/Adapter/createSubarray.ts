interface ISubarray<T> {
  readonly length: number
  get(index: number): T | undefined
  set(index: number, value: T): void
  subarray(start: number, end: number): ISubarray<T>
}

const EMPTY_SUBARRAY: ISubarray<any> = {
  length: 0,
  get: () => undefined,
  set: () => {},
  subarray: () => EMPTY_SUBARRAY
}

class ArraySubarrayAdapter<T> implements ISubarray<T> {
  private readonly _data: T[]
  private readonly _start: number
  private readonly _end: number

  constructor(data: T[], start = 0, end = data.length) {
    this._data = data
    this._start = start
    this._end = end
  }

  get(index: number): T | undefined {
    const n = this.length
    if (index < 0) index += n
    if (index < 0 || index >= n) return undefined
    return this._data[this._start + index]
  }

  set(index: number, value: T): void {
    const n = this.length
    if (index < 0) index += n
    if (index < 0 || index >= n) return
    this._data[this._start + index] = value
  }

  subarray(start: number, end: number): ISubarray<T> {
    const n = this.length
    if (start < 0) start += n
    if (start < 0) start = 0
    if (end < 0) end += n
    if (end > n) end = n
    if (start >= end) return EMPTY_SUBARRAY
    return new ArraySubarrayAdapter(this._data, this._start + start, this._start + end)
  }

  get length(): number {
    return this._end - this._start
  }
}

function createSubarray<T>(array: T[], start = 0, end = array.length): ISubarray<T> {
  const n = array.length
  if (start < 0) start += n
  if (start < 0) start = 0
  if (end < 0) end += n
  if (end > n) end = n
  return new ArraySubarrayAdapter(array, start, end)
}

export { ISubarray, createSubarray }

if (require.main === module) {
  const arr = [1, 2, 3, 4, 5]
  const subarray = createSubarray(arr, 1, 4)
  console.log(subarray.get(0), subarray.get(1), subarray.get(2))
  subarray.set(1, 10)
  console.log(arr)
}
