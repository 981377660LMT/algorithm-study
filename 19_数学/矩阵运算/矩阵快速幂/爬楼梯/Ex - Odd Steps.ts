import * as fs from 'fs'

function useInput(debugCase?: string) {
  const data = debugCase == void 0 ? fs.readFileSync(process.stdin.fd, 'utf8') : debugCase
  const dataIter = _makeIter(data)

  function input(): string {
    return dataIter.next().value.trim()
  }

  function* _makeIter(str: string): Generator<string, string, any> {
    yield* str.trim().split(/\r\n|\r|\n/)
    return ''
  }

  return {
    input
  }
}

type BigIntMatrix = BigInt64Array[]

function mul(m1: BigIntMatrix, m2: BigIntMatrix, mod: bigint): BigIntMatrix {
  const [ROW, COL] = [m1.length, m2[0].length]
  const res = Array.from({ length: ROW }, () => new BigInt64Array(COL))
  for (let r = 0; r < ROW; r++) {
    for (let c = 0; c < COL; c++) {
      for (let i = 0; i < ROW; i++) {
        res[r][c] = m1[r][i] * m2[i][c] + res[r][c]
        res[r][c] %= mod
      }
    }
  }

  return res
}

/**
 * @description 矩阵快速幂
 * @param base 矩阵(N*N)
 * @param exp 幂次
 * @param mod 每个数最后取的模
 */
function matqpow(base: BigIntMatrix, exp: bigint, mod: bigint): BigIntMatrix {
  const N = base.length
  let res = Array.from({ length: N }, () => new BigInt64Array(N))
  for (let r = 0; r < N; r++) res[r][r] = 1n

  while (exp) {
    if (exp & 1n) res = mul(res, base, mod)
    exp >>= 1n
    base = mul(base, base, mod)
  }

  return res
}

const MOD = 998244353n
const { input } = useInput()
const [n, s] = input()
  .split(' ')
  .map(v => BigInt(v))
const bad = [
  0n,
  ...input()
    .split(' ')
    .map(v => BigInt(v))
]

let res = [new BigInt64Array([1n]), new BigInt64Array([0n]), new BigInt64Array([0n])]
const trans = [
  new BigInt64Array([1n, 0n, 1n]),
  new BigInt64Array([1n, 0n, 1n]),
  new BigInt64Array([0n, 1n, 0n])
]

for (let i = 1; i < bad.length; i++) {
  const [pre, cur] = [bad[i - 1], bad[i]]
  const tmp = matqpow(trans, cur - pre, MOD)
  res = mul(tmp, res, MOD)
  res[0][0] = 0n
}

const tmp = matqpow(trans, s - bad[bad.length - 1], MOD)
res = mul(tmp, res, MOD)
console.log(Number(res[0][0]))

export {}
