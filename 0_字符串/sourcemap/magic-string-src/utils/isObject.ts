const toString = Object.prototype.toString

export default function isObject(thing: unknown): boolean {
  return toString.call(thing) === '[object Object]'
}
