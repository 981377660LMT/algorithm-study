/**
 * !01矩阵(方阵)乘法 C = A * B.
 * 输入矩阵中的元素只包含0或1,然后进行正常的矩阵乘法。
 * 也就是说输出矩阵元素值的范围为0~n.
 * 一个直观意义是, C[i][j]代表A的第i行和B的第j列的`公共元素个数`.
 * @param mat1 方阵A.边长n<=3000.
 * @param mat2 方阵B.边长n<=3000.
 * @returns 方阵C.每个元素的值范围为0~n(在uint16范围内).
 * @complexity O(n^3/32)
 * `1000*1000 => 100ms`.
 * `2000*2000 => 670ms`.
 * `3000*3000 => 2.2s`.
 * `4000*4000 => 5.4s`.
 * `5000*5000 => 10s`.
 */
function matMul01(mat1: ArrayLike<number>[], mat2: ArrayLike<number>[]): Uint16Array[] {
  if (!mat1.length || !mat2.length) return []
  if (mat1.length !== mat2.length) throw new Error('mat1.length !== mat2.length')

  const n = mat1.length
  const blockCount = (n + 31) >>> 5
  const block1 = new Uint32Array(blockCount * n)
  const block2 = new Uint32Array(blockCount * n)
  for (let i = 0; i < n; ++i) {
    const cache1 = mat1[i]
    const cache2 = mat2[i]
    for (let j = 0; j < n; ++j) {
      block1[i * blockCount + (j >> 5)] |= cache1[j] << (j & 31)
      block2[j * blockCount + (i >> 5)] |= cache2[j] << (i & 31)
    }
  }

  const res = Array<Uint16Array>(n)
  for (let i = 0; i < n; ++i) {
    const row = new Uint16Array(n)
    res[i] = row
    for (let j = 0; j < n; ++j) {
      let sum = 0
      for (let k = 0; k < blockCount; ++k) {
        sum += bitCount32(block1[i * blockCount + k] & block2[j * blockCount + k])
      }
      row[j] = sum
    }
  }

  return res

  function bitCount32(uint32: number): number {
    uint32 -= (uint32 >>> 1) & 0x55555555
    uint32 = (uint32 & 0x33333333) + ((uint32 >>> 2) & 0x33333333)
    return (((uint32 + (uint32 >>> 4)) & 0x0f0f0f0f) * 0x01010101) >>> 24
  }
}

export { matMul01 }

if (require.main === module) {
  const mat1 = [
    [1, 0, 1],
    [0, 1, 0],
    [1, 0, 1]
  ]
  const mat2 = [
    [1, 0, 1],
    [0, 1, 0],
    [1, 0, 1]
  ]
  console.log(matMul01(mat1, mat2))

  // time
  const n = 5000
  const get01Matrix = (n: number) => {
    const mat = Array(n)
    for (let i = 0; i < n; ++i) {
      const row = Array(n)
      mat[i] = row
      for (let j = 0; j < n; ++j) row[j] = Math.random() > 0.5 ? 1 : 0
    }
    return mat
  }

  const m1 = get01Matrix(n)
  const m2 = get01Matrix(n)

  console.time('matMul01')
  matMul01(m1, m2)
  console.timeEnd('matMul01')
}
