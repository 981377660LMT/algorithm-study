import { resolve } from 'path'
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

interface ICrawler {
  crawl(): Promise<void>
}

abstract class Crawler implements ICrawler {
  private url: string
  private xPathExpression: string
  private dirName?: string
  private extName?: `.${string}`

  /**
   *
   * @param url
   * @param xPathExpression
   * @param dirName 默认值process.cwd()
   * @param extName 默认值'.md'
   */
  constructor(options: Options)
  /**
   *
   * @param url
   * @param xPathExpression
   * @param dirName 默认值process.cwd()
   * @param extName 默认值'.md'
   */
  constructor(url: string, xPathExpression: string, dirName?: string, extName?: `.${string}`)
  constructor(...args: any[]) {
    const paramIsOptions = args.length === 1 && typeof args[0] === 'object'
    const paramIsNotOptions = args.length >= 2 && args.length <= 4

    if (paramIsOptions) {
      const { url, xPathExpression, dirName = process.cwd(), extName = '.md' } = args[0] as Options
      this.url = url
      this.xPathExpression = xPathExpression
      this.dirName = dirName
      this.extName = extName
    } else if (paramIsNotOptions) {
      const [url, xPathExpression, dirName, extName] = args
      this.url = url
      this.xPathExpression = xPathExpression
      this.dirName = dirName ?? process.cwd()
      this.extName = extName ?? '.md'
    } else {
      throw new Error('invalid input')
    }
  }

  public async crawl(): Promise<void> {
    const html = await this.getSourceFrom(this.url)
    const fileNames = await this.parseFromSource(html, this.xPathExpression)
    const normalizedFileNames = await this.normalizeData(fileNames)
    await this.generateFiles(normalizedFileNames)
  }

  /**
   *
   * @param path 从path提取html
   */
  protected async getSourceFrom(url: string): Promise<string> {
    return (
      await axios.get<string>(url, {
        headers: {
          'User-Agent':
            'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.81 Safari/537.36',
        },
      })
    ).data
  }

  /**
   *
   * @param source html字符串
   * @param xpathExpression xpath表达式
   * @description
   * 从html提取文件名
   */
  protected async parseFromSource(source: string, xPathExpression: string): Promise<string[]> {
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
  protected async normalizeData(data: string[]): Promise<string[]> {
    return data
      .filter(title => /^\d+/.test(title))
      .map(title => title.replace(/[/\\?%*:|"<>]/g, '-')) // Illegal characters  https://en.wikipedia.org/wiki/Filename#Reserved_characters_and_words
      .reverse()
  }

  protected async generateFiles(fileNames: string[]): Promise<void> {
    for (const fileName of fileNames) {
      const destination = resolve(this.dirName ?? process.cwd(), fileName) + this.extName ?? '.md'
      await writeFile(destination, '')
    }
  }
}

export { Crawler }
