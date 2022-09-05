// !/dev/stdinとprocess.stdin.fdに速度面での違いはほぼなさそうですね。
// それに対してreadlineはかなり時間がかかるようです。

// 解析Number
// parseInt 要解析基数 所以会慢一些
// !parseFloat、Number、+itemが最も高速なようです

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

/// ////////////////////////////////////////
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
