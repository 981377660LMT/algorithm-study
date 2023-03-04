/* eslint-disable no-shadow */
/* eslint-disable no-inner-declarations */
/* eslint-disable no-console */
/* eslint-disable prefer-destructuring */
/* eslint-disable no-param-reassign */
/* eslint-disable new-cap */

import assert from 'assert'

const ARRAYTYPE_RECORD = {
  INT8: Int8Array,
  UIN8: Uint8Array,
  INT16: Int16Array,
  UINT16: Uint16Array,
  INT32: Int32Array,
  UINT32: Uint32Array,
  FLOAT32: Float32Array,
  FLOAT64: Float64Array
}

type DataType = keyof typeof ARRAYTYPE_RECORD

interface Options {
  initialCapacity?: number
  arrayLike?: ArrayLike<number>
}

/**
 * 支持增删改查的typedArray
 *
 * @param dataType {@link DataType}
 * @param options {@link Options}
 */
function useMutableTypedArray(dataType: DataType, options?: Options) {
  const { initialCapacity = 1 << 4, arrayLike = [] } = options || {}

  const arrayType = ARRAYTYPE_RECORD[dataType]
  let _elementData = new arrayType(initialCapacity)
  _ensureCapacity(arrayLike.length)
  _elementData.set(arrayLike)
  let _length = arrayLike.length

  function at(index: number): number {
    const newIndex = _normalizeIndex(index)
    _rangeCheck(newIndex)
    return _elementData[newIndex]
  }

  function set(index: number, value: number): void {
    const newIndex = _normalizeIndex(index)
    _rangeCheck(newIndex)
    _elementData[newIndex] = value
  }

  function pop(index = -1): number {
    const newIndex = _normalizeIndex(index)
    _rangeCheck(newIndex)
    const popped = _elementData[newIndex]
    _elementData.copyWithin(newIndex, newIndex + 1, _length)
    _length--
    return popped
  }

  function popleft(): number {
    return pop(0)
  }

  function insert(index: number, value: number): void {
    const newIndex = _normalizeIndex(index)
    _rangeCheckForAdd(newIndex)
    _ensureCapacity(_length + 1)
    _elementData.copyWithin(newIndex + 1, newIndex, _length)
    _elementData[newIndex] = value
    _length++
  }

  function append(value: number): void {
    insert(_length, value)
  }

  function appendleft(value: number): void {
    insert(0, value)
  }

  function slice(start: number, end?: number): InstanceType<typeof ARRAYTYPE_RECORD[DataType]> {
    return _elementData.slice(start, end)
  }

  function subarray(start: number, end?: number): InstanceType<typeof ARRAYTYPE_RECORD[DataType]> {
    return _elementData.subarray(start, end)
  }

  return {
    at,
    set,
    pop,
    insert,
    append,
    popleft,
    appendleft,
    slice,
    subarray,
    get length(): number {
      return _length
    },
    toString(): string {
      return _elementData.subarray(0, _length).toString()
    }
  }

  function _normalizeIndex(index: number): number {
    if (index < 0) {
      index += _length
    }

    return index
  }

  function _rangeCheck(index: number): void {
    if (index < 0 || index >= _length) {
      throw new RangeError(`index ${index} is out of range ${[0, _length - 1]}`)
    }
  }

  function _rangeCheckForAdd(addedIndex: number): void {
    if (addedIndex < 0 || addedIndex > _length) {
      throw new RangeError(`added index ${addedIndex} is out of range ${[0, _length]}`)
    }
  }

  function _ensureCapacity(minCapacity: number): void {
    if (minCapacity > _elementData.length) {
      _grow(minCapacity)
    }
  }

  function _grow(minCapacity: number): void {
    let newLength = _elementData.length << 1
    if (newLength < minCapacity) {
      newLength = minCapacity
    }

    const newElementData = new arrayType(newLength)
    newElementData.set(_elementData)
    _elementData = newElementData
  }
}

/**
 * typedArray实现的SortedList
 *
 * @param dataType {@link DataType}
 * @param options {@link Options}
 */
function useSortedTypedArray(dataType: DataType, options?: Options) {
  const { initialCapacity = 1 << 4, arrayLike = [] } = options || {}

  const arrayType = ARRAYTYPE_RECORD[dataType]
  let _elementData = new arrayType(initialCapacity)
  _ensureCapacity(arrayLike.length)
  let _length = 0
  for (let i = 0; i < arrayLike.length; i++) {
    add(arrayLike[i])
  }

  function at(index: number): number {
    const newIndex = _normalizeIndex(index)
    _rangeCheck(newIndex)
    return _elementData[newIndex]
  }

  function pop(index = -1): number {
    const newIndex = _normalizeIndex(index)
    _rangeCheck(newIndex)
    const popped = _elementData[newIndex]
    _elementData.copyWithin(newIndex, newIndex + 1, _length)
    _length--
    return popped
  }

  function add(value: number): void {
    const pos = bisectLeft(value)
    _ensureCapacity(_length + 1)
    _elementData.copyWithin(pos + 1, pos, _length)
    _elementData[pos] = value
    _length++
  }

  function bisectLeft(value: number): number {
    let left = 0
    let right = _length - 1
    while (left <= right) {
      const mid = Math.floor((left + right) / 2)
      const midElement = _elementData[mid]
      if (midElement < value) {
        left = mid + 1
      } else {
        right = mid - 1
      }
    }

    return left
  }

  function bisectRight(value: number): number {
    let left = 0
    let right = _length - 1
    while (left <= right) {
      const mid = Math.floor((left + right) / 2)
      const midElement = _elementData[mid]
      if (midElement <= value) {
        left = mid + 1
      } else {
        right = mid - 1
      }
    }

    return left
  }

  function slice(start: number, end?: number): InstanceType<typeof ARRAYTYPE_RECORD[DataType]> {
    return _elementData.slice(start, end)
  }

  function subarray(start: number, end?: number): InstanceType<typeof ARRAYTYPE_RECORD[DataType]> {
    return _elementData.subarray(start, end)
  }

  return {
    at,
    add,
    pop,
    bisectLeft,
    bisectRight,
    slice,
    subarray,
    get length(): number {
      return _length
    },
    toString(): string {
      return _elementData.subarray(0, _length).toString()
    }
  }

  function _normalizeIndex(index: number): number {
    if (index < 0) {
      index += _length
    }

    return index
  }

  function _rangeCheck(index: number): void {
    if (index < 0 || index >= _length) {
      throw new RangeError(`index ${index} is out of range ${[0, _length - 1]}`)
    }
  }

  function _ensureCapacity(minCapacity: number): void {
    if (minCapacity > _elementData.length) {
      _grow(minCapacity)
    }
  }

  function _grow(minCapacity: number): void {
    let newLength = _elementData.length << 1
    if (newLength < minCapacity) {
      newLength = minCapacity
    }

    const newElementData = new arrayType(newLength)
    newElementData.set(_elementData)
    _elementData = newElementData
  }
}

/**
 * 数组实现的SortedList
 *
 * @param iterable {@link Iterable}
 * @param compareFn SortedList的比较函数
 */
function useSortedList<E>(
  iterable: Iterable<E> = [],
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  compareFn: (a: E, b: E) => number = (a: any, b: any) => a - b
) {
  const _elementData: E[] = []
  for (const item of iterable) {
    add(item)
  }

  function at(index: number): E {
    const newIndex = _normalizeIndex(index)
    _rangeCheck(newIndex)
    return _elementData[newIndex]
  }

  function pop(index = -1): E {
    const newIndex = _normalizeIndex(index)
    _rangeCheck(newIndex)
    const popped = _elementData[newIndex]
    _elementData.splice(newIndex, 1)
    return popped
  }

  function add(value: E): void {
    const pos = bisectLeft(value)
    _elementData.splice(pos, 0, value)
    _elementData[pos] = value
  }

  function bisectLeft(value: E): number {
    let left = 0
    let right = _elementData.length - 1
    while (left <= right) {
      const mid = Math.floor((left + right) / 2)
      const midElement = _elementData[mid]
      if (compareFn(midElement, value) < 0) {
        left = mid + 1
      } else {
        right = mid - 1
      }
    }

    return left
  }

  function bisectRight(value: E): number {
    let left = 0
    let right = _elementData.length - 1
    while (left <= right) {
      const mid = Math.floor((left + right) / 2)
      const midElement = _elementData[mid]
      if (compareFn(midElement, value) <= 0) {
        left = mid + 1
      } else {
        right = mid - 1
      }
    }

    return left
  }

  return {
    at,
    add,
    pop,
    bisectLeft,
    bisectRight,
    get length(): number {
      return _elementData.length
    },
    toString(): string {
      return _elementData.toString()
    }
  }

  function _normalizeIndex(index: number): number {
    if (index < 0) {
      index += _elementData.length
    }

    return index
  }

  function _rangeCheck(index: number): void {
    if (index < 0 || index >= _elementData.length) {
      throw new RangeError(`index ${index} is out of range ${[0, _elementData.length - 1]}`)
    }
  }
}

export { useMutableTypedArray, useSortedTypedArray, useSortedList }

if (require.main === module) {
  const arr = Array.from({ length: 2e5 }, () => ~~(Math.random() * 100000))
  const sl2 = useSortedTypedArray('UINT32')
  console.time('add')
  for (const x of arr) {
    sl2.add(x)
  }
  console.timeEnd('add')
}

if (require.main === module) {
  // 2071. 你可以安排的最多任务数目
  function maxTaskAssign(
    tasks: number[],
    workers: number[],
    pills: number,
    strength: number
  ): number {
    tasks.sort((a, b) => a - b)
    workers.sort((a, b) => a - b)
    let left = 0
    let right = Math.min(tasks.length, workers.length)
    while (left <= right) {
      const mid = (left + right) >> 1
      if (check(mid)) {
        left = mid + 1
      } else {
        right = mid - 1
      }
    }

    return right

    function check(mid: number): boolean {
      let remain = pills
      const sl = useSortedTypedArray('INT32', { arrayLike: workers.slice(-mid) })
      // const wls = useSortedList(workers.slice(-mid))
      for (let i = mid - 1; i >= 0; i--) {
        const t = tasks[i]
        if (sl.at(sl.length - 1) >= t) {
          sl.pop()
        } else {
          if (remain === 0) {
            return false
          }
          const cand = sl.bisectLeft(t - strength)
          if (cand === sl.length) {
            return false
          }
          remain -= 1
          sl.pop(cand)
        }
      }

      return true
    }
  }

  assert.strictEqual(maxTaskAssign([3, 2, 1], [0, 3, 3], 1, 1), 3)
}
