/* eslint-disable no-confusing-arrow */
/* eslint-disable implicit-arrow-linebreak */
// https://leetcode.cn/problems/number-of-atoms/

import {
  NumberLiteral,
  Parser,
  lazy,
  oneOf,
  regExp,
  seqOf,
  str,
  zeroOrMore,
  zeroOrOne
} from './Parser'

// 直接按照题目描述翻译的话，很容易写出这样的语法：
// formula -> Atom Count? | formula+ | '(' formula ')' Count?
// !然而这是不行的。原因是formula -> formula+这样的产生式是左递归的。直接实现并运行的话会递归调用到栈溢出。
// 这跟手写递归下降算法遇到的左递归问题是一样的.

// !解决的方法还是改写语法：
// atomGroup -> Atom Count?
// formulaGroup -> "(" formula ")" Count?
// formula -> (atomGroup | formulaGroup)*

const atom = regExp(/^[A-Z][a-z]*/)
const count = NumberLiteral
const atomGroup = seqOf(atom, zeroOrOne(count)).map(
  ([atom, count]) => new Map([[atom, count ? +count : 1]])
)

const formulaGroup: Parser<Map<string, number>> = lazy(() =>
  seqOf(str('('), formula, str(')'), zeroOrOne(count)).map(([_, res, __, count]) => {
    const counter = new Map<string, number>()
    const curCount = count ? +count : 1
    for (const mp of res) {
      for (const [atom, cnt] of mp) {
        counter.set(atom, (counter.get(atom) || 0) + cnt * curCount)
      }
    }
    return counter
  })
)

const formula: Parser<Map<string, number>[]> = lazy(() =>
  zeroOrMore(oneOf(atomGroup, formulaGroup))
)

// 返回所有原子的数量，格式为：第一个（按字典序）原子的名字，跟着它的数量（如果数量大于 1），
// 然后是第二个原子的名字（按字典序），跟着它的数量（如果数量大于 1），以此类推。
function countOfAtoms(s: string): string {
  const res = formula.parse(s).result!
  const counter = new Map<string, number>()
  for (const mp of res) {
    for (const [atom, cnt] of mp) {
      counter.set(atom, (counter.get(atom) || 0) + cnt)
    }
  }

  const atoms = [...counter.keys()].sort()
  return atoms.map(atom => `${atom}${counter.get(atom) === 1 ? '' : counter.get(atom)}`).join('')
}

if (require.main === module) {
  console.log(countOfAtoms('H2O'))
  console.log(countOfAtoms('Mg(OH)2'))
  console.log(countOfAtoms('K4(ON(SO3)2)2'))
}
