/* eslint-disable no-useless-constructor */
/* eslint-disable max-len */
/* eslint-disable no-constant-condition */
/* eslint-disable no-empty */
/* eslint-disable no-inner-declarations */
/* eslint-disable implicit-arrow-linebreak */

// 也叫ParseC.把语法规则和词法规则翻译成了一系列parser的 “组合”.
// https://qszhu.github.io/2021/08/22/parser-combinators.html
// https://github.com/francisrstokes/arcsecond
// https://time.geekbang.org/column/intro/436
// 缺点：各种地方都慢，因为是暴力匹配ll(k)，绝大多数情况ll(1)顶多ll(2)就够用了

// 状态载体：Parser
// 基本单元：str,regExp,Whitespace,Ignored,StringLiteral,NumberLiteral
// 带空格的基本单元：token,regExpToken
// 修饰符：zeroOrMore,zeroOrOne,oneOrMore,oneOf,seqOf
// 避免循环依赖：lazy
// 工具函数: between,sepBy...

interface IParserState {
  /** 原始的输入字符串. */
  origin: string
  /** 当前的解析位置. */
  index: number
  /** 解析过程中是否出错.不使用try-catch来加速. */
  hasError: boolean
  /** 当前的解析结果. */
  result?: unknown
}

/** 解析器函数.输一个状态，返回一个 **新的** 状态. */
type ParserFn = (state: IParserState) => IParserState

class Parser {
  /**
   * @param parserFn 解析器函数.输一个状态，返回一个 **新的** 状态.
   */
  constructor(readonly parserFn: ParserFn) {}

  parse(s: string): IParserState {
    const initState = { origin: s, index: 0, hasError: false }
    const nextState = this.parserFn(initState)
    if (nextState.index !== s.length) nextState.hasError = true
    return nextState
  }

  map(f: (parsedResult: any) => unknown): Parser {
    return new Parser(state => {
      const nextState = this.parserFn(state)
      if (nextState.hasError) return nextState
      nextState.result = f(nextState.result)
      return nextState
    })
  }
}

function str(s: string): Parser {
  return new Parser(state => {
    if (state.hasError) return state
    const { origin, index } = state
    if (!origin.startsWith(s, index)) return { ...state, hasError: true }
    return { ...state, index: index + s.length, result: s }
  })
}

/**
 * 支持正则匹配的parser构造函数.
 */
function regExp(pattern: RegExp): Parser {
  return new Parser(state => {
    if (state.hasError) return state
    // !如果正则表达式不是以^开头的话，则可能在某个中间位置匹配到结果，那样就相当于跳过了某些字符做了匹配
    if (!pattern.source.startsWith('^')) {
      throw new Error(`regExp: "${pattern}" should start with "^"`)
    }
    const { origin, index } = state
    const matching = origin.slice(index).match(pattern)
    if (matching == null) {
      return { ...state, hasError: true }
    }
    const [res] = matching
    return { ...state, index: index + res.length, result: res }
  })
}

/**
 * 对应正则表达式中的 `*`，表示匹配 0 次或多次.
 * 不会抛出错误，不能匹配到的话只会返回空的结果。
 */
function zeroOrMore(parser: Parser): Parser {
  return new Parser(state => {
    if (state.hasError) return state
    const res: unknown[] = []
    let nextState = state
    while (true) {
      const testState = parser.parserFn(nextState)
      if (testState.hasError) break
      nextState = testState
      res.push(nextState.result)
    }
    return { ...nextState, result: res }
  })
}

/**
 * 对应正则表达式中的 `?`，表示匹配 0 次或 1 次.
 * 不会抛出错误，不能匹配到的话只会返回空的结果。
 */
function zeroOrOne(parser: Parser): Parser {
  return new Parser(state => {
    if (state.hasError) return state
    const nextState = parser.parserFn(state)
    if (nextState.hasError) return { ...state, result: undefined }
    return nextState
  })
}

/**
 * 对应正则表达式中的 `+`，表示匹配 1 次或多次.
 * 至少匹配一次，否则抛出错误.
 */
function oneOrMore(parser: Parser): Parser {
  return new Parser(state => {
    if (state.hasError) return state
    const res: unknown[] = []
    let nextState = state
    while (true) {
      nextState = parser.parserFn(nextState)
      if (nextState.hasError) break
      res.push(nextState.result)
    }
    if (!res.length) return { ...state, hasError: true }
    return { ...nextState, result: res }
  })
}

/**
 * 对应正则表达式中的 `|`.
 * 至少匹配一次，否则抛出错误.
 */
function oneOf(...parsers: Parser[]): Parser {
  return new Parser(state => {
    if (state.hasError) return state
    for (let i = 0; i < parsers.length; i++) {
      const p = parsers[i]
      const nextState = p.parserFn(state)
      if (!nextState.hasError) return nextState
    }
    return { ...state, hasError: true }
  })
}

/**
 * 构造匹配一个`序列`的parser.
 * 传入的parser必须依次匹配成功，否则会抛出错误.
 */
function seqOf(...parsers: Parser[]): Parser {
  return new Parser(state => {
    if (state.hasError) return state
    const res: unknown[] = []
    let nextState = state
    for (let i = 0; i < parsers.length; i++) {
      const p = parsers[i]
      nextState = p.parserFn(nextState)
      if (nextState.hasError) return { ...state, hasError: true }
      res.push(nextState.result)
    }
    return { ...nextState, result: res }
  })
}

/**
 * @see {@link https://github.com/francisrstokes/arcsecond}
 */
const between = (left: Parser, right: Parser) => (content: Parser) =>
  seqOf(left, content, right).map(([_, res]) => res)

/**
 * @see {@link https://github.com/francisrstokes/arcsecond}
 */
const sepBy = (sep: Parser) => (value: Parser) =>
  new Parser(state => {
    if (state.hasError) return state
    const res: unknown[] = []
    let nextState = state
    while (true) {
      const valueState = value.parserFn(nextState)
      if (valueState.hasError) break
      res.push(valueState.result)
      nextState = valueState
      const sepState = sep.parserFn(nextState)
      if (sepState.hasError) break
      nextState = sepState
    }
    return { ...nextState, result: res }
  })

/**
 * 惰性求值来处理语法中的循环依赖.
 */
const lazy = (thunk: () => Parser) => new Parser(state => thunk().parserFn(state))

const Whitespace = regExp(/^\s/)
const Ignored = zeroOrMore(Whitespace)

/**
 * 如果遇到前面有空白符号，则在匹配之后丢弃.
 */
const token = (s: string): Parser => seqOf(Ignored, str(s)).map(([_, res]) => res)

/**
 * 如果遇到前面有空白符号，则在匹配之后丢弃.
 */
const regexToken = (pattern: RegExp): Parser =>
  seqOf(Ignored, regExp(pattern)).map(([_, res]) => res)

/** 匹配小括号()间的内容. */
const betweenParens = between(token('('), token(')'))
/** 匹配中括号[]间的内容. */
const betweenBrackets = between(token('['), token(']'))
/** 匹配大括号{}间的内容. */
const betweenBraces = between(token('{'), token('}'))

/** 终结符Identifier */
const Identifier = regexToken(/^[a-zA-Z_][a-zA-Z0-9_]*/)

/** 只支持双引号(“)，单引号(‘)和反引号(`)不支持. */
const StringLiteral = regexToken(/^"[^"]*"/)

/** 数字为十进制，支持正负和小数. */
const NumberLiteral = regexToken(/^[+-]?[0-9]+(\.[0-9]*)?/)

/** 小写英文字母. */
const LowercaseLiteral = regexToken(/^[a-z]+/)

/** 大写英文字母. */
const UpperCaseLiteral = regexToken(/^[A-Z]+/)

export {
  Parser,
  str,
  regExp,
  zeroOrMore,
  zeroOrOne,
  oneOrMore,
  oneOf,
  seqOf,
  between,
  sepBy,
  lazy,
  Whitespace,
  Ignored,
  token,
  regexToken,
  betweenParens,
  betweenBrackets,
  betweenBraces,
  Identifier,
  StringLiteral,
  NumberLiteral,
  LowercaseLiteral,
  UpperCaseLiteral
}

if (require.main === module) {
  // 1. prog = (functionDecl | functionCall)* ;
  // 2. functionDecl: "function" Identifier "(" ")"  functionBody;
  // 3. functionBody : '{' functionCall* '}' ;
  // 4. functionCall : Identifier '(' parameterList? ')' ;
  // 5. parameterList : StringLiteral (',' StringLiteral)* ;
  // 6. expression: primary (binOP primary)* ;
  // 7. primary: StringLiteral | DecimalLiteral | IntegerLiteral | functionCall | '(' expression ')' ;

  // !无论我们如何调整这些规则实现的先后顺序，都会遇到“在定义前使用”的编译错误。这可以通过惰性求值来避免：
  // 1. prog = (functionDecl | functionCall)* ;
  const prog = lazy(() => zeroOrMore(oneOf(functionDecl, functionCall))).map(stmts => ({
    type: 'prog',
    stmts
  }))

  // 2. functionDecl: "function" Identifier "(" ")"  functionBody;
  const functionDecl = lazy(() =>
    seqOf(token('function'), Identifier, token('('), token(')'), functionBody)
  ).map(([_, name, _lp, _rp, body]) => ({ type: 'functionDecl', name, body }))

  // 3. functionBody : '{' functionCall* '}' ;
  const functionBody = lazy(() => seqOf(token('{'), zeroOrMore(functionCall), token('}'))).map(
    ([_lb, calls, _rb]) => calls
  )

  // 4. functionCall : Identifier '(' parameterList? ')' ;
  const functionCall = lazy(() =>
    seqOf(Identifier, token('('), zeroOrOne(parameterList), token(')'))
  ).map(([name, _lp, params, _rp]) => ({ type: 'functionCall', name, params }))

  // 5. parameterList : StringLiteral (',' StringLiteral)* ;
  const parameterList = lazy(() =>
    seqOf(StringLiteral, zeroOrMore(seqOf(token(','), StringLiteral)))
  ).map(([param, params]) => [param, ...params.map(([_comma, param]: unknown[]) => param)])

  test()
  function test() {
    const res3 = prog.parse(`
    function foo() {
      println("in foo...")
    }
    
    function bar() {
      println("in bar...")
      foo()
    }
    
    bar()`)
    console.log(JSON.stringify(res3, null, 2))
  }

  testSepBy()
  function testSepBy(): void {
    const p = sepBy(token(','))
    const res = p(StringLiteral).parse('"a",  "b", "c"')
    console.log(res)
  }
}
