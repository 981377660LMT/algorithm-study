/* eslint-disable max-len */
/* eslint-disable implicit-arrow-linebreak */

import {
  NumberLiteral,
  Parser,
  StringLiteral,
  lazy,
  oneOf,
  oneOrMore,
  regexToken,
  seqOf,
  token,
  zeroOrMore,
  zeroOrOne
} from '../Parser'

// prog = statementList* EOF;
const prog = lazy(() => zeroOrMore(statementList)).map(stmts => ({ type: 'prog', stmts }))

// statementList = (variableDecl | functionDecl | expressionStatement)+ ;
const statementList = lazy(() => oneOrMore(oneOf(variableDecl, functionDecl, expressionStatement)))

// statement: block | expressionStatement | returnStatement | ifStatement | forStatement
//         | emptyStatement | functionDecl | variableDecl ;
const statement = lazy(() =>
  oneOf(
    block,
    expressionStatement,
    returnStatement,
    ifStatement,
    forStatement,
    emptyStatement,
    functionDecl,
    variableDecl
  )
)

// block : '{' statementList? '}' ;
const block = lazy(() => seqOf(token('{'), zeroOrOne(statementList), token('}')))

// ifStatement : 'if' '(' expression ')' statement ('else' statement)? ;
const ifStatement: Parser = lazy(() =>
  seqOf(
    token('if'),
    token('('),
    expression,
    token(')'),
    statement,
    zeroOrOne(seqOf(token('else'), statement))
  )
)

// forStatement : 'for' '(' (expression | 'let' variableDecl)? ';' expression? ';' expression? ')' statement ;
const forStatement: Parser = lazy(() =>
  seqOf(
    token('for'),
    token('('),
    zeroOrOne(oneOf(expression, seqOf(token('let'), variableDecl))),
    token(';'),
    zeroOrOne(expression),
    token(';'),
    zeroOrOne(expression),
    token(')'),
    statement
  )
)

// variableStatement : 'let' variableDecl ';';
const variableStatement = lazy(() => seqOf(token('let'), variableDecl, token(';')))

// variableDecl : Identifier typeAnnotation？ ('=' expression)? ;
const variableDecl = lazy(() =>
  seqOf(Identifier, zeroOrOne(typeAnnotation), zeroOrOne(seqOf(token('='), expression)))
)

// typeAnnotation : ':' Identifier;
const typeAnnotation = lazy(() => seqOf(token(':'), Identifier))

// functionDecl: "function" Identifier callSignature  block ;
const functionDecl: Parser = lazy(() =>
  seqOf(token('function'), Identifier, callSignature, block)
).map(([_, name, __, body]) => ({ type: 'functionDecl', name, body }))

// callSignature: '(' parameterList? ')' typeAnnotation? ;
const callSignature = lazy(() =>
  seqOf(token('('), zeroOrOne(parameterList), token(')'), zeroOrOne(typeAnnotation))
)

// returnStatement: 'return' expression? ';' ;
const returnStatement = lazy(() => seqOf(token('return'), zeroOrOne(expression), token(';')))

// emptyStatement: ';' ;
const emptyStatement = lazy(() => token(';'))

// expressionStatement: expression ';' ;
const expressionStatement = lazy(() => seqOf(expression, token(';')))

// expression: assignment;
const expression: Parser = lazy(() => assignment)

// assignment: binary (assignmentOp binary)* ;
const assignment = lazy(() => seqOf(binary, zeroOrMore(seqOf(assignmentOp, binary))))

// binary: unary (binOp unary)* ;
const binary = lazy(() => seqOf(unary, zeroOrMore(seqOf(binOp, unary))))

// unary: primary | prefixOp unary | primary postfixOp ;
const unary: Parser = lazy(() => oneOf(primary, seqOf(prefixOp, unary), seqOf(primary, postfixOp)))

// primary: StringLiteral | DecimalLiteral | IntegerLiteral | functionCall | '(' expression ')' ;
const primary = lazy(() =>
  oneOf(StringLiteral, NumberLiteral, functionCall, seqOf(token('('), expression, token(')')))
)

// assignmentOp = '=' | '+=' | '-=' | '*=' | '/=' | '>>=' | '<<=' | '>>>=' | '^=' | '|=' ;
const assignmentOp = oneOf(
  ...['=', '+=', '-=', '*=', '/=', '>>=', '<<=', '>>>=', '^=', '|='].map(token)
)

// binOp: '+' | '-' | '*' | '/' | '==' | '!=' | '<=' | '>=' | '<'
//      | '>' | '&&'| '||'|...;
const binOp = oneOf(
  ...['+', '-', '*', '/', '==', '!=', '<=', '>=', '<', '>', '&&', '||'].map(token)
)

// prefixOp = '+' | '-' | '++' | '--' | '!' | '~';
const prefixOp = oneOf(...['+', '-', '++', '--', '!', '~'].map(token))

// postfixOp = '++' | '--';
const postfixOp = oneOf(...['++', '--'].map(token))

// functionCall : Identifier '(' argumentList? ')' ;
const functionCall = lazy(() =>
  seqOf(Identifier, token('('), zeroOrOne(argumentList), token(')'))
).map(([name, _lparen, params, _rparen]) => ({ type: 'functionCall', name, params }))

// argumentList : expression (',' expression)* ;
const argumentList = lazy(() => seqOf(expression, zeroOrMore(seqOf(token(','), expression))))

const Identifier = regexToken(/^[a-zA-Z_][a-zA-Z0-9_]*/)

// parameterList : parameter (',' parameter)* ;
const parameterList = lazy(() => seqOf(parameter, zeroOrMore(seqOf(token(','), parameter)))).map(
  ([param, params]) => [param, ...params.map(([_comma, param]: unknown[]) => param)]
)

// parameter : Identifier typeAnnotation? ;
const parameter = lazy(() => seqOf(Identifier, zeroOrOne(typeAnnotation)))

//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
// !运算符优先级：确保运算符优先级的传统方法之一是修改语法
// primary : Num
// product : primary ('*' primary)*
// sum : product ('+' product)*

const Num = oneOf(token('0'), regexToken(/^[1-9][0-9]*/)).map(Number)

// product : primary ('*' primary)*
const product = lazy(() => seqOf(primary, zeroOrMore(seqOf(token('*'), Num)))).map(
  ([lhs, rest]) => [lhs, ...rest].reduce((lhs, [op, rhs]) => [lhs, op, rhs]) // reduce来将匹配到的结果左结合
)

// sum : product ('+' product)*
const sum = lazy(() => seqOf(product, zeroOrMore(seqOf(token('+'), product)))).map(([lhs, rest]) =>
  [lhs, ...rest].reduce((lhs, [op, rhs]) => [lhs, op, rhs])
)

console.log(sum.parse('1+2*3').result)

//
//
// !右结合的运算符要怎么写呢,比如求幂运算**。诀窍就在于右递归.
// primary : Num
// power : primary '**' power | primary
// product : power ('*' power)*
// sum : product ('+' product)*
// 此外因为其优先级比*要高，所以要插在*的规则上面
//
//
//
// !但缺点是每遇到一个运算符就要写一个语法规则出来
//
//
//
// !reduce和curry化
// 1) primary = Num
// 2) product = f('*', primary)
// 3) sum = f('+', product)
// 1)代入2)：
// 2) product = f('*', Num)
// 3) sum = f('+', product)
// 2)代入3)：
// !3) sum = f('+', f('*', Num)) <=> sum = ['*', '+'].reduce((term, op) => f(op, term), Num)
// !也就是说我们只要把运算符按照优先级从高到低放在一个数组里，再reduce一下就好了。
// reducer 本质是 幺半群(e+op)
const infix = (operator: Parser) => (nextTerm: Parser) =>
  seqOf(nextTerm, zeroOrMore(seqOf(operator, nextTerm)))
//   const sum = ...
//   seqOf(product, zeroOrMore(seqOf(token('+'), product))))
// //      ^^^^^^^                   ^^^^^^^^^^  ^^^^^^^
//     .map(...)

// const infix = (operator: Parser, nextTerm: Parser) =>
//   seqOf(nextTerm, zeroOrMore(seqOf(operator, nextTerm)))
// //      ^^^^^^^^                   ^^^^^^^^  ^^^^^^^^
//     .map(...)

const expr = [
  infix(oneOf(token('*'), token('/'), token('%'))),
  infix(oneOf(token('+'), token('-')))
].reduce((term, op) => op(term), Num)
console.log(expr.parse('1+2*3').result)
//
// 简化右结合运算符
const infixRight = (operator: Parser) => (nextTerm: Parser) => {
  const parser: Parser = lazy(() => oneOf(seqOf(nextTerm, operator, parser), nextTerm))
  return parser
}
const exprOps = [
  infixRight(oneOf(token('**'))),
  infix(oneOf(token('*'), token('/'), token('%'))),
  infix(oneOf(token('+'), token('-')))
]
const expr2 = exprOps.reduce((term, op) => op(term), Num)
console.dir(expr2.parse('1 + 2 * 3 ** 4 + 5').result, { depth: null })

// !这样，增加新的运算符和调整运算符优先级，都只要修改上面的exprOps数组就可以了。
// !前缀和后缀表达式
// 前缀运算符需要向右结合，后缀运算符需要向左结合，所以分别参照中缀表达式中左结合和右结合运算符的实现修改即可
// !后缀和左结合：
const postfix = (operator: Parser) => (nextTerm: Parser) => seqOf(nextTerm, zeroOrMore(operator))
const infix2 = (operator: Parser) => (nextTerm: Parser) =>
  seqOf(nextTerm, zeroOrMore(seqOf(operator, nextTerm)))
// !前缀和右结合：
const prefix = (operator: Parser) => (nextTerm: Parser) => {
  const parser: Parser = lazy(() => oneOf(seqOf(operator, parser), nextTerm))
  return parser
}
const infixRight2 = (operator: Parser) => (nextTerm: Parser) => {
  const parser: Parser = lazy(() => oneOf(seqOf(nextTerm, operator, parser), nextTerm))
  return parser
}
const exprOps2 = [
  postfix(oneOf(token('++'), token('--'))),
  prefix(oneOf(token('++'), token('--'), token('+'), token('-'))),
  infixRight(token('**')),
  infix(oneOf(token('*'), token('/'), token('%'))),
  infix(oneOf(token('+'), token('-')))
]
console.log(exprOps2.reduce((term, op) => op(term), Num).parse('1 + 2 * 3 ** 4 + 5').result)

// !确定性
// 到目前为止我们写的parser能支持的语法都是无二义性的
// 如果相同的表达式能解析出不同的结果，则说明语法是有二义性的。
// 你可以把parser combinators理解为动态生成的递归下降算法的实现，
// 不加特别处理的话是没有回溯的过程的，就只能解析出一个确定的结果（或是解析出错）。
// 对于程序设计语言来说我觉得这不是一件坏事。

if (require.main === module) {
  const source = `let a: number = 1 + 2 * 3;`
  const { stmts } = prog.parse(source).result!
  console.dir(stmts, { depth: null })
  /*
  { type: 'prog',
    stmts:
    [ { type: 'let',
        decl:
          { name: 'a', type: 'number', val: [ 1, [ '+', 2 ], [ '+', 3 ] ] } } ] }
  */
}
