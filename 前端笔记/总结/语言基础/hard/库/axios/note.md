Axios 是一个基于 Promise 的 HTTP 客户端，该库拥有以下特性：

- 支持 Promise API；
- 能够拦截请求和响应；
- 能够转换请求和响应数据；
- 客户端支持防御 CSRF 攻击；
- 同时支持浏览器和 Node.js 环境；
- 能够取消请求及自动转换 JSON 数据。

axios 与 scrapy 很多相似的地方

1. http 拦截器
   Axios 的作用是用于发送 HTTP 请求，请求拦截器和响应拦截器分别对应于 HTTP 请求的不同阶段，它们的本质是一个实现特定功能的函数。这时我们就可以按照功能把发送 HTTP 请求拆解成不同类型的子任务，比如有 用于处理请求配置对象的子任务，用于发送 HTTP 请求的子任务 和 用于处理响应对象的子任务。当我们按照指定的顺序来执行这些子任务时，就可以完成一次完整的 HTTP 请求。
   既然已经提到了任务，我们就会联想到任务管理系统的基本功能：`任务注册`、`任务编排`（优先级排序）和`任务调度`等。因此我们就可以考虑从 任务注册、任务编排和任务调度 `三个方面`来分析 Axios 拦截器的实现。

   1. 任务注册
      注册到 Manager
   2. 任务编排
      FIFO/LIFO
   3. 任务调度

2. HTTP 适配器

```JS
// lib/core/dispatchRequest.js
module.exports = function dispatchRequest(config) {
  // 省略部分代码
  var adapter = config.adapter || defaults.adapter;

  return adapter(config).then(function onAdapterResolution(response) {
    // 省略部分代码
    return response;
  }, function onAdapterRejection(reason) {
    // 省略部分代码
    return Promise.reject(reason);
  });
};

// lib/defaults.js
var defaults = {
  adapter: getDefaultAdapter(),
  xsrfCookieName: 'XSRF-TOKEN',
  xsrfHeaderName: 'X-XSRF-TOKEN',
  //...
}

function getDefaultAdapter() {
  var adapter;
  if (typeof XMLHttpRequest !== 'undefined') {
    // For browsers use XHR adapter
    adapter = require('./adapters/xhr');
  } else if (typeof process !== 'undefined' &&
    Object.prototype.toString.call(process) === '[object process]') {
    // For node use HTTP adapter
    adapter = require('./adapters/http');
  }
  return adapter;
}


```

自定义适配器
其实除了默认的适配器外，我们还可以自定义适配器。那么如何自定义适配器呢

```JS
var settle = require('./../core/settle');
module.exports = function myAdapter(config) {
  // 当前时机点：
  //  - config配置对象已经与默认的请求配置合并
  //  - 请求转换器已经运行
  //  - 请求拦截器已经运行

  // 使用提供的config配置对象发起请求
  // 根据响应对象处理Promise的状态
  return new Promise(function(resolve, reject) {
    var response = {
      data: responseData,
      status: request.status,
      statusText: request.statusText,
      headers: responseHeaders,
      config: config,
      request: request
    };

    settle(resolve, reject, response);

    // 此后:
    //  - 响应转换器将会运行
    //  - 响应拦截器将会运行
  });
}

```

那么自定义适配器有什么用呢
通过自定义适配器，让开发者可以轻松地模拟请求

3.  CSRF 防御
    CSRF 的两个特点：CSRF（通常）发生在第三方域名。CSRF 攻击者`不能获取到Cookie等信息`，只是使用。
    CSRF 防御措施:检查 Referer 字段(注意 referer 可修改)/CSRF token/`双重 Cookie 防御`
    `双重 Cookie 防御`利用 CSRF 攻击不能获取到用户 Cookie 的特点，我们可以要求 Ajax 和表单请求携带一个 Cookie 中的值。 就是将 token 设置在 Cookie 中，在提交（POST、PUT、PATCH、DELETE）等请求时提交 Cookie，并通过请求头或请求体带上 Cookie 中已设置的 token，服务端接收到请求后，再进行对比校验。
    双重 Cookie 采用以下流程：

        1. 在用户访问网站页面时，向请求域名注入一个 Cookie，内容为随机字符串（例如 csrfcookie=v8g9e4ksfhw）。
        2. 在前端向后端发起请求时，取出 Cookie，并添加到 URL 的参数中（接上例 POST https://www.a.com/comment?csrfcookie=v8g9e4ksfhw）。
        3. 后端接口验证 Cookie 中的字段与 URL 参数中的字段是否一致，不一致则拒绝。

    利用 CSRF 攻击不能获取到用户 Cookie 的特点，我们可以要求 Ajax 和表单请求携带一个 Cookie 中的值。
    Axios 提供了 xsrfCookieName 和 xsrfHeaderName 两个属性来分别设置 `CSRF 的 Cookie 名称和 HTTP 请求头的名称`，它们的默认值如下所示：

```JS
// lib/defaults.js
var defaults = {
  adapter: getDefaultAdapter(),

  // 省略部分代码
  xsrfCookieName: 'XSRF-TOKEN',
  xsrfHeaderName: 'X-XSRF-TOKEN',
};

// lib/adapters/xhr.js
module.exports = function xhrAdapter(config) {
  return new Promise(function dispatchXhrRequest(resolve, reject) {
    var requestHeaders = config.headers;

    var request = new XMLHttpRequest();
    // 省略部分代码

    // 添加xsrf头部
    if (utils.isStandardBrowserEnv()) {
      var xsrfValue = (config.withCredentials || isURLSameOrigin(fullPath)) && config.xsrfCookieName ?
        cookies.read(config.xsrfCookieName) :
        undefined;

      if (xsrfValue) {
        requestHeaders[config.xsrfHeaderName] = xsrfValue;
      }
    }

    request.send(requestData);
  });
};


一重是浏览器自动附加的 cookie，另一重就是在页面代码中通过其他手段（自定义请求头、请求体、URL查询参数）传递的 cookie。
攻击者发起 CSRF 攻击时，请求中只会包含浏览器附加的 cookie 请求头，把这种只有一重 cookie 的请求拦掉就实现了防御。

```

Axios 内部是使用` 双重 Cookie 防御` 的方案来防御 CSRF 攻击。

4. Axios 如何取消重复请求？
   https://segmentfault.com/a/1190000021290514
   对于浏览器环境来说，Axios 底层是利用 XMLHttpRequest 对象来发起 HTTP 请求。如果要取消请求的话，我们可以通过调用 `XMLHttpRequest 对象上的 abort` 方法来取消请求：
   而对于 Axios 来说，我们可以通过 Axios 内部提供的`CancelToken`来取消请求：

   - 判断重复请求:
     请求 URL 地址和请求参数来生成一个唯一的 key，同时为每个请求创建一个专属的 CancelToken，然后把 `key 和 cancel 函数`以键值对的形式保存到 Map 对象（pendingRequest）中

   ```JS
    import qs from 'qs'

    const pendingRequest = new Map();
    // GET -> params；POST -> data
    const requestKey = [method, url, qs.stringify(params), qs.stringify(data)].join('&');
    const cancelToken = new CancelToken(function executor(cancel) {
      if(!pendingRequest.has(requestKey)){
        pendingRequest.set(requestKey, cancel);
      }
    })

   ```

   - 取消重复请求
     前置拦截:检查是否存在重复请求，若存在则取消已发的请求；把当前请求信息添加到 pendingRequest 对象中
     后置拦截:从 pendingRequest 对象中移除请求

     如果一个页面有两个一样的请求，那就会被拦截，这种情况怎么处理呢
     手动在请求里`加一个多余的参数`，这样两个请求就会不一样了；或者`加一个字段标识`，识别到有这个字段则不加入 pending

   - CancelToken 的工作原理

   ```JS
       // lib/adapters/xhr.js
     if (config.cancelToken) {
       config.cancelToken.promise.then(function onCanceled(cancel) {
         if (!request) { return; }
         request.abort(); // 取消请求
         reject(cancel);
         request = null;
       });
     }

   ```

   需要注意的是已取消的请求可能已经达到服务端，针对这种情形，`服务端的对应接口需要进行幂等控制`
   cancel 的还是会发到后端，只是浏览器不会处理;取消上次请求并不能中断上次请求发送到后端
   有可能是在请求发出去的过程中中断，也有可能是在响应回来的路上中断，但我只能保证前端获取的数据正常

5. Axios 如何实现请求重试？
6. Axios 如何缓存请求数据？
