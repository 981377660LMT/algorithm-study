/**
 * @param {string[]} source
 * @return {string[]}
 */
const removeComments = (source: string[]): string[] => {
  // 注意 不可以加flag 's'即点匹配换行符   因为只有/* */允许多行匹配 而// 只允许单行匹配
  // 要用[\s\S]*?表示包括换行符的所有字符
  // 为了消除// 后所有字符 我们不能使用非贪婪匹配?
  const str = source.join('\n').replace(/(\/\*)([\s\S]*?)(\*\/)|\/\/(.*)/g, '')
  return str.split('\n').filter(v => v)
}

// 输出: ["int main()","{ ","  ","int a, b, c;","a = b + c;","}"]
console.log(
  removeComments([
    '/*Test program */',
    'int main()',
    '{ ',
    '  // variable declaration ',
    'int a, b, c;',
    '/* This is a test',
    '   multiline  ',
    '   comment for ',
    '   testing */',
    'a = b + c;',
    '}'
  ])
)

// 输出: ["ab"]
// console.log(removeComments(['a/*comment', 'line', 'more_comment*/b']))

// /* /* test */ */ 不属于正确的代码语法

export {}
