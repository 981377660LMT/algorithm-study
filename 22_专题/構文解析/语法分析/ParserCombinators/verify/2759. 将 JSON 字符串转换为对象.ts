/* eslint-disable no-param-reassign */
/* eslint-disable no-return-assign */
/* eslint-disable implicit-arrow-linebreak */
import {
  lazy,
  oneOf,
  regexToken,
  token,
  seqOf,
  sepBy,
  betweenBrackets,
  betweenBraces
} from './Parser'

// https://leetcode.cn/problems/convert-json-string-to-object/solution/yong-parser-combinatorsmiao-by-qinsi-kmuq/
// JSON的语法规则:
// JSON可以是布尔、null、字符串、数字、数组或对象

// jsonLit -> booleanLit | nullLit | stringLit | numberLit | arrayLit | objectLit ;
const jsonLit = lazy(() => oneOf(booleanLit, nullLit, stringLit, numberLit, arrayLit, objectLit))

// booleanLit -> "true" | "false" ;
const booleanLit = oneOf(token('true'), token('false')).map(res => res === 'true')

// nullLit -> "null" ;
const nullLit = token('null').map(() => null)

// 字符串被""包围
// stringLit -> '"' '"' | '"' chars '"' ;
const stringLit = regexToken(/^"[^"]*"/).map(res => res.slice(1, res.length - 1))

// 数字为十进制，支持正负和小数
// numberLit : /^[+-]?[0-9]+(\.[0-9]*)?/ ;
const numberLit = regexToken(/^[+-]?[0-9]+(\.[0-9]*)?/).map(Number)

// 数组是,分隔的JSON列表，并被[]包围
// arrayLit -> "[" (jsonLit ("," jsonLit)*)? "]" ;
const arrayLit = lazy(() => betweenBrackets(_sepByComma(jsonLit)))

// 对象是,分割的键值对列表，并被{}包围
// objectLit -> "{" (kvPair ("," kvPair)*)? "}" ;
const objectLit = lazy(() => betweenBraces(_sepByComma(_kvPair))).map(
  (res: [key: string, value: string]) =>
    res.reduce((pre, [key, val]) => ((pre[key] = val), pre), Object.create(null))
)

const _sepByComma = sepBy(token(','))
// kvPair -> stringLit ":" jsonLit ;
const _kvPair = lazy(() => seqOf(stringLit, token(':'), jsonLit)).map(([key, _, value]) => [
  key,
  value
])

function jsonParse(str: string): any {
  return jsonLit.parse(str).result
}

export {}

if (require.main === module) {
  console.log(jsonParse('true'))
}
