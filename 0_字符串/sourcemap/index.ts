// https://juejin.cn/post/7299477531640561727#heading-9
// magic-string是一个用于处理字符串的JavaScript库。它可以让你在字符串中进行插入、删除、替换等操作，并且能够生成准确的sourcemap
// 这个库特别适用于需要对源代码进行轻微修改并保存sourcemap的情况，比如替换字符、添加内容等操作。通过 magic-string，你可以确保在字符串操作的同时，sourcemap能够保持准确，不会因为操作而失真。

import MagicString from 'magic-string'
import fs from 'fs'
import path from 'path'

const s = new MagicString(
  `
  The quick brown fox
  jumped over 
  the lazy dog.
  `
)

s.update(0, 8, 'answer')
s.toString() // 'answer = 99'

s.update(11, 13, '42') // character indices always refer to the original string
s.toString() // 'answer = 42'

s.prepend('var ').append(';') // most methods are chainable
s.toString() // 'var answer = 42;'

const map = s.generateMap({
  source: 'source.js',
  file: 'converted.js.map',
  includeContent: true
}) // generates a v3 sourcemap

// 使用数组形式的原始映射生成 sourcemap 对象，而不是 base64vlq编码后的数组
console.dir(s.generateDecodedMap(), { depth: null })

// 在字符串的每一行前面加上前缀
s.indent('haha', { exclude: [0, 1] })

console.log(map.toUrl()) // !inline -> 将.map 作为 DataURI 嵌入，不单独生成.map 文件

fs.writeFileSync(path.resolve(__dirname, './converted.js'), s.toString())
fs.writeFileSync(path.resolve(__dirname, './converted.js.map'), map.toString())
