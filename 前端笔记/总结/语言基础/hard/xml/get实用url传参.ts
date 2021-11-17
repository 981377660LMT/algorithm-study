const http = new XMLHttpRequest()
const endpoint = 'https://api.example.com/endpoint'
const params = {
  a: 1,
  b: 2,
  c: 3,
}
const url = endpoint + formatParams(params) // 实际应用中需要判断endpoint是否已经包含查询参数
// => "https://api.example.com/endpoint?a=1&b=2&c=3";

http.open('GET', url, true)
http.onreadystatechange = function () {
  if (http.readyState == 4 && http.status == 200) {
    alert(http.responseText)
  }
}
http.send(null) // 请求方法是GET或HEAD时，设置请求体为空

function formatParams(params: Record<any, any>) {
  return (
    '?' +
    Object.entries(params)
      .map(([key, value]) => `${key}=${value}`)
      .join('&')
  )
}

export {}
