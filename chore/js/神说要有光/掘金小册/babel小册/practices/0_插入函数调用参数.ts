// https://juejin.cn/book/6946117847848321055/section/6945997926376144899

import * as parser from '@babel/parser'
import traverse from '@babel/traverse'
import generate from '@babel/generator'
import * as types from '@babel/types'

const sourceCode = `
console.log(1);

function func() {
    console.info(2);
}

export default class Clazz {
    say() {
        console.debug(3);
    }
    render() {
        return <div>{console.error(4)}</div>
    }
}
`

/**
 * 设置为 unambiguous，让 babel 根据内容是否包含 import、export 来自动设置 module/script.
 * 用到了 jsx 的语法，所以 parser 要开启 jsx 的 plugin.
 */
const ast = parser.parse(sourceCode, { sourceType: 'unambiguous', plugins: ['jsx'] })

console.dir(ast, { depth: null })

const targetCalleeName = ['log', 'info', 'warn', 'error', 'debug'].map(v => `console.${v}`)

traverse(ast, {
  /** 当 callee 部分是成员表达式(object.property)，并且是 console.xxx 时，那在参数中插入文件名和行列号，行列号从 AST 的公共属性 loc 上取. */
  CallExpression(path, state) {
    const callee = path.node.callee
    const isConsoleFunc =
      types.isMemberExpression(callee) && callee.object.name === 'console' && ['log', 'info', 'warn', 'error', 'debug'].includes(callee.property.name)

    if (isConsoleFunc) {
      const startPos = path.node.loc?.start
      if (startPos) {
        const { line, column } = startPos
        path.node.arguments.unshift(types.stringLiteral(`filename: (${line}, ${column})`))
      }
    }
  }
})

const { code, map } = generate(ast)

console.log({ code, map })
