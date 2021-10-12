// 如果你在项目中使用了css-loader，
// 你可以像这样通过localIdentName来变换class name。
// localIdentName: "[path][name]__[local]--[hash:base64:5]",

// 52进制转换
class ClassNameGenerator {
  private static readonly chars = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ'
  private position: number

  constructor() {
    this.position = 1
  }

  next() {
    const sb: string[] = []
    let pos = this.position
    this.position++

    while (pos) {
      const [div, mod] = [~~((pos - 1) / 52), (pos - 1) % 52]
      pos = div
      sb.push(ClassNameGenerator.chars[mod])
    }

    return sb.reverse().join('')
  }

  reset() {
    this.position = 1
  }
}

const classNameGenerator = new ClassNameGenerator()
// 仅使用字母: a - z , A - Z
function getUniqueClassName(): string {
  // your code here
  return classNameGenerator.next()
}

getUniqueClassName.reset = function () {
  // your code here
  classNameGenerator.reset()
}

export {}

if (require.main === module) {
  console.log(getUniqueClassName())
  console.log(getUniqueClassName())
  console.log(getUniqueClassName())
  console.log(getUniqueClassName())
  // 'a'

  getUniqueClassName()
  // 'b'

  getUniqueClassName()
  // 'c'

  // skip cases till 'Y'

  getUniqueClassName()
  // 'Z'

  getUniqueClassName()
  // 'aa'

  getUniqueClassName()
  // 'ab'

  getUniqueClassName()
  // 'ac'

  // skip more cases

  getUniqueClassName()
  // 'ZZ'

  getUniqueClassName()
  // 'aaa'

  console.log(getUniqueClassName())
  // 'aab'

  getUniqueClassName()
  // 'aac'

  getUniqueClassName.reset()

  console.log(getUniqueClassName())
  // 'a'
}
