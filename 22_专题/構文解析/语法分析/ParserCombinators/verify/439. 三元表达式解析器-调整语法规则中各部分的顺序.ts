/* eslint-disable no-confusing-arrow */
/* eslint-disable implicit-arrow-linebreak */

// 439. 三元器-调整语法规则中各部分的顺序
// https://leetcode.cn/problems/ternary-expression-parser/

import { NumberLiteral, Parser, lazy, oneOf, seqOf, str } from './Parser'

// 语法:
// cond -> "T" | "F"
// ternary -> cond "?" expr ":" expr
// expr -> Num | ternary | cond
// !注意parserCombinators没有回溯，oneOf会返回第一个匹配到的结果.如果在expr中cond写在ternary的前面
// !那么在原本可以匹配到ternary的地方，就会先匹配到cond并返回。这可能并不是我们想要的。
// 遇到这种情况，只需要调整相对顺序即可：

const t = str('T').map(() => true)
const f = str('F').map(() => false)
const cond = oneOf(t, f)
const ternary = lazy(() =>
  seqOf(cond, str('?'), expr, str(':'), expr).map(([cond, _, trueExpr, __, falseExpr]) => (cond ? trueExpr : falseExpr))
)

const expr: Parser<string | boolean> = lazy(() => oneOf(NumberLiteral, ternary, cond))

function parseTernary(expression: string): string {
  const num = expr.parse(expression).result!
  if (num === true) return 'T'
  if (num === false) return 'F'
  return num
}

if (require.main === module) {
  console.log(parseTernary('T?2:3')) // "2"
}
