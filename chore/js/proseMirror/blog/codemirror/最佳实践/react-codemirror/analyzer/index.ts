import * as acorn from 'acorn'
import Mustache from 'mustache'

// #region types
export interface JsViolation {
  type:
    | 'REQUIRE_IMPORT'
    | 'VARIABLE_ASSIGNMENT'
    | 'SET_VALUE_CALL'
    | 'TRIGGER_CALL'
    | 'MUSTACHE_EXPRESSION'
  message: string
  /** 在原始 script 中的起始字符偏移（供编辑器 lint 标记位置使用）. */
  from?: number
  /** 在原始 script 中的结束字符偏移. */
  to?: number
}

export interface CodeSnippet {
  script: string
  label: string
}
// #endregion

// #region extract JS code from task
/**
 * 从 task 的 flowList 中提取需要校验的 JS 代码片段.
 * - type === 'js_query' → content: { script: string }
 * - type === 'transformer' → content: { content: { name: 'script', code: string } }
 */
export function extractJsCode(task): CodeSnippet[] {
  const snippets: CodeSnippet[] = []

  for (const flow of task.flowList ?? []) {
    if (!flow.content) continue

    try {
      const parsed = JSON.parse(flow.content)

      if (flow.type === 'js_query' && typeof parsed?.script === 'string' && parsed.script.trim()) {
        snippets.push({
          script: parsed.script,
          label: `接口 "${task.name}" 的 JsQuery 代码`
        })
      }

      if (
        flow.type === 'transformer' &&
        parsed?.content?.name === 'script' &&
        typeof parsed?.content?.code === 'string' &&
        parsed.content.code.trim()
      ) {
        snippets.push({
          script: parsed.content.code,
          label: `接口 "${task.name}" 的 Transformer 代码`
        })
      }
    } catch {}
  }

  return snippets
}
// #endregion

// #region collect global names from schema
/**
 * 收集 schema 中所有全局名称（组件名、state 名、task 名、模块出入参名）.
 */
export function collectGlobalNames(schema: Schema): Set<string> {
  const names = new Set<string>()
  ;(schema?.schema?.children || []).forEach(item => {
    const name = item.data?.value?.data?.name
    if (name) names.add(name)
  })
  ;(schema?.states || []).forEach(item => {
    if (item.name) names.add(item.name)
  })
  ;(schema?.tasks || []).forEach(item => {
    if (item.name) names.add(item.name)
  })
  ;(schema?.landerModuleInput || []).forEach(item => {
    if (item.name) names.add(item.name)
  })
  return names
}
// #endregion

// #region AST analysis
/**
 * 解析 JS 代码并检测违规模式.
 * 如果代码存在语法错误，静默跳过不报错.
 *
 * 检测策略：
 * 1. 检测 Mustache 表达式 {{ }}：JsQuery / Transformer 中不应包含 {{ }}，它们会被状态管理层误识别为模板表达式，导致代码被错误求值
 * 2. 正则检测静态 import 声明（避免 AST 解析 import 与 return 共存导致的兼容问题）
 * 3. 剥离 import 行后，将剩余代码包裹在 function 中以 script 模式解析 AST，
 *    检测 require()、动态 import()、全局变量赋值、.setValue()、.trigger() 等违规
 */
export function analyzeJsCode(script: string, globalNames: Set<string>): JsViolation[] {
  const violations: JsViolation[] = []

  // ---- 检测 Mustache 表达式 ----
  checkMustacheExpressions(script, violations)

  // ---- 第一步：正则检测静态 import 声明 ----
  checkStaticImports(script, violations)

  // ---- 第二步：把 import 行替换为同长空格（保留偏移），用 AST 检测其余违规 ----
  const strippedScript = script.replace(/^[ \t]*import\b(?!\s*\().*$/gm, match =>
    ' '.repeat(match.length)
  )

  if (strippedScript.trim()) {
    try {
      const wrapperPrefix = '(function() {\n'
      const offset = wrapperPrefix.length // AST 偏移 → 原 script 偏移：sub `offset`
      const ast = acorn.parse(`${wrapperPrefix}${strippedScript}\n})()`, {
        ecmaVersion: 2020,
        sourceType: 'script'
      })
      walkNode(ast, violations, globalNames, offset)
    } catch {}
  }

  return violations
}

/**
 * UI 层便捷入口：接收一个全局变量名的可迭代集合（如 `Object.keys(globalContext)`），
 * 自动包装为 `Set<string>` 后调用 `analyzeJsCode`.
 *
 * @param script 待校验的 JS 代码片段
 * @param globalNames 全局变量名（数组、Set 或其它可迭代结构）
 * @param options.excludeNames 需要排除的名字（如当前正在编辑的 query 自身名称，避免误报"对自身赋值"）
 */
export function analyzeJsCodeWithGlobalContext(
  script: string,
  globalNames: Iterable<string>,
  options?: { excludeNames?: Iterable<string> }
): JsViolation[] {
  const set = new Set(globalNames)
  if (options?.excludeNames) {
    for (const name of options.excludeNames) set.delete(name)
  }
  return analyzeJsCode(script, set)
}

/**
 * 检测 JS 代码中的 Mustache 表达式 {{ }}.
 * JsQuery / Transformer 中不应包含 {{ }}，它们会被状态管理层误识别为模板表达式，
 * 导致代码被错误求值。应直接引用变量名（如 query1.data）而非 {{query1.data}}。
 */
function checkMustacheExpressions(script: string, violations: JsViolation[]): void {
  try {
    // 先移除注释，避免误报注释中的 {{ }}
    const scriptWithoutComments = script
      .replace(/\/\*[\s\S]*?\*\//g, '') // 移除多行注释
      .replace(/\/\/.*$/gm, '') // 移除单行注释

    const tokens = Mustache.parse(scriptWithoutComments)
    if (!tokens.some(token => token[0] !== 'text')) {
      return
    }

    // 在原始 script 上 regex 定位，以便产出编辑器可用的 from/to.
    // 注释里的 {{ }} 会被一并匹配——这是 Mustache 引擎实际的求值范围，与原行为一致.
    const regex = /\{\{([\s\S]+?)\}\}/g
    let match: RegExpExecArray | null
    while ((match = regex.exec(script)) !== null) {
      const expr = match[1].trim()
      const snippet = expr.length > 40 ? expr.slice(0, 40) + '...' : expr
      violations.push({
        type: 'MUSTACHE_EXPRESSION',
        message: `禁止在 JS 代码中使用 Mustache 表达式 {{${snippet}}}，请直接引用变量（如 ${
          snippet.split('.')[0] || 'variable'
        }）`,
        from: match.index,
        to: match.index + match[0].length
      })
    }
  } catch {}
}

/**
 * 正则检测静态 import 声明.
 * 逐行扫描，跳过注释行，匹配以 `import` 开头且非动态 import() 的语句.
 */
function checkStaticImports(script: string, violations: JsViolation[]): void {
  let lineStart = 0
  for (const line of script.split('\n')) {
    const trimmed = line.trim()
    const lineEnd = lineStart + line.length

    // 跳过单行注释和块注释起始行
    if (!trimmed.startsWith('//') && !trimmed.startsWith('/*')) {
      // 匹配静态 import（排除动态 import()）
      if (/^import\s/.test(trimmed) && !/^import\s*\(/.test(trimmed)) {
        const fromMatch = trimmed.match(/from\s+['"]([^'"]+)['"]/)
        const sideEffectMatch = trimmed.match(/^import\s+['"]([^'"]+)['"]/)
        const moduleName = fromMatch?.[1] || sideEffectMatch?.[1] || '未知模块'
        const leadingSpace = line.length - line.trimStart().length
        violations.push({
          type: 'REQUIRE_IMPORT',
          message: `禁止使用 import 语句（导入模块："${moduleName}"）`,
          from: lineStart + leadingSpace,
          to: lineEnd
        })
      }
    }

    lineStart = lineEnd + 1 // +1 for the \n
  }
}

/**
 * 递归遍历 AST 节点，检测违规模式.
 * @param offset AST 节点的 start/end 减去该值即为原 script 中的字符偏移.
 */
function walkNode(
  node: any,
  violations: JsViolation[],
  globalNames: Set<string>,
  offset: number
): void {
  if (!node || typeof node !== 'object') return

  switch (node.type) {
    case 'ImportExpression': {
      // 动态 import：import('xxx')
      const source = node.source
      const modulePath =
        source?.type === 'Literal' && typeof source.value === 'string'
          ? `"${source.value}"`
          : '动态路径'
      violations.push({
        type: 'REQUIRE_IMPORT',
        message: `禁止使用动态 import()（导入模块：${modulePath}）`,
        from: nodeFrom(node, offset),
        to: nodeTo(node, offset)
      })
      break
    }

    case 'CallExpression':
      checkCallExpression(node, violations, globalNames, offset)
      break

    case 'AssignmentExpression':
      checkAssignmentExpression(node, violations, globalNames, offset)
      break

    case 'UpdateExpression':
      checkUpdateExpression(node, violations, globalNames, offset)
      break
  }

  for (const key of Object.keys(node)) {
    if (key === 'type') continue
    const child = node[key]
    if (Array.isArray(child)) {
      child.forEach(item => {
        if (item && typeof item === 'object' && item.type) {
          walkNode(item, violations, globalNames, offset)
        }
      })
    } else if (child && typeof child === 'object' && child.type) {
      walkNode(child, violations, globalNames, offset)
    }
  }
}

function nodeFrom(node: any, offset: number): number {
  return Math.max(0, (node?.start ?? 0) - offset)
}
function nodeTo(node: any, offset: number): number {
  return Math.max(0, (node?.end ?? 0) - offset)
}

/**
 * 检查 CallExpression 节点：
 * - require('xxx') → 禁止
 * - globalName.setValue(...) / globalName.xxx.setValue(...) → 禁止
 * - globalName.trigger(...) / globalName.xxx.trigger(...) → 禁止
 */
function checkCallExpression(
  node: any,
  violations: JsViolation[],
  globalNames: Set<string>,
  offset: number
): void {
  const callee = node.callee

  // require('xxx')
  if (callee?.type === 'Identifier' && callee.name === 'require') {
    const firstArg = node.arguments?.[0]
    const modulePath =
      firstArg?.type === 'Literal' && typeof firstArg.value === 'string'
        ? `"${firstArg.value}"`
        : '动态路径'
    violations.push({
      type: 'REQUIRE_IMPORT',
      message: `禁止使用 require()（导入模块：${modulePath}）`,
      from: nodeFrom(node, offset),
      to: nodeTo(node, offset)
    })
    return
  }

  // xxx.setValue() / xxx.trigger()
  if (callee?.type !== 'MemberExpression') return

  const propName = getMemberPropertyName(callee)
  if (propName === 'setValue') {
    pushGlobalMethodViolation(
      node,
      callee,
      globalNames,
      violations,
      'SET_VALUE_CALL',
      'setValue',
      offset
    )
  } else if (propName === 'trigger') {
    pushGlobalMethodViolation(
      node,
      callee,
      globalNames,
      violations,
      'TRIGGER_CALL',
      'trigger',
      offset
    )
  }
}

/**
 * 若 MemberExpression 的根标识符落在全局变量集合中，则记录一条对应类型的违规.
 */
function pushGlobalMethodViolation(
  callNode: any,
  callee: any,
  globalNames: Set<string>,
  violations: JsViolation[],
  type: 'SET_VALUE_CALL' | 'TRIGGER_CALL',
  methodName: 'setValue' | 'trigger',
  offset: number
): void {
  const rootName = getRootIdentifierName(callee.object)
  if (!rootName || !globalNames.has(rootName)) return
  const callChain = getMemberExpressionString(callee)
  violations.push({
    type,
    message: `禁止调用 ${methodName} 方法（表达式：${callChain}(...)，变量 "${rootName}" 是全局变量）`,
    from: nodeFrom(callNode, offset),
    to: nodeTo(callNode, offset)
  })
}

/**
 * 检查 AssignmentExpression 节点：
 * - globalName = xxx → 禁止
 * - globalName.xxx = xxx → 禁止
 */
function checkAssignmentExpression(
  node: any,
  violations: JsViolation[],
  globalNames: Set<string>,
  offset: number
): void {
  const left = node.left

  if (left?.type === 'Identifier') {
    if (globalNames.has(left.name)) {
      violations.push({
        type: 'VARIABLE_ASSIGNMENT',
        message: `禁止对全局变量 "${left.name}" 直接赋值（${left.name} = ...）`,
        from: nodeFrom(node, offset),
        to: nodeTo(node, offset)
      })
    }
  } else if (left?.type === 'MemberExpression') {
    const rootName = getRootIdentifierName(left)
    if (rootName && globalNames.has(rootName)) {
      const fullPath = getMemberExpressionString(left)
      violations.push({
        type: 'VARIABLE_ASSIGNMENT',
        message: `禁止对全局变量 "${rootName}" 的属性赋值（${fullPath} = ...）`,
        from: nodeFrom(node, offset),
        to: nodeTo(node, offset)
      })
    }
  }
}

/**
 * 检查 UpdateExpression 节点：
 * - globalName++ / globalName-- → 禁止
 * - ++globalName / --globalName → 禁止
 */
function checkUpdateExpression(
  node: any,
  violations: JsViolation[],
  globalNames: Set<string>,
  offset: number
): void {
  const argument = node.argument

  if (argument?.type === 'Identifier') {
    if (globalNames.has(argument.name)) {
      const operator = node.operator
      const expr = node.prefix ? `${operator}${argument.name}` : `${argument.name}${operator}`
      violations.push({
        type: 'VARIABLE_ASSIGNMENT',
        message: `禁止对全局变量 "${argument.name}" 使用自增/自减运算符（${expr}）`,
        from: nodeFrom(node, offset),
        to: nodeTo(node, offset)
      })
    }
  } else if (argument?.type === 'MemberExpression') {
    const rootName = getRootIdentifierName(argument)
    if (rootName && globalNames.has(rootName)) {
      const operator = node.operator
      const fullPath = getMemberExpressionString(argument)
      const expr = node.prefix ? `${operator}${fullPath}` : `${fullPath}${operator}`
      violations.push({
        type: 'VARIABLE_ASSIGNMENT',
        message: `禁止对全局变量 "${rootName}" 的属性使用自增/自减运算符（${expr}）`,
        from: nodeFrom(node, offset),
        to: nodeTo(node, offset)
      })
    }
  }
}

/**
 * 获取 MemberExpression 的属性名.
 * 支持 obj.prop（Identifier）和 obj['prop']（Literal）.
 */
function getMemberPropertyName(node: any): string | undefined {
  if (!node || node.type !== 'MemberExpression') return undefined
  const prop = node.property
  if (!node.computed && prop?.type === 'Identifier') return prop.name
  if (node.computed && prop?.type === 'Literal' && typeof prop.value === 'string') return prop.value
  return undefined
}

/**
 * 递归展开 MemberExpression 链，返回最根部的 Identifier 名称.
 * e.g. state1.nested.deep → 'state1'
 */
function getRootIdentifierName(node: any): string | undefined {
  if (!node) return undefined
  if (node.type === 'Identifier') return node.name
  if (node.type === 'MemberExpression') return getRootIdentifierName(node.object)
  return undefined
}

/**
 * 递归重建 MemberExpression 链的完整表达式字符串.
 * e.g. state1.nested['deep'] → 'state1.nested["deep"]'
 */
function getMemberExpressionString(node: any): string {
  if (!node) return '?'
  if (node.type === 'Identifier') return node.name
  if (node.type === 'MemberExpression') {
    const obj = getMemberExpressionString(node.object)
    const prop = node.property
    if (!node.computed && prop?.type === 'Identifier') return `${obj}.${prop.name}`
    if (node.computed && prop?.type === 'Literal') return `${obj}[${JSON.stringify(prop.value)}]`
    return `${obj}[?]`
  }
  return '?'
}
// #endregion
