/**
 * @param {string} str
 * @return {Generator}
 * 输入的字符串只包括非负整数字符和+、-、 *、 /、 (、 ) 和空格，空格需要忽略。
 */
function* tokenize(str: string): Generator<string> {
  const tokens = str.match(/(\d+)|[\+\-\*\/\(\)]/g)
  if (tokens) {
    for (const token of tokens) {
      yield token
    }
  }
}

const tokens = tokenize(' 1 * (20 -   300      ) ')
for (let token of tokens) {
  console.log(token)
}
