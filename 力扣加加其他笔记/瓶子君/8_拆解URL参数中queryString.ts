const url = 'http://sample.com/?a=1&b=2&c=xx&d=2#hash'
const queryString = (str: string) => {
  const obj = Object.create(null)
  str.replace(/([^?&=]+)=([^&]+)/g, (_, k, v) => (obj[k] = v))
  return obj
}

console.log(queryString(url))

export {}
