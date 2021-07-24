// 请你找出一个最大整数 m ，以满足 str = [str2, m] 可以从 str1 获得(子序列)。
const getMaxRepetion = (s1: string, n1: number, s2: string, n2: number) => {
  // const a = s1.repeat(n1)

  let index = 0
  let heat = 0

  for (let i = 0; i < n1; i++) {
    for (const letter of s1) {
      if (letter === s2[index]) {
        index++
      }

      if (index === s2.length) {
        index = 0
        heat++
      }
    }
  }

  return Math.floor(heat / n2)
}

console.log(getMaxRepetion('acb', 4, 'ab', 32))
// 输出：2
export {}
