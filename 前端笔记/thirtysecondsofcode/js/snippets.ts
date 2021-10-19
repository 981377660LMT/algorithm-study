import { Stream } from 'stream'
import { Func } from '../../bigfrontend/代码题/typings'

export {}

// toHSLObject
// 色相，饱和度，亮度
const toHSLObject = (hslStr: string) => {
  const [hue, saturation, lightness] = (hslStr.match(/\d+/g) ?? []).map(Number)
  return { hue, saturation, lightness }
}
toHSLObject('hsl(50, 10%, 10%)') // { hue: 50, saturation: 10, lightness: 10 }

const parseCookie = (cookie: string): Record<string, string> => {
  return cookie
    .split(';')
    .map(str => str.split('='))
    .reduce<Record<string, string>>((obj, entry) => {
      obj[decodeURIComponent(entry[0].trim())] = decodeURIComponent(entry[1].trim())
      return obj
    }, {})
}
parseCookie('foo=bar; equation=E%3Dmc%5E2')
// { foo: 'bar', equation: 'E=mc^2' }

const pipeAsyncFunctions =
  (...fns: Func[]) =>
  (arg: any) =>
    fns.reduce((p, f) => p.then(f), Promise.resolve(arg))
const sum = pipeAsyncFunctions(
  x => x + 1,
  x => new Promise(resolve => setTimeout(() => resolve(x + 2), 1000)),
  x => x + 3,
  async x => (await x) + 4
)
;(async () => {
  console.log(await sum(5)) // 15 (after one second)
})()

const queryStringToObject = (url: string) =>
  [...new URLSearchParams(url.split('?')[1])].reduce((a, [k, v]) => ((a[k] = v), a), {})
queryStringToObject('https://google.com?page=1&count=10')
// {page: '1', count: '10'}

const formatSeconds = (s: number) => {
  const [hour, minute, second, sign] =
    s > 0 ? [s / 3600, (s / 60) % 60, s % 60, ''] : [-s / 3600, (-s / 60) % 60, -s % 60, '-']

  return sign + [hour, minute, second].map(v => `${Math.floor(v)}`.padStart(2, '0')).join(':')
}
formatSeconds(200) // '00:03:20'
formatSeconds(-200) // '-00:03:20'
formatSeconds(99999) // '27:46:39'

const getSiblings = (el: HTMLHeadElement) =>
  [...el.parentNode.childNodes].filter(node => node !== el)
getSiblings(document.querySelector('head')!) // ['body']

// 欧几里得距离
const euclideanDistance = (a: number[], b: number[]) =>
  // hypot:计算一直角三角形的斜边长度
  Math.hypot(...Object.keys(a).map(k => b[Number(k)] - a[Number(k)]))
euclideanDistance([1, 1], [2, 3]) // ~2.2361
euclideanDistance([1, 1, 1], [2, 3, 2]) // ~2.4495

const cycleGenerator = function* <T>(arr: T[]) {
  let i = 0
  while (true) {
    yield arr[i % arr.length]
    i++
  }
}

const binaryCycle = cycleGenerator([0, 1])
binaryCycle.next() // { value: 0, done: false }

const isAnagram = (str1: string, str2: string) => {
  const normalize = (str: string) =>
    str
      .toLowerCase()
      .replace(/[^a-z0-9]/gi, '')
      .split('')
      .sort()
      .join('')
  return normalize(str1) === normalize(str2)
}
isAnagram('iceman', 'cinema') // true

const findLastKey = (obj: Record<string, any>, fn: (...args: any[]) => any) =>
  Object.keys(obj)
    .reverse()
    .find(key => fn(obj[key], key, obj))
findLastKey(
  {
    barney: { age: 36, active: true },
    fred: { age: 40, active: false },
    pebbles: { age: 1, active: true },
  },
  x => x['active']
) // 'pebbles'

const objectToQueryString = (queryParameters: Record<string, any>): string => {
  return queryParameters
    ? Object.entries(queryParameters).reduce((pre, [key, val], index) => {
        const symbol = pre.length === 0 ? '?' : '&'
        pre += typeof val === 'string' ? `${symbol}${key}=${val}` : ''
        return pre
      }, '')
    : ''
}
objectToQueryString({ page: '1', size: '2kg', key: undefined })
// '?page=1&size=2kg'

const isPalindrome = (str: string) => {
  const s = str.toLowerCase().replace(/[\W_]/g, '')
  return s === [...s].reverse().join('')
}
isPalindrome('taco cat') // true

const words = (str: string, pattern = /[^a-zA-Z-]+/) => str.split(pattern).filter(Boolean)
words('I love javaScript!!') // ['I', 'love', 'javaScript']
words('python, javaScript & coffee') // ['python', 'javaScript', 'coffee']

// @ts-ignore
const isDateValid = (...val: any[]) => !Number.isNaN(new Date(...val).valueOf())
isDateValid('December 17, 1995 03:24:00') // true
isDateValid('1995-12-17T03:24:00') // true
isDateValid('1995-12-17 T03:24:00') // false
isDateValid('Duck') // false
isDateValid(1995, 11, 17) // true
isDateValid(1995, 11, 17, 'Duck') // false
isDateValid({}) // false

const isGeneratorFunction = (val: any): val is Generator<any, any, any> =>
  Object.prototype.toString.call(val) === '[object GeneratorFunction]'
isGeneratorFunction(function () {}) // false
isGeneratorFunction(function* () {}) // true

const isPromiseLike = (obj: any): obj is PromiseLike<any> =>
  obj !== null &&
  (typeof obj === 'object' || typeof obj === 'function') &&
  typeof obj.then === 'function'

isPromiseLike({
  then: function () {
    return ''
  },
}) // true
isPromiseLike(null) // false
isPromiseLike({}) // false

const isAsyncFunction = (val: any): val is (...args: any[]) => Promise<any> =>
  Object.prototype.toString.call(val) === '[object AsyncFunction]'
isAsyncFunction(function () {}) // false
isAsyncFunction(async function () {}) // true

const isArrayLike = (obj: any): obj is ArrayLike<any> =>
  obj != null && typeof obj[Symbol.iterator] === 'function'
isArrayLike([1, 2, 3]) // true
isArrayLike(document.querySelectorAll('.className')) // true
isArrayLike('abc') // true
isArrayLike(null) // false

const isBrowser = () => ![typeof window, typeof document].includes('undefined')
isBrowser() // true (browser)
isBrowser() // false (Node)

const isNode = () => typeof process !== 'undefined' && !!process.versions && !!process.versions.node
isNode() // true (Node)
isNode() // false (browser)

const isStream = (val: any) =>
  val !== null && typeof val === 'object' && typeof val.pipe === 'function'

// 返回第一个参数为 true 的合并(combined)函数。
const coalesceFactory =
  (valid: (v: any) => boolean) =>
  (...args: any[]) =>
    args.find(valid)
const customCoalesce = coalesceFactory((v: any) => ![null, undefined, '', NaN].includes(v))
customCoalesce(undefined, null, NaN, '', 'Waldo') // 'Waldo'

// unix 时间精确到 s
const getTimestamp = (date = new Date()) => Math.floor(date.getTime() / 1000)
getTimestamp() // 1602162242

// 在字符串中规范行结束。
const normalizeLineEndings = (str: string, normalized = '\r\n') => str.replace(/\r?\n/g, normalized)
normalizeLineEndings('This\r\nis a\nmultiline\nstring.\r\n')
// 'This\r\nis a\r\nmultiline\r\nstring.\r\n'
normalizeLineEndings('This\r\nis a\nmultiline\nstring.\r\n', '\n')
// 'This\nis a\nmultiline\nstring.\n'
// 同理
const splitLines = (str: string) => str.split(/\r?\n/)
splitLines('This\nis a\nmultiline\nstring.\n')
// ['This', 'is a', 'multiline', 'string.' , '']

const isalnum = (str: string) => /^[a-z0-9]+$/gi.test(str)
isalnum('hello123') // true
isalnum('123') // true
isalnum('hello 123') // false (space character is not alphanumeric)
isalnum('#$hello') // false

const serializeCookie = (name: string, val: string) =>
  `${encodeURIComponent(name)}=${encodeURIComponent(val)}`
serializeCookie('foo', 'bar') // 'foo=bar'

// 从字符串中移除 HTML/XML 标记。
const stripHTMLTags = (str: string) => str.replace(/<[^>].*?>/g, '')
stripHTMLTags('<p><em>lorem</em> <strong>ipsum</strong></p>') // 'lorem ipsum'

// isValidJSON
const isValidJSON = (str: any) => {
  try {
    JSON.parse(str)
    return true
  } catch (e) {
    return false
  }
}
isValidJSON('{"name":"Adam","age":20}') // true
isValidJSON('{"name":"Adam",age:"20"}') // false
isValidJSON(null) // true

const untildify = (str: string) => str.replace(/^~($|\\|\/)/, `${require('os').homedir()}$1`)
// 将波浪线路径转换为绝对路径。
untildify('~/node') // '/Users/aUser/node'

const getType = (v: any) =>
  v === undefined ? 'undefined' : v === null ? 'null' : v.constructor.name
getType(new Set([1, 2, 3])) // 'Set'

const isPlainObject = (val: any) => !!val && typeof val === 'object' && val.constructor === Object
isPlainObject({ a: 1 }) // true
isPlainObject(new Map()) // false

const isPrimitive = (val: any) => Object(val) !== val
isPrimitive(null) // true
isPrimitive(undefined) // true
isPrimitive(50) // true
isPrimitive('Hello!') // true
isPrimitive(false) // true
isPrimitive(Symbol()) // true
isPrimitive([]) // false
isPrimitive({}) // false

const isObject = (obj: any): obj is object => obj === Object(obj)
isObject([1, 2, 3, 4]) // true
isObject([]) // true
isObject(['Hello!']) // true
isObject({ a: 1 }) // true
isObject({}) // true
isObject(true) // false

const isTravisCI = () => 'TRAVIS' in process.env && 'CI' in process.env
isTravisCI() // true (if code is running on Travis CI)

const round = (n: number, decimals = 0) => Number(`${Math.round(`${n}e${decimals}`)}e-${decimals}`)
round(1.005, 2) // 1.01

const cloneRegExp = (regExp: RegExp) => new RegExp(regExp.source, regExp.flags)
const regExp = /lorem ipsum/gi
const regExp2 = cloneRegExp(regExp) // regExp !== regExp2

// isLocalStorageEnabled
const isLocalStorageEnabled = () => {
  try {
    const key = `__storage__test`
    window.localStorage.setItem(key, '')
    window.localStorage.removeItem(key)
    return true
  } catch (e) {
    return false
  }
}
isLocalStorageEnabled() // true, if localStorage is accessible

// 检查是否支持触摸事件
const supportsTouchEvents = () => window && 'ontouchstart' in window
supportsTouchEvents() // true

const logBase = (n: number, base: number) => Math.log(n) / Math.log(base)
logBase(10, 10) // 1
logBase(100, 10) // 2

// 检查给定的元素是否被聚焦
const elementIsFocused = (el: Element) => el === document.activeElement
elementIsFocused(document.querySelector('#app')!) // true if the element is focused

// 检查页面的浏览器选项卡是否聚焦
const isBrowserTabFocused = () => !document.hidden
isBrowserTabFocused() // true

const toSafeInteger = (num: number) => {
  return Math.round(Math.max(Math.min(num, Number.MAX_SAFE_INTEGER), Number.MIN_SAFE_INTEGER))
}
toSafeInteger(Infinity) // 9007199254740991
