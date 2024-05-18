/** ISlice. */
interface ISequence<T> {
  readonly length: number
  at(index: number): T | undefined
  set(index: number, value: T): void
  subsequence(start: number, end: number): ISequence<T>
}

interface MutableArrayLike<T> {
  readonly length: number
  [index: number]: T
}

const EMPTY_SEQUENCE: ISequence<any> = {
  length: 0,
  at: () => undefined,
  set: () => {},
  subsequence: () => EMPTY_SEQUENCE
}

/**
 * SliceWrapper.
 */
class ArraySequenceAdapter<T> implements ISequence<T> {
  private readonly _data: MutableArrayLike<T>
  private readonly _start: number
  private readonly _end: number

  constructor(data: MutableArrayLike<T>, start = 0, end = data.length) {
    this._data = data
    this._start = start
    this._end = end
  }

  at(index: number): T | undefined {
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

  subsequence(start: number, end: number): ISequence<T> {
    const n = this.length
    if (start < 0) start += n
    if (start < 0) start = 0
    if (end < 0) end += n
    if (end > n) end = n
    if (start >= end) return EMPTY_SEQUENCE
    return new ArraySequenceAdapter(this._data, this._start + start, this._start + end)
  }

  get length(): number {
    return this._end - this._start
  }
}

class FunctionSequenceAdapter<T> implements ISequence<T> {
  private readonly _f: (i: number) => T
  private readonly _start: number
  private readonly _end: number

  constructor(f: (i: number) => T, start: number, end: number) {
    this._f = f
    this._start = start
    this._end = end
  }

  at(index: number): T | undefined {
    const n = this.length
    if (index < 0) index += n
    if (index < 0 || index >= n) return undefined
    return this._f(this._start + index)
  }

  set(): void {
    throw new Error('Not supported')
  }

  subsequence(start: number, end: number): ISequence<T> {
    const n = this.length
    if (start < 0) start += n
    if (start < 0) start = 0
    if (end < 0) end += n
    if (end > n) end = n
    if (start >= end) return EMPTY_SEQUENCE
    return new FunctionSequenceAdapter(this._f, this._start + start, this._start + end)
  }

  get length(): number {
    return this._end - this._start
  }
}

function createSequence<T>(data: MutableArrayLike<T>, start?: number, end?: number): ISequence<T>
function createSequence<T>(f: (i: number) => T, start: number, end: number): ISequence<T>
function createSequence<T>(
  dataOrF: ArrayLike<T> | ((i: number) => T),
  start?: number,
  end?: number
): ISequence<T> {
  if (typeof dataOrF === 'function') {
    return new FunctionSequenceAdapter(dataOrF, start!, end!)
  } else {
    if (start == undefined) start = 0
    if (end == undefined) end = dataOrF.length
    return new ArraySequenceAdapter(dataOrF, start, end)
  }
}

export { ISequence, createSequence }

if (require.main === module) {
  const seq = createSequence(i => i * i, 0, 10)
  console.log(seq.at(3))
  console.log(seq.at(10))
  console.log(seq.at(-1))
  const subSeq = seq.subsequence(2, 5)
  console.log(subSeq.length)
  console.log(subSeq.at(0))
  console.log(subSeq.at(1))
  console.log(subSeq.at(2))

  const seq2 = createSequence('hello')
  console.log(seq2.at(0))
  console.log(seq2.at(1))
  const subSeq2 = seq2.subsequence(1, 4)
  console.log(subSeq2.length)
  console.log(subSeq2.at(0))
}
