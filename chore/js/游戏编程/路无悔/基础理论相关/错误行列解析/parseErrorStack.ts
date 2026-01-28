import ErrorStackParser from 'error-stack-parser'
import { cloneDeep } from 'lodash-es'

try {
  // 模拟一个错误调用链
  function inner() {
    throw new Error('Something went wrong')
  }
  function outer() {
    inner()
  }
  outer()
} catch (e: any) {
  const frames = ErrorStackParser.parse(e)

  console.log('Total Frames:', frames.length)

  // 遍历堆栈帧
  frames.forEach((frame, index) => {
    console.log(`Frame ${index}:`)
    console.log(`  Function: ${frame.functionName}`) // 函数名 (例如 'inner')
    console.log(`  File:     ${frame.fileName}`) // 文件路径
    console.log(`  Line:     ${frame.lineNumber}`) // 行号
    console.log(`  Column:   ${frame.columnNumber}`) // 列号
    console.log(`  Source:   ${frame.source}`) // 原始堆栈行信息 (如果有)
  })
}

{
  /**
   * 获取错误的顶层调用位置信息
   */
  function getErrorOrigin(error: Error) {
    try {
      const frames = ErrorStackParser.parse(error)
      if (frames && frames.length > 0) {
        // 通常 frames[0] 是错误的发生点
        const topFrame = frames[0]
        return {
          loc: `${topFrame.fileName}:${topFrame.lineNumber}:${topFrame.columnNumber}`,
          func: topFrame.functionName || '<anonymous>'
        }
      }
    } catch (parseError) {
      console.error('Failed to parse error stack', parseError)
    }
    return null
  }

  // 测试工具函数
  const locInfo = getErrorOrigin(new Error('Test Error'))
  console.log('Location Info:', locInfo)
}

{
  /**
   * 模拟低代码执行器
   * @param userCode 用户编写的代码字符串
   * @param context 注入的上下文变量
   */
  function executeLowCode(userCode: string, context: Record<string, any>) {
    // 1. 计算偏移量：
    // 假设我们的 wrapper 结构如下：
    // Line 1: with(ctx) {  <-- 偏移 1 行
    // Line 2:   userCode...
    // Line 3: }
    const preamble = 'with(ctx) {\n'
    const preambleLineOffset = 1 // 根据实际拼接字符串的换行数量确定

    // 2. 注入 sourceURL：这让 DevTools 能看到名为 UserScript.js 的文件，不仅是 anonymous
    const sourceUrl = `\n//# sourceURL=UserScript_${Date.now()}.js`

    try {
      // 构造最终执行的函数体
      const fnBody = preamble + userCode + '\n}' + sourceUrl
      const fn = new Function('ctx', fnBody)
      fn(context)
    } catch (e: any) {
      const frames = ErrorStackParser.parse(e)

      // 通常 frames[0] 就是生成的 new Function 内部的堆栈
      const topFrame = frames[0]

      if (topFrame && typeof topFrame.lineNumber === 'number') {
        // 【核心逻辑】还原行号：
        // 堆栈行号 - 前置包裹行号 = 用户代码行号
        // 注意：不同浏览器对 new Function 行号起始定义可能不同（通常从 1 开始），需实测微调
        const realLine = topFrame.lineNumber - preambleLineOffset

        console.group('🚨 [低代码引擎] 捕获运行时错误')
        console.log(`错误信息: ${e.message}`)
        console.log(`原始位置: Line ${topFrame.lineNumber}`)
        console.log(`修正位置: Line ${realLine} (对应用户代码编辑器)`)

        // 可选：直接打印出错的那一行代码
        const userCodeLines = userCode.split('\n')
        // realLine 从 1 开始，数组下标从 0 开始
        if (userCodeLines[realLine - 1]) {
          console.log(
            `错误代码: "%c${userCodeLines[realLine - 1].trim()}%c"`,
            'color: red; font-weight: bold',
            ''
          )
        }
        console.groupEnd()
      }

      // 记得再次抛出或上报，不要吞掉错误
      throw e
    }
  }

  // --- 测试 ---
  const badUserCode = `
console.log('Start execution');
const a = 10;
// 这一行会报错，因为 doNotExist 未定义，且位于用户代码第 4 行
doNotExist(); 
console.log('End');
`

  executeLowCode(badUserCode, { console })
}

{
  interface ErrorLocation {
    line: number | null
    column: number | null
  }

  /**
   * 解析错误堆栈，提取用户代码中错误发生的行列号
   * @param error 错误对象
   * @param userCode 用户代码字符串
   * @param codeOffset 代码偏移量（用于补偿包装代码的行数）
   * @returns 错误位置信息
   */
  function parseErrorLocation(error: Error, userCode: string, codeOffset = 0): ErrorLocation {
    const stack = error.stack || ''
    const message = error.message || ''
    const maxLine = userCode.split('\n').length

    // 匹配模式：[正则, 匹配源, 是否应用偏移]
    const patterns: Array<[RegExp, string, boolean]> = [
      [/<anonymous>:(\d+):(\d+)/, stack, true], // Chrome/V8: <anonymous>:行:列
      [/Function:(\d+):(\d+)/, stack, true], // 某些浏览器: Function:行:列
      [/at line (\d+),?\s*column (\d+)/i, message, false] // 语法错误消息
    ]

    for (const [regex, source, applyOffset] of patterns) {
      const match = source.match(regex)
      if (match) {
        const line = parseInt(match[1], 10) - (applyOffset ? codeOffset : 0)
        const column = parseInt(match[2], 10)
        if (line >= 1 && line <= maxLine) {
          return { line, column }
        }
      }
    }

    return { line: null, column: null }
  }

  /**
   * 格式化错误消息，包含行列号信息
   * @param message 原始错误消息
   * @param location 错误位置信息
   * @returns 格式化后的错误消息
   */
  function formatErrorMessage(message: string, location: ErrorLocation): string {
    if (location.line !== null) {
      const columnInfo = location.column !== null ? `, 列 ${location.column}` : ''
      return `${message} (行 ${location.line}${columnInfo})`
    }
    return message
  }

  type AnyObject = Record<string, any>

  const shim = {
    structuredClone: (value: unknown) => cloneDeep(value)
  }

  function execJavaScript(expression: string, ctx: AnyObject, superCtx = {}) {
    const fn = new Function(
      'context',
      'superCtx',
      'shim',
      `with(shim) { with(superCtx) { with(context) { return ${expression}; } } }`
    )
    const result = fn(ctx, superCtx, shim)
    return result
  }

  function runJs(expression: string, ctx: AnyObject, superCtx = {}) {
    try {
      const content = execJavaScript(`(() => {\n${expression};\n })()`, ctx, superCtx)
      return {
        status: 'success',
        content,
        error: null
      }
    } catch (e) {
      // 转换后结构:
      // 行1: (function anonymous(context,superCtx,shim
      // 行2: ) {
      // 行3: with(shim) { with(superCtx) { with(context) { return (() => {
      // 行4: 用户代码第1行
      // 用户代码行号 = 堆栈行号 - 3
      const error = e as Error
      const errorLocation = parseErrorLocation(error, expression, 3)
      const formattedError = formatErrorMessage(error.message, errorLocation)
      return {
        status: 'failure',
        error: formattedError,
        content: undefined,
        errorLocation
      }
    }
  }
}
