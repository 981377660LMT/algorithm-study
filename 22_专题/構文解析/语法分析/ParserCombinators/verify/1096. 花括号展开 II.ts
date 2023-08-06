/* eslint-disable implicit-arrow-linebreak */
// https://leetcode.cn/problems/brace-expansion-ii/
// 1096. 花括号展开 II
// TODO WA

import { enumerateProduct } from '../../../../../13_回溯算法/itertools/product'
import { LowercaseLiteral, Parser, lazy, oneOrMore, oneOf, seqOf, str } from './Parser'

// union -> "{" expr ("," expr)+ "}"
// expr -> (Letter | union)+

const letter = LowercaseLiteral.map(res => [res])

const union = lazy(() =>
  seqOf(str('{'), expr, oneOrMore(seqOf(str(','), expr).map(([_, res]) => res)), str('}')).map(
    ([_lhs, first, rest, _rhs]) => {
      const member = [...first, ...rest.flat()]
      return member
    }
  )
)

const expr: Parser<string[]> = lazy(() =>
  oneOrMore(oneOf(letter, union)).map(res => {
    const product: string[] = []
    enumerateProduct(res, group => {
      product.push(group.join(''))
    })
    return product
  })
)

function braceExpansionII(expression: string): string[] {
  const res = expr.parse(expression).result!
  return [...new Set(res)].sort()
}

if (require.main === module) {
  console.log(braceExpansionII('{{a,z},a{b,c},{ab,z}}')) // ["a","ab","ac","z"]
  console.log(braceExpansionII('{{d,e},f{g,h}}')) // [ 'd', 'e', 'fg', 'fh' ]
  console.log()

  // !WA
  //   {a{b,c}}{{d,e},f{g,h}}
  // 只能解析出前面的{a{b,c}}
  console.log(braceExpansionII('{a{b,c}}{{d,e},f{g,h}}')) // ["abdg","abdh","acdg","acdh","ade","adf","aeg","aeh"]
}
