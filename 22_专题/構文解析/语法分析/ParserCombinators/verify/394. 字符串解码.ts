/* eslint-disable no-shadow */
/* eslint-disable implicit-arrow-linebreak */

import {
  LowercaseLiteral,
  NumberLiteral,
  betweenBrackets,
  lazy,
  oneOf,
  seqOf,
  zeroOrMore,
  Parser
} from './Parser'

// encoded -> (Str | Num "[" encoded "]")*
const encoded: Parser = lazy(() =>
  zeroOrMore(oneOf(LowercaseLiteral, seqOf(NumberLiteral, betweenBrackets(encoded))))
)

type Encoded = string | [string, Encoded[]]

function decodeString(s: string): string {
  const res = encoded.parse(s).result as Encoded[]
  return res.map(_normalize).join('')

  function _normalize(encoded: Encoded): string {
    if (typeof encoded === 'string') return encoded
    const [repeat, children] = encoded
    return children.map(_normalize).join('').repeat(+repeat)
  }
}

if (require.main === module) {
  console.log(decodeString('3[a]2[bc]'))
  // abc3[cd]xyz
  console.log(decodeString('abc3[cd]xyz'))
}
