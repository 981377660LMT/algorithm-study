/* eslint-disable max-len */
/* eslint-disable no-inner-declarations */
/* eslint-disable implicit-arrow-linebreak */

import { lazy, oneOf, oneOrMore, seqOf, zeroOrMore, zeroOrOne, Identifier, StringLiteral, token } from '../lib'

// prog = statementList* EOF;
const prog = lazy(() => zeroOrMore(statementList))

// statementList = (variableDecl | functionDecl | expressionStatement)+ ;
const statementList = lazy(() => oneOrMore(oneOf(variableStatement, functionDecl, expressionStatement)))

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
const ifStatement = lazy(() =>
  seqOf(token('if'), token('('), expression, token(')'), statement, zeroOrOne(seqOf(token('else'), statement)))
)

// forStatement :
// 'for' '(' (expression | 'let' variableDecl)? ';' expression? ';' expression? ')' statement ;
const forStatement = lazy(() =>
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
const variableDecl = lazy(() => seqOf(Identifier, zeroOrOne(typeAnnotation), zeroOrOne(seqOf(token('='), expression))))

// typeAnnotation : ':' Identifier;
const typeAnnotation = lazy(() => seqOf(token(':'), Identifier))

// functionDecl: "function" Identifier callSignature  block ;
const functionDecl = lazy(() => seqOf(token('function'), Identifier, callSignature, block))

// callSignature: '(' parameterList? ')' typeAnnotation? ;
const callSignature = lazy(() => seqOf(token('('), zeroOrOne(parameterList), token(')'), zeroOrOne(typeAnnotation)))

// returnStatement: 'return' expression? ';' ;
const returnStatement = lazy(() => seqOf(token('return'), zeroOrOne(expression), token(';')))

// emptyStatement: ';' ;
const emptyStatement = lazy(() => token(';'))

// expressionStatement: expression ';' ;
const expressionStatement = lazy(() => seqOf(expression, token(';')))

// expression: assignment;
const expression = lazy(() => assignment)

// assignment: binary (assignmentOp binary)* ;
const assignment = lazy(() => seqOf(binary, zeroOrMore(seqOf(assignmentOp, binary))))

// binary: unary (binOp unary)* ;
const binary = lazy(() => seqOf(unary, zeroOrMore(seqOf(binOp, unary))))

// unary: primary | prefixOp unary | primary postfixOp ;
const unary = lazy(() => oneOf(primary, seqOf(prefixOp, unary), seqOf(primary, postfixOp)))

// primary: StringLiteral |  functionCall | '(' expression ')' ;
const primary = lazy(() => oneOf(StringLiteral, functionCall, seqOf(token('('), expression, token(')'))))

// assignmentOp = '=' | '+=' | '-=' | '*=' | '/=' | '>>=' | '<<=' | '>>>=' | '^=' | '|=' ;
const assignmentOp = oneOf(...['=', '+=', '-=', '*=', '/=', '>>=', '<<=', '>>>=', '^=', '|='].map(token))

// binOp: '+' | '-' | '*' | '/' | '==' | '!=' | '<=' | '>=' | '<'
//      | '>' | '&&'| '||'|...;
const binOp = oneOf(...['+', '-', '*', '/', '==', '!=', '<=', '>=', '<', '>', '&&', '||'].map(token))

// prefixOp = '+' | '-' | '++' | '--' | '!' | '~';
const prefixOp = oneOf(...['+', '-', '++', '--', '!', '~'].map(token))

// postfixOp = '++' | '--';
const postfixOp = oneOf(...['++', '--'].map(token))

// functionCall : Identifier '(' argumentList? ')' ;
const functionCall = lazy(() => seqOf(Identifier, token('('), zeroOrOne(argumentList), token(')')))

// argumentList : expression (',' expression)* ;
const argumentList = lazy(() => seqOf(expression, zeroOrMore(seqOf(token(','), expression))))

// parameterList : parameter (',' parameter)* ;
const parameterList = lazy(() => seqOf(parameter, zeroOrMore(seqOf(token(','), parameter))))

// parameter : Identifier typeAnnotation? ;
const parameter = lazy(() => seqOf(Identifier, zeroOrOne(typeAnnotation)))

if (require.main === module) {
  async function main() {
    const source = 'let a: number = 1 + 2 + 3;'
    const res = await prog.parseToEnd(source)
    // Parsing end at 0: "let a: number = 1 + 2 + 3;"
    console.log(res)
  }

  main()
}
