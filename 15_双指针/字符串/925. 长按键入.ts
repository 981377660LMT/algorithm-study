/**
 * @param {string} name
 * @param {string} typed
 * @return {boolean}
 * 你的朋友正在使用键盘输入他的名字 name。偶尔，在键入字符 c 时，按键可能会被长按，而字符可能被输入 1 次或多次。
   你将会检查键盘输入的字符 typed。如果它对应的可能是你的朋友的名字（其中一些字符可能被长按），那么就返回 True。
 */
const isLongPressedName = function (name: string, typed: string): boolean {
  const match1 = Array.from(name.matchAll(/(\w)\1*/g)).map(item => item[0])
  const match2 = Array.from(typed.matchAll(/(\w)\1*/g)).map(item => item[0])
  if (match1.length !== match2.length) return false
  for (let i = 0; i < match1.length; i++) {
    if (match1[i].length > match2[i].length || match1[i][0] !== match2[i][0]) return false
  }
  return true
}

// const isLongPressedName = function (name: string, typed: string): boolean {
//   let i = 0
//   let j = 0
//   while (j < typed.length) {
//     if (i < name.length && name[i] === typed[j]) i++
//     else if()
//   }
// }

console.log(isLongPressedName('alex', 'aaleex'))
console.log(isLongPressedName('saeed', 'ssaaedd'))
