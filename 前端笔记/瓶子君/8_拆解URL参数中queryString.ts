import { URL } from 'url'

const url = 'http://sample.com/?a=1&b=2&c=xx&d=2#hash'
const queryString1 = (str: string) => {
  const obj = Object.create(null)
  str.replace(/([^?&=]+)=([^&]+)/g, (_, g1, g2) => (obj[g1] = g2))
  return obj
}

// 解析params参数
const getParamByKey = (url: string, key: string) => {
  const params = new URL(url).searchParams
  return Object.fromEntries(params.entries())[key]
}

console.log(queryString1(url))
console.log(getParamByKey('http://sample.com/?a=1&b=2&c=xx&d=2#hash', 'a'))

export {}
