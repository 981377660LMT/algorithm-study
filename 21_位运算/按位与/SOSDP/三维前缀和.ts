/* eslint-disable new-cap */
/* eslint-disable func-call-spacing */
/* eslint-disable no-spaced-func */
/* eslint-disable no-shadow */
/* eslint-disable implicit-arrow-linebreak */

import assert from 'assert'

type Matrix3D = ArrayLike<ArrayLike<ArrayLike<number>>>
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

/**
 * 三维前缀和
 */
class MatrixPreSum3D {
  private readonly _preSum: Matrix3D
  private readonly _xSize: number
  private readonly _ySize: number
  private readonly _zSize: number

  constructor(matrix: Matrix3D, dataType?: keyof typeof ARRAYTYPE_RECORD) {
    const xSize = matrix.length
    const ySize = matrix[0].length
    const zSize = matrix[0][0].length

    const arrayType = dataType ? ARRAYTYPE_RECORD[dataType] : Array
    const preSum = Array(xSize + 1)
    for (let i = 0; i < xSize + 1; i++) {
      preSum[i] = Array(ySize + 1)
      for (let j = 0; j < ySize + 1; j++) {
        preSum[i][j] = new arrayType(zSize + 1).fill(0)
      }
    }

    for (let x = 1; x < xSize + 1; x++) {
      for (let y = 1; y < ySize + 1; y++) {
        for (let z = 1; z < zSize + 1; z++) {
          preSum[x][y][z] = matrix[x - 1][y - 1][z - 1]
        }
      }
    }

    for (let x = 1; x < xSize + 1; x++) {
      for (let y = 1; y < ySize + 1; y++) {
        for (let z = 1; z < zSize + 1; z++) {
          preSum[x][y][z] += preSum[x - 1][y][z]
        }
      }
    }

    for (let x = 1; x < xSize + 1; x++) {
      for (let y = 1; y < ySize + 1; y++) {
        for (let z = 1; z < zSize + 1; z++) {
          preSum[x][y][z] += preSum[x][y - 1][z]
        }
      }
    }

    for (let x = 1; x < xSize + 1; x++) {
      for (let y = 1; y < ySize + 1; y++) {
        for (let z = 1; z < zSize + 1; z++) {
          preSum[x][y][z] += preSum[x][y][z - 1]
        }
      }
    }

    this._preSum = preSum
    this._xSize = xSize
    this._ySize = ySize
    this._zSize = zSize
  }

  /**
   * 查询 sum(A[x1:x2+1][y1:y2+1][z1:z2+1]) 的值
   *
   * @param x1 0 <= x1 < x2 < {@link _xSize}
   * @param y1 0 <= y1 < y2 < {@link _ySize}
   * @param z1 0 <= z1 < z2 < {@link _zSize}
   * @param x2 0 <= x1 < x2 < {@link _xSize}
   * @param y2 0 <= y1 < y2 < {@link _ySize}
   * @param z2 0 <= z1 < z2 < {@link _zSize}
   *
   * @example
   * ```ts
   * matrixPreSum3D.query(0, 0, 0, 1, 1, 1) // 查询A[0:2][0:2][0:2]的和
   * ```
   */
  query(x1: number, y1: number, z1: number, x2: number, y2: number, z2: number): number {
    return (
      this._preSum[x2 + 1][y2 + 1][z2 + 1] -
      this._preSum[x1][y2 + 1][z2 + 1] -
      this._preSum[x2 + 1][y1][z2 + 1] -
      this._preSum[x2 + 1][y2 + 1][z1] +
      this._preSum[x1][y1][z2 + 1] +
      this._preSum[x1][y2 + 1][z1] +
      this._preSum[x2 + 1][y1][z1] -
      this._preSum[x1][y1][z1]
    )
  }
}

if (require.main === module) {
  const matrix3d = [
    [
      [1, 2, 3],
      [4, 5, 6],
      [7, 8, 9]
    ],

    [
      [10, 11, 12],
      [13, 14, 15],
      [16, 17, 18]
    ],
    [
      [19, 20, 21],
      [22, 23, 24],
      [25, 26, 27]
    ]
  ]

  const preSum3d = new MatrixPreSum3D(matrix3d)
  assert.strictEqual(preSum3d.query(1, 1, 1, 1, 1, 1), 14)
  assert.strictEqual(preSum3d.query(0, 0, 0, 1, 1, 1), 60)
  assert.strictEqual(
    preSum3d.query(0, 0, 0, 2, 2, 2),
    matrix3d.reduce(
      (pre, cur) => pre + cur.reduce((pre, cur) => pre + cur.reduce((pre, cur) => pre + cur, 0), 0),
      0
    )
  )
  console.log('test ok')

  // 512*512*512
  // const matrix3d2 = Array.from({ length: 512 }, () =>
  //   Array.from({ length: 512 }, () => new Uint8Array(512).fill(1))
  // )
  // console.time('preSum3d')
  // for (let i = 1; i < 512; i++) {
  //   for (let j = 1; j < 512; j++) {
  //     for (let k = 1; k < 512; k++) {
  //       matrix3d2[i][j][k] += matrix3d2[i - 1][j][k]
  //     }
  //   }
  // }
  // console.timeEnd('preSum3d') // preSum3d: 409.034ms
}

export { MatrixPreSum3D }
