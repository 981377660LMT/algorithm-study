1.  Vue 的优点及缺点
    响应式和组件化
    缺点：基于对象配置文件的写法，也就是 options 写法，开发时不利于对一个属性的查找。另外一些缺点，在小项目中感觉不太出什么，vuex 的魔法字符串，对 ts 的支持。兼容性上存在一些问题
    不利于 seo,初次加载时耗时多
2.  Vue 中 hash 模式和 history 模式的区别
    hash 模式的 URL 中会夹杂着#号，而 history 没有
    hash 模式是依靠 **onhashchange** 事件(监听 location.hash 的改变)，而 history 模式是主要是依靠的 HTML5 history 中新增的两个方法，**pushState**()可以改变 url 地址且不会发送请求，**replaceState**()可以读取历史记录栈,还可以对浏览器记录进行修改。
    利用 pushState, replaceState 可以改变 url 同时浏览器不刷新，并且通过 **onpopstate** 监听浏览器历史记录的方式，完成一系列的异步动作。

    ```JS
    window.onhashchange = function(event){
    // location.hash 获取到的是包括#号的，如"#heading-3"
    // 所以可以截取一下
    let hash = location.hash.slice(1);
    }


    ```

    现在 router 还新增了 abstract 模式,为了支持移动端、后端等各个环境。

3.  要获取当前时间你会放到 computed 还是 methods 里？
    放在 methods 中，因为 computed 会有惰性，并不能知道 new Date()的改变。
4.  MVVM
    MVVM 是 Model-View-ViewModel 缩写，**也就是把 MVC 中的 Controller 演变成 ViewModel(Vue 对象 就是 ViewModel)**。Model 层代表数据模型，View 代表 UI 组件，ViewModel 是 View 和 Model 层的桥梁，数据会绑定到 viewModel 层并自动将数据渲染到页面中，视图变化的时候会通知 viewModel 层更新数据。
5.  MVC 和 MVVM 区别
    MVC 的思想：一句话描述就是 Controller 负责将 Model 的数据用 View 显示出来，换句话说就是在 Controller 里面把 Model 的数据赋值给 View，解耦。
    MVVM 与 MVC 最大的区别就是：它实现了 View 和 Model 的自动同步，也就是当 Model 的属性改变时，我们不用再自己手动操作 Dom 元素，来改变 View 的显示，而是改变属性后该属性对应 View 层显示会自动改变（对应 Vue 数据驱动的思想）
    为什么官方要说 Vue 没有完全遵循 MVVM 思想呢？
    严格的 MVVM 要求 View 不能和 Model 直接通信，而 Vue 提供了$refs 这个属性，让 Model 可以直接操作 View，违反了这一规定，所以说 Vue 没有完全遵循 MVVM。

6.  nextTick 知道吗，实现原理是什么？
    它主要是为了解决：例如一个 data 中的数据它的改变会导致视图的更新，而在某一个很短的时间被改变了很多次，假如是 1000 次，每一次的改变如果都都将促发数据中的 setter 并按流程跑下来直到修改真实 DOM，那 DOM 就会被更新 1000 次，这样的做法肯定是非常低效的。
    Vue.js 源码中分别用 Promise、setTimeout、setImmediate 等方式定义了一个异步方法 nextTick，它接收的是一个回调函数，多次调用 nextTick 会将传入的回调函数存入队列中，当当前栈的任务都执行完毕之后才来执行这个队列中刚刚存储的那些回调函数，并且通过这个异步方法清空当前队列。
7.  接口请求一般放在哪个生命周期中？
    如果异步请求不需要依赖 Dom 推荐在 created 钩子函数中调用异步请求，因为在 created 钩子函数中调用异步请求有以下优点：
    能更快获取到服务端数据，减少页面 loading 时间；
    ssr 不支持 beforeMount 、mounted 钩子函数，所以放在 created 中有助于一致性；
    也可以放到 mounted 中，但需要注意的是服务端渲染时不支持 mounted

8.  Vue 模版编译原理三步

```JS
 const template = `<p>{{message}}</p>` 转成
 with(this){return createElement('p',[createTextVNode(toString(message))])}
 //  其中this是vue实例
```

    Vue 模版编译，也就是 compile 阶段，它其实就是将 **template** 转化为 **render 函数** 的过程，它会经过以下三个阶段：

    1.  **parse** 阶段将 template 字符串通过各种正则表达式生成一颗抽象语法树 AST，在解析的过程中是通过 while 不断循环这个字符串，每解析完一个标签指针向下移动；并且用栈来建立节点间的层级关系，也就是用来保存解析好的标签头。

    2.  **optimize** 阶段将生成的抽象语法树 AST 进行静态标记，这些被标记为静态的节点在后面的 patch 过程中会被跳过对比，从而达到优化的目的。标记的主要过程是为每个节点设置类似于 **static** 这样的属性，或者给根节点设置一个 **staticRoot** 属性表明这是不是一个静态根。这些被标记为静态的节点在后面的 patch 过程中会被跳过对比，从而达到优化的目的。

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


        vue3 中的优化：
        1. patchFlag
        2. hoistStatic
        3. cacheHandler

    3. 在进入到 **generate** 阶段之前，说明已经生成了被静态标记过的 AST，而 generate 就是将 AST 转化为 render function 字符串。

1.  with 语法
    使用 with 改变 {} 作用域内自由变量的查找规则，当作 obj 属性来找
    找不到报错

```JS
const obj={
    a:1,
    b:2
}
with(obj){
    console.log(a)
    console.log(b)
}
```

9.  keep-alive 了解吗
    keep-alive 中运用了 LRU 算法。可以实现组件缓存，当组件切换时不会对当前组件进行卸载
10. SSR 了解吗
    远古技术
    SSR 也就是服务端渲染，也就是将 Vue 在客户端把标签渲染成 HTML 的工作放在服务端完成，然后再把 html 直接返回给客户端。
    SSR 有着更好的 SEO、并且首屏加载速度更快等优点。不过它也有一些缺点，比如我们的开发条件会受到限制，服务器端渲染只**支持 beforeCreate 和 created 两个钩子**，当我们需要一些外部扩展库时需要特殊处理，服务端渲染应用程序也需要处于 Node.js 的运行环境。还有就是服务器会有更大的负载需求。
11. Vue 的 diff 算法
    简单来说，diff 算法有以下过程

    1. **patch**:先同级比较再比较子节点。如果相同 vnodenode(通过 tag 和 key 判断) ,走 **patchVNode**,否则重建。
    2. **patchVNode**:先判断一方有子节点和一方没有子节点的情况。如果新的一方有子节点，旧的一方没有，相当于新的子节点替代了原来没有的节点；同理，如果新的一方没有子节点，旧的一方有，相当于要把老的节点删除。再来比较都有子节点的情况，直接进 **updateChildren**。
    3. **updateChildren** : **如果是相同节点的话才会进入到 updateChildren 阶段**。在这个阶段核心是采用**双指针**，双端比较的算法，同时从新旧节点的两端进行比较，在这个过程中，会用到模版编译时的静态标记配合 key 来跳过对比静态节点，如果不是的话再进行其它的比较。
       整体的执行思路如下：

       1. vnode 头对比 oldVnode 头
          vnode 尾对比 oldVnode 尾
          vnode 头对比 oldVnode 尾
          vnode 尾对比 oldVnode 头

       只要符合一种情况就进行 **patchVNode**，移动节点，移动下标等操作

       2. 都不对再在 老节点的 **KeyToIndexMap** 中找一个 key 和 newStart 相同的节点
          找不到，新建一个。
          找到，获取这个节点，判断它和 newStartVnode 是不是同一个节点(还要比较 type)
          如果是相同节点，进行 **patchVNode** 然后将这个老节点，newStart 下标继直接搬过来，继续移动
          如果不是相同节点，需要执行 createElm 创建新元素

       为什么会有头对尾、尾对头的操作？
       **可以快速检测出 reverse 操作，加快 diff 效率**

12. nextTick 中的 waiting 是什么时候变为 true 的呢

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

由于 setTimeout 的原因，`pending 变为了 true 之后就不会执行 if 里的代码了`，而是等定时器执行了之后才变回来

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
    v-if 在编译过程中会被转化成三元表达式,条件不满足时不渲染此节点。
    v-show 会被编译成指令，条件不满足时控制样式将对应节点隐藏 （display:none）
    display:none、visibility:hidden 和 opacity:0 之间的区别？
    display:none 不占位
    属性|是否占据空间|事件绑定是否触发
    -----|-----|-----
    display:none|x|x
    visibility:hidden|√|x
    opacity:0|√|√
19. 为什么 for if 不能连用
    **vue-for 的优先级高于 vue-if**
    每次渲染都会`先循环再进行条件判断`(就是我会把所有的代码**先渲染出来**在进行条件判断，这样就造成了性能的浪费)
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
23. v-model

    - 在普通标签上是:
      value+@input 的语法糖，并且会处理拼音输入法的问题，
      text 和 textarea 元素使用 value property 和 input 事件；
      checkbox 和 radio 使用 checked property 和 change 事件；
      select 字段将 value 作为 prop 并将 change 作为事件。

    - 在组件上是:value+@input 的语法糖

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

    ```JS
       // src/init.js

       Vue.prototype.$mount = function (el) {
       const vm = this;
       const options = vm.$options;
       el = document.querySelector(el);

       // 如果不存在render属性
       if (!options.render) {
           // 如果存在template属性
           let template = options.template;

           if (!template && el) {
           // 如果不存在render和template 但是存在el属性 直接将模板赋值到el所在的外层html结构（就是el本身 并不是父元素）
           template = el.outerHTML;
           }

           // 最终需要把tempalte模板转化成render函数
           if (template) {
           const render = compileToFunctions(template);
           options.render = render;
           }
       }

       // 将当前组件实例挂载到真实的el节点上面
       return mountComponent(vm, el);
       };

    ```

    5. 接下来调用 runtime/index.js 中的 $mount, 重新获取 el 并调用 mountComponent() 方法。
       mountComponent 用于触发 beforeMount，定义 updateComponent,创建 watcher 实例，触发 mounted,并最终返回 vm 实例。

    6. 创建完 watcher 的实例后会调用一次 watcher.get() 方法，该方法会调用 updateComponent(), updateComponent()又会调用 vm.render() 以及 **vm.update()**,vm.\_update()会调用 vm.**patch**()挂载真实 dom,并将真实 dom 记录于 vm.$el 中。

    总结：

    - 解析模板为 render 函数(或在开发环境已完成,vue-loader)
    - 触发；响应式 vm.render 函数会触发 getter
    - 执行 render 生成 vnode 并 patch

39. Vue 更新过程(详细过程见 Vue 源码解析:响应式原理)
    总结：
    触发 setter
    重新 render ,patch
    nextTick 异步渲染
40. 请简述虚拟 DOM 中 Key 的作用和好处。

解析：

​ 作用： **标识节点在当前层级的唯一性**。
​ 好处： 在执行 updateChildren 对比新旧 Vnode 的子节点差异时，通过设置 key 可以进行更高效的比较，便于复用节点。 降低创建销毁节点成本，从而减少 dom 操作，提升更新 dom 的性能。

40. 如何理解 MVVM(**数据驱动视图**)
    很久以前就有组件化 **ejs 模板引擎的 include**
    View:DOM
    ViewModel:Vue
    Model:Object
    **解耦了 V 和 M 层**
41. Proxy 与 Object.defineProperty 的区别

    1. Object.defineProperty
       深度监听需要递归到底
       无法监听新增属性(Vue.set Vue.delete)
       不能监听数组变化

    2. Proxy
       lazy 监听
       可监听新属性
       可监听数组变化
       无法兼容所有浏览器，无法 polyfill

42. JS 模拟 DOM 结构

```HTML
    <div id="div1" class="container">
      <p>vdom</p>
      <ul style="font-size: 20px">
        <li>a</li>
      </ul>
    </div>
```

```TS
简易版虚拟DOM
interface VirtualDom {
  type: keyof HTMLElementTagNameMap
  props: {
    children?: Children[] | Children
    [attr: string]: any
  }
}

type Children = VirtualDom | string

const vdom: VirtualDom = {
  type: 'div1',
  props: {
    className: 'container',
    id: 'div1',
    children: [
      { type: 'p', props: { children: 'vdom' } },
      {
        type: 'ul',
        props: {
          style: 'font-size:20px',
          children: {
            type: 'li',
            props: {
              children: 'a',
            },
          },
        },
      },
    ],
  },
}

```

43. hash 路由特点
    1. 触发网页跳转
    2. 不刷新页面
    3. 不提交到 server
44. Vue3 比 Vue2 的优势/新功能
    - createApp 方法挂载到 Vue 上=>方法挂载到根实例 app 上
    - emits 属性 emits 时**子组件声明 emits 名字** setup 里 emit
    - 生命周期
    - 多事件
    - Fragment
    - 移除.sync 变为 v-model

```JS
<Component v-model:title='pageTitle' />
是以下的缩写:
<Component :title='pageTitle' @update:title='pageTitle=$event' />

```

    - 异步组件写法
    - 移除 filter
    - teleport
    - Suspense Suspense 内部有一个具名插槽 fallback
    - CompositionAPI

45. Vue3 生命周期
    Options API 生命周期 改名

    ```JS
    beforeCreate() {
        console.log('beforeCreate')
    },
    created() {
        console.log('created')
    },
    beforeMount() {
        console.log('beforeMount')
    },
    mounted() {
        console.log('mounted')
    },
    beforeUpdate() {
        console.log('beforeUpdate')
    },
    updated() {
        console.log('updated')
    },
    // beforeDestroy 改名
    beforeUnmount() {
        console.log('beforeUnmount')
    },
    // destroyed 改名
    unmounted() {
        console.log('unmounted')
    }
    ```

    Composition API 生命周期

    ```JS
        // 等于 beforeCreate 和 created
        setup() {
            console.log('setup')

            onBeforeMount(() => {
                console.log('onBeforeMount')
            })
            onMounted(() => {
                console.log('onMounted')
            })
            onBeforeUpdate(() => {
                console.log('onBeforeUpdate')
            })
            onUpdated(() => {
                console.log('onUpdated')
            })
            onBeforeUnmount(() => {
                console.log('onBeforeUnmount')
            })
            onUnmounted(() => {
                console.log('onUnmounted')
            })
        },
    ```

46. CompositionAPI

    - 代码组织
    - 逻辑复用
    - 类型推导

47. CompositionAPI 与 reactHooks 区别

    - 前者 **setup 只被调一次** 而后者函数会被调很多次
    - 前者无需 useMemo/useCallback 因为 setup 只调一次
    - 前者无需顾虑顺序 后者需要保证 hooks 调用顺序
    - 前者 ref 加 reactive 比 useState 难以理解

48. 为什么需要 ref
    - 返回值类型会丢失响应式
    - setep/computed/合成函数 都有可能返回值类型
    - 如果 vue 不定义 ref 用户会自造 ref，反而混乱
49. 为什么 ref 需要.value
    - ref 是一个对象，vaue 存值
    - 通过.value 的 get 和 set 实现响应式
    - 用于模板和 reactive 不需要.value 其余需要
50. 为什么需要 toRef toRefs
    不创造响应式，延续响应式
51. proxy 模拟

    ```JS
    // 创建响应式
    function reactive(target = {}) {
    if (typeof target !== 'object' || target == null) {
        // 不是对象或数组，则返回
        return target
    }

    // 代理配置
    const proxyConf = {
        get(target, key, receiver) {
        // 只处理本身（非原型的）属性
        const ownKeys = Reflect.ownKeys(target)
        const result = Reflect.get(target, key, receiver)
        // 深度监听:惰性的监听，只有被用到才递归监听
        // 性能如何提升的？
        return reactive(result)
        },
        set(target, key, val, receiver) {
        // 重复的数据，不处理
        if (val === target[key]) {
            return true
        }
        const ownKeys = Reflect.ownKeys(target)
        const result = Reflect.set(target, key, val, receiver)
        return result // 是否设置成功
        },
        deleteProperty(target, key) {
        const result = Reflect.deleteProperty(target, key)
        return result // 是否删除成功
        },
    }

    // 生成代理对象
    const observed = new Proxy(target, proxyConf)
    return observed
    }
    ```

52. setup 内如何获取实例
    CompositionAPI 里没有 this
    **getCurrentInstance**
53. Vue3 为什么比 Vue2 快
    https://vue-next-template-explorer.netlify.app/

    1. proxy
    2. patchFlag
       编译模板时动态节点做标记,分为 Text PROPS 等类型
       diff 算法时可以区分动静态结点/不同动态节点
       **输入做了标记** 从而 diff 性能得到提高
       优化并不只是 diff 算法 而是整个流程
    3. hoistStatic
       将静态节点定义提升到父作用域缓存
       多个相邻的静态节点会被合并(相邻的静态节点多到一定程度会被合并)
       空间换时间

       ```HTML
       <div>Hello World!</div>
       <div>{{ma}}</div>
       ```

       ```JS

        import { createElementVNode as _createElementVNode, toDisplayString as _toDisplayString, Fragment as _Fragment, openBlock as _openBlock, createElementBlock as _createElementBlock } from "vue"

        const _hoisted_1 = /*#__PURE__*/_createElementVNode("div", null, "Hello World!", -1 /* HOISTED */)

        export function render(_ctx, _cache, $props, $setup, $data, $options) {
            return (_openBlock(), _createElementBlock(_Fragment, null, [
                _hoisted_1,
                _createElementVNode("div", null, _toDisplayString(_ctx.ma), 1 /* TEXT */)
            ], 64 /* STABLE_FRAGMENT */))
        }

       ```

    4. cacheHandler
       缓存事件

       ```HTML
        <div>Hello World!</div>
        <div @click="clickHandler">haha</div>
       ```

       ```JS
       不开cacheHandler
        const _hoisted_1 = /*#__PURE__*/_createElementVNode("div", null, "Hello World!", -1 /* HOISTED */)
        const _hoisted_2 = ["onClick"]  // 缓存事件

        export function render(_ctx, _cache, $props, $setup, $data, $options) {
        return (_openBlock(), _createElementBlock(_Fragment, null, [
            _hoisted_1,
            _createElementVNode("div", { onClick: _ctx.clickHandler }, "haha", 8 /* PROPS */, _hoisted_2)
        ], 64 /* STABLE_FRAGMENT */))
        }

        开cacheHandler
        const _hoisted_1 = /*#__PURE__*/_createElementVNode("div", null, "Hello World!", -1 /* HOISTED */)

        export function render(_ctx, _cache, $props, $setup, $data, $options) {
        return (_openBlock(), _createElementBlock(_Fragment, null, [
            _hoisted_1,
            _createElementVNode("div", {
            onClick: _cache[0] || (_cache[0] = (...args) => (_ctx.clickHandler && _ctx.clickHandler(...args)))
            }, "haha")
        ], 64 /* STABLE_FRAGMENT */))
        }
       ```

    5. SSR 优化
       静态节点绕过 vdom
       动态节点还是要
    6. tree-shaking
       根据模板里的属性动态 import 处理函数

54. jsx 与 template
    jsx 已经是 ES 规范
    template 还是 Vue 自家规范，但是可读性更好
    **jsx 写 slot 更加方便**
    组件更推荐使用 jsx
55. 怎样理解 Vue 的单向数据流
    数据总是从父组件传到子组件，子组件没有权利修改父组件传过来的数据，只能请求父组件对原始数据进行修改。这样会防止从子组件意外改变父级组件的状态，从而导致你的应用的数据流向难以理解。
    如果实在要改变父组件的 prop 值 可以再 data 里面定义一个变量 并用 prop 的值初始化它 之后用$emit 通知父组件去修改
56. Vue 的父子组件生命周期钩子函数执行顺序
    加载渲染过程
    父 beforeCreate->父 created->父 beforeMount->子 beforeCreate->子 created->子 beforeMount->子 mounted->父 mounted
    子组件更新过程
    父 beforeUpdate->子 beforeUpdate->子 updated->父 updated
    销毁过程
    父 beforeDestroy->子 beforeDestroy->子 destroyed->父 destroyed
57. vue-router 动态路由是什么 有什么问题

```JS
const router = new VueRouter({
  routes: [
    // 动态路径参数 以冒号开头
    { path: "/user/:id", component: User },
  ],
});
```

58. vue-router 组件复用导致路由参数失效怎么办？
    通过 watch 监听路由参数再发请求

```JS
watch: { //通过watch来监听路由变化
 "$route": function(){
 this.getData(this.$route.params.xxx);
 }
}
```

59. 谈一下对 vuex 的个人理解
    vuex 是专门为 vue 提供的全局状态管理系统，用于多个组件中数据共享、数据缓存等。（无法持久化、内部核心原理是通过**创造一个全局实例 new Vue**）
    State：定义了应用状态的数据结构，可以在这里设置默认的初始状态。
    Getter：允许组件从 Store 中获取数据，mapGetters 辅助函数仅仅是将 store 中的 getter 映射到局部计算属性。
    Mutation：是唯一更改 store 中状态的方法，且必须是同步函数。
    Action：用于提交 mutation，而不是直接变更状态，可以包含任意异步操作。
    Module：允许将单一的 Store 拆分为多个 store 且同时保存在单一的状态树中。
    之所以拆成 mutation 和 action 是为了实现时间旅行功能(devtools 记录状态变化)
60. Vuex 页面刷新数据丢失怎么解决
    推荐使用 vuex-persist 插件，它就是为 Vuex 持久化存储而生的一个插件。不需要你手动存取 storage ，而是直接将状态保存至 cookie 或者 localStorage 中
61. Vuex 为什么要分模块并且加命名空间
    抽离
62. 写过自定义指令吗 原理是什么
    1. 在生成 ast 语法树时，遇到指令会给当前元素添加 directives 属性
    2. 通过 genDirectives 生成指令代码
    3. 在 patch 前将指令的钩子提取到 cbs 中,在 patch 过程中调用对应的钩子
    4. 当执行指令对应钩子函数时，调用对应指令定义的方法
63. Vue 修饰符有哪些
    .stop 阻止事件继续传播
    .prevent 阻止标签默认行为
    .capture 使用事件捕获模式,即元素自身触发的事件先在此处处理，然后才交由内部元素进行处理
    .self 只当在 event.target 是当前元素自身时触发处理函数
    .once 事件将只会触发一次
    .passive 告诉浏览器你不想阻止事件的默认行为
64. 生命周期钩子是如何实现的
    对应阶段 callHook

```JS
export function callHook(vm, hook) {
  // 依次执行生命周期对应的方法
  const handlers = vm.$options[hook];
  if (handlers) {
    for (let i = 0; i < handlers.length; i++) {
      handlers[i].call(vm); //生命周期里面的this指向当前实例
    }
  }
}

// 调用的时候
Vue.prototype._init = function (options) {
  const vm = this;
  vm.$options = mergeOptions(vm.constructor.options, options);
  callHook(vm, "beforeCreate"); //初始化数据之前
  // 初始化状态
  initState(vm);
  callHook(vm, "created"); //初始化数据之后
  if (vm.$options.el) {
    vm.$mount(vm.$options.el);
  }
};
```

65. vue-router 中路由方法 pushState 和 replaceState 能否触发 popSate 事件
    答案是：不能
    注意:用 history.pushState()或者 history.replaceState()不会触发 popstate 事件
    只有在做出浏览器动作时，才会触发该事件，如用户点击浏览器的回退按钮（或者在 Javascript 代码中调用 history.back()）
    注意:仅改变网址,网页不会真的跳转,也不会获取到新的内容,本质上网页还停留在原页面
    popstate 事件会在点击后退、前进按钮(或调用 history.back()、history.forward()、history.go()方法)时触发

66. Vue2 中改变数组的陷阱是什么?
    在以下两种情况下，Vue 无法检测数组的更改,

    ```JS
    vm.todos[indexOfTodo] = newTodo
    vm.todos.length = todosLength
    使用 set 和 splice 方法克服这两个问题,

    // Vue.set
    Vue.set(vm.todos, indexOfTodo, newTodoValue)
    // splice
    vm.todos.splice(todosLength)
    ```

67. 描述一下 props 的验证？

```JS
Vue.component('user-profile', {
  props: {
    // Basic type check (`null` matches any type)
    age: Number,
    // Multiple possible types
    identityNumber: [String, Number],
    // Required string
    email: {
      type: String,
      required: true
    },
    // Number with a default value
    minBalance: {
      type: Number,
      default: 10000
    },
    // Object with a default value
    message: {
      type: Object,
      // Object or array defaults must be returned from
      // a factory function
      default: function () {
        return { message: 'Welcome to Vue' }
      }
    },
    // Custom validator function
    location: {
      validator: function (value) {
        // The value must match one of these strings
        return ['India', 'Singapore', 'Australia'].includes(value)
      }
    }
  }
})

```

68. 什么是自定义指令？
    所有的钩子都有 el、 binding 和 vnode 作为参数
    binding.value 是传给指令的值

```HTML
<input v-focus>
<div v-avatar="{ width: 500, height: 400, url: 'path/logo', text: 'Iron Man' }"></div>

```

```JS

// Register a global custom directive called `v-focus`
Vue.directive('focus', {
  // When the bound element is inserted into the DOM...
  inserted: function (el) {
    // Focus the element
    el.focus()
  }
})

Vue.directive('avatar', function (el, binding) {
  console.log(binding.value.width) // 500
  console.log(binding.value.height)  // 400
  console.log(binding.value.url) // path/logo
  console.log(binding.value.text)  // "Iron Man"
})

```

69. 什么是动态组件？
    动态组件用于使用元素动态切换多个组件，并将数据传递给 v-bind: is 属性。

```JS
new Vue({
  el: '#app',
  data: {
    currentPage: 'home'
  },
  components: {
    home: {
      template: "<p>Home</p>"
    },
    about: {
      template: "<p>About</p>"
    },
    contact: {
      template: "<p>Contact</p>"
    }
  }
})

<div id="app">
   <component v-bind:is="currentPage">
       <!-- component changes when currentPage changes! -->
       <!-- output: Home -->
   </component>
</div>
```

70. 什么是递归组件？
    可以在自己的模板中递归调用自己的组件称为递归组件。

```JS
Vue.component('recursive-component', {
  template: `<!--Invoking myself!-->
             <recursive-component></recursive-component>`
});

```

递归组件对于在博客、嵌套菜单或基本上任何父子内容相同的地方显示评论非常有用，尽管内容不同。
请记住，递归组件可能导致无限循环，最大堆栈大小超过错误，因此请确保递归调用是有条件的(例如 v-if 指令)。

71. 在 vuex 中使用严格模式的目的是什么？
    Vuex 的 state 在 mutation 之外被改变时，就会抛出一个错误。

```JS
const store = new Vuex.Store({
  // ...
  strict: process.env.NODE_ENV !== 'production'
})

```

72. 什么是 vuex 插件？

```JS
const myPlugin = store => {
  // called when the store is initialized
  store.subscribe((mutation, state) => {
    // called after every mutation.
    // The mutation comes in the format of `{ type, payload }`.
  })
}

const store = new Vuex.Store({
  // ...
  plugins: [myPlugin]
})

```

73. vuex 不直接更新状态的原因是什么？
    我们希望显式地跟踪应用程序状态，以便实现可以记录每次变异、获取状态快照、甚至执行时间旅行调试的工具。所以我们需要进行 mutation，而不是直接改变存储的状态。
74. 什么是 vuex 的模块？

```JS
const moduleOne = {
  state: { ... },
  mutations: { ... },
  actions: { ... },
  getters: { ... }
}

const moduleTwo = {
  state: { ... },
  mutations: { ... },
  actions: { ... },
  getters: { ... }
}

const store = new Vuex.Store({
  modules: {
    one: moduleOne,
    two: moduleTwo
  }
})

store.state.one // -> `moduleOne's state
store.state.two // -> `moduleTwo's state

```

75. vue 中怎么重置 data

```JS
Object.assign(this.$data, this.$options.data())
this.$data获取当前状态下的data
this.$options.data()获取该组件初始状态下的data。
```

76. 怎么解决 vue 动态设置 img 的 src 不生效的问题？

```HTML
<img :src="require('../../../assets/images/xxx.png')" />
```

77. vue 为什么要求组件模板只能有一个根元素？
    1. 在全局组件生命时，Vue 组件需要指定 el 元素的 id 为挂载对象， 这也就明确了根元素只能有一个
    2. 在单文件组件中， template 最终渲染时，不会出现在页面中，此时 template 中的元素挂载时也需要唯一入口
    3. 在虚拟 dom patch 时， 新旧节点自上而下，逐级比较；有子节点则向下递归。所以树的根节点只能有一个
78. 你知道 vue 的模板语法用的是哪个 web 模板引擎的吗？说说你对这模板引擎的理解
    期初 vue 用的 jade 模板，后来由于商标原因改成了 pug,只是换个名字，语法都与 jade 一样，
79. style 加 scoped 属性的用途和原理
    用途：防止全局同名 CSS 污染
    原理：在标签加上 v-data-something 属性，再在选择器时加上对应[v-data-something]，即 CSS 带属性选择器，以此完成类似作用域的选择方式
    一、scoped 会在元素上添加唯一的属性（data-v-x 形式），css 编译后也**会加上属性选择器**，从而达到限制作用域的目的。
    缺点：
    （1）由于只是通过属性限制，**类还是原来的类**，所以在其他地方对类设置样式还是可以造成污染。
    （2）添加了属性选择器，对于 CSS 选择器的**权重加重了(10)**。
    （3）外层组件包裹子组件，会给子组件的根节点添加 data 属性。**在外层组件中无法修改子组件中除了根节点以外的节点的样式**。比如子组件中有 box 类，在父节点中设置样式，会被编译为
    .box[data-v-x]的形式，但是 box 类所在的节点上没有添加 data 属性，因此无法修改样式。
    可以**使用/deep/或者>>>穿透 CSS**，这样外层组件设置的 box 类编译后的就为[data-v-x] .box 了，就可以进行修改。
    二、可以使用 CSS Module
    CSS Module 没有添加唯一属性，而是通过**修改类名限制作用域**。这样类发生了变化，在其他地方设置样式无法造成污染，也没有使 CSS 选择器的权重增加。
80. 从 0 到 1 自己构架一个 vue 项目，说说有哪些步骤、哪些重要插件、目录结构你会怎么组织
    vue-cli 实际上已经很成熟了，目录除了脚手架默认的，
    1、一般会额外创建 views，components，api，utils，stores 等；
    2、下载重要插件，比如 axios，dayjs（moment 太大），其他的会根据项目需求补充；
    3、封装 axios，统一 api 调用风格和基本配置；
    4、如果有国际化需求，配置国际化；
    5、开发，测试和正式环境配置，一般不同环境 API 接口基础路径会不一样；
81. 使用 vue 渲染大量数据时应该怎么优化
    1. 虚拟列表：vue-virtual-scroll-list，vue-virtual-scroller……
    2. 冻结属性，**让不必要的属性不响应**：Object.freeze, 或者使用 Object.defineProperty 将对象属性的 configurable 设置为 false，源码：
    3. 分页
82. vue 自定义事件中父组件怎么接收子组件的多个参数

```JS
this.$emit("eventName",data)
data为一个对象
```

83. vue 的优缺点?

- 优点

  ① 易用性：vue 提供数据响应式、基于配置的组件系统以及大量的指令等，这些让开发者只需关心核心业务即可。

  ② 灵活性：如果我们的应用足够小，可以只使用 vue 的核心库即可；随着应用的规模不断扩大，可以根据需求引入 vue-router、vuex、vue-cli 等其他工具。

  ③ 高效性：vue 操作的是虚拟 DOM，采用 diff 算法更新 DOM，比传统的 DOM 操作更加的高效。

- 缺点：

  ① 不支持 IE8 及以下版本

  ② 这不是单纯的 vue 的缺点，其他框架如 react、angular 均有，那就是不利于 SEO，不适合做需要浏览器抓取信息的电商网站，比较适合做后台管理系统。当然 vue 也有相关的 SSR 服务端渲染方式，也有针对 vue 服务端渲染的库，nuxt.js 来提高 SEO。

84. vue-cli 默认是单页面的，那要弄成多页面该怎么办呢
    在 vue.config.js 配置文件下的 pages 配置项中配置多个页面，每个页面有一个对应的 js 入口
85. vue 变量名如果以\_、$开头的属性会发生什么问题？怎么访问到它们的值？
以 _ 或 $ 开头的属性 不会 被 Vue 实例代理，因为它们可能和 Vue 内置的属性、API 方法冲突。你可以使用例如 vm.$data.\_property 的方式访问这些属性。

86. vuex 怎么知道 state 是通过 mutation 修改还是外部直接修改的？

```
默认严格模式下：Vuex 中修改 state 的唯一渠道就是执行 commit('xx', payload) 方法，其底层通过执行 this._withCommit(fn) 设置_committing `标志变量为 true`，然后才能修改 state，修改完毕还需要还原_committing 变量。外部修改虽然能够直接修改 state，但是并没有修改_committing 标志位，所以只要 watch 一下 state，state change 时判断是否_committing 值为 true，即可判断修改的合法性。s
```

87. 组件中写 name 选项有什么作用
    https://www.jb51.net/article/140702.htm
    1. DOM 做递归组件时需要调用自身 name
    2. vue-devtools 调试工具里显示的组见名称是由 vue 中组件 name 决定的
    3. 项目使用 keep-alive 时，可搭配组件 name 进行缓存过滤
    4. 组件安装
88. 说下你对指令的理解？
    指令就是把操作 DOM 的方法糅合为一条命令。
89. vue-router 有哪几种导航钩子（ 导航守卫 ）？
    3 种，全局的、组件的、单个路由独享的
90. 使用 vue 写一个 tab 切换

```HTML
<template>
  <div class="tab-wrap">
    <div class="tab">
      <button
        v-for="i in list"
        :key="i.id"
        @click="choose = i.id"
      >{{i.title}}</button>
    </div>
    <div class="tab-content">{{list[choose].content}}</div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      list: [
        {id: 0, title: 'A', content: 'aaaaaaaaaaaaaa'},
        {id: 1, title: 'B', content: 'bbbbbbbbbbbbbb'},
        {id: 2, title: 'C', content: 'cccccccccccccc'}
      ],
      choose: 0
    }
  }
}
</script>

```

91. vue-router 如何响应路由参数的变化？
    使用 watch 监听

```JS
watch: {
    $route(to, from){
        if(to != from) {
            console.log("监听到路由变化，做出相应的处理");
        }
    }
}

```

92. 切换到新路由时，页面要滚动到顶部或保持原先的滚动位置怎么做呢
    通过 router 的 meta 来记录需要保存滚动条的位置,在 new VueRouter()时调用 **scrollBehavior**(to, from, savedPosition) {return { x: 0, y: 0 }}的方法
93. route 和 router 有什么区别
    route：代表当前路由信息对象，可以获取到当前路由的信息参数
    router：代表路由实例的对象，包含了路由的跳转方法，钩子函数等
94. vue 项目本地开发完成后部署到服务器后报 404 是什么原因呢
    1. 检查 nginx 配置，是否正确设置了资源映射条件；
    2. 检查 vue.config.js 中是否配置了 publicPath，若有则检查是否和项目资源文件在服务器摆放位置一致。
95. 说说组件的命名规范
    定义组件名有两种方式：
    1.kebab-case（短横线分隔命名），引用时必须也采用 kebab-case；
    2.PascalCase（首字母大写命名），引用时既可以采用 PascalCase 也可以使用 kebab-case；
    但在 DOM 中使用只有 kebab-case 是有效的
96. e.target 和 e.currentTarget 有什么区别？
    event.currentTarget 指向事件所绑定的元素，而 event.target 始终指向事件发生时的元素。
97. v-show 指令算是重排吗？
    v-show 本质是通过元素 css 的 display 属性来控制是否显示，在 DOM 渲染时仍然会先渲染元素，然后才会进行判断是否显示（通过 display 属性），而对于重排的定义是渲染树中的节点信息发生了**大小、边距等改变，要重新计算各节点和 css 具体的大小和位置**。
    当用 display 来控制元素的显示和隐藏时，会改变节点的大小和渲染树的布局，导致发生重排，因此 v-show 指令算是重排。
98. 跟 keep-alive 有关的生命周期是哪些？描述下这些生命周期
    activated 和 deactivated

        keep-alive的生命周期
        1.activated： 页面第一次进入的时候，钩子触发的顺序是created->mounted->activated
        2.deactivated: 页面退出的时候会触发deactivated，当再次前进或者后退的时候只触发activated

99. vue 要做权限管理该怎么做？如果控制到按钮级别的权限怎么做？
    封装了一个指令权限，能简单快速的实现按钮级别的权限判断。 v-permission

```HTML
<template>
  <!-- Admin can see this -->
  <el-tag v-permission="['admin']">admin</el-tag>

  <!-- Editor can see this -->
  <el-tag v-permission="['editor']">editor</el-tag>

  <!-- Editor can see this -->
  <el-tag v-permission="['admin','editor']">Both admin or editor can see this</el-tag>
</template>

<script>
// 当然你也可以为了方便使用，将它注册到全局
import permission from '@/directive/permission/index.js' // 权限判断指令
export default{
  directives: { permission }
}
</script>
```

100. vue 使用 v-for 遍历对象时，是按什么顺序遍历的？如何保证顺序？
     1、会先判断是否有 iterator 接口，如果有循环执行 next()方法
     2、没有 iterator 的情况下，会调用 Object.keys()方法，在不同浏览器中，JS 引擎不能保证输出顺序一致
     3、保证对象的输出顺序可以把对象放在数组中，作为数组的元素
     ![顺序](https://user-images.githubusercontent.com/42334454/62844542-9fc46080-bcf4-11e9-9f52-9213e84f6536.png)
101. ue 的 extend（构造器）的理解，它主要是用来做什么的？
     extend 的作用是继承当前的 Vue 类，传入一个 extendOption 生成一个新的构造函数。在 extend 的时候会进行 mergeOption，融合 Vue 原型上的 baseOption，所以 extend 出来的子类也能使用 v-model、keep-alive 等全局性的组件。


    作用是生成组件类。**在挂载全局组件和设置了 components 属性的时候会使用到**。在生成 DOM 的时候会 new 实例化挂载。

102. 手写一个自定义指令及写出如何调用
103. vue 中什么是递归组件？举个例子说明下
     通过组件的 name 属性，调用自身。例如**生成树型菜单**。
104. 怎么修改 vue 打包后生成文件路径
     webpack：output.path
     vue-cli3: outputDir
105. 如果一个元素上同时存在 class 和:class 可以吗
     v-bind:class 指令也可以与普通的 class attribute 共存。:class 是动态的添加一个 class，与标签中再多添加一个 class 类似，所以不会出现错误。
106. 在 vue 中 watch 和 created 哪个先执行？为什么？
     官网的生命周期图中，init reactivity 是晚于 beforeCreate 但是早于 created 的。
     watch 加了 immediate，应当同 init reactivity 周期一同执行，早于 created。
     而正常的 watch，则是 mounted 周期后触发 data changes 的周期执行，晚于 created。
107. 如何解决 vue 打包 vendor 过大的问题

     1. 路由懒加载
     2. 在 webpack.base.conf.js 新增 externals 配置，表示不需要打包的文件，然后在 index.html 中通过 CDN 引入

        ```JS
        externals: {
            "vue": "Vue",
            "vue-router": "VueRouter",
            "vuex": "Vuex",
            "element-ui": "ELEMENT",
            "BMap": "BMap"
        }

        ```

108. 为什么我们写组件的时候可以写在.vue 里呢？可以是别的文件名后缀吗？
     配合相应的 loader 想咋写就咋写
109. 说说你对单向数据流和双向数据流的理解
     父传子的 props 就是单向数据流，v-model 就是双向数据流
110. vue 部署上线前需要做哪些准备工作？
     router 是不是 hash 是否需要配置 nginx , publicPath , 是不是要配置 cdn

     **publicPath**: 代表打包后，**index.html 里面引用资源的相对地址**
     后台访问时要加上 publicPath 地址 dist，即 http://localhost:8080/dist/index.html#
     publicPath

111. axios 同时请求多个接口，如果当 token 过期时，怎么取消后面的请求？
     Axios 取消请求 CancelToken
112. 说说你对 vue 的错误处理的了解
     分为 errorCaptured 与 errorHandler。
     errorCaptured 是组件内部钩子，可捕捉本组件与子孙组件抛出的错误，接收 error、vm、info 三个参数，return false 后可以阻止错误继续向上抛出。
     errorHandler 为全局钩子，使用 Vue.config.errorHandler 配置，接收参数与 errorCaptured 一致，2.6 后可捕捉 v-on 与 promise 链的错误，可用于统一错误处理与错误兜底。
113. 说说你对 vue 的表单修饰符.lazy 的理解
     v-model 默认的触发条件是 input 事件,加了.lazy 修饰符之后,v-model 会在 change 事件触发的时候去监听
114. 预渲染和 SSR(服务端渲染)有什么区别
     预渲染的使用场景更多是我们所说的静态页面的形式。服务端渲染适用于大型的、页面数据处理较多且较为复杂的、与服务端有数据交互的功能型网站，一个明显的使用场景就是电商网站。
115. 请解释下 Vue 中 slot 和 slot-scope 两者的区别
     slot 插槽分为默认插槽和具名插槽
     slot-scope 是作用域插槽(render props)，父组件可以直接使用子组件中定义 data 数据
116. 说说你对选项 el,template,render 的理解
     el: 把当前实例挂载在元素上
     template: 实例模版, 可以是.vue 中的 template, 也可以是 template 选项, 最终会编译成 render 函数
     render: 不需要通过编译的可执行函数

template 和 render, 开发时各有优缺点, 不过在线上尽量不要有 template

117. 为什么 Vue 被称为“渐进框架”？
     区别于 angular
     用户可以按需使用，使用的过程中，随着开发者的需求的丰富，从简单到复杂，vue 都能提供覆盖，降低了上手成本，随着业务量复杂了，就会使用上各种配套的全家桶，引入官方或第三方插件，但如果业务量不复杂，那就只用最核心的 api 即可，这样整个库的使用就很轻巧。
118. 说下你的 vue 项目的目录结构，如果是大型项目你该怎么划分结构和划分组件呢
     views 目录存放一级路由的组件，即视图组件
     Components 目录存放组件
     Store 存放 vuex 相关文件
     Router 目录存放路由相关文件
     Untils 目录存放工具 js 文件
     API 目录存放封装好的与后端交互的逻辑
     Assets 存放静态文件
119. provide 和 inject 原理
     在初始化阶段， Vue 会遍历当前组件 inject 选项中的所有键名，拿这个键名在上级组件中的 `_provided` 属性里面进行查找，如果找到了就使用对应的值，如果找不到则再向上一级查找，直到查找完根组件为止，最终如果没有找到则使用默认值，通过这样层层向上查找的方式最终实现 provide 和 inject 的数据传递机制。
     这里没有讲到 provide 的初始化，因为比较简单，就是将 provide 选项挂载在了 `_provided` 属性上。

以上就是 provide 和 inject 的实现原理。
