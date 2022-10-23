/* eslint-disable no-shadow */
/* eslint-disable no-console */

// !将 1-33 选择9个数的所有组合输出到文件

import { createReadStream, createWriteStream, PathLike } from 'fs'
import { createInterface } from 'readline'
import { Readable } from 'stream'

import { combinations } from './combinations'

/**
 * 将可迭代对象写入文件
 *
 * @param iterable 可迭代对象
 * @param outputPath 输出文件
 */
function output(iterable: Iterable<unknown>, outputPath: PathLike): Promise<void> {
  const writeStream = createWriteStream(outputPath)
  return new Promise((resolve, reject) => {
    const readStream = Readable.from(iterable)
    readStream.pipe(writeStream)
    readStream.on('end', () => {
      resolve()
    })
    readStream.on('error', err => {
      reject(err)
    })
  })
}

/**
 * 分割文件
 *
 * @param inputPath 源文件
 * @param lineLimit 每个文件的行数最大值
 */
function split(inputPath: PathLike, lineLimit = 1e7): Promise<void> {
  const readStream = createReadStream(inputPath)
  const lineReader = createInterface({
    input: readStream,
    crlfDelay: Infinity
  })

  let count = 0
  let writeStream = createWriteStream(`${inputPath}.${count}.txt`)
  return new Promise((resolve, reject) => {
    lineReader.on('line', line => {
      writeStream.write(`${line}\n`)
      count++
      if (count % lineLimit === 0) {
        const lineRecord = count // !注意回调里count会改变
        writeStream.end(() => {
          console.log(`split line : ${lineRecord}`)
        })
        writeStream = createWriteStream(`${inputPath}.${count}.txt`)
      }
    })

    lineReader.on('close', () => {
      writeStream.end(() => {
        console.log(`split done , all line : ${count}`)
        resolve()
      })
    })

    lineReader.on('error', err => {
      reject(err)
    })
  })
}

interface WrapperOptions {
  /**
   * 每个元素的长度
   */
  itemEndPaddingLength?: number

  /**
   * 元素间的分隔符
   */
  itemSplt?: string

  /**
   * 每组元素的长度
   */
  groupEndPaddingLength?: number

  /**
   * 每行的组数
   */
  groupCountEachRow?: number

  /**
   * 组间的分隔符
   */
  groupSplit?: string

  /**
   * 换行符
   */
  crlf?: 'CR' | 'LF' | 'CRLF'
}

function* wrapper(iterable: Iterable<unknown[]>, options?: WrapperOptions): Generator<string> {
  const {
    itemEndPaddingLength = 0,
    itemSplt = ' ',
    groupEndPaddingLength = 0,
    groupCountEachRow = 1,
    groupSplit = ' ',
    crlf = 'CRLF'
  } = options ?? {}

  const crlfMap = new Map([
    ['CR', '\r'],
    ['LF', '\n'],
    ['CRLF', '\r\n']
  ])
  const crlfEscape = crlfMap.get(crlf) ?? '\r\n'

  let mod = 0
  for (const data of iterable) {
    mod = (mod + 1) % groupCountEachRow
    const group = data
      .map(item => String(item).padEnd(itemEndPaddingLength))
      .join(itemSplt)
      .padEnd(groupEndPaddingLength)
    const groupEndPadding = mod === 0 ? '' : groupSplit
    yield `${group}${groupEndPadding}`
    if (mod === 0) yield crlfEscape
  }
}

if (require.main === module) {
  const nums = new Uint8Array(33).map((_, i) => i + 1)
  const comb = combinations(nums, 9)
  const wrappedComb = wrapper(comb, {
    itemEndPaddingLength: 2,
    groupCountEachRow: 4,
    groupSplit: '\t|  ',
    crlf: 'CRLF'
  })

  const outputPath = 'output.txt'
  output(wrappedComb, outputPath)
    .then(() => split(outputPath, 1e6))
    .catch(console.error)
}
