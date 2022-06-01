// https://github.com/chrisakroyd/bit-vec/blob/main/src/index.js

function countBits(count) {
  let n = count
  n = n - ((n >> 1) & 0x55555555)
  n = (n & 0x33333333) + ((n >> 2) & 0x33333333)
  return (((n + (n >> 4)) & 0xf0f0f0f) * 0x1010101) >> 24
}

class BitVector {
  /**
   * BitVector constructor.
   *
   * @constructor
   * @param {Number} size -> Size of the array in bits.
   */
  constructor(size) {
    this.array = new Uint8Array(Math.ceil(size / 8))
    this.bitsPerElement = this.array.BYTES_PER_ELEMENT * 8
  }

  /**
   *  .bits() + .length() are semi-dynamic properties that may change frequently,
   *  therefore these are computed on the fly via getters.
   */
  get bits() {
    return this.bitsPerElement * this.array.length
  }

  get length() {
    return this.array.length
  }

  get bitVector() {
    return this.array
  }

  set bitVector(bitArray) {
    this.array = bitArray
  }

  /**
   * Clears the bit at the given index.
   *
   * @param {Number} index -> Number for index: 0 <= index < bitVec.bits.
   * @throws {RangeError} Throws range error if index is out of range.
   */
  rangeCheck(index) {
    if (!(index < this.bits) || index < 0) {
      throw new RangeError(`Given index ${index} out of range of bit vector length ${this.bits}`)
    }
  }

  /**
   * `bitVec.get(index)`
   * Performs a get operation on the given index, retrieving the stored value (0 or 1).
   *
   * @param {Number} index -> Number for index: 0 <= index < bitVec.bits.
   * @return {Number} Returns number, 1 if set, 0 otherwise.
   */
  get(index) {
    this.rangeCheck(index)
    const byteIndex = Math.floor(index / this.bitsPerElement)
    const bitIndex = index % this.bitsPerElement

    return (this.array[byteIndex] & (1 << bitIndex)) > 0 ? 1 : 0
  }

  /**
   * `bitVec.set(index)`
   * Performs a set operation on the given index, setting the value to either 0 or 1.
   *
   * @param {Number} index -> Number for index: 0 <= index < bitVec.bits.
   * @param {Number} value -> Number, 0 or 1, defaults to 1.
   * @return {BitVector} Returns `BitVector` for chaining with the bit cleared.
   */
  set(index, value = 1) {
    this.rangeCheck(index)
    const byteIndex = Math.floor(index / this.bitsPerElement)
    const bitIndex = index % this.bitsPerElement

    if (value) {
      this.array[byteIndex] |= 1 << bitIndex
    } else {
      this.array[byteIndex] &= ~(1 << bitIndex)
    }

    return this
  }

  /**
   * `bitVec.clear(index)`
   * Clears the bit at the given index.
   *
   * @param {Number} index -> Number for index: 0 <= index < bitVec.bits.
   * @return {BitVector} Returns `BitVector` for chaining with the bit cleared.
   */
  clear(index) {
    return this.set(index, 0)
  }

  /**
   * `bitVec.flip(index)`
   * Flips the bit at the given index.
   *
   * @param {Number} index -> Number for index: 0 <= index < bitVec.bits.
   * @return {BitVector} Returns `BitVector` for chaining with the bit cleared.
   */
  flip(index) {
    this.rangeCheck(index)
    const byteIndex = Math.floor(index / this.bitsPerElement)
    const bitIndex = index % this.bitsPerElement
    this.array[byteIndex] ^= 1 << bitIndex
    return this
  }

  /**
   * `bitVec.test(index)`
   * Tests whether the given index is set to 1.
   *
   * @param {Number} index -> Number for index: 0 <= index < bitVec.bits.
   * @return {Boolean} Returns Boolean `true` if index is set, `false` otherwise .
   */
  test(index) {
    return this.get(index) === 1
  }

  /**
   * `bitVec.count()`
   * Counts the number of set bits in the bit vector.
   *
   * @return {Number} Number of indices currently set to 1.
   */
  count() {
    let c = 0
    for (let i = 0; i < this.array.length; i += 1) {
      c += countBits(this.array[i])
    }
    return c
  }

  /**
   * `bitVec.setRange(begin, end, value = 1)`
   * Sets a range of bits from begin to end.
   *
   * @param {Number} begin -> Number for index: 0 <= index < bitVec.bits.
   * @param {Number} end -> Number for index: 0 <= index < bitVec.bits.
   * @param {Number} value -> The value to set the index to (0 or 1).
   * @return {BitVector} Returns `BitVector` for chaining with the bits set.
   */
  setRange(begin, end, value = 1) {
    for (let i = begin; i < end; i += 1) {
      this.set(i, value)
    }
    return this
  }

  /**
   * `bitVec.clearRange(begin, end)`
   * Clears a range of bits from begin to end.
   *
   * @param {Number} begin -> Number for index: 0 <= index < bitVec.bits.
   * @param {Number} end -> Number for index: 0 <= index < bitVec.bits.
   * @return {BitVector} Returns `BitVector` for chaining with the bits set.
   */
  clearRange(begin, end) {
    this.setRange(begin, end, 0)
    return this
  }

  /**
   * `bitVec.shortLong(bitVec)`
   *  Useful function allowing for the comparison of two differently sized BitVector's.
   *  Simply returns the short and long arrays.
   *
   * @param {BitVector} bitVec -> BitVector, instance of BitVector class.
   * @return {Object} Returns object with two keys of type `BitVector`,
   *                  short = shorter bit vector, long = longer bit vector.
   */
  shortLong(bitVec) {
    let short
    let long

    if (bitVec.length < this.length) {
      short = bitVec.array
      long = this.array
    } else {
      short = this.array
      long = bitVec.array
    }

    return { short, long }
  }

  /**
   * `bitVec.or(bitVec)`
   * Performs the bitwise or operation between two BitVectors and returns the result as a
   * new BitVector object.
   *
   * @param {BitVector} bitVec -> BitVector, instance of BitVector class.
   * @return {BitVector} Returns new `BitVector` object with the result of the operation.
   */
  or(bitVec) {
    // Get short and long arrays, assign correct variables -> for ops between two diff sized arrays.
    const { short, long } = this.shortLong(bitVec)
    const array = new Uint8Array(long.length)

    // Perform operation over shorter array.
    for (let i = 0; i < short.length; i += 1) {
      array[i] = short[i] | long[i]
    }

    // Fill in the remaining unchanged numbers from the longer array.
    for (let j = short.length; j < long.length; j += 1) {
      array[j] = long[j]
    }

    // Return a new BitVector object.
    return BitVector.fromArray(array)
  }

  /**
   * `bitVec.xor(bitVec)`
   * Performs the bitwise xor operation between two BitVectors and returns the result as a
   * new BitVector object.
   *
   * @param {BitVector} bitVec -> BitVector, instance of BitVector class.
   * @return {BitVector} Returns new `BitVector` object with the result of the operation.
   */
  xor(bitVec) {
    // Get short and long arrays, assign correct variables -> for ops between two diff sized arrays.
    const { short, long } = this.shortLong(bitVec)
    const array = new Uint8Array(long.length)

    // Perform operation over shorter array.
    for (let i = 0; i < short.length; i += 1) {
      array[i] = short[i] ^ long[i]
    }

    // Fill in the remaining numbers from the longer array.
    for (let j = short.length; j < long.length; j += 1) {
      array[j] = 0 ^ long[j]
    }

    // Return a new BitVector object.
    return BitVector.fromArray(array)
  }

  /**
   * `bitVec.and(bitVec)`
   * Performs the bitwise and operation between two BitVectors and returns the result as a
   * new BitVector object.
   *
   * @param {BitVector} bitVec -> BitVector, instance of BitVector class.
   * @return {BitVector} Returns new `BitVector` object with the result of the operation.
   */
  and(bitVec) {
    // Get short and long arrays, assign correct variables -> for ops between two diff sized arrays.
    const { short, long } = this.shortLong(bitVec)
    const array = new Uint8Array(long.length)

    // Perform operation over shorter array.
    for (let i = 0; i < short.length; i += 1) {
      array[i] = short[i] & long[i]
    }

    // Fill in the remaining unchanged numbers from the longer array.
    for (let j = short.length; j < long.length; j += 1) {
      array[j] = long[j]
    }

    // Return a new BitVector object.
    return BitVector.fromArray(array)
  }

  /**
   * `bitVec.equals(otherBitVec)`
   * Determines if two bit vectors are equal.
   *
   * @param {BitVector} bitVec -> BitVector, instance of BitVector class.
   * @return {Boolean} Returns Boolean `true` if the two bit vectors are equal, `false` otherwise.
   */
  equals(bitVec) {
    const { short, long } = this.shortLong(bitVec)

    for (let i = 0; i < short.length; i += 1) {
      if (short[i] !== long[i]) {
        return false
      }
    }

    // If the longer array is all 0 then they are equal, if not then they are not.
    // equiv to padding shorter bit array to larger array length and comparing.
    // Allows comparisons along vecs of different length.
    for (let j = short.length; j < long.length; j += 1) {
      if (long[j] !== 0) {
        return false
      }
    }

    return true
  }

  /**
   * `bitVec.notEquals(otherBitVec)`
   * Determines if two bit vectors are not equal.
   *
   * @param {BitVector} bitVec -> BitVector, instance of BitVector class.
   * @return {Boolean} Returns Boolean `true` if the two bit vectors are not equal,
   *                   `false` otherwise.
   */
  notEquals(bitVec) {
    return !this.equals(bitVec)
  }

  /**
   * `bitVec.not()`
   * Performs the bitwise not operation on this BitVector and returns the result as a
   * new BitVector object.
   *
   * @return {BitVector} Returns new `BitVector` object with the result of the operation.
   */
  not() {
    const array = new Uint8Array(this.array.length)

    for (let i = 0; i < this.array.length; i += 1) {
      array[i] = ~this.array[i]
    }

    return BitVector.fromArray(array)
  }

  /**
   * `bitVec.invert()`
   *
   * Inverts this BitVector, alias of .not().
   *
   * @return {BitVector} Returns new `BitVector` object with the result of the operation.
   */
  invert() {
    this.array = this.not().array
    return this
  }

  /**
   * `bitVec.orEqual(bitVec)`
   * Performs the bitwise or operation between two BitVectors and assigns the result to
   * this BitVector.
   *
   * @param {BitVector} bitVec -> BitVector, instance of BitVector class.
   * @return {BitVector} Returns `this` for chaining with the bits set.
   */
  orEqual(bitVec) {
    this.array = this.or(bitVec).array
    return this
  }

  /**
   * `bitVec.xorEqual(bitVec)`
   * Performs the bitwise xor operation between two BitVectors and assigns the result to
   * this BitVector.
   *
   * @param {BitVector} bitVec -> BitVector, instance of BitVector class.
   * @return {BitVector} Returns `this` for chaining with the bits set.
   */
  xorEqual(bitVec) {
    this.array = this.xor(bitVec).array
    return this
  }

  /**
   * `bitVec.andEqual(bitVec)`
   * Performs the bitwise and operation between two BitVectors and assigns the result to
   * this BitVector.
   *
   * @param {BitVector} bitVec -> BitVector, instance of BitVector class.
   * @return {BitVector} Returns `this` for chaining with the bits set.
   */
  andEqual(bitVec) {
    this.array = this.and(bitVec).array
    return this
  }

  /**
   * `bitVec.notEqual(bitVec)`
   * Performs the bitwise not operation between two BitVectors and assigns the result to
   * this BitVector.
   *
   * @return {BitVector} Returns `this` for chaining with the bits set.
   */
  notEqual() {
    this.array = this.not().array
    return this
  }

  /**
   * `bitVec.isEmpty()`
   * Tests whether this BitVector has any set bits.
   *
   * @return {Boolean} Returns Boolean `true` if the bit vector has no set bits, `false` otherwise.
   */
  isEmpty() {
    for (let i = 0; i < this.array.length; i += 1) {
      if (this.array[i] !== 0) {
        return false
      }
    }
    return true
  }

  toArray() {
    return this.array
  }

  static fromArray(bitVec) {
    const newBitVec = new BitVector(0)
    newBitVec.bitVector = bitVec
    return newBitVec
  }
}

export default BitVector
