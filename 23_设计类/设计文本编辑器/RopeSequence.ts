// https://github.com/marijnh/rope-sequence

const GOOD_LEAF_SIZE = 200

/**
 * A rope sequence is a persistent sequence data structure
 * that supports appending, prepending, and slicing without doing a
 * full copy. It is represented as a mostly-balanced tree.
 */
class RopeSequence<T> {
  /**
   * The length of the rope.
   */
  length: number

  /**
   * Append an array or other rope to this one, returning a new rope.
   */
  append(other: T[] | RopeSequence<T>): RopeSequence<T> {
    if (!other.length) return this
    other = RopeSequence.from(other)

    return (
      (!this.length && other) ||
      (other.length < GOOD_LEAF_SIZE && this.leafAppend(other)) ||
      (this.length < GOOD_LEAF_SIZE && other.leafPrepend(this)) ||
      this.appendInner(other)
    )
  }

  /**
   * Prepend an array or other rope to this one, returning a new rope.
   */
  prepend(other: T[] | RopeSequence<T>): RopeSequence<T> {
    if (!other.length) return this
    return RopeSequence.from(other).append(this)
  }

  /**
   * @internal
   * Inner implementation of append
   */
  protected appendInner(other: RopeSequence<T>): RopeSequence<T> {
    return new Append<T>(this, other)
  }

  /**
   * Create a rope repesenting a sub-sequence of this rope.
   */
  slice(from: number = 0, to: number = this.length): RopeSequence<T> {
    if (from >= to) return RopeSequence.empty as RopeSequence<T>
    return this.sliceInner(Math.max(0, from), Math.min(this.length, to))
  }

  /**
   * @internal
   * Inner implementation of slice
   */
  protected sliceInner(from: number, to: number): RopeSequence<T> {
    throw new Error('Abstract method sliceInner not implemented')
  }

  /**
   * Retrieve the element at the given position from this rope.
   */
  get(i: number): T | undefined {
    if (i < 0 || i >= this.length) return undefined
    return this.getInner(i)
  }

  /**
   * @internal
   * Inner implementation of get
   */
  protected getInner(i: number): T {
    throw new Error('Abstract method getInner not implemented')
  }

  /**
   * Call the given function for each element between the given
   * indices. This tends to be more efficient than looping over the
   * indices and calling `get`, because it doesn't have to descend the
   * tree for every element.
   */
  forEach(
    f: (element: T, index: number) => boolean | void,
    from: number = 0,
    to: number = this.length
  ): void {
    if (from <= to) this.forEachInner(f, from, to, 0)
    else this.forEachInvertedInner(f, from, to, 0)
  }

  /**
   * @internal
   * Inner implementation of forEach for forward traversal
   */
  protected forEachInner(
    f: (element: T, index: number) => boolean | void,
    from: number,
    to: number,
    start: number
  ): boolean | void {
    throw new Error('Abstract method forEachInner not implemented')
  }

  /**
   * @internal
   * Inner implementation of forEach for inverted traversal
   */
  protected forEachInvertedInner(
    f: (element: T, index: number) => boolean | void,
    from: number,
    to: number,
    start: number
  ): boolean | void {
    throw new Error('Abstract method forEachInvertedInner not implemented')
  }

  /**
   * Map the given functions over the elements of the rope, producing
   * a flat array.
   */
  map<U>(f: (element: T, index: number) => U, from: number = 0, to: number = this.length): U[] {
    let result: U[] = []
    this.forEach(
      (elt, i) => {
        result.push(f(elt, i))
      },
      from,
      to
    )
    return result
  }

  /**
   * @internal
   * Helper method for leaf append
   */
  protected leafAppend(other: RopeSequence<T>): RopeSequence<T> | undefined {
    return undefined
  }

  /**
   * @internal
   * Helper method for leaf prepend
   */
  protected leafPrepend(other: RopeSequence<T>): RopeSequence<T> | undefined {
    return undefined
  }

  /**
   * Return the content of this rope as an array.
   */
  flatten(): T[] {
    throw new Error('Abstract method flatten not implemented')
  }

  /**
   * The depth of the rope tree
   */
  get depth(): number {
    throw new Error('Abstract property depth not implemented')
  }

  /**
   * Create a rope representing the given array, or return the rope
   * itself if a rope was given.
   */
  static from<U>(values: U[] | RopeSequence<U>): RopeSequence<U> {
    if (values instanceof RopeSequence) return values
    return values && values.length ? new Leaf(values) : (RopeSequence.empty as RopeSequence<U>)
  }

  /**
   * The empty rope sequence.
   */
  static empty: RopeSequence<any>
}

class Leaf<T> extends RopeSequence<T> {
  private values: T[]

  constructor(values: T[]) {
    super()
    this.values = values
  }

  override flatten(): T[] {
    return this.values
  }

  protected override sliceInner(from: number, to: number): RopeSequence<T> {
    if (from == 0 && to == this.length) return this
    return new Leaf<T>(this.values.slice(from, to))
  }

  protected override getInner(i: number): T {
    return this.values[i]
  }

  protected override forEachInner(
    f: (element: T, index: number) => boolean | void,
    from: number,
    to: number,
    start: number
  ): boolean | void {
    for (let i = from; i < to; i++) if (f(this.values[i], start + i) === false) return false
    return undefined
  }

  protected override forEachInvertedInner(
    f: (element: T, index: number) => boolean | void,
    from: number,
    to: number,
    start: number
  ): boolean | void {
    for (let i = from - 1; i >= to; i--) if (f(this.values[i], start + i) === false) return false
    return undefined
  }

  protected override leafAppend(other: RopeSequence<T>): RopeSequence<T> | undefined {
    if (this.length + other.length <= GOOD_LEAF_SIZE)
      return new Leaf<T>(this.values.concat(other.flatten()))
    return undefined
  }

  protected override leafPrepend(other: RopeSequence<T>): RopeSequence<T> | undefined {
    if (this.length + other.length <= GOOD_LEAF_SIZE)
      return new Leaf<T>(other.flatten().concat(this.values))
    return undefined
  }

  override get length(): number {
    return this.values.length
  }

  override get depth(): number {
    return 0
  }
}

// Initialize the empty rope sequence
RopeSequence.empty = new Leaf<never>([])

class Append<T> extends RopeSequence<T> {
  private left: RopeSequence<T>
  private right: RopeSequence<T>
  override readonly length: number
  override readonly depth: number

  constructor(left: RopeSequence<T>, right: RopeSequence<T>) {
    super()
    this.left = left
    this.right = right
    this.length = left.length + right.length
    this.depth = Math.max(left.depth, right.depth) + 1
  }

  override flatten(): T[] {
    return this.left.flatten().concat(this.right.flatten())
  }

  protected override getInner(i: number): T {
    return i < this.left.length ? this.left.get(i)! : this.right.get(i - this.left.length)!
  }

  protected override forEachInner(
    f: (element: T, index: number) => boolean | void,
    from: number,
    to: number,
    start: number
  ): boolean | void {
    let leftLen = this.left.length
    if (from < leftLen && this.left.forEachInner(f, from, Math.min(to, leftLen), start) === false)
      return false
    if (
      to > leftLen &&
      this.right.forEachInner(
        f,
        Math.max(from - leftLen, 0),
        Math.min(this.length, to) - leftLen,
        start + leftLen
      ) === false
    )
      return false
    return undefined
  }

  protected override forEachInvertedInner(
    f: (element: T, index: number) => boolean | void,
    from: number,
    to: number,
    start: number
  ): boolean | void {
    let leftLen = this.left.length
    if (
      from > leftLen &&
      this.right.forEachInvertedInner(
        f,
        from - leftLen,
        Math.max(to, leftLen) - leftLen,
        start + leftLen
      ) === false
    )
      return false
    if (
      to < leftLen &&
      this.left.forEachInvertedInner(f, Math.min(from, leftLen), to, start) === false
    )
      return false
    return undefined
  }

  protected override sliceInner(from: number, to: number): RopeSequence<T> {
    if (from == 0 && to == this.length) return this
    let leftLen = this.left.length
    if (to <= leftLen) return this.left.slice(from, to)
    if (from >= leftLen) return this.right.slice(from - leftLen, to - leftLen)
    return this.left.slice(from, leftLen).append(this.right.slice(0, to - leftLen))
  }

  protected override leafAppend(other: RopeSequence<T>): RopeSequence<T> | undefined {
    let inner = this.right.leafAppend(other)
    if (inner) return new Append<T>(this.left, inner)
    return undefined
  }

  protected override leafPrepend(other: RopeSequence<T>): RopeSequence<T> | undefined {
    let inner = this.left.leafPrepend(other)
    if (inner) return new Append<T>(inner, this.right)
    return undefined
  }

  protected override appendInner(other: RopeSequence<T>): RopeSequence<T> {
    if (this.left.depth >= Math.max(this.right.depth, other.depth) + 1)
      return new Append<T>(this.left, new Append<T>(this.right, other))
    return new Append<T>(this, other)
  }
}

export { RopeSequence }
