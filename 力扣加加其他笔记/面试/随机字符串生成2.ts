// 随机生成一个长度为 8 的字符串，要求只能是小写字母和数字，字母和数字可重复，但是生成的随机字符串不能重复
const useGetRandomStr = () => {
  const visited = new Set<string>()

  function* genStr(length: number) {
    const nums = Array.from({ length: 10 }, (_, i) => i)
    const lowercases = Array.from({ length: 26 }, (_, i) => String.fromCodePoint(97 + i))
    // // @ts-ignore
    // console.log({
    //   ...Object.fromEntries(nums.entries()),
    //   ...Object.fromEntries(lowercases.entries()),
    // })
    const map = [...nums, ...lowercases]
    for (let i = 0; i < length; i++) {
      yield map[~~(Math.random() * 36)]
    }
  }

  const randomStr = (length: number, failbackTimes: number = Infinity) => {
    let res = [...genStr(length)].join('')

    for (let i = 0; i < failbackTimes; i++) {
      if (!visited.has(res)) break
      res = [...genStr(length)].join('')
    }

    visited.add(res)
    return res
  }

  return randomStr
}

export {}

const getRandomStr = useGetRandomStr()
console.log(getRandomStr(8))
