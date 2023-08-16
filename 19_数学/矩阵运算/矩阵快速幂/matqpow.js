/* eslint-disable no-inner-declarations */

// 矩阵快速幂
// 长度较小时，Uint32Array的优势不明显，所以这里的二维矩阵直接用number[][]。
// 如果可以将矩阵压缩成一维数组，则使用Uint32Array优势更明显。

function newMatrix(row, col) {
  const res = Array(row)
  for (let i = 0; i < row; i++) {
    res[i] = Array(col).fill(0)
  }
  return res
}

function eye(n) {
  const res = Array(n)
  for (let i = 0; i < n; i++) {
    const row = Array(n).fill(0)
    row[i] = 1
    res[i] = row
  }
  return res
}

function copy(raw) {
  const row = raw.length
  const res = Array(row)
  for (let i = 0; i < row; i++) {
    res[i] = raw[i].slice()
  }
  return res
}

/**
 * uint32乘法.
 */
function mulUint32(num1, num2, mod = 1e9 + 7) {
  return (((Math.floor(num1 / 65536) * num2) % mod) * 65536 + (num1 & 65535) * num2) % mod
}

/**
 * uint32矩阵乘法.
 */
function matMul(mat1, mat2, mod = 1e9 + 7) {
  const row1 = mat1.length
  const row2 = mat2.length
  const col2 = mat2[0].length
  const res = newMatrix(row1, col2)
  for (let i = 0; i < row1; i++) {
    const resRow = res[i]
    const m1Row = mat1[i]
    for (let k = 0; k < row2; k++) {
      const m2Row = mat2[k]
      for (let j = 0; j < col2; j++) {
        resRow[j] = (resRow[j] + mulUint32(m1Row[k], m2Row[j], mod)) % mod
        if (resRow[j] < 0) resRow[j] += mod
      }
    }
  }
  return res
}

/**
 * uint32矩阵快速幂.
 */
function matPow(base, exp, mod = 1e9 + 7) {
  if (base.length !== base[0].length) throw new Error('base is not a square matrix')
  const n = base.length
  let res = eye(n)
  let baseCopy = copy(base)
  while (exp) {
    if (exp & 1) res = matMul(res, baseCopy, mod)
    baseCopy = matMul(baseCopy, baseCopy, mod)
    exp = Math.floor(exp / 2)
  }
  return res
}
