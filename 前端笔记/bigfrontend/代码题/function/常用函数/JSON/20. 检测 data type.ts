// 对于JavaScript中的所有基础数据类型，请实现一个方法进行检测。
/**
 * @param {any} data
 * @return {string}
 */
function detectType(data: any): string {
  // if (data instanceof FileReader) return 'object'
  return Object.prototype.toString.call(data).slice(1, -1).split(' ')[1].toLowerCase()
}

console.log(detectType(1))
