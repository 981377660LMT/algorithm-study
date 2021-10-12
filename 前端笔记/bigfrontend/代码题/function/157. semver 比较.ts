// semver:语义化版本号
// Semver是一个专门分析Semantic Version（语义化版本）的工具，
// “semver”其实就是这两个单词的缩写。
// Npm使用了该工具来处理版本相关的工作。
// X 是主版本号、Y 是次版本号、而Z 为修订号。

function compare(v1: string, v2: string): 0 | 1 | -1 {
  const getVersion = (str: string) => str.split('.').map(Number)

  const version1 = getVersion(v1)
  const version2 = getVersion(v2)
  for (let i = 0; i < 3; i++) {
    if (version1[i] > version2[i]) return 1
    if (version1[i] < version2[i]) return -1
  }

  return 0
}

console.log(compare('12.1.0', '12.0.9'))
console.log(compare('12.1.0', '12.1.2'))
console.log(compare('5.0.1', '5.0.1'))
