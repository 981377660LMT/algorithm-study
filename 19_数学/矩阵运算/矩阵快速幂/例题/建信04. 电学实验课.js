// https://leetcode.cn/contest/ccbft-2021fall/problems/lSjqMF/
// # 实验目标要求同学们用导线连接所有「目标插孔」，
// # 即从任意一个「目标插孔」沿导线可以到达其他任意「目标插孔」
// # 一条导线可连接相邻两列的且行间距不超过 1 的两个插孔
// # 每一列插孔中最多使用其中一个插孔（包括「目标插孔」）
// # 若实验目标可达成，请返回使用导线数量最少的连接所有目标插孔的方案数；否则请返回 0。

// # 1 <= row <= 20
// # 3 <= col <= 10^9
// # 1 < position.length <= 1000

// # O(row^3*log(col)*position.length) = 20^3*log(10^9)*1000
// 超时

const MOD = 1e9 + 7

function electricityExperiment(row, col, position) {
  const T = Array.from({ length: row }, () => Array(row).fill(0))
  for (let r = 0; r < row; r++) {
    T[r][r] = 1
    if (r !== 0) T[r][r - 1] = 1
    if (r !== row - 1) T[r][r + 1] = 1
  }

  const n = position.length
  position.sort((a, b) => a[1] - b[1])
  for (let i = 0; i < n - 1; i++) {
    const [row1, col1] = position[i]
    const [row2, col2] = position[i + 1]
    const rowDiff = Math.abs(row1 - row2)
    const colDiff = Math.abs(col1 - col2)
    if (rowDiff > colDiff) return 0
  }

  let res = 1
  for (let i = 0; i < n - 1; i++) {
    const [row1, col1] = position[i]
    const [row2, col2] = position[i + 1]
    const colDiff = Math.abs(col1 - col2)
    res = mulUint32(res, cal(row1, row2, colDiff), MOD)
  }

  return Number(res)

  function cal(row1, row2, k) {
    const resT = matPow(T, k, MOD)
    return resT[row1][row2]
  }
}

function newMatrix(row, col) {
  const res = Array(row)
  for (let i = 0; i < row; i++) {
    res[i] = Array(col).fill(0)
  }
  return res
}

/**
 * 单位矩阵.
 * @example
 * ```ts
 * eye(3)
 * // [
 * //   [1, 0, 0],
 * //   [0, 1, 0],
 * //   [0, 0, 1]
 * // ]
 * ```
 */
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
