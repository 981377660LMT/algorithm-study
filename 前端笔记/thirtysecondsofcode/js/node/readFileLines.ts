// 使用 fs.readFileSync ()从文件创建缓冲区。
// buf.toString (encoding)函数将 buffer 转换为 string。
import { PathOrFileDescriptor, readFileSync } from 'fs'

const readFileLines = (filename: PathOrFileDescriptor) =>
  readFileSync(filename).toString('utf-8').split('\n')
/*
contents of test.txt :
  line1
  line2
  line3
  ___________________________
*/
const arr = readFileLines('test.txt')
console.log(arr) // ['line1', 'line2', 'line3']
