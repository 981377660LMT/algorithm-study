// DimensionIndex/IndexDimension
// dimension: 是一个数组，例如 [2, 3, 4]
// index: 是一个数字，例如 8

class DimensionIndex {
  private readonly _dimensions: number[]
  private readonly _size: number

  constructor(...dimension: number[]) {
    this._dimensions = dimension
    let size = 1
    for (let i = 0; i < dimension.length; i++) size *= dimension[i]
    this._size = size
  }

  index(...dimension: number[]): number {
    let res = dimension[0]
    for (let i = 1; i < dimension.length; i++) {
      res = res * this._dimensions[i] + dimension[i]
    }
    return res
  }

  dimension(index: number): number[] {
    const res = Array(this._dimensions.length).fill(0)
    for (let i = this._dimensions.length - 1; i >= 0; i--) {
      res[i] = index % this._dimensions[i]
      index = Math.floor(index / this._dimensions[i])
    }
    return res
  }

  isValidDimension(...dimension: number[]): boolean {
    for (let i = 0; i < dimension.length; i++) {
      if (dimension[i] < 0 || dimension[i] >= this._dimensions[i]) return false
    }
    return true
  }

  indexOfSpecifiedDimension(index: number, d: number): number {
    for (let i = this._dimensions.length - 1; i >= 0; i--) {
      if (i === d) return index % this._dimensions[i]
      index = Math.floor(index / this._dimensions[i])
    }
    throw new Error('Invalid dimension')
  }

  get size(): number {
    return this._size
  }
}

export { DimensionIndex }

if (require.main === module) {
  const di = new DimensionIndex(2, 3, 4)
  console.log(di.index(0, 2, 1))
  console.log(di.dimension(9))
  console.log(di.isValidDimension(1, 2, 1))
  console.log(di.isValidDimension(1, 2, 4))
  console.log(di.indexOfSpecifiedDimension(9, 0))
  console.log(di.indexOfSpecifiedDimension(9, 1))
  console.log(di.indexOfSpecifiedDimension(9, 2))
  console.log(di.size)
}
