// 请编写代码检查s2是否为s1旋转而成
function isFlipedString(s1: string, s2: string): boolean {
  // 两个都一样
  return s1.length === s2.length && s1.repeat(2).includes(s2)
  // return s1.length === s2.length && s2.repeat(2).includes(s1)
}

console.log(isFlipedString('waterbottle', 'erbottlewat'))
