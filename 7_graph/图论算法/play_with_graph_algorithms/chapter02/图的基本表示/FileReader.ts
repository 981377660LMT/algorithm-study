import { createReadStream } from 'fs'
import { createInterface } from 'readline'

class FileReader<FileData extends unknown[] = string[][]> {
  private constructor(public readonly fileData: FileData) {}

  static async asyncBuild(fileName: string) {
    const fileData = await this.prototype.readLines(fileName)
    return new FileReader(fileData)
  }

  // 可以模仿Java的Scanner.nextInt()
  private async readLines(fileName: string) {
    try {
      const fileStream = createReadStream(fileName, { encoding: 'utf-8' })
      const lines = createInterface({
        input: fileStream,
        crlfDelay: Infinity,
      })
      const buffer: string[][] = []

      for await (const line of lines) {
        buffer.push(line.split(/\s+/))
      }

      return buffer
    } catch (error) {
      throw new Error(`读取文件${fileName}失败 ${error}`)
    }
  }
}

export { FileReader }
