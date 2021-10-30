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
7. Vue 模版编译原理三步
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

vue3 中的优化：11

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
   4. **如果是相同节点的话才会进入到 patchVNode 阶段**。在这个阶段核心是采用**双指针**，双端比较的算法，同时从新旧节点的两端进行比较，在这个过程中，会用到模版编译时的静态标记配合 key 来跳过对比静态节点，如果不是的话再进行其它的比较。
      整体的执行思路如下：

      1. vnode 头对比 oldVnode 头
         vnode 尾对比 oldVnode 尾
         vnode 头对比 oldVnode 尾
         vnode 尾对比 oldVnode 头

      只要符合一种情况就进行 patch，移动节点，移动下标等操作

      2. 都不对再在 oldChild 中找一个 key 和 newStart 相同的节点
         找不到，新建一个。
         找到，获取这个节点，判断它和 newStartVnode 是不是同一个节点
         如果是相同节点，进行 patch 然后将这个节点插入到 oldStart 之前，newStart 下标继续移动
         如果不是相同节点，需要执行 createElm 创建新元素

      为什么会有头对尾、尾对头的操作？
      **可以快速检测出 reverse 操作，加快 diff 效率**

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
18. vue if show 的实现
    最终编译出的虚拟 dom 带有指令
19. 为什么 for if 不能连用
    **vue-for 的优先级高于 vue-if**
    每次渲染都会先循环再进行条件判断(就是我会把所有的代码**先渲染出来**在进行条件判断，这样就造成了性能的浪费)
    如果避免出现这种情况，则在外层嵌套 template（页面渲染不生成 dom 节点），在这一层进行 v-if 判断，然后在内部进行 v-for 循环
    ```Vue
    <template v-if="isShow">
        <p v-for="item in items">
    </template>
    ```
    如果条件出现在循环内部，可通过计算属性 computed 提前过滤掉那些不需要显示的项
20. 组件渲染和更新过程(前端开发最重要的 2 个工作，一个是把数据渲染到页面，另一个是处理用户交互)

    1. 第一步，解析 template 模板
       把模板中用到的 data 中的属性，都变成 JS 变量
       把模板中的 Vue 指令（v-for、v-model 等）绑定了相应的 JS 逻辑
       最终，把模板作为参数传入 with 函数中，返回 render 函数，**render 函数最终返回 Vnode 对象**。
       在创建组件的时候会调用一个 Vue.extend({}) 方法,创建完成之后会给组件加上一些 Hooks（钩子函数）
       render((h)=>h(...))
    2. 第二步，创建响应式，开始监听数据变化
       Object.defineProperty + 发布-订阅模式，实现响应式，
       再结合 DOM 事件的监听，
       最终实现双向绑定。

    3. 第三步，首次渲染，显示页面
       首先，初次渲染模版，会访问到 data 里的数据，触发 get 方法，创建订阅者，加入订阅中心。
       然后执行 通过 `vm._update` 完成的 函数，将 vnode 渲染成 DOM，初次渲染完成
    4. 第四步，data 属性变化，更新页面
       修改属性，被响应式的 set 监听到
       set 中执行 updateComponent
       updateComponent 重新执行，触发`vm._update`函数
       通过 patch 算法对比新旧 Vnode 对象，从而进行打补丁，
       修改原有的 dom 结构，更新页面

21. 为什么 vue 的 data 是函数
    否则多个组件共用一份 data
22. vue 事件绑定原理
    @click.native 与@click
    生成虚拟 dom 时，组件中的@click 变成 **on** @click.native 变成 **nativeOn**
23. v-model 是:value+@input 的语法糖

24. v-html 问题 :XSS 攻击/可能会替换掉标签内的子元素
    `<img src="" onerror="alert(1)" />`

25. 父子组件顺序
    父子子父 洋葱模型
26. 组件通信
    **父子**
    props/$on$emt
    $parent $children
    provide/inject
    $ref
    **兄弟**
    eventBus(Vue.prototype.$bus=new Vue())
    vuex
27. provide 与 inject 实现
    const provide=vm.$options.provide
    将 provide 挂载到 vm 上
    **inject 则不断寻找父亲的 provide 属性**
    并将 provide 进行 defineReactive 到自己身上
28. vue 中相同逻辑如何抽离
29. 异步组件
    异步组件实现的本质**是 2 次渲染**，除了 0 delay 的高级异步组件第一次直接渲染成 loading 组件外，其它都是第一次渲染生成一个注释节点，当异步获取组件成功后，再通过 forceRender 强制重新渲染，这样就能正确渲染出我们异步加载的组件了。
30. 什么是作用域插槽
    **作用域插槽允许你传递一个模板而不是已经渲染好的元素给插槽**
    之所以叫做”作用域“插槽，是因为模板虽然是在父级作用域中渲染的，却能拿到子组件的数据
    即：作用域插槽**会被解析成函数**而不是孩子节点 (React 传的函数 props=>...)
    被应用于表格中
31. 对 keep-alive 的了解
    组件缓存

    1. 2 个属性：include/exclude **字符串或正则表达式**。只有匹配的组件会被缓存/任何匹配的组件都不会被缓存
    2. 2 个生命周期：actived/deactived(含在 keep-alive 中创建的组件，会多出两个生命周期的钩子)
    3. 一个算法：LRU
       keep-alive 是一个抽象组件

    ```JS
        props:{
            include,
            exclude,
            max
        }
    ```

    - 被 keepalive 包含的组件不会被再次初始化，也就意味着不会重走生命周期函数
    - 组件一旦被 缓存，再次渲染就不会执行 created、mounted 生命周期钩子函数
    - 不会在函数式组件中正常工作，因为它们没有缓存实例。
    - activated 和 deactivated 是配合 keep-alive 一起使用的
      activated 和 deactivated 没有 keep-alive 的时候是不会被触发的
      在存在 keep-alive 的时候可以将 activated 当作 created 进行使用
      deactivated 是组件销毁的时候触发，此时的 destory 是不执行的
      **keep-alive 实现原理**
      组件通过插槽，获取第一个子节点。根据 include、exclude 判断是否需要缓存，通过组件的 key，判断是否命中缓存。利用 LRU 算法，更新缓存以及对应的 keys 数组。根据 max 控制缓存的最大组件数量。

```Vue
<div id="app" class='wrapper'>
    <keep-alive>
        <!-- 需要缓存的视图组件 -->
        <router-view v-if="$route.meta.keepAlive"></router-view>
    </keep-alive>
    <!-- 不需要缓存的视图组件 -->
    <router-view v-if="!$route.meta.keepAlive"></router-view>
</div>

```

32. 什么阶段才能访问 DOM:Mounted

```JS
callHook(vm, 'beforeCreate')
// 初始化 inject
// 初始化 props、methods、data、computed 和 watch
// 初始化 provide
callHook(vm, 'created')
// 挂载实例 vm.$mount(vm.$options.el)
```

为什么 created 之后才挂载实例
`vm.$mount(vm.$options.el)`
created 钩子函数中可以访问到数据，在 mounted 钩子函数中可以访问到 DOM，在 destroy 钩子函数中可以做一些定时器销毁工作

33. Vue 性能优化
    1. 编码优化
       1. 不要把数据全放在 data 中
       2. vue-for 事件代理
       3. keep-alive
       4. v-if 代替 v-show
       5. key 保证唯一性
       6. Object.freeze 避免劫持(大数据列表)
       7. 合理路由懒加载和异步组件
       8. 尽量用 runtime 版本
       9. 防抖节流
    2. Vue 加载性能优化
       1. 第三方模块按需导入
       2. 图片懒加载(h5 原生，img 的 src 替换)
       3. 滚动到可视区加载(IntersectionObserver) vue-virtual-scroll-list
    3. 用户体验
       1. 骨架屏
       2. app-shell
    4. SEO 优化 ssr
    5. 打包优化 cdn happypack 多线程打包 splitchunks sourcemap
    6. 缓存压缩 gzip
34. Vue3 改进
    1. proxy
    2. diff 只更新 vdom 动态数据的部分
35. hash 路由和 history 路由
    1. onhashchange
    2. history.pushState/replaceState
36. Vue 的入口
    Vue.js 是一个跨平台的 MVVM 框架，它可以跑在 web 上，也可以配合 weex 跑在 native 客户端上。**platform 是 Vue.js 的入口**，2 个目录代表 2 个主要入口，分别打包成运行在 web 上和 weex 上的 Vue.js。
37. Vue SSR 大体思路
    服务端渲染主要的工作是**把组件渲染为服务器端的 HTML 字符串**，将它们直接发送到浏览器，最后将静态标记"混合"为客户端上完全交互的应用程序
38. Vue 首次渲染过程
    解析：

    1. Vue 初始化，添加实例成员、静态成员，并在原型上挂载**patch**方法和$mount 方法。

    2. 初始化结束，调用 new Vue()。在 new Vue()的过程中，调用 **this.init()**方法, 给 vue 的实例挂载各种功能。

    3. 在 this.init() 内部最终会调用 entry-runtime-with-compiler.js 中的 vm.**$mount()**,用于获取 render 函数。

    4. $mount 获取 render 过程: 如果用户没有传入 render,会将 template 编译为 **render**，如果 template 也没有，则将 el 中的内容作为模版，通过 compileToFunctions() 生成 render。

    5. 接下来调用 runtime/index.js 中的 $mount, 重新获取 el 并调用 mountComponent() 方法。
       mountComponent 用于触发 beforeMount，定义 updateComponent,创建 watcher 实例，触发 mounted,并最终返回 vm 实例。

    6. 创建完 watcher 的实例后会调用一次 watcher.get() 方法，该方法会调用 updateComponent(), updateComponent()又会调用 vm.render() 以及 **vm.update()**,vm.\_update()会调用 vm.**patch**()挂载真实 dom,并将真实 dom 记录于 vm.$el 中。

39. 请简述虚拟 DOM 中 Key 的作用和好处。

解析：

​ 作用： **标识节点在当前层级的唯一性**。
​ 好处： 在执行 updateChildren 对比新旧 Vnode 的子节点差异时，通过设置 key 可以进行更高效的比较，便于复用节点。 降低创建销毁节点成本，从而减少 dom 操作，提升更新 dom 的性能。

40. 如何理解 MVVM(**数据驱动视图**)
