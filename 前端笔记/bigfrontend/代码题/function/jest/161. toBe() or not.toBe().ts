interface IMatcher {
  toBe(data: any): void
}

function myExpect(input: any): IMatcher & { not: IMatcher } {
  return new Matcher(input)
}

class Matcher implements IMatcher {
  private input: any
  private isReversed: boolean

  constructor(input: any) {
    this.input = input
    this.isReversed = false
  }

  toBe(data: any): boolean {
    const isIdentical = Object.is(this.input, data)
    if ((isIdentical && !this.isReversed) || (!isIdentical && this.isReversed)) return true
    throw new Error('foo')
  }

  get not() {
    this.isReversed = !this.isReversed
    return this
  }
}

if (require.main === module) {
  console.log(myExpect(3).toBe(3)) // ✅
  console.log(myExpect(4).not.toBe(4)) // ❌
}
