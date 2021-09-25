class StringIterator {
  private nextCharIndex: number // compressedString 中下一个要生成的字符
  private remain: number // 当前字符剩余的数目
  private compressedString: string
  private curChar: string
  constructor(compressedString: string) {
    this.compressedString = compressedString
    this.nextCharIndex = 0
    this.remain = 0
    this.curChar = ' '
  }

  next(): string {
    if (!this.hasNext()) return ' '
    this.upgradeQuery()
    this.remain--
    return this.curChar
  }

  hasNext(): boolean {
    return this.nextCharIndex < this.compressedString.length || this.remain !== 0
  }

  private upgradeQuery() {
    if (this.remain === 0) {
      this.curChar = this.compressedString.charAt(this.nextCharIndex++)
      while (
        this.nextCharIndex < this.compressedString.length &&
        this.isDigit(this.compressedString.charAt(this.nextCharIndex))
      ) {
        this.remain = this.remain * 10 + Number(this.compressedString.charAt(this.nextCharIndex++))
      }
    }
  }

  private isDigit(x: any) {
    return !isNaN(parseFloat(x)) && isFinite(x)
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
