class StringIterator {
  private raw: string
  private cur: string
  private curRemain: number // 当前字符剩余的数目
  private nextCharIndex: number // compressedString 中下一个要生成的字符
  constructor(compressedString: string) {
    this.raw = compressedString
    this.cur = ' '
    this.curRemain = 0
    this.nextCharIndex = 0
  }

  next(): string {
    if (!this.hasNext()) return ' '
    this.mayMove()
    this.curRemain--
    return this.cur
  }

  hasNext(): boolean {
    return this.nextCharIndex < this.raw.length || this.curRemain !== 0
  }

  private mayMove(): void {
    if (this.curRemain === 0) {
      this.cur = this.raw.charAt(this.nextCharIndex++)
      while (
        this.nextCharIndex < this.raw.length &&
        this.isDigit(this.raw.charAt(this.nextCharIndex))
      ) {
        this.curRemain = this.curRemain * 10 + Number(this.raw.charAt(this.nextCharIndex++))
      }
    }
  }

  private isDigit(s: string): boolean {
    const code = s.codePointAt(0)!
    return code >= 48 && code <= 57
  }
}

const iter = new StringIterator('L1e2t1C1o1d1e1')
console.log(iter.next())
console.log(iter.next())
console.log(iter.next())
console.log(iter.next())
console.log(iter.next())
export {}

// 类似于251. 展开二维向量
