// 三维树状数组
// 单点修改，区间查询

class BIT3D {
  private readonly _x: number
  private readonly _y: number
  private readonly _z: number
  private readonly _hash1: number
  private readonly _hash2: number
  private readonly _data: Float64Array

  constructor(x: number, y: number, z: number) {
    x++
    y++
    z++
    this._x = x
    this._y = y
    this._z = z
    this._hash1 = y * z
    this._hash2 = z
    this._data = new Float64Array(x * y * z)
  }

  /**
   * 0<=x<X, 0<=y<Y, 0<=z<Z.
   */
  add(x: number, y: number, z: number, v: number): void {
    for (let i = x; i < this._x; i |= i + 1) {
      for (let j = y; j < this._y; j |= j + 1) {
        for (let k = z; k < this._z; k |= k + 1) {
          this._data[this._index(i, j, k)] += v
        }
      }
    }
  }

  /**
   * 0<=x<X, 0<=y<Y, 0<=z<Z.
   */
  queryPrefix(x: number, y: number, z: number): number {
    x--
    y--
    z--
    let res = 0
    for (let i = x; i >= 0; i = (i & (i + 1)) - 1) {
      for (let j = y; j >= 0; j = (j & (j + 1)) - 1) {
        for (let k = z; k >= 0; k = (k & (k + 1)) - 1) {
          res += this._data[this._index(i, j, k)]
        }
      }
    }
    return res
  }

  /**
   * 0<=x1<=x2<X, 0<=y1<=y2<Y, 0<=z1<=z2<Z.
   */
  queryRange(x1: number, y1: number, z1: number, x2: number, y2: number, z2: number): number {
    return (
      this.queryPrefix(x2, y2, z2) -
      this.queryPrefix(x1, y2, z2) -
      this.queryPrefix(x2, y1, z2) -
      this.queryPrefix(x2, y2, z1) +
      this.queryPrefix(x1, y1, z2) +
      this.queryPrefix(x1, y2, z1) +
      this.queryPrefix(x2, y1, z1) -
      this.queryPrefix(x1, y1, z1)
    )
  }

  private _index(x: number, y: number, z: number): number {
    return x * this._hash1 + y * this._hash2 + z
  }
}

export { BIT3D }

if (require.main === module) {
  const bit3d = new BIT3D(10, 10, 10)
  bit3d.add(1, 1, 1, 1)
  bit3d.add(1, 2, 1, 7)
  console.log(bit3d.queryRange(0, 0, 0, 10, 10, 10))
}
