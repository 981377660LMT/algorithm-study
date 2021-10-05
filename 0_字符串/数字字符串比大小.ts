function compare(str1: string, str2: string) {
  return str1.length === str2.length ? str1 > str2 : str1.length > str2.length
}

// 数组的隐式转换会先调用join
// console.log(['128', '9'] > ['129', '0'])
// 即'1289'>'1290'

// 怎么判断相等
// console.log(!(['a', 'b'] > ['a', 'b'])&&!(['a', 'b'] < ['a', 'b']))
