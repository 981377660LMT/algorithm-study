//  短网址
//  TinyURL是一种URL简化服务， 比如：
//  当你输入一个URL https://leetcode.com/problems/design-tinyurl 时，
//  它将返回一个简化的URL http://tinyurl.com/4e9iAk.

//  你的加密和解密算法如何设计和运作是没有限制的，
//  你只需要保证一个URL可以被加密成一个TinyURL，
//  并且这个TinyURL可以用解密方法恢复成原本的URL。
//////////////////////////////////////////////////////////////////////////////////////////////////
// 短网址原理
// 当我们在浏览器里输入 http://tinyurl.com/4e9iAk 时
// DNS首先解析获得 http://tinyurl.com 的 IP 地址
// 当 DNS 获得 IP 地址以后（比如：74.125.225.72），会向这个地址发送 HTTP GET 请求，查询短码 4e9iAk
// http://t.cn 服务器会通过短码 4e9iAk 获取对应的长 URL
// 请求通过 HTTP 301 转到对应的长 URL https://leetcode.com/problems/design-tinyurl 。
// 301 是永久重定向，302 是临时重定向。短地址一经生成就不会变化，所以用 301 是符合 http 语义的。

const chars = '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz'
// tinyURL 到 longURL映射
const codeDB = new Map<string, string>()
// longURL 到 tinyURL映射
const urlDB = new Map<string, string>()
const getCode = () => {
  const code = Array.from({ length: 6 }, () => chars.charAt(~~(Math.random() * 62)))
  return 'http://tinyurl.com/' + code.join('')
}

/**
 * Encodes a URL to a shortened URL.
 */
function encode(longUrl: string): string {
  if (urlDB.has(longUrl)) return urlDB.get(longUrl)!
  let code = getCode()
  while (codeDB.has(code)) code = getCode()
  codeDB.set(code, longUrl)
  urlDB.set(longUrl, code)
  return code
}

/**
 * Decodes a shortened URL to its original URL.
 */
function decode(shortUrl: string): string {
  return codeDB.get(shortUrl)!
}

/**
 * Your functions will be called as such:
 * decode(encode(strs));
 */
