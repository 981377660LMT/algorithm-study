interface JSONPOptions {
  url: string
  callback: (...args: any[]) => any
  params?: Record<string, any>
  callbackKey?: string
}

function JSONP(JSONPOptions: JSONPOptions) {
  const { url, callback, params = {}, callbackKey = 'cb' } = JSONPOptions
  // 定义本地的一个callback的名称
  const callbackName = 'JSONP_CALLBACK'
  // 把这个名称加入到参数中: 'cb=JSONP_CALLBACK' 告诉后端这是callback名称
  params[callbackKey] = callbackName
  // @ts-ignore 把这个callback加入到window对象中，这样就能执行这个回调了
  window[callbackName] = callback

  // 得到'id=1&cb=JSONP_CALLBACK'
  const paramString = Object.keys(params)
    .map(key => `${key}=${params[key]}`)
    .join('&')

  const script = document.createElement('script')
  script.setAttribute('src', `${url}?${paramString}`)
  document.body.appendChild(script)
}

JSONP({
  url: 'http://localhost:8080/api/jsonp',
  params: { id: 1 },
  callbackKey: 'cb',
  callback(res) {
    console.log(res)
  },
})
// 后端返回字符串
// JSONP_CALLBACK(res)

// 前端构造一个恶意页面，请求JSONP接口，收集服务端的敏感信息。如果JSONP接口还涉及一些敏感操作或信息（比如登录、删除等操作），那就更不安全了。
// 解决方法：验证JSONP的调用来源（Referer），服务端判断Referer是否是白名单，或者部署随机Token来防御。
// JSONP只支持GET请求，CORS支持所有类型的HTTP请求
// JSONP的优势在于支持老式浏览器，以及可以向不支持CORS的网站请求数据
