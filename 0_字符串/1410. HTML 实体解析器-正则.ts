// 1410. HTML 实体解析器
// https://leetcode.cn/problems/html-entity-parser/description/
// !将字符实体转换为对应的特殊字符
//
// 「HTML 实体解析器」 是一种特殊的解析器，它将 HTML 代码作为输入，并用字符本身替换掉所有这些特殊的字符实体。
//
// HTML 里这些特殊字符和它们对应的字符实体包括：
//
// 双引号：字符实体为 &quot; ，对应的字符是 " 。
// 单引号：字符实体为 &apos; ，对应的字符是 ' 。
// 与符号：字符实体为 &amp; ，对应对的字符是 & 。
// 大于号：字符实体为 &gt; ，对应的字符是 > 。
// 小于号：字符实体为 &lt; ，对应的字符是 < 。
// 斜线号：字符实体为 &frasl; ，对应的字符是 / 。
// 给你输入字符串 text ，请你实现一个 HTML 实体解析器，返回解析器解析后的结果。

const RULES = [
  {
    entityRegexp: /&quot;/g,
    character: '"'
  },
  {
    entityRegexp: /&apos;/g,
    character: "'"
  },
  {
    entityRegexp: /&gt;/g,
    character: '>'
  },
  {
    entityRegexp: /&lt;/g,
    character: '<'
  },
  {
    entityRegexp: /&frasl;/g,
    character: '/'
  },

  // &的替换需要放在最后
  {
    entityRegexp: /&amp;/g,
    character: '&'
  }
]

function entityParser(text: string): string {
  for (const rule of RULES) {
    const { entityRegexp, character } = rule
    text = text.replace(entityRegexp, character)
  }

  return text
}

console.log(entityParser('&amp; is an HTML entity but &ambassador; is not.'))

export {}
