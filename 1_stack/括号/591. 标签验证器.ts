/**
 * @param {string} code
 * @return {boolean}
 */
const isValid = function (code: string): boolean {
  // 去除cdata
  while (/<!\[CDATA\[.*?\]\]>/.test(code)) {
    code = code.replace(/<!\[CDATA\[.*?\]\]>/g, '$TMP$')
  }
  console.log(code)
  while (/<([A-Z]{1,9})>([^<]*)<\/(\1)>/.test(code)) {
    code = code.replace(/<([A-Z]{1,9})>([^<]*)<\/(\1)>/g, '$VALID$')
  }

  return code === '$VALID$'
}

console.log(isValid('<DIV>This is the first line <![CDATA[<div>]]></DIV>'))
console.log(isValid('<DIV>>>  ![cdata[]] <![CDATA[<div>]>]]>]]>>]</DIV>'))

console.log(isValid('<A>  <B> </A>   </B>'))
console.log(isValid('<DIV>  div tag is not closed  <DIV>'))
console.log(isValid('<DIV>  unmatched <  </DIV>'))
console.log(isValid('<DIV> closed tags with invalid tag name  <b>123</b> </DIV>'))
console.log(
  isValid('<DIV> unmatched tags with invalid tag name  </1234567890> and <CDATA[[]]>  </DIV>')
)
console.log(isValid('<DIV>  unmatched start tag <B>  and unmatched end tag </C>  </DIV>'))
console.log(isValid('<![CDATA[ABC]]><TAG>sometext</TAG>'))

// 真的没必要
