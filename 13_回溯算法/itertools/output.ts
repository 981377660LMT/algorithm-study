// Readable.from 的用法
// 将 1-33 选择9个数的所有组合输出到文件
// todo nodejs 文件读取最佳实践

import { createReadStream, createWriteStream, PathLike } from 'fs'
import { createInterface } from 'readline'
import { Readable } from 'stream'

import { combinations } from './combinations'

/**
 * 将 1-33 选择9个数的所有组合输出到文件
 */
function output(to = 'output.txt') {
  const writeStream = createWriteStream(to)
  const nums = new Uint8Array(33).map((_, i) => i + 1)
  const gen = wrapper(combinations(nums, 9))
  Readable.from(gen).pipe(writeStream)

  function* wrapper<E>(iter: Iterable<E[]>): Generator<string> {
    for (const item of iter) {
      yield `${item.join(' ')}\n`
    }
  }
}

/**
 * 分割文件
 *
 * @param from 源文件
 * @param lineLimit 每个文件的行数最大值
 */
function split(from: PathLike, lineLimit = 1e7) {
  const readStream = createReadStream(from)
  const lineReader = createInterface({
    input: readStream,
    crlfDelay: Infinity
  })

  let count = 0
  let writeStream = createWriteStream(`${from}.${count}`)
  lineReader.on('line', line => {
    writeStream.write(`${line}\n`)
    count++
    if (count % lineLimit === 0) {
      writeStream.end()
      writeStream = createWriteStream(`${from}.${count}`)
    }
  })
}

split('out.txt', 1e7)
