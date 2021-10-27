1.  Vue 的优点及缺点
    响应式和组件化
    缺点：基于对象配置文件的写法，也就是 options 写法，开发时不利于对一个属性的查找。另外一些缺点，在小项目中感觉不太出什么，vuex 的魔法字符串，对 ts 的支持。兼容性上存在一些问题
    不利于 seo,初次加载时耗时多
2.  Vue 中 hash 模式和 history 模式的区别
    hash 模式的 URL 中会夹杂着#号，而 history 没有
    hash 模式是依靠 **onhashchange** 事件(监听 location.hash 的改变)，而 history 模式是主要是依靠的 HTML5 history 中新增的两个方法，**pushState**()可以改变 url 地址且不会发送请求，**replaceState**()可以读取历史记录栈,还可以对浏览器记录进行修改。
    利用 pushState, replaceState 可以改变 url 同时浏览器不刷新，并且通过 **popstate** 监听浏览器历史记录的方式，完成一系列的异步动作。

---

    ````JS
    window.onhashchange = function(event){
    // location.hash 获取到的是包括#号的，如"#heading-3"
    // 所以可以截取一下
    let hash = location.hash.slice(1);
    }

        ```
    ````

3. 要获取当前时间你会放到 computed 还是 methods 里？
   放在 methods 中，因为 computed 会有惰性，并不能知道 new Date()的改变。
4. MVVM
   MVVM 是 Model-View-ViewModel 缩写，**也就是把 MVC 中的 Controller 演变成 ViewModel(Vue 对象 就是 ViewModel)**。Model 层代表数据模型，View 代表 UI 组件，ViewModel 是 View 和 Model 层的桥梁，数据会绑定到 viewModel 层并自动将数据渲染到页面中，视图变化的时候会通知 viewModel 层更新数据。
5. 4.nextTick 知道吗，实现原理是什么？
   它主要是为了解决：例如一个 data 中的数据它的改变会导致视图的更新，而在某一个很短的时间被改变了很多次，假如是 1000 次，每一次的改变如果都都将促发数据中的 setter 并按流程跑下来直到修改真实 DOM，那 DOM 就会被更新 1000 次，这样的做法肯定是非常低效的。
   Vue.js 源码中分别用 Promise、setTimeout、setImmediate 等方式定义了一个异步方法 nextTick，它接收的是一个回调函数，多次调用 nextTick 会将传入的回调函数存入队列中，当当前栈的任务都执行完毕之后才来执行这个队列中刚刚存储的那些回调函数，并且通过这个异步方法清空当前队列。
6. 接口请求一般放在哪个生命周期中？
   接口请求一般放在 mounted 中，但需要注意的是服务端渲染时不支持 mounted，需要放到 created 中。
7. Vue 模版编译原理
   Vue 模版编译，也就是 compilte 阶段，它其实就是将 template 转化为 render function 的过程，它会经过以下三个阶段：

   1. **parse** 阶段将 template 字符串通过各种正则表达式生成一颗抽象语法树 AST，在解析的过程中是通过 while 不断循环这个字符串，每解析完一个标签指针向下移动；并且用栈来建立节点间的层级关系，也就是用来保存解析好的标签头。

   2. **optimize** 阶段将生成的抽象语法树 AST 进行静态标记，这些被标记为静态的节点在后面的 patch 过程中会被跳过对比，从而达到优化的目的。标记的主要过程是为每个节点设置类似于 **static** 这样的属性，或者给根节点设置一个 **staticRoot** 属性表明这是不是一个静态根。这些被标记为静态的节点在后面的 patch 过程中会被跳过对比，从而达到优化的目的。

```JS
https://github.com/LinDaiDai/niubility-coding-js/blob/master/%E6%A1%86%E6%9E%B6-%E5%BA%93/Vue/vuejs%E6%BA%90%E7%A0%81-%E6%A8%A1%E7%89%88%E7%BC%96%E8%AF%91%E5%8E%9F%E7%90%86.md
// 判断是否是静态节点
function isStatic (node) {
  const { type } = node
  if (type === 2) { // 表达式节点
    return false
  } else if (type === 3) { // 文本节点
    return true
  }
  return (!node.if && !node.for) // 不存在if和for
}

```

3. 在进入到 **generate** 阶段之前，说明已经生成了被静态标记过的 AST，而 generate 就是将 AST 转化为 render function 字符串。

4. keep-alive 了解吗
   keep-alive 中运用了 LRU 算法。可以实现组件缓存，当组件切换时不会对当前组件进行卸载
5. SSR 了解吗
   远古技术
   SSR 也就是服务端渲染，也就是将 Vue 在客户端把标签渲染成 HTML 的工作放在服务端完成，然后再把 html 直接返回给客户端。
   SSR 有着更好的 SEO、并且首屏加载速度更快等优点。不过它也有一些缺点，比如我们的开发条件会受到限制，服务器端渲染只**支持 beforeCreate 和 created 两个钩子**，当我们需要一些外部扩展库时需要特殊处理，服务端渲染应用程序也需要处于 Node.js 的运行环境。还有就是服务器会有更大的负载需求。
6. Vue 的 diff 算法
   简单来说，diff 算法有以下过程

   1. 先同级比较再比较子节点
   2. 先判断一方有子节点和一方没有子节点的情况。如果新的一方有子节点，旧的一方没有，相当于新的子节点替代了原来没有的节点；同理，如果新的一方没有子节点，旧的一方有，相当于要把老的节点删除。
   3. 再来比较都有子节点的情况，这里是 diff 的核心。首先会通过判断两个节点的 key、tag、isComment、data 同时定义或不定义以及当标签类型为 input 的时候 type 相不相同来确定两个节点是不是相同的节点，如果不是的话就将新节点替换旧节点。
   4. **如果是相同节点的话才会进入到 patchVNode 阶段**。在这个阶段核心是采用双端比较的算法，同时从新旧节点的两端进行比较，在这个过程中，会用到模版编译时的静态标记配合 key 来跳过对比静态节点，如果不是的话再进行其它的比较。

7. nextTick 中的 waiting 是什么时候变为 true 的呢

**在下次 DOM 更新循环结束之后**
nextTick 的实现
Vue.js 实现了一个 nextTick 函数，传入一个 cb ，这个 cb 会被存储到一个队列中，在下一个 tick 时触发队列中的所有 cb 事件。

```JS
let callbacks = [];
let pending = false;

function nextTick (cb) {
    callbacks.push(cb);

    if (!pending) {
        pending = true;
        setTimeout(flushCallbacks, 0);
    }
}

function flushCallbacks () {
    pending = false;
    const copies = callbacks.slice(0);
    callbacks.length = 0;
    for (let i = 0; i < copies.length; i++) {
        copies[i]();
    }
}

```

这样做的目的是：

例如在一个时间点内，一直调用 nextTick

nextTick(cb)
nextTick(cb)
nextTick(cb)

由于 setTimeout 的原因，pending 变为了 true 之后就不会执行 if 里的代码了，而是等定时器执行了之后才变回来

8. Proxy 只会代理对象的第一层，那么 Vue3 又是怎样处理这个问题的呢
   判断当前 Reflect.get 的返回值是否为 Object，如果是则再通过 reactive 方法做代理， 这样就实现了深度观测。
9. Vue3 做了哪些优化
   响应式的原理不同
   更加细致的静态标记

10. Vue 实现响应式
    - initData(vm)
    - new Observer(value) // 观测过了就不会重复观测(看 value 里的`__ob__`是不是 instanceof Observer)
    - walk(obj)
    - defineReactive(obj) // 定义响应式
    - Object.defineProperty  
      get 中收集依赖 **dep.depend()**
      set 中更新时 **dep.notify()**
11. Vue 如何检测数组变化
    函数劫持 重写了 **data 里的数组原型 上的 7 个数组方法** `data.arr.__proto__=MyArray`
    里面会手动 dep.notify()通知视图更新
    数组里的对象继续递归观测
12. Vue 为何异步渲染
    防止频繁更新
    - dep.notify() 通知 Watcher 更新
    - subs[i].update() 依次调用 Watcher 的 update
    - queueWatcher 将 watcher 去重(根据 watcher 的 id 去重)后放到队列
    - nextTick（flushCallbacks) 异步清空 watcher 队列
13. nextTick 实现原理
14. Computed 实现
    - initComputed(vm,opts.computed)
    - new Watcher
    - watcher.**dirty** 为 true 则重新计算值
15. watch 中的 deep:true 如何实现的
    递归 好性能
16. Vue 的生命周期
17. ajax 请求放哪里
    created 中 dom 还未渲染出来
    **非 SSR**:mounted
    **SSR**:created
