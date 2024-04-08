abstract class Block<S, U> {
  abstract fullyUpdate(update: U): void
  abstract partialUpdate(index: number, update: U): void
  abstract fullyQuery(sum: S): void
  abstract partialQuery(index: number, sum: S): void

  beforePartialUpdate() {}
  afterPartialUpdate() {}
  beforePartialQuery() {}
  afterPartialQuery() {}
}

class BlockSplit<S, U> {
  private readonly _b: number
  private readonly _n: number
  private readonly _blocks: Block<S, U>[]

  constructor(n: number, blockSize: number, supplier: (left: number, right: number) => Block<S, U>) {
    this._b = blockSize
    this._n = n
    this._blocks = Array(Math.ceil(n / blockSize))
    for (let i = 0; i < this._blocks.length; i++) {
      const l = i * blockSize
      let r = l + blockSize - 1
      r = Math.min(r, n - 1)
      this._blocks[i] = supplier(l, r)
    }
  }

  updateRange(left: number, right: number, update: U): void {
    for (let i = 0; i < this._blocks.length; i++) {
      const l = i * this._b
      let r = l + this._b - 1
      r = Math.min(r, this._n - 1)
      if (this._leave(left, right, l, r)) {
        continue
      } else if (this._enter(left, right, l, r)) {
        this._blocks[i].fullyUpdate(update)
      } else {
        this._blocks[i].beforePartialUpdate()
        for (let j = Math.max(l, left), to = Math.min(r, right); j <= to; j++) {
          this._blocks[i].partialUpdate(j - l, update)
        }
        this._blocks[i].afterPartialUpdate()
      }
    }
  }

  updatePoint(index: number, update: U): void {
    const blockId = Math.floor(index / this._b)
    this._blocks[blockId].partialUpdate(index - blockId * this._b, update)
    this._blocks[blockId].afterPartialUpdate()
  }

  queryRange(left: number, right: number, sum: S): void {
    for (let i = 0; i < this._blocks.length; i++) {
      const l = i * this._b
      let r = l + this._b - 1
      r = Math.min(r, this._n - 1)
      if (this._leave(left, right, l, r)) {
        continue
      } else if (this._enter(left, right, l, r)) {
        this._blocks[i].fullyQuery(sum)
      } else {
        this._blocks[i].beforePartialQuery()
        for (let j = Math.max(l, left), to = Math.min(r, right); j <= to; j++) {
          this._blocks[i].partialQuery(j - l, sum)
        }
        this._blocks[i].afterPartialQuery()
      }
    }
  }

  private _enter(L: number, R: number, l: number, r: number): boolean {
    return L <= l && R >= r
  }

  private _leave(L: number, R: number, l: number, r: number): boolean {
    return L > r || R < l
  }
}

export {}
