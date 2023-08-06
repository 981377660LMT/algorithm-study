class Parser {
  constructor(parserFn) {
    this.parserFn = parserFn
  }

  parse(target) {
    const initState = { target, index: 0, isError: false }
    const nextState = this.parserFn(initState)
    if (nextState.index !== target.length) {
      nextState.hasError = true
    }
    return nextState
  }

  map(fn) {
    return new Parser(state => {
      const nextState = this.parserFn(state)
      if (nextState.isError) return nextState
      nextState.result = fn(nextState.result)
      return nextState
    })
  }
}

const str = s =>
  new Parser(state => {
    if (state.isError) return state
    const { target, index } = state
    if (!target.startsWith(s, index)) return { ...state, isError: true }
    return { ...state, index: index + s.length, result: s }
  })

const regExp = pattern =>
  new Parser(state => {
    if (state.isError) return state
    const { target, index } = state
    const slicedTarget = target.slice(index)
    const match = slicedTarget.match(pattern)
    if (!match) return { ...state, isError: true }
    const [result] = match
    return { ...state, index: index + result.length, result }
  })

const zeroOrMore = parser =>
  new Parser(state => {
    if (state.isError) return state
    const results = []
    let nextState = state
    while (true) {
      const testState = parser.parserFn(nextState)
      if (testState.isError) break
      nextState = testState
      results.push(nextState.result)
    }
    return { ...nextState, result: results }
  })

const oneOf = (...parsers) =>
  new Parser(state => {
    if (state.isError) return state
    for (const parser of parsers) {
      const nextState = parser.parserFn(state)
      if (!nextState.isError) return nextState
    }
    return { ...state, isError: true }
  })

const seqOf = (...parsers) =>
  new Parser(state => {
    if (state.isError) return state
    const results = []
    let nextState = state
    for (const parser of parsers) {
      nextState = parser.parserFn(nextState)
      if (nextState.isError) return { ...state, isError: true }
      results.push(nextState.result)
    }
    return { ...nextState, result: results }
  })

const between = (left, right) => content =>
  seqOf(left, content, right).map(([_left, content, _right]) => content)

const sepBy = sep => value =>
  new Parser(state => {
    if (state.isError) return state
    const results = []
    let nextState = state
    while (true) {
      const valueState = value.parserFn(nextState)
      if (valueState.isError) break
      results.push(valueState.result)
      nextState = valueState
      const sepState = sep.parserFn(nextState)
      if (sepState.isError) break
      nextState = sepState
    }
    return { ...nextState, result: results }
  })

const lazy = thunk => new Parser(state => thunk().parserFn(state))

const whitespace = regExp(/^\s/)
const ignored = zeroOrMore(whitespace)

const token = s => seqOf(ignored, str(s)).map(([_, res]) => res)
const regexToken = pattern => seqOf(ignored, regExp(pattern)).map(([_, res]) => res)

const betweenBrackets = between(token('['), token(']'))
const betweenBraces = between(token('{'), token('}'))
const sepByComma = sepBy(token(','))

const stringLit = regexToken(/^"[^"]*"/).map(res => res.slice(1, res.length - 1))

const numberLit = regexToken(/^[+-]?[0-9]+(\.[0-9]*)?/).map(res => Number(res))

const nullLit = token('null').map(() => null)

const booleanLit = oneOf(token('true'), token('false')).map(res => (res === 'true' ? true : false))

const arrayLit = lazy(() => betweenBrackets(sepByComma(jsonLit)))

const kvPair = lazy(() => seqOf(stringLit, token(':'), jsonLit)).map(([k, _colon, v]) => [k, v])

const objectLit = lazy(() => betweenBraces(sepByComma(kvPair))).map(res =>
  res.reduce((acc, [key, val]) => ((acc[key] = val), acc), {})
)

const jsonLit = lazy(() => oneOf(booleanLit, nullLit, stringLit, numberLit, arrayLit, objectLit))

/**
 * @param {string} str
 * @return {*}
 */
var jsonParse = function (str) {
  return jsonLit.parse(str).result
}

console.log(jsonParse(JSON.stringify({ a: 1 })))
