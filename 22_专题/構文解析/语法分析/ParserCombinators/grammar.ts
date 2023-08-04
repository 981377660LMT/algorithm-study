/* eslint-disable max-len */
/* eslint-disable no-confusing-arrow */
/* eslint-disable import/export */
/* eslint-disable implicit-arrow-linebreak */

import { Parser, lazy, oneOf, oneOrMore, regExp, seqOf, str, zeroOrMore, zeroOrOne } from './lib'

const whiteSpace = regExp(/^\s/)
const ignored = zeroOrMore(whiteSpace)
const token = (s: string) => seqOf(ignored, str(s)).map(([, res]) => res)
const regexToken = (r: RegExp) => seqOf(ignored, regExp(r)).map(([, res]) => res)
const Prologue = regexToken(/^```.*```/s).map(res => res.slice(3, res.length - 3)) // s: 启用 “dotall” 模式，允许点 . 匹配换行符 \n
const Regex = regexToken(/^[/][^].*[/](?=\s*[;.])/)
const MapCode = regexToken(/^\.map\(.*?\)(?=\s*;)/s)
const Terminal = regexToken(/^[A-Z][A-Za-z_]*/)
const NonTerminal = regexToken(/^[a-z][A-Za-z_]*/)
const Literal = regexToken(/^"[^"]*"/).map(res => `token(${res})`)

const primary: Parser = lazy(() => oneOf(Terminal, NonTerminal, Literal, seqOf(token('('), choice, token(')')))).map(
  res => {
    if (Array.isArray(res)) return res[1]
    return res
  }
)

const qualified: Parser = lazy(() =>
  oneOf(seqOf(primary, token('?')), seqOf(primary, token('*')), seqOf(primary, token('+')), primary)
).map(res => {
  if (Array.isArray(res)) {
    const [primary, quantifier] = res
    if (quantifier === '?') return `zeroOrOne(${primary})`
    if (quantifier === '*') return `zeroOrMore(${primary})`
    if (quantifier === '+') return `oneOrMore(${primary})`
  }
  return res
})

const sequence: Parser = lazy(() => oneOrMore(qualified)).map(res =>
  res.length > 1 ? `seqOf(${res.join(', ')})` : res[0]
)

const choice: Parser = lazy(() => seqOf(sequence, zeroOrMore(seqOf(token('|'), sequence)))).map(([first, rest]) => {
  if (Array.isArray(rest) && rest.length > 0) {
    rest = rest.map(r => r[1])
    return `oneOf(${[first, ...rest].join(', ')})`
  }
  return first
})

const syntax: Parser = lazy(() => seqOf(NonTerminal, token('->'), choice, zeroOrOne(MapCode), token(';'))).map(
  ([head, _arrow, body, code, _semi]) => `export const ${head}: Parser = lazy(() => ${body})${code || ''};\n`
)

const lexical: Parser = lazy(() => seqOf(Terminal, token(':'), Regex, zeroOrOne(MapCode), token(';'))).map(
  ([head, _arrow, body, code, _semi]) => `export const ${head} = regexToken(${body})${code || ''};\n`
)

export const grammar: Parser = lazy(() => seqOf(zeroOrOne(Prologue), zeroOrMore(oneOf(syntax, lexical)))).map(
  ([prologue, rules]) => `${prologue || ''}${rules.join('\n')}`
)
