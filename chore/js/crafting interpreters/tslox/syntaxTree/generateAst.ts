/* eslint-disable no-useless-concat */

/**
 * Lox Syntax Tree.
 *
 * @see {@link https://craftinginterpreters.com/appendix-ii.html}
 *
 * 1. Expr:
 * - Assign       : name IToken, value Expr
 * - Binary       : left Expr, operator IToken, right Expr
 * - Call         : callee Expr, paren IToken, args Expr[]
 * - Get          : obj Expr, name IToken
 * - Grouping     : expression Expr
 * - Literal      : value any
 * - Logical      : left Expr, operator IToken, right Expr
 * - Set          : obj Expr, name IToken, value Expr
 * - Super        : keyword IToken, method IToken
 * - This         : keyword IToken
 * - Unary        : operator IToken, right Expr
 * - VariableExpr : name IToken
 *
 * 2. Stmt:
 * - Block        : statements Stmt[]
 * - Class        : name IToken, superclass VariableExpr|undefined, methods FunctionStmt[]
 * - Expression   : expression Expr
 * - Function     : name IToken, params IToken[], body Stmt[]
 * - IfStmt       : condition Expr, thenBranch Stmt, elseBranch Stmt|undefined
 * - Print        : expression Expr
 * - Return       : keyword IToken, value Expr|undefined
 * - VariableDecl : name IToken, initializer Expr|undefined
 * - WhileStmt    : condition Expr, body Stmt
 */

import { join } from 'path'
import { existsSync, mkdirSync, writeFileSync } from 'fs'
import { capitalize } from '../utils'

/**
 * Generate the AST classes for the expressions.
 */
function generateAst(outputDir: string): void {
  defineAst(
    'Expr',
    [
      'Assign       : name IToken, value Expr',
      'Binary       : left Expr, operator IToken, right Expr',
      'Call         : callee Expr, paren IToken, args Expr[]',
      'Get          : obj Expr, name IToken',
      'Grouping     : expression Expr',
      'Literal      : value any',
      'Logical      : left Expr, operator IToken, right Expr',
      'SetExpr      : obj Expr, name IToken, value Expr',
      'SuperExpr    : keyword IToken, method IToken',
      'ThisExpr     : keyword IToken',
      'Unary        : operator IToken, right Expr',
      'VariableExpr : name IToken'
    ],
    {
      import: "import { type IToken } from '../types'"
    }
  )

  defineAst(
    'Stmt',
    [
      'Block          : statements Stmt[]',
      'ClassStmt      : name IToken, superclass VariableExpr|undefined, methods Func[]',
      'Expression     : expression Expr',
      'Func           : name IToken, params IToken[], body Stmt[]',
      'IfStmt         : condition Expr, thenBranch Stmt, elseBranch Stmt|undefined',
      'Print          : expression Expr',
      'ReturnStmt     : keyword IToken, value Expr|undefined',
      'VariableDecl   : name IToken, initializer Expr|undefined',
      'WhileStmt      : condition Expr, body Stmt'
    ],
    {
      import: "import { type Expr, type VariableExpr } from './Expr' " + '\n' + "import { type IToken } from '../types'"
    }
  )

  function defineAst(
    baseClassName: string,
    classPropsStrings: string[],
    options?: {
      import?: string
    }
  ): void {
    if (!existsSync(outputDir)) {
      mkdirSync(outputDir, { recursive: true })
    }
    const path = join(outputDir, `${baseClassName}.ts`)

    writeFileSync(path, '') // clear the file

    const wl = (content: string) => writeFileSync(path, `${content}\n`, { flag: 'a' })

    wl('// Generated by generateAst.ts')
    wl('')
    options?.import && wl(options.import)
    wl('')
    defineVisitor(classPropsStrings)
    wl('')
    wl(`export abstract class ${baseClassName} {`)
    wl(`  abstract accept<R>(visitor: ${capitalize(baseClassName)}Visitor<R>): R`)
    wl('}')
    wl('')
    for (const v of classPropsStrings) {
      const [className, propsString] = v.split(':').map(s => s.trim())
      defineType(className, propsString)
    }

    function defineVisitor(types: string[]): void {
      wl(`export interface ${capitalize(baseClassName)}Visitor<R> {`)
      for (const t of types) {
        const typeName = t.split(':')[0].trim()
        wl(`  visit${typeName}${baseClassName}(${typeName.toLowerCase()}: ${typeName}): R`)
      }
      wl('}')
    }

    function defineType(className: string, propsString: string): void {
      wl(`export class ${className} extends ${baseClassName} {`)

      const props = propsString.split(', ').map(s => s.split(' '))
      for (const [name, type] of props) {
        wl(`  readonly ${name}: ${type}`)
      }
      wl('')

      const constructorParameters = props.map(([name, type]) => `${name}: ${type}`).join(', ')
      wl(`  constructor(${constructorParameters}) {`)
      wl('    super()')
      for (const [name] of props) {
        wl(`    this.${name} = ${name}`)
      }
      wl('  }')
      wl('')

      wl(`  override accept<R>(visitor: ${capitalize(baseClassName)}Visitor<R>): R {`)
      wl(`    return visitor.visit${className}${baseClassName}(this)`)
      wl('  }')
      wl('}')
      wl('')
    }
  }
}

if (typeof require !== 'undefined' && typeof module !== 'undefined' && require.main === module) {
  generateAst(__dirname)
}
