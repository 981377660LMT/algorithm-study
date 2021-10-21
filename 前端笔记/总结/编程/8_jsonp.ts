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
