// !/dev/stdinとprocess.stdin.fdに速度面での違いはほぼなさそうですね。
// それに対してreadlineはかなり時間がかかるようです。

// 解析Number
// parseInt 要解析基数 所以会慢一些
// !parseFloat、Number、+itemが最も高速なようです

import * as fs from 'fs'

declare global {
  interface Array<T> {
    at(index: number): T | undefined
    count(item: T): number
  }
}

Array.prototype.at = function (index: number): any {
  if (index < 0) index += this.length
  return this[index]
}

Array.prototype.count = function (item: any): number {
  let res = 0
  for (let i = 0; i < this.length; i++) {
    if (this[i] === item) res++
  }
  return res
}

function makeArray1<T>(size: number, initValue?: T): T[] {
  return Array(size).fill(initValue)
}

function makeArray2<T>(row: number, col: number, initValue?: T): T[][] {
  return Array.from({ length: row }, () => Array(col).fill(initValue))
}

function useInput(debugCase?: string) {
  const data = debugCase == void 0 ? fs.readFileSync(process.stdin.fd, 'utf8') : debugCase
  const dataIter = _makeIter(data)

  function input(): string {
    return dataIter.next().value.trim()
  }

  function nextNum(): number {
    return Number(input())
  }

  function nextNums(): number[] {
    return input().split(' ').map(Number)
  }

  function* _makeIter(str: string): Generator<string, string, any> {
    yield* str.trim().split(/\r\n|\r|\n/)
    return ''
  }

  return {
    input,
    nextNum,
    nextNums,
    next: input,
  }
}

if (require.main === module) {
  const { input } = useInput(
    `
    1
    2 3
    test
    `
  )
  const a = Number(input())
  const [b, c] = input().split(' ').map(Number)
  const s = input()
  console.log(`${a + b + c} ${s}`)
}

export { useInput }
