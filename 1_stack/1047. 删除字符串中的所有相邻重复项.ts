/**
 * @param {string} s
 * @return {string}
 */
// var removeDuplicates = function (s) {
//   const duplicateRegexp = /(\w)\1+/g
//   const remove = s => {
//     const shouldRemove = duplicateRegexp.test(s)
//     return shouldRemove ? remove(s.replace(duplicateRegexp, '')) : s
//   }
//   return remove(s)
// }
var removeDuplicates = function (s: string): string {
  const stack: string[] = []
  for (const char of s) {
    stack[stack.length - 1] === char ? stack.pop() : stack.push(char)
  }

  return stack.join('')
}

console.log(removeDuplicates('abbaca'))
console.log(removeDuplicates('aaaaaa'))
export {}
