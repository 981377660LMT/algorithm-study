/* eslint-disable implicit-arrow-linebreak */
// exprList -> "(" expr ("," expr)* ")"
// notExpr -> "!" exprList
// andExpr -> "&" exprList
// orExpr -> "|" exprList
// expr -> "t" | "f" | notExpr | andExpr | orExpr

import assert from 'assert'
import { Parser, lazy, oneOf, seqOf, str, zeroOrMore } from './Parser'

const expr: Parser<boolean> = lazy(() => oneOf(t, f, notExpr, andExpr, orExpr))
const t = str('t').map(() => true)
const f = str('f').map(() => false)
const orExpr = lazy(() => seqOf(str('|'), exprList)).map(([_, res]) => {
  let newRes = false
  for (let i = 0; i < res.length; i++) {
    const boleans = res[i]
    for (let j = 0; j < boleans.length; j++) {
      newRes = newRes || boleans[j]
    }
  }
  return newRes
})
const andExpr = lazy(() => seqOf(str('&'), exprList)).map(([_, res]) => {
  let newRes = true
  for (let i = 0; i < res.length; i++) {
    const boleans = res[i]
    for (let j = 0; j < boleans.length; j++) {
      newRes = newRes && boleans[j]
    }
  }
  return newRes
})
const notExpr = lazy(() => seqOf(str('!'), exprList)).map(([_, res]) => !res[0][0])
const exprList = lazy(() =>
  seqOf(
    str('('),
    expr.map(res => [res]),
    zeroOrMore(seqOf(str(','), expr).map(([_, res]) => res)),
    str(')')
  )
).map<boolean[][]>(res => res.slice(1, -1) as boolean[][])

function parseBoolExpr(expression: string): boolean {
  return expr.parse(expression).result as boolean
}

// console.log(parseBoolExpr('&(|(f))'))
// // console.log(parseBoolExpr('|(f,f,f,t)'))
// // console.log(parseBoolExpr('!(&(f,t))'))
// // "&(t,f)"
// console.log(parseBoolExpr('&(t,&(t,f))'))
// // "&(|(f))"
// console.log(parseBoolExpr('&(|(f))'))
// "!(&(f,t))"

if (require.main === module) {
  assert(parseBoolExpr('!(&(f,t))'))
}
