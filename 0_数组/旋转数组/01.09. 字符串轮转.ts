// 请编写代码检查s2是否为s1旋转而成
function isFlipedString(s1: string, s2: string): boolean {
  return s1.length === s2.length && s1.repeat(2).includes(s2)
}

console.log(isFlipedString('waterbottle', 'erbottlewat'))
