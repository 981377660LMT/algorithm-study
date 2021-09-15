// 编写一种方法，将字符串中的空格全部替换为%20。
// 假定该字符串尾部有足够的空间存放新增字符，并且知道字符串的“真实”长度。
function replaceSpaces(S: string, length: number): string {
  // return encodeURI(S.substring(0, length))
  return S.slice(0, length).replace(/\s/g, '%20')
}

console.log(replaceSpaces(' asa  as ', 9))
