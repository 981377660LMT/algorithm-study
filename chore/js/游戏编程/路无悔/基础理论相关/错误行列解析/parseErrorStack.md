`error-stack-parser` 是一个跨浏览器的 JavaScript 库，用于解析 `Error` 对象中的堆栈跟踪信息。由于不同浏览器（Chrome, Firefox, Safari, IE/Edge）对 `Error.stack` 字符串的格式实现不一致，这个库将它们统一解析为标准的 `StackFrame` 对象数组。

以下是关于 `ErrorStackParser` 的深入讲解和使用示例：

### 1. 核心功能

它主要解决以下痛点：

- **标准化**：将不同浏览器杂乱的堆栈字符串转换为统一的对象结构。
- **提取关键信息**：轻松获取函数名、文件名、行号和列号。

### 2. 基本使用

你当前的代码已经展示了最基础的用法。`ErrorStackParser.parse(error)` 返回一个 `StackFrame[]` 数组。

```typescript
import ErrorStackParser from 'error-stack-parser'

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
```

### 3. StackFrame 对象详解

`ErrorStackParser` 解析出的每个元素都是一个对象（通常兼容 `stackframe` 库的接口），包含以下常用属性：

| 属性名         | 类型     | 说明                                                   |
| :------------- | :------- | :----------------------------------------------------- |
| `functionName` | `string` | 报错所在的函数名称。如果是匿名函数可能为空或推断名称。 |
| `fileName`     | `string` | 报错所在的文件 URL 或路径。                            |
| `lineNumber`   | `number` | 报错所在的行号。                                       |
| `columnNumber` | `number` | 报错所在的列号。                                       |
| `args`         | `any[]`  | 函数参数（通常很少能解析出来，取决于浏览器）。         |

### 4. 常见应用场景

1.  **错误监控系统 (Sentry 等原理)**：
    在全局 `window.onerror` 或 `Vue.config.errorHandler` 中捕获错误，使用此库解析堆栈，然后上报给后端进行聚合分析。

2.  **日志格式化**：
    在控制台输出更易读的错误日志，而不是一长串凌乱的字符串。

3.  **Source Map 还原前的准备**：
    获取到压缩代码的 `lineNumber` 和 `columnNumber` 后，结合 `source-map` 库，可以定位到源码的具体位置。

### 5. 注意事项

- **跨域脚本**：如果报错的脚本是跨域的（CDN 资源），且没有设置 CORS 头，浏览器可能会隐藏详细的行列号（Script Error），此时无法解析出有效信息。
- **兼容性**：虽然它支持大多数浏览器，但极个别老旧环境可能解析不全，不过在现代开发（如你现在的 TypeScript/Node 环境或主流浏览器）中通常表现良好。

### 6. 扩展代码示例：封装一个错误解析工具

```typescript
// ...existing code...

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
```

### 7. 进阶示例：过滤系统堆栈 (Filtering Noise)

在实际开发中，堆栈跟踪往往包含大量 `node_modules` 或框架内部（如 Webpack 包装代码）的调用，这些信息对定位业务逻辑 bug 帮助不大。我们可以利用解析后的对象轻松过滤它们。

```typescript
import ErrorStackParser from 'error-stack-parser'

function getAppStackFrames(error: Error) {
  const frames = ErrorStackParser.parse(error)

  // 过滤规则：只保留不包含 'node_modules' 且有文件名的帧
  return frames.filter(frame => {
    return (
      frame.fileName &&
      !frame.fileName.includes('node_modules') &&
      !frame.fileName.includes('webpack-internal:')
    )
  })
}

try {
  // 模拟业务逻辑错误
  throw new Error('Business Logic Error')
} catch (e: any) {
  const appFrames = getAppStackFrames(e)

  console.log('=== 精简后的堆栈 ===')
  appFrames.forEach(f => {
    console.log(`${f.functionName || 'Top Level'} -> ${f.fileName}:${f.lineNumber}`)
  })
}
```

### 8. 进阶示例：集成到 Window 全局监听 (简易监控 SDK)

这是一个前端监控 SDK 的核心原型。通过监听 `error` 事件，自动捕获未处理的异常，解析堆栈，并整理成标准格式准备上报。

```typescript
import ErrorStackParser from 'error-stack-parser'

interface ErrorReport {
  type: string
  message: string
  stackTrace: object[]
  timestamp: number
  url: string
}

// 初始化全局监控
function initGlobalMonitor() {
  window.addEventListener('error', (event: ErrorEvent) => {
    // 忽略没有 error 对象的 Script Error（通常是跨域脚本问题）
    if (!event.error) return

    try {
      const frames = ErrorStackParser.parse(event.error)

      const report: ErrorReport = {
        type: 'JsRunTimeError',
        message: event.message,
        // 提取前 5 帧即可，避免包体过大
        stackTrace: frames.slice(0, 5).map(f => ({
          function: f.functionName,
          file: f.fileName,
          line: f.lineNumber,
          col: f.columnNumber
        })),
        timestamp: Date.now(),
        url: window.location.href
      }

      console.log('🚨 捕获到异常，准备上报:', report)
      // TODO: navigator.sendBeacon('/api/log', JSON.stringify(report))
    } catch (parseErr) {
      console.warn('ErrorStackParser 解析失败', parseErr)
    }
  })
}
```

### 9. 结合格式化美化输出

有时候我们需要将错误对象转换为类似于 Python 风格的简洁 Log 字符串，以便写入日志文件。

```typescript
function formatError(error: Error): string {
  const frames = ErrorStackParser.parse(error)
  const stackString = frames
    .map(f => `  at ${f.functionName || '<anonymous>'} (${f.fileName}:${f.lineNumber})`)
    .join('\n')

  return `Error: ${error.message}\nTraceback (most recent call last):\n${stackString}`
}
```

---

低代码平台通常使用 `new Function` 配合 `with` 作用域来执行用户编写的代码片段。在这种场景下，直接查看堆栈会有两个主要问题：

1.  **文件名丢失**：报错信息通常显示为 `<anonymous>` 或 `eval`。
2.  **行号偏移**：由于平台包裹了外层代码（如 `with(ctx) { ... }`），堆栈中的行号是包含包裹代码的，无法对应到用户编写的“第 X 行”。

我们可以利用 `ErrorStackParser` 解析出行号，然后减去“包裹代码的前置行数”来还原用户视角的位置。此外，利用 `sourceURL` 可以让浏览器调试器更友好。

### 10. 特殊场景：低代码/沙箱环境 (`with` / `new Function`) 的错误定位

低代码平台常用 `new Function('ctx', 'with(ctx) { ...userCode... }')` 来执行代码。这会导致堆栈行号包含外层包裹代码，且缺乏文件名。

**解决方案：**

1.  **修正行号**：解析堆栈，减去平台注入的 Wrapper 前置行数。
2.  **SourceURL**：注入 `//# sourceURL=xxx` 让浏览器 DevTools 识别虚拟文件。

```typescript
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
```
