### 1. 为什么需要 JSON5？（痛点分析）

传统的 JSON 非常轻量且易于机器解析，但作为**人工配置语言**时，有几个显著的痛点：

1. **不支持注释**：无法在配置中解释某个字段的作用。
2. **过于严格的引号**：键名必须显式使用双引号 `""`，即使它是合法的标识符。
3. **不支持尾逗号**：在数组或对象末尾添加、删除项时，经常因为多写或漏写逗号导致解析报错，这对版本控制（Git Diff）也不友好。
4. **数字表达受限**：不支持十六进制、不支持 `NaN` 或 `Infinity`。

`JSON5 诞生正是为了解决这些问题，它让配置文件的书写体验更接近于写普通的 JavaScript 对象。`

---

### 2. JSON5 的核心语法特性

JSON5 在 JSON 的基础上增加了以下 ES5 特性：

#### 2.1 对象（Objects）

- **允许键名不加引号**：只要键名是合法的 ES5 标识符。
- **支持单引号**：键名和字符串不仅可以使用双引号 `""`，也可以使用单引号 `''`。
- **支持尾逗号（Trailing commas）**。

#### 2.2 数组（Arrays）

- **支持尾逗号**。

#### 2.3 字符串（Strings）

- **支持单引号引用的字符串**。
- **支持多行字符串**：通过在行尾添加反斜杠 `\` 来实现字面上的多行折行。
- 允许转义字符。

#### 2.4 数字（Numbers）

- **支持十六进制**：如 `0x1A`。
- **支持小数点开头或结尾**：如 `.867` 或 `867.`。
- **支持正号**：如 `+1`。
- **支持 IEEE 754 特殊值**：支持 `Infinity`、`-Infinity` 和 `NaN`。

#### 2.5 注释（Comments）

- **支持单行注释**：`//`
- **支持多行注释**：`/* */`

---

### 3. JSON 与 JSON5 的代码对比

**标准 JSON：**

```json
{
  "name": "example",
  "version": "1.0.0",
  "description": "It is strictly formatted.",
  "dependencies": ["lodash", "react"],
  "price": 12.5
}
```

**JSON5 示例：**

```json5
{
  // 1. 键名可以不加引号
  name: 'example',

  // 2. 字符串支持单引号
  version: '1.0.0',

  /*
   * 3. 多行注释
   * 多行字符串
   */
  description: 'It is strictly \
formatted but easy to read.',

  // 4. 数组和对象支持尾部逗号 (有利于 git diff)
  dependencies: ['lodash', 'react'],

  // 5. 丰富的数字表示
  price: +12.5,
  hexColor: 0xff0000,
  infinityValue: Infinity,
  notANumber: NaN
}
```

---

### 4. 工具生态与使用 (Node.js / 浏览器)

原生 JavaScript 提供的是 `JSON.parse` 和 `JSON.stringify`，不支持 JSON5。你需要借助于三方库 [`json5`](https://json5.org/)。

#### 安装

```bash
npm install json5
```

#### 基本使用

它的 API 设计与原生 JSON 保持绝对一致：

```javascript
const JSON5 = require('json5')

const json5Str = `{
  foo: 'bar', // commented
  arr: [1, 2, ],
}`

// 解析 JSON5 字符串
const obj = JSON5.parse(json5Str)
console.log(obj.foo) // 'bar'

// 序列化回 JSON5（可以保留某些格式化特效）
const outStr = JSON5.stringify(obj, null, 2)
```

#### 在构建工具中使用

许多现代构建工具和框架已经在底层默认或可通过插件支持 JSON5：

- **Babel**：Babel 的配置文件常常使用 `.babelrc` (实际上是用 JSON5 解析的)。
- **Next.js / Nuxt / Vite** 等框架：在处理特定类型的配置时部分支持。
- **Webpack**：可以通过 `json5-loader` 去直接 import `.json5` 文件。

---

### 5. 配置文件格式大比拼（JSON 家族及相关）

在选用配置文件时，经常会拿 JSON5 与其他格式对比：

| 特性 / 格式            | 原生 JSON                    | JSONC (JSON with Comments)                    | JSON5                     | YAML                                | TOML                               |
| :--------------------- | :--------------------------- | :-------------------------------------------- | :------------------------ | :---------------------------------- | :--------------------------------- |
| **注释**               | ❌                           | ✅                                            | ✅                        | ✅                                  | ✅                                 |
| **尾逗号**             | ❌                           | ✅                                            | ✅                        | ❌(无需)                            | ✅                                 |
| **单引号/无引号键名**  | ❌ (仅双引号)                | ❌ (仅双引号)                                 | ✅                        | ✅                                  | ✅                                 |
| **多行字符串**         | ❌                           | ❌                                            | ✅ (反斜杠折行)           | ✅                                  | ✅                                 |
| **无穷大/NaN/Hex数字** | ❌                           | ❌                                            | ✅                        | ✅                                  | ✅                                 |
| **常见使用场景**       | 数据传输 (API)、package.json | VS Code 配置 (`settings.json`), tsconfig.json | 复杂配置文件，Babel配置等 | CI/CD (GitHub Actions), Docker, K8s | Rust (Cargo.toml)、Python 项目配置 |
| **解析速度复杂度**     | 极快 / 极低                  | 快 / 低                                       | 较快 / 低                 | 极慢 / 极高（因缩进敏感）           | 快 / 中                            |

- **JSONC**：微软主导的规范，仅仅是在严格的 JSON 基础上**放开了注释和尾逗号**。VS Code 的配置文件（如 `settings.json`）和 TypeScript 的 tsconfig.json 用的其实都是 JSONC，而不是完整的 JSON5（不支持单引号和无引号键）。
- **YAML**：功能异常强大，书写最为极简，但规则非常复杂（比如著名的“挪威国家代码 `NO` 被解析为 False”问题），对缩进极其敏感。
- **TOML**：语义清晰，天生为“配置文件”而生，在后端生态愈发流行。

### 6. 总结建议

1. **API 数据传输**：永远使用**原生 JSON**。不要使用 JSON5，因为各语言的原生库并不原生包含对 JSON5 的支持，强行使用会导致联调障碍。
2. **需要人类频繁干预的配置文件**：**推荐使用 JSON5** 或 JSONC。避免了传统 JSON 中的小坑（如多一个逗号使服务崩溃），并且保持了 JS 开发者极其熟悉的语法直觉。
3. **如果你只是在处理 TS/JS 项目**：将配置文件命名为 `.json5`，或者在普通无后缀名的 rc 文件中使用，它能显著提升基建配置的阅读体验。

---

- RPC 编辑器(JSEditor) 相关
  - 目标语法本质上是 JSON5 + {{ ... }} 作为一等表达式，无法被任何标准格式（JSON/JSON5/JS）原生解析

  ```js
  {
    "a": "1",
    b: 2,
    "c": {{sql1.value}},
    "d": '{{sql2.value}}'
  }
  ```

  - 采用字符串存储(request_body)，传给后端时正则替换运行时的值
  - 后段parse：整体感觉是将字符串当作JSON5解析
    单一正则 placeholder 替换 + JSON5.parse + JSON5 序列化时正则还原

    ```ts
    import JSON5 from 'json5'

    /**
     * 单次扫描同时匹配「字符串字面量 / 行注释 / 块注释 / `{{ ... }}`」。
     * 字符串和注释整段被吞掉、原样返回；`{{ ... }}` 只在出现在字符串/注释外时才被替换。
     * 这样即使字符串内含 `{{ ... }}`（如 `"{{1}}"`）也不会被破坏。
     如果只用最朴素的 /\{\{[\s\S]*?\}\}/g 全局替换，会把字符串字面量内部的 {{...}} 也替换掉
     */
    const TOKEN_RE =
      /"(?:\\.|[^"\\])*"|'(?:\\.|[^'\\])*'|\/\/[^\n]*|\/\*[\s\S]*?\*\/|\{\{[\s\S]*?\}\}/g

    /** 匹配 JSON 序列化后带引号的 sentinel 字符串。 */
    const SENTINEL_RE = /"__LANDER_EXPR_(\d+)__"/g

    /**
     * 解析 "JSON5 + `{{ ... }}` 表达式" 文本，返回顶层对象的扁平 KV 表。
     *
     * 步骤：
     *  1. tokenize：把 `{{ ... }}` 替换为 sentinel 字符串字面量，让文本变成合法 JSON5。
     *  2. parse：交给 `JSON5.parse`（无 `new Function`，无任意代码执行风险）。
     *  3. render：对每个值 `JSON.stringify`，让所有 sentinel 都带引号，再用一个正则把
     *     `"sentinel"` 整体替换回原始 `{{ ... }}`（去外层引号 → 表达式 token）。
     *     顶层字符串（json 仍以 `"` 开头）再 `JSON.parse` 剥掉一层 JSON 引号。
     */
    export const parseLanderJson = (src: string): Record<string, string> => {
      const placeholders: string[] = []
      const protectedText = src.replace(TOKEN_RE, m => {
        if (!m.startsWith('{{')) {
          return m
        }
        placeholders.push(m)
        return JSON.stringify(`__LANDER_EXPR_${placeholders.length - 1}__`)
      })

      let parsed: unknown
      try {
        parsed = JSON5.parse(protectedText)
      } catch (e) {
        throw new Error((e as Error)?.message || 'JSON5 解析失败')
      }
      if (!parsed || typeof parsed !== 'object' || Array.isArray(parsed)) {
        throw new Error('请输入对象格式（顶层必须是 `{ ... }`）')
      }

      const render = (v: unknown): string => {
        if (v === null || v === undefined) {
          return ''
        }
        const json = JSON.stringify(v).replace(SENTINEL_RE, (_, i) => placeholders[Number(i)] ?? '')
        return json.startsWith('"') ? JSON.parse(json) : json
      }

      const result: Record<string, string> = {}
      Object.entries(parsed as Record<string, unknown>).forEach(([k, v]) => {
        result[render(k)] = render(v)
      })
      return result
    }
    ```

---

## 核心难题

输入文本 **不是** 任何标准格式：

- 不是 JSON（key 可不加引号、可有单引号、可有注释）
- 不是 JSON5（多了个非法语法 `{{ ... }}`）
- 不是 JS（`{{x}}` 在 JS 里会被解析成「块语句包裹一个标识符」，语义完全错）

直接喂给任何现成 parser 都会报错。所以必须做**预处理**，把它变成合法格式。

## 思路：替身 + 还原（"sentinel" 模式）

把麻烦的 `{{ ... }}` **暂时换掉**，让文本变合法 → 用现成 parser 解析 → 解析完后**换回来**。

第 1 步把 `{{sql1.value}}` 替换成什么？必须满足：

1. 在 JSON5 中合法
2. 唯一可识别（不会和用户内容冲突）
3. 能记住"原本是什么"（按顺序记到数组里）

选择：**编号字符串字面量** `"__LANDER_EXPR_0__"`、`"__LANDER_EXPR_1__"`……

```
原文:  { c: {{sql1.value}}, d: {{sql2.x}} }
                ↓ 替换
中间:  { c: "__LANDER_EXPR_0__", d: "__LANDER_EXPR_1__" }
       (合法 JSON5！)
       placeholders = ["{{sql1.value}}", "{{sql2.x}}"]
```

JSON5 解析得到对象：`{ c: "__LANDER_EXPR_0__", d: "__LANDER_EXPR_1__" }`。最后遍历每个值，看见 `"__LANDER_EXPR_N__"` 就用 `placeholders[N]` 还原。

## 第一个坑：字符串内的 `{{ }}` 不该替换

如果用户写 `"page_size": "{{1}}"`，这个 `{{1}}` 在字符串字面量内部，是字面文本，不是表达式。

天真的全局正则 `/\{\{.*?\}\}/g` 不分上下文，会把它换成 `"page_size": ""__LANDER_EXPR_0__""` —— 双重引号，JSON5 直接报错。**这就是用户报的 bug。**

解决办法：让正则在替换前**先识别"字符串字面量"这个上下文**，遇到字符串就整段跳过。

## 关键技巧：alternation 正则 + 回调判断

```js
const TOKEN_RE = /"(?:\\.|[^"\\])*"|'(?:\\.|[^'\\])*'|\/\/[^\n]*|\/\*[\s\S]*?\*\/|\{\{[\s\S]*?\}\}/g
```

这是 5 个候选用 `|` 拼起来：

| 分支                 | 匹配什么                     |
| -------------------- | ---------------------------- |
| `"(?:\\.\|[^"\\])*"` | 双引号字符串（含 `\"` 转义） |
| `'(?:\\.\|[^'\\])*'` | 单引号字符串                 |
| `\/\/[^\n]*`         | 行注释                       |
| `\/\*[\s\S]*?\*\/`   | 块注释                       |
| `\{\{[\s\S]*?\}\}`   | 表达式                       |

**正则引擎从左到右扫描，每次只匹配最早出现的那一个候选**。所以遇到 `"page_size": "{{1}}"` 时：

1. 引擎走到 `"page_size"`，匹配到第 1 个分支，**整段** `"page_size"` 被吞掉
2. 跳过冒号空格
3. 走到 `"{{1}}"`，又匹配第 1 个分支，**整段** `"{{1}}"` 被吞掉（连同里面的 `{{1}}`）
4. 第 5 个分支 `\{\{...\}\}` 根本没机会匹配字符串内的 `{{1}}`

回调里只对 `{{...}}` 分支做替换，其它分支原样返回：

```js
src.replace(TOKEN_RE, matched => {
  if (!matched.startsWith('{{')) return matched // 字符串/注释原样
  placeholders.push(matched)
  return JSON.stringify(`__LANDER_EXPR_${placeholders.length - 1}__`)
})
```

这其实就是一个**最小词法分析器**——只是借用正则引擎实现，没自己写状态机。

## 第二个坑：嵌套对象里的还原

顶层 `{ c: {{x}} }` 简单：解析得 `c: "__LANDER_EXPR_0__"`，是个字符串，整体替换即可。

但 `{ a: { x: {{x}}, y: 1 } }` 里 `{{x}}` 嵌在子对象中。我们要把整个子对象作为**一段字符串**塞进 KVItem.value，最终结果应该是：

```
{ key: "a", value: '{"x":{{x}},"y":1}' }
                       ↑ 注意：{{x}} 没有引号，是表达式 token
```

代码里这一步：

```js
value = restoreInJsonString(JSON.stringify(v))
```

- `v = { x: "__LANDER_EXPR_0__", y: 1 }`
- `JSON.stringify(v)` → `'{"x":"__LANDER_EXPR_0__","y":1}'`
- `restoreInJsonString` 用正则 `/"__LANDER_EXPR_(\d+)__"/g` 把**带引号**的 sentinel 整体替换为原始 `{{x}}`（去掉外层引号 → 变回表达式 token）
- 结果：`'{"x":{{x}},"y":1}'`

为什么必须连引号一起替换？因为在 JSON 序列化里，普通字符串是带引号的（`"hello"`），而我们想让 `{{x}}` 变成"裸 token"——所以替换正则必须包含外层 `"..."`，一并消掉。

## 整个流程完整跑一遍

输入：

```
{ c: {{sql1.value}}, a: { x: {{foo}} }, p: "{{1}}" }
```

**Step 1: tokenize 替换**

- 走第 5 分支：`{{sql1.value}}` → `"__LANDER_EXPR_0__"`，记录 `placeholders[0]`
- 走第 5 分支：`{{foo}}` → `"__LANDER_EXPR_1__"`，记录 `placeholders[1]`
- 走第 1 分支：`"{{1}}"` 整段保留，`{{1}}` 不替换
- 中间文本：`{ c: "__LANDER_EXPR_0__", a: { x: "__LANDER_EXPR_1__" }, p: "{{1}}" }`

**Step 2: JSON5.parse**

```js
{ c: '__LANDER_EXPR_0__', a: { x: '__LANDER_EXPR_1__' }, p: '{{1}}' }
```

**Step 3: 遍历顶层 entries**

- `c`：值是字符串，匹配 sentinel → 还原 `{{sql1.value}}`
- `a`：值是对象 → `JSON.stringify` 得 `'{"x":"__LANDER_EXPR_1__"}'` → 正则替换带引号 sentinel → `'{"x":{{foo}}}'`
- `p`：值是字符串 `'{{1}}'`，**不匹配 sentinel 正则**（sentinel 是 `__LANDER_EXPR_N__` 格式）→ 原样

最终：

```js
{
  c: '{{sql1.value}}',
  a: '{"x":{{foo}}}',
  p: '{{1}}'
}
```

## 为什么这么做就够好

| 问题                          | 解决方式                                                                                    |
| ----------------------------- | ------------------------------------------------------------------------------------------- |
| 任意代码执行风险              | `JSON5.parse` 替代 `new Function`——纯数据解析，零副作用                                     |
| 字符串内 `{{` 误匹配          | alternation 正则让字符串/注释 token **先于** `{{` 匹配                                      |
| 嵌套对象的表达式 token 序列化 | 用"带引号 sentinel"模式，反向替换时一并消掉引号                                             |
| 不重复造轮子                  | JSON5 处理所有 unquoted key/单引号/注释/尾逗号；正则引擎处理 tokenize；自己只写约 30 行胶水 |

核心思想就一句话：**用正则的 alternation 实现"上下文敏感的替换"，把非法语法暂时换成合法占位符，让现成解析器去干重活，最后按记录还原。**
