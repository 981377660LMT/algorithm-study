// 读取语法文件，生成parser

import fs from 'fs'
import { resolve } from 'path'
import { grammar } from './grammar'

async function main(fn: string) {
  const source = fs.readFileSync(fn).toString().trim()
  const res = await grammar.parse(source)
  fs.writeFileSync(resolve(__dirname, '_output.ts'), res.result as string)
}

if (require.main === module) {
  main(resolve(__dirname, 'grammar.grammar')).catch(console.error)
}
