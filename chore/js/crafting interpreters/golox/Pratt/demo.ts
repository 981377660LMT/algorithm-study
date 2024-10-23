import { createParser } from './parser'
import { createScanner } from './scanner'

const text = '1 + 2 * 3'
const scanner = createScanner(text)
const parser = createParser(scanner)
const ast = parser.parseProgram()
console.log(JSON.stringify(ast, null, 2))
