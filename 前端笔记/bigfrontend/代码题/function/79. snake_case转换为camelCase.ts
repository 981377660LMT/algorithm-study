import _ from 'lodash'

/**
 * @param {string} str
 * @return {string}
 */
function snakeToCamel(str: string): string {
  const snakePattern = /([^_])_([^_])/g
  // 先匹配到位置 再转换
  // s I
  // S o
  // r A
  return str.replace(snakePattern, (_, g1: string, g2: string) => `${g1}${g2.toUpperCase()}`)
}

// 连续的下划线__，打头的下划线 _a和结尾的下划线a_需要被保留。
// console.log(snakeToCamel('_double__underscore_'))
// '_double__underscore_'
console.log(snakeToCamel('is_IOS_or_Android'))
// 'isIOSOrAndroid'
console.log(_.camelCase('is_IOS_or_Android'))
export { snakeToCamel }
