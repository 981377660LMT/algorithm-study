const a = getUrlParams('a', 'http://lucifer.ren?a=1&b=2&a=3')
// const b = getUrlParams('b', 'http://lucifer.ren?a=1&b=2&a=3')
// const c = getUrlParams('c', 'http://lucifer.ren?a=1&b=2&a=3')

console.log(a)
// console.log(b)
// console.log(c)
// 查找key是否在url的查询字符串中， 如果在就返回，如果不在返回null，如果存在多个就返回数组。
function getUrlParams(key: string, href: string) {
  const match = [...href.matchAll(/([^?&=]+)=([^&]+)/g)]
  const res = match.filter(g => g[1] === key).map(g => g[2])
  return res.length ? res : null
}

export {}
