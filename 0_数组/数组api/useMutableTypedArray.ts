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

export { useMutableTypedArray }

if (require.main === module) {
  const nums = useMutableTypedArray('INT32', { arrayLike: [1, 2, 3, 4, 5] })
  assert.strictEqual(nums.length, 5)
  assert.strictEqual(nums.at(0), 1)
  assert.strictEqual(nums.at(1), 2)
  assert.strictEqual(nums.popleft(), 1)
  assert.strictEqual(nums.toString(), '2,3,4,5')
  assert.strictEqual(nums.pop(), 5)
  assert.strictEqual(nums.toString(), '2,3,4')
  console.log(nums.toString(), nums.length)
  nums.append(6)
  nums.appendleft(1)
  assert.strictEqual(nums.toString(), '1,2,3,4,6')
  nums.insert(2, 3)
  assert.strictEqual(nums.toString(), '1,2,3,3,4,6')
  nums.set(2, 2)
  assert.strictEqual(nums.toString(), '1,2,2,3,4,6')
  nums.append(7)
  assert.strictEqual(nums.toString(), '1,2,2,3,4,6,7')

  console.time('useMutableTypedArray')
  const nums2 = useMutableTypedArray('UINT32', {
    // initialCapacity: 2e5,
    arrayLike: Array(1e5).map((_, i) => i)
  })

  for (let i = 0; i < 5e4; i++) {
    const pos = Math.random() * nums2.length
    nums2.append(i)
    nums2.appendleft(i)
    nums2.pop(pos)
    nums2.pop(pos)
    nums2.at(pos)
    nums2.insert(pos, i)
    nums2.set(pos, i)
  }

  console.timeEnd('useMutableTypedArray')
}
