// 非常酷的类型编程
// https://playfulprogramming.com/posts/fun-with-types
// https://github.com/ronami/meta-typing
// https://github.com/ronami/HypeScript
// https://github.com/codemix/ts-sql
// https://www.learningtypescript.com/articles/extreme-explorations-of-typescripts-type-system

type ExtractFromDelimiters<S extends string, L extends string, R extends string> = string extends S
  ? string[] // 如果字符串不是字符串字面量，类型不会出错
  : S extends ''
  ? []
  : S extends `${infer _Head}${L}${infer U}${R}${infer Tail}`
  ? [U, ...ExtractFromDelimiters<Tail, L, R>]
  : []

type Test = ExtractFromDelimiters<'Hello, {{name}}!', '{{', '}}'> // Expected: ['name']
type Test2 = ExtractFromDelimiters<'{{name}} {{age}}', '{{', '}}'> // Expected: ['name', 'age']

const template = 'Hello, {{name}}! Your age is {{age}}.'
type Extracted = ExtractFromDelimiters<typeof template, '{{', '}}'> // Expected: ['name', 'age']

export {}
