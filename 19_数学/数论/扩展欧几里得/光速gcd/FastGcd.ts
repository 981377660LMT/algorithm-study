import { gcd } from '../gcd'

/**
 * O(值域)时间预处理，O(1)时间查询gcd.
 * @deprecated 使用普通的gcd更快.
 */
class FastGcd {
  private readonly _max: number
  private readonly _sqrt: number
  private readonly _preGcd: Uint32Array
  private readonly _fac: Uint32Array
  private readonly _isPrime: Uint8Array
  private readonly _primes: Uint32Array
  private _total = 0

  constructor(max: number) {
    max++
    this._max = max
    this._sqrt = Math.floor(Math.sqrt(max))
    this._preGcd = new Uint32Array((this._sqrt + 1) * (this._sqrt + 1))
    this._fac = new Uint32Array((max + 1) * 3)
    this._isPrime = new Uint8Array(max + 1)
    this._primes = new Uint32Array(max + 1)
    this._build()
  }

  /**
   * alias of {@link query}.
   */
  gcd(a: number, b: number): number {
    return this.query(a, b)
  }

  query(a: number, b: number): number {
    if (a < 0) a = -a
    if (b < 0) b = -b
    if (a > this._max) throw new RangeError(`a(${a}) is out of range.`)
    if (b > this._max) throw new RangeError(`b(${b}) is out of range.`)
    let res = 1
    for (let i = 0; i < 3; i++) {
      let tmp: number
      const f = this._fac[a * 3 + i]
      if (f > this._sqrt) {
        tmp = b % f ? 1 : f
      } else {
        tmp = this._preGcd[f * (this._sqrt + 1) + (b % f)]
      }
      b = Math.floor(b / tmp)
      res *= tmp
    }
    return res
  }

  private _build(): void {
    const { _max, _sqrt, _preGcd, _fac, _isPrime, _primes } = this
    _fac[3] = 1
    _fac[4] = 1
    _fac[5] = 1
    for (let i = 2; i <= _max; i++) {
      if (!_isPrime[i]) {
        _fac[i * 3] = 1
        _fac[i * 3 + 1] = 1
        _fac[i * 3 + 2] = i
        _primes[++this._total] = i
      }
      for (let j = 1; _primes[j] * i <= _max; j++) {
        const tmp = _primes[j] * i
        _isPrime[tmp] = 1
        _fac[tmp * 3] = _fac[i * 3] * _primes[j]
        _fac[tmp * 3 + 1] = _fac[i * 3 + 1]
        _fac[tmp * 3 + 2] = _fac[i * 3 + 2]
        if (_fac[tmp * 3] > _fac[tmp * 3 + 1]) {
          _fac[tmp * 3] ^= _fac[tmp * 3 + 1]
          _fac[tmp * 3 + 1] ^= _fac[tmp * 3]
          _fac[tmp * 3] ^= _fac[tmp * 3 + 1]
        }
        if (_fac[tmp * 3 + 1] > _fac[tmp * 3 + 2]) {
          _fac[tmp * 3 + 1] ^= _fac[tmp * 3 + 2]
          _fac[tmp * 3 + 2] ^= _fac[tmp * 3 + 1]
          _fac[tmp * 3 + 1] ^= _fac[tmp * 3 + 2]
        }
        if (i % _primes[j] === 0) {
          break
        }
      }
    }

    for (let i = 0; i <= _sqrt; i++) {
      _preGcd[i] = i
      _preGcd[i * (_sqrt + 1)] = i
    }
    for (let i = 1; i <= _sqrt; i++) {
      for (let j = 1; j <= i; j++) {
        _preGcd[i * (_sqrt + 1) + j] = _preGcd[j * (_sqrt + 1) + (i % j)]
        _preGcd[j * (_sqrt + 1) + i] = _preGcd[i * (_sqrt + 1) + j]
      }
    }
  }
}

export {}

if (require.main === module) {
  const fastGcd = new FastGcd(1e5)
  console.time('FastGcd')
  for (let i = 0; i < 1e5; ++i) {
    fastGcd.gcd(i, i)
  }
  console.timeEnd('FastGcd')

  console.time('naive gcd')
  for (let i = 0; i < 1e5; ++i) {
    gcd(i, i)
  }
  console.timeEnd('naive gcd')
}
