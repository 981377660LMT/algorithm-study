import { oneOf, oneOrMore, seqOf, str } from './Parser'

// oneOrMore+oneOf
const p = oneOrMore(oneOf(str('sss'), str('s')))
console.log(p.parse('ss'))
// { origin: 'ss', index: 2, hasError: true, result: [ 'ss' ] }
// 处理index==origin.length的情况
