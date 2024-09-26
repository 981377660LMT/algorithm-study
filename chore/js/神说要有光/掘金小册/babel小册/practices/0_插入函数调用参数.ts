// https://juejin.cn/book/6946117847848321055/section/6945997926376144899

import * as parser from '@babel/parser'
import traverse from '@babel/traverse'
import generate from '@babel/generator'
import * as types from '@babel/types'
import template from '@babel/template'
import { transformSync } from '@babel/core'

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

const targetCalleeName = ['log', 'info', 'warn', 'error', 'debug'].map(v => `console.${v}`)

function addLineColumnToConsoleCall(sourceCode: string) {
  /**
   * 设置为 unambiguous，让 babel 根据内容是否包含 import、export 来自动设置 module/script.
   * 用到了 jsx 的语法，所以 parser 要开启 jsx 的 plugin.
   */
  const ast = parser.parse(sourceCode, { sourceType: 'unambiguous', plugins: ['jsx'] })

  console.dir(ast, { depth: null })

  traverse(ast, {
    /** 当 callee 部分是成员表达式(object.property)，并且是 console.xxx 时，那在参数中插入文件名和行列号，行列号从 AST 的公共属性 loc 上取. */
    CallExpression(path, state) {
      const calleeName = path.get('callee').toString()

      if (targetCalleeName.includes(calleeName)) {
        const startPos = path.node.loc?.start
        if (startPos) {
          const { line, column } = startPos
          path.node.arguments.unshift(types.stringLiteral(`filename: (${line}, ${column})`))
        }
      }
    }
  })

  const { code, map } = generate(ast)
  return { code, map }
}

/** 在当前 console.xx 的 AST 之前插入一个 console.log 的 AST. */
function addLineColumnBeforeConsoleCall(sourceCode: string) {
  const ast = parser.parse(sourceCode, { sourceType: 'unambiguous', plugins: ['jsx'] })
  traverse(ast, {
    CallExpression(path, state) {
      // !要跳过新的节点的处理，就需要在节点上加一个标记，如果有这个标记的就跳过。
      if (path.node.$isNew) return
      const calleeName = path.get('callee').toString()
      if (targetCalleeName.includes(calleeName)) {
        const startPos = path.node.loc?.start
        if (startPos) {
          const { line, column } = startPos
          const newNode = template.expression(`console.log('filename: (${line}, ${column})')`)()
          newNode.$isNew = true
          // 判断要替换的节点是否在 JSXElement 下，所以要用 findParent 的 api 顺着 path 查找是否有 JSXElement 节点
          if (path.findParent(p => p.isJSXElement())) {
            path.replaceWith(types.arrayExpression([newNode, path.node]))
            path.skip() // !跳过新节点的遍历
          } else {
            path.insertBefore(newNode)
          }
        }
      }
    }
  })

  const { code, map } = generate(ast)
  return { code, map }
}

/**
 * 将转换功能改造成babel插件.
 * 然后通过 @babel/core 的 transformSync 方法来编译代码，并引入上面的插件.
 */
function babelPluginStyle() {
  const plugin = ({ types, template }) => ({
    visitor: {
      CallExpression(path, state) {
        if (path.node.$isNew) return
        const calleeName = path.get('callee').toString()
        if (targetCalleeName.includes(calleeName)) {
          const startPos = path.node.loc?.start
          if (startPos) {
            const { line, column } = startPos
            const newNode = template.expression(`console.log('filename: (${line}, ${column})')`)()
            newNode.$isNew = true
            if (path.findParent(p => p.isJSXElement())) {
              path.replaceWith(types.arrayExpression([newNode, path.node]))
              path.skip()
            } else {
              path.insertBefore(newNode)
            }
          }
        }
      }
    }
  })

  const res = transformSync(sourceCode, {
    plugins: [plugin],
    parserOpts: {
      sourceType: 'unambiguous',
      plugins: ['jsx']
    }
  })
  return res?.code
}

if (require.main === module) {
  // console.log(addLineColumnToConsoleCall(sourceCode))
  // console.log(addLineColumnBeforeConsoleCall(sourceCode))
  console.log(babelPluginStyle())
}
