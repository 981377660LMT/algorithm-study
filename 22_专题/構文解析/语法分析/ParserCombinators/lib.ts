/* eslint-disable max-len */
/* eslint-disable no-constant-condition */
/* eslint-disable no-empty */
/* eslint-disable no-await-in-loop */
/* eslint-disable no-inner-declarations */
/* eslint-disable implicit-arrow-linebreak */

// 也叫ParseC.把语法规则和词法规则翻译成了一系列parser的 “组合”.
// https://qszhu.github.io/2021/08/22/parser-combinators.html
// https://github.com/francisrstokes/arcsecond
// https://time.geekbang.org/column/intro/436

// 状态载体：Parser
// 基本单元：str,regExp, whitespace,identifier,stringLiteral
// 带空格的基本单元：token,regExpToken
// 修饰符：zeroOrMore,zeroOrOne,oneOf,seqOf
// 避免循环依赖：lazy

interface IParserState {
  /**
   * 原始的输入字符串.
   */
  origin: string

  /**
   * 当前的解析位置.
   */
  index: number

  /**
   * 当前的解析结果.
   */
  result?: unknown
}

type ParserFn = (state: IParserState) => Promise<IParserState>

class Parser {
  private readonly _parserFn: ParserFn

  constructor(parserFn: ParserFn) {
    this._parserFn = parserFn
  }

  async parse(s: string): Promise<IParserState> {
    return this._parserFn({ origin: s, index: 0 })
  }

  /**
   * 解析过程结束后仍有未匹配到输入的情况定义为出错.
   */
  async parseToEnd(s: string): Promise<IParserState> {
    const res = await this.parse(s)
    if (res.index !== s.length) {
      throw new Error(`Parsing end at ${res.index}: "${_peek(res)}"`)
    }
    return res
  }

  map(f: (arg: any) => unknown): Parser {
    return new Parser(async state => {
      const nextState = await this.parseFn(state)
      return { ...nextState, result: f(nextState.result) }
    })
  }

  get parseFn(): ParserFn {
    return this._parserFn
  }
}

function str(s: string): Parser {
  return new Parser(async state => {
    const { origin, index } = state
    if (index >= origin.length) {
      throw new Error(`str: Tried to match "${s}", but got unexpected EOF.`)
    }

    if (origin.startsWith(s, index)) {
      return { origin, index: index + s.length, result: s }
    }

    throw new Error(`str: Tried to match "${s}", but got "${_peek(state)}" at index ${index}.`)
  })
}

/**
 * 支持正则匹配的parser构造函数.
 */
function regExp(pattern: RegExp): Parser {
  return new Parser(async state => {
    // 如果正则表达式不是以^开头的话，则可能在某个中间位置匹配到结果，那样就相当于跳过了某些字符做了匹配
    console.assert(pattern.source.startsWith('^'), 'regExp should start with "^"')

    const { origin, index } = state
    if (index >= origin.length) {
      throw new Error(`regExp: Tried to match "${pattern}", but got unexpected EOF.`)
    }

    const matching = origin.slice(index).match(pattern)
    if (matching != null) {
      const [res] = matching
      return { ...state, index: index + res.length, result: res }
    }

    throw new Error(
      `regExp: Tried to match "${pattern}", but got "${_peek(state)}" at index ${index}.`
    )
  })
}

/**
 * 对应正则表达式中的 `*`，表示匹配 0 次或多次.
 * 不会抛出错误，不能匹配到的话只会返回空的结果。
 */
function zeroOrMore(parser: Parser): Parser {
  return new Parser(async state => {
    const res: unknown[] = []
    let nextState = state

    try {
      while (true) {
        nextState = await parser.parseFn(nextState)
        res.push(nextState.result)
      }
    } catch (ignore) {}

    return { ...nextState, result: res }
  })
}

/**
 * 对应正则表达式中的 `?`，表示匹配 0 次或 1 次.
 * 不会抛出错误，不能匹配到的话只会返回空的结果。
 */
function zeroOrOne(parser: Parser): Parser {
  return new Parser(async state => {
    try {
      const nextState = await parser.parseFn(state)
      return { ...nextState, result: nextState.result }
    } catch (ignore) {}

    return { ...state, result: undefined }
  })
}

/**
 * 对应正则表达式中的 `+`，表示匹配 1 次或多次.
 * 至少匹配一次，否则抛出错误.
 */
function oneOrMore(parser: Parser): Parser {
  return new Parser(async state => {
    const res: unknown[] = []
    let nextState = state

    try {
      while (true) {
        nextState = await parser.parseFn(nextState)
        res.push(nextState.result)
      }
    } catch (ignore) {}

    if (!res.length) {
      throw new Error(`oneOrMore: Unable to match parser at index ${state.index}: "${_peek(state)}`)
    }

    return { ...nextState, result: res }
  })
}

/**
 * 对应正则表达式中的 `|`.
 * 至少匹配一次，否则抛出错误.
 * 构造的parser会返回传入的一连串parser中第一个匹配成功的结果.
 */
function oneOf(...parsers: Parser[]): Parser {
  return new Parser(async state => {
    for (const parser of parsers) {
      try {
        return await parser.parseFn(state)
      } catch (ignore) {}
    }
    throw new Error(`oneOf: Unable to match any parser at index ${state.index}: "${_peek(state)}"`)
  })
}

/**
 * 构造匹配一个`序列`的parser.
 * 传入的parser必须依次匹配成功，否则会抛出错误.
 */
function seqOf(...parsers: Parser[]): Parser {
  return new Parser(async state => {
    const res: unknown[] = []
    let nextState = state
    for (const parser of parsers) {
      nextState = await parser.parseFn(nextState)
      res.push(nextState.result)
    }
    return { ...nextState, result: res }
  })
}

/**
 * 惰性求值来避免循环引用.
 */
function lazy(thunk: () => Parser): Parser {
  return new Parser(state => {
    const parser = thunk()
    return parser.parseFn(state)
  })
}

/**
 * 如果遇到前面有空白符号，则在匹配之后丢弃.
 */
function token(s: string): Parser {
  return seqOf(Whitespace, str(s)).map(([_, res]) => res)
}

/**
 * 如果遇到前面有空白符号，则在匹配之后丢弃.
 */
function regExpToken(pattern: RegExp): Parser {
  return seqOf(Whitespace, regExp(pattern)).map(([_, res]) => res)
}

function _peek(state: IParserState) {
  const { origin, index } = state
  return origin.slice(index, index + 30)
}

const Whitespace = regExp(/^\s*/)

/**
 * 终结符Identifier.
 */
const Identifier = regExpToken(/^[a-zA-Z_][a-zA-Z0-9_]*/)

/**
 * 只支持双引号(“)，单引号(‘)和反引号(`)就不支持了.
 */
const StringLiteral = regExpToken(/^"[^"]*"/)

export {
  Parser,
  str,
  regExp,
  token,
  zeroOrMore,
  zeroOrOne,
  oneOrMore,
  oneOf,
  seqOf,
  lazy,
  Identifier,
  StringLiteral
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
  async function test() {
    const res = await prog.parse('foo()')
    console.log(JSON.stringify(res, null, 2))
    // const res2 = await prog.parseToEnd('println("hello", "world")')
    // console.log(JSON.stringify(res2, null, 2))
    const res3 = await prog.parseToEnd(`
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

  // testStr()
  // testZeroOrMore()
  // testZeroOrOne()
  // testOneOf()
  // testZeroOrMoreAndOneOf()
  // testSeqOf()

  async function testStr() {
    const parser = token('hello')
    const res = await parser.parse('hello world')
    console.log(res)
    const res2 = await parser.parse('world')
    console.log(res2)
  }

  async function testZeroOrMore() {
    console.log('testZeroOrMore')
    const parser = zeroOrMore(token('hello'))
    const res = await parser.parse('hellohellohello world')
    console.log(res)
    const res2 = await parser.parse('world')
    console.log(res2)
  }

  async function testZeroOrOne() {
    console.log('testZeroOrOne')
    const parser = zeroOrOne(token('hello'))
    const res = await parser.parse('hello world')
    console.log(res)
    const res2 = await parser.parse('world')
    console.log(res2)
  }

  async function testOneOf() {
    console.log('testOneOf')
    const parser = oneOf(token('hello'), token('world'))
    const res = await parser.parse('hello world')
    console.log(res)
    const res2 = await parser.parse('world')
    console.log(res2)
  }

  async function testZeroOrMoreAndOneOf() {
    console.log('testZeroOrMoreAndOneOf')
    const parser = zeroOrMore(oneOf(token('hello'), token('world')))
    const res = await parser.parse('helloworld')
    console.log(res)
    const res2 = await parser.parse('worldhello')
    console.log(res2)
  }

  async function testSeqOf() {
    console.log('testSeqOf')
    const parser = seqOf(token('hello'), token('world'))
    const res = await parser.parse('helloworld')
    console.log(res)
    const res2 = await parser.parse('worldhello')
    console.log(res2)
  }
}
