// http://www.qq.com  // 通过

// http://www.qq.com.cn  // 不通过

// http://www.qq.com/a/b  // 通过

// http://www.qq.com?a=1  // 通过

// http://www.123qq.com?a=1  // 不通过

const check = (url: string): boolean => {
  const regexp = /^https?:\/\/(www\.)?qq\.com(\/.*)?$/
  return regexp.test(url)
}

console.log(
  [
    'http://www.qq.com',
    'http://www.qq.com.cn',
    'http://www.qq.com/a/b',
    'http://www.qq.com?a=1',
    'http://www.123qq.com?a=1',
    'http://www.baidu.com?redirect=http://www.qq.com/a'
  ].map(check)
)

export {}
