/**
 * @param {string[]} source
 * @return {string[]}
 */
const removeComments = function (source: string[]): string[] {}

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
    '}',
  ])
)

// 输出: ["int main()","{ ","  ","int a, b, c;","a = b + c;","}"]
console.log(removeComments(['a/*comment', 'line', 'more_comment*/b']))
// 输出: ["ab"]

// /* /* test */ */ 不属于正确的代码语法
