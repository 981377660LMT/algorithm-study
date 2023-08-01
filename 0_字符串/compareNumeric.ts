// 字符串数字比大小

/**
 * @param str1 字符串数字
 * @param str2 字符串数字
 * @returns 如果str1大于str2则返回true
 */
function compareNumeric(str1: string, str2: string): boolean {
  return str1.length === str2.length ? str1 > str2 : str1.length > str2.length
}

if (require.main === module) {
  console.log(compareNumeric('9', '18'))
  console.log(compareNumeric('91', '18'))
}

export { compareNumeric }
