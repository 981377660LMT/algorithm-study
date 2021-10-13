import { dirname, resolve } from 'path'
import { writeFile } from 'fs/promises'

import axios from 'axios'
import { select } from 'xpath'
import { DOMParser } from 'xmldom'

interface Options {
  url: string
  xPathExpression: string
  dirName?: string
  extName?: `.${string}`
}

// todo
// interface FileGenerator {
//   getSourceFrom: (url: string) => Promise<string>
//   parseFromSource: <T>(source: string, xPathExpression: string) => Promise<T>
//   normalizeData: <T>(data: T) => Promise<string[]>
//   generateFiles: (fileNames: string[], dir?: string, extName?: `.${string}`) => Promise<void>
// }

/**
 *
 * @param path 从path提取html
 */
async function getSourceFrom(url: string): Promise<string> {
  return (await axios.get<string>(url)).data
}

/**
 *
 * @param source html字符串
 * @param xpathExpression xpath表达式
 * @description
 * 从html提取文件名
 */
async function parseFromSource(source: string, xPathExpression: string): Promise<string[]> {
  const dom = new DOMParser().parseFromString(source)
  const textNodes = select(xPathExpression, dom) as Text[]
  return textNodes.map(textNode => textNode.data)
}

/**
 *
 * @param data 需要整理的文件名
 * @description
 * 去除不需要的的文件和违法字符，并排序
 */
async function normalizeData(data: string[]): Promise<string[]> {
  return data
    .filter(title => /^\d+/.test(title))
    .map(title => title.replace(/[/\\?%*:|"<>]/g, '-')) // Illegal characters  https://en.wikipedia.org/wiki/Filename#Reserved_characters_and_words
    .reverse()
}

async function generateFiles(
  fileNames: string[],
  dir = dirname(process.cwd()),
  extName = '.md'
): Promise<void> {
  for (const fileName of fileNames) {
    const destination = resolve(dir, fileName) + extName
    await writeFile(destination, '')
  }
}

/**
 * @description
 * generate files
 */
async function crawl(options: Options): Promise<void> {
  const { url, xPathExpression, dirName, extName } = options
  const html = await getSourceFrom(url)
  const fileNames = await parseFromSource(html, xPathExpression)
  const normalizedFileNames = await normalizeData(fileNames)
  await generateFiles(normalizedFileNames, dirName, extName)
}

if (require.main === module) {
  crawl({
    url: 'https://bigfrontend.dev/zh/question',
    xPathExpression: "//ul[@class='List__ListItems-sc-1p5i700-1 kUISA-D']//li//text()",
    dirName: __dirname,
    extName: '.md',
  }).catch(e => console.log(e))
}

export { crawl }
