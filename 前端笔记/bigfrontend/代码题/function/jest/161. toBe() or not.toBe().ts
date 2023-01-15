// 161. toBe() or not.toBe()

interface IMatcher {
  toBe(data: unknown): void
}

function myExpect(input: unknown): IMatcher & { not: IMatcher } {
  return new Matcher(input)
}

class Matcher implements IMatcher {
  private _isReversed = false
  private readonly _input: unknown

  constructor(input: unknown) {
    this._input = input
  }

  toBe(data: unknown): boolean {
    const isSame = Object.is(this._input, data)
    if (isSame === this._isReversed) {
      throw new Error(`expect ${this._input} to be ${data}`)
    }
    return true
  }

  get not() {
    this._isReversed = !this._isReversed
    return this
  }
}

if (require.main === module) {
  console.log(myExpect(3).toBe(3)) // ✅
  console.log(myExpect(4).not.toBe(4)) // ❌
}
