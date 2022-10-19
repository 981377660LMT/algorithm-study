/* eslint-disable no-param-reassign */
/* eslint-disable no-console */
/* eslint-disable @typescript-eslint/no-non-null-assertion */

/**
 * A rope is an efficient data structure for storing and manipulating very large mutable strings
 * It reduces memory reallocation and data copy overhead for applications
 * that are constantly operating on very large strings by
 * splitting them into multiple smaller strings transparently.
 * Efficient random access is achieved via a binary tree.
 *
 * @see {@link https://github.com/component/rope}
 * @see {@link https://github.com/component/rope/issues/2}
 */
class Rope {
  /**
   * The threshold used to split a leaf node into two child nodes.
   */
  private static readonly _SPLIT_LENGTH = 1 << 15

  /**
   * The threshold used to join two child nodes into one leaf node.
   */
  private static readonly _JOIN_LENGTH = 1 << 14

  /**
   * The threshold used to trigger a tree node rebuild when rebalancing the rope.
   */
  private static readonly _REBALANCE_RATIO = 1 << 1

  length: number

  private _left: Rope | undefined

  private _right: Rope | undefined

  private _value: string | undefined

  /**
   * Creates a rope data structure
   *
   * @param str - String to populate the rope.
   * @complexity O(n)
   */
  constructor(str: string) {
    this._value = str
    this.length = str.length
    this._adjust()
  }

  /**
   * Inserts text into the rope on the specified position.
   *
   * @param  position - Where to insert the text
   * @param  value - Text to be inserted on the rope
   * @complexity O(k)
   */
  insert(position: number, value: string) {
    if (position < 0) position += this.length
    if (position < 0) {
      position = 0
    } else if (position > this.length) {
      position = this.length
    }

    if (typeof this._value !== 'undefined') {
      this._value = this._value.slice(0, position) + value.toString() + this._value.slice(position)
      this.length = this._value.length
    } else {
      const leftLength = this._left!.length
      if (position < leftLength) {
        this._left!.insert(position, value)
      } else {
        this._right!.insert(position - leftLength, value)
      }
      this.length = this._left!.length + this._right!.length
    }

    this._adjust()
  }

  /**
   * Removes text from the rope between the `start` and `end` positions.
   * The character at `start` gets removed, but the character at `end` is
   * not removed.
   *
   * @param  start - Initial position (inclusive)
   * @param  end - Final position (not-inclusive)
   * @complexity O(k)
   */
  remove(start: number, end: number): void {
    if (start < 0 || start > this.length) throw new RangeError('Start is not within rope bounds.')
    if (end < 0 || end > this.length) throw new RangeError('End is not within rope bounds.')
    if (start > end) throw new RangeError('Start is greater than end.')

    if (typeof this._value !== 'undefined') {
      this._value = this._value.substring(0, start) + this._value.substring(end)
      this.length = this._value.length
    } else {
      const leftLength = this._left!.length
      const leftStart = Math.min(start, leftLength)
      const leftEnd = Math.min(end, leftLength)
      const rightLength = this._right!.length
      const rightStart = Math.max(0, Math.min(start - leftLength, rightLength))
      const rightEnd = Math.max(0, Math.min(end - leftLength, rightLength))
      if (leftStart < leftLength) {
        this._left!.remove(leftStart, leftEnd)
      }
      if (rightEnd > 0) {
        this._right!.remove(rightStart, rightEnd)
      }
      this.length = this._left!.length + this._right!.length
    }

    this._adjust()
  }

  /**
   * Returns text from the rope between the `start` and `end` positions.
   * The character at `start` gets returned, but the character at `end` is
   * not returned.
   *
   * @param  start - Initial position (inclusive)
   * @param  end - Final position (not-inclusive)
   * @complexity O(k)
   */
  slice(start?: number, end?: number): string {
    if (typeof start === 'undefined') {
      start = 0
    }

    if (typeof end === 'undefined') {
      end = this.length
    }

    if (start < 0) {
      start = 0
    } else if (start > this.length) {
      start = this.length
    }

    if (end < 0) {
      end = 0
    } else if (end > this.length) {
      end = this.length
    }

    if (typeof this._value !== 'undefined') {
      return this._value.slice(start, end)
    }

    const leftLength = this._left!.length
    const leftStart = Math.min(start, leftLength)
    const leftEnd = Math.min(end, leftLength)
    const rightLength = this._right!.length
    const rightStart = Math.max(0, Math.min(start - leftLength, rightLength))
    const rightEnd = Math.max(0, Math.min(end - leftLength, rightLength))

    if (leftStart !== leftEnd) {
      if (rightStart !== rightEnd) {
        return this._left!.slice(leftStart, leftEnd) + this._right!.slice(rightStart, rightEnd)
      }
      return this._left!.slice(leftStart, leftEnd)
    }

    if (rightStart !== rightEnd) {
      return this._right!.slice(rightStart, rightEnd)
    }

    return ''
  }

  /**
   * Returns the character at `position`
   *
   * @param  position
   * @complexity O(logn)
   */
  at(position: number): string {
    return this.slice(position, position + 1)
  }

  /**
   * Converts the rope to a JavaScript String.
   */
  toString(): string {
    if (typeof this._value !== 'undefined') {
      return this._value
    }
    return this._left!.toString() + this._right!.toString()
  }

  /**
   * Rebuilds the entire rope structure, producing a balanced tree.
   */
  private _rebuild(): void {
    if (typeof this._value === 'undefined') {
      this._value = this._left!.toString() + this._right!.toString()
      this._left = undefined
      this._right = undefined
      this._adjust()
    }
  }

  /**
   * Finds unbalanced nodes in the tree and rebuilds them.
   */
  private _rebalance(): void {
    if (typeof this._value === 'undefined') {
      if (
        this._left!.length / this._right!.length > Rope._REBALANCE_RATIO ||
        this._right!.length / this._left!.length > Rope._REBALANCE_RATIO
      ) {
        this._rebuild()
      } else {
        this._left!._rebalance()
        this._right!._rebalance()
      }
    }
  }

  /**
   * Adjusts the tree structure, so that very long nodes are split
   * and short ones are joined
   */
  private _adjust(): void {
    if (typeof this._value !== 'undefined') {
      if (this.length > Rope._SPLIT_LENGTH) {
        const divide = Math.floor(this.length / 2)
        this._left = new Rope(this._value.slice(0, divide))
        this._right = new Rope(this._value.slice(divide))
        this._value = undefined
      }
    } else if (this.length < Rope._JOIN_LENGTH) {
      this._value = this._left!.toString() + this._right!.toString()
      this._left = undefined
      this._right = undefined
    }
  }
}

if (require.main === module) {
  console.time('Rope')
  const rope = new Rope('Hello World'.repeat(1e5))
  for (let _ = 0; _ < 1000; _++) {
    rope.insert(5, '!')
    rope.slice(0, 10)
  }
  console.timeEnd('Rope') // Rope: 21.849ms

  console.time('String')
  let string = 'Hello World'.repeat(1e5)
  for (let _ = 0; _ < 1000; _++) {
    string = `${string.slice(0, 5)}!${string.slice(5)}`
  }
  console.timeEnd('String') // String: 5.334s
}

export { Rope }
