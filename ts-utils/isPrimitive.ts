/**
 * 判断一个对象是否为基本类型.
 */
function isPrimitive(
  o: unknown
): o is number | string | boolean | symbol | bigint | null | undefined {
  return o === null || (typeof o !== 'object' && typeof o !== 'function')
}

export { isPrimitive }
