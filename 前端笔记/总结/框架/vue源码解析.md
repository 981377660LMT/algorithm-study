三个关键：响应式原理+虚拟 dom+编译

1. runtime 版本是什么
   通常我们利用 vue-cli 去初始化我们的 Vue.js 项目的时候会询问我们用 Runtime Only 版本的还是 Runtime + Compiler 版本
   Runtime Only:需要借助如 webpack 的 vue-loader 工具把 .vue 文件编译成 JavaScript，因为是在编译阶段做的，所以它只包含运行时的 Vue.js 代码
   Runtime + Compiler:我们如果没有对代码做预编译，但又使用了 Vue 的 template 属性并传入一个字符串，则需要在客户端编译模板

   ```JS
       // 需要编译器的版本
       new Vue({
           template: '<div>{{ hi }}</div>'
       })

       // 这种情况不需要
       new Vue({
           render (h) {
               return h('div', this.hi)
           }
       })
   ```

   Vue.js 2.0 中，最终渲染都是通过 render 函数，如果写 template 属性，则需要编译成 render 函数，那么这个编译过程会发生运行时，所以需要带有编译器的版本。

2. new Vue 发生了什么
   `src/core/instance/init.js`

```JS
function Vue (options) {
  if (process.env.NODE_ENV !== 'production' &&
    !(this instanceof Vue)
  ) {
    warn('Vue is a constructor and should be called with the `new` keyword')
  }
  this._init(options)
}
```

```JS
  vm._self = vm
  initLifecycle(vm)
  initEvents(vm)
  initRender(vm)
  callHook(vm, 'beforeCreate')
  initInjections(vm) // resolve injections before data/props
  initState(vm)
  initProvide(vm) // resolve provide after data/props
  callHook(vm, 'created')
  ..
   if (vm.$options.el) {
    vm.$mount(vm.$options.el)  // 如果有 el 属性，则调用 vm.$mount 方法挂载 vm，挂载的目标就是把模板渲染成最终的 DOM
  }
```

Vue 初始化主要就干了几件事情:调用`Vue.__init(options)`,合并配置，初始化生命周期，初始化事件中心，初始化渲染，初始化 data、props、computed、watcher 等等。
Vue 的初始化逻辑写的非常清楚，把不同的功能逻辑拆成一些单独的函数执行，让主线逻辑一目了然，这样的编程思想是非常值得借鉴和学习的

3. Vue 实例挂载的实现(mountComponent:`vm._render`+`vm._update`)

   1. 首先，它对 el 做了限制，Vue 不能挂载在 body、html 这样的根节点上
   2. 如果没有定义 render 方法，则会把 el 或者 template 字符串转换成 render 方法(在 线 编 译)
      这里我们要牢记，在 Vue 2.0 版本中，**所有 Vue 的组件的渲染最终都需要 render 方法**
      vm.\_render 最终是通过执行 createElement 方法并返回的是 vnode，它是一个虚拟 Node。
   3. 就是要把这个 VNode 渲染成一个真实的 DOM 并渲染出来，这个过程是通过 `vm._update` 完成的

4. Vue 的 `_update`调用时机:一个是首次渲染，一个是数据更新的时候
   `_update` 的核心就是调用 `vm.__patch__`
   `Vue.prototype.__patch__ = inBrowser ? patch : noop`
   在 web 平台上，是否是服务端渲染也会对这个方法产生影响。因为在服务端渲染中，没有真实的浏览器 DOM 环境，所以不需要把 VNode 最终转换成 DOM，因此是一个空函数，而在浏览器端渲染中，它指向了 patch 方法
   为何 Vue.js 源码绕了这么一大圈，把相关代码分散到各个目录
   因为前面介绍过，patch 是平台相关的，在 Web 和 Weex 环境，它们把虚拟 DOM 映射到 “平台 DOM” 的方法是不同的

![img](https://ustbhuangyi.github.io/vue-analysis/assets/new-vue.png)
init=>$mount=>compile=>render=>vNode=>patch=>Dom

5. Vue 响应式系统三大件
   Vue 响应系统，其核心有三点：observe、watcher、dep：
   **observe**:遍历 data 中的属性，使用 Object.defineProperty 的 get/set 方法对其进行数据劫持
   **dep**：每个属性拥有自己的消息订阅器 dep，用于存放所有订阅了该属性的观察者对象
   **watcher**：观察者（对象），通过 dep 实现对响应属性的监听，监听到结果后，主动触发自己的回调进行响应

   为什么 Vue.js 不能兼容 IE8 及以下浏览器
   ES5 的 Object.defineProperty

   `src/core/instance/state.js`
   在 Vue 的初始化阶段，`_init` 方法执行的时候，会执行 initState(vm) 方法

```JS
export function initState (vm: Component) {
  vm._watchers = []
  const opts = vm.$options
    if (opts.props) initProps(vm, opts.props)
    if (opts.methods) initMethods(vm, opts.methods)
    if (opts.data) {
    initData(vm)
  } else {
    observe(vm._data = {}, true /* asRootData */)
  }
    if (opts.computed) initComputed(vm, opts.computed)
    if (opts.watch && opts.watch !== nativeWatch) {
    initWatch(vm, opts.watch)
    }
}
```

initProps:遍历的过程主要做两件事情
一个是调用 **defineReactive** 方法把每个 prop 对应的值变成响应式
另一个是通过 **proxy** 把 `vm._props.xxx` 的访问代理到 vm.xxx 上

defineReactive 的功能就是定义一个响应式对象

```JS
export function defineReactive (
  obj: Object,
  key: string,
  val: any,
  customSetter?: ?Function,
  shallow?: boolean
) {
  const dep = new Dep()
  ...
```

initData:遍历的过程主要做两件事情
一个是调用 **observe** 方法，里面进而调用**defineReactive**把每个 data 也变成响应式
另一个是通过 **proxy** 把 `vm._data.xxx` 的访问代理到 vm.xxx 上

proxy 方法的实现很简单，通过 Object.defineProperty 把 target[sourceKey][key] 的读写变成了对 target[key] 的读写
observe 方法的作用就是给非 VNode 的对象类型数据添加一个 Observer，如果已经添加过则直接返回，否则在满足一定条件下去实例化一个 Observer 对象实例。
Observer 是一个类，它的作用是给对象的属性添加 getter 和 setter，用于依赖收集和派发更新：

```JS
export class Observer {
  value: any;
  dep: Dep;
  vmCount: number; // number of vms that has this object as root $data
  ...
```

## 响应式对象(Observer 先给对象的属性遍历添加 getter 和 setter,然后再用于依赖收集和派发更新,配置依赖观测)

响应式对象，核心就是利用 Object.defineProperty 给数据添加了 getter 和 setter，目的就是为了在我们访问数据以及写数据的时候能自动执行一些逻辑：**getter 做的事情是依赖收集，setter 做的事情是派发更新**
Dep:Subject Watcher:Observer

## 依赖收集(观察者 watcher 注册到 可观察对象 Dep.target)：

Dep 是整个 getter 依赖收集的核心
每一个数据都有自己的 Dep 类实例，用来管理依赖数据的 Watcher 类实例
getter 方法:对象什么时候被访问?=>**Vue 的 mount 过程**执行 `vm._render()`,这个时候就触发了数据对象的 getter,watcher 注册到 可观察对象 Dep.target,这样实际上已经完成了一个依赖收集的过程。

```JS
export default class Dep {
  static target: ?Watcher;
  id: number;
  subs: Array<Watcher>;
}
```

Dep 是一个 Class，它定义了一些属性和方法，**这里需要特别注意的是它有一个静态属性 target，这是一个全局唯一 Watcher**，这是一个非常巧妙的设计，因为在同一时间只能有一个全局的 Watcher 被计算，另外它的自身属性 subs 也是 Watcher 的数组。

## 派发更新(可观察对象 Dep.notify() 触发所有的 观察者 watcher 的 update())：

实际上就是当数据发生变化的时候，触发 setter 逻辑，watcher 对象在进行更新执行 update，内部主要执行了一个 queueWatcher 函数判重操作，相同的 watcher 对象只会被加入到 queue 队列一次。在 **nextTick 传入回调后执行所有 watcher 的 run，进而触发 Dep.notify()**。

Watcher 和 Dep 就是一个非常经典的观察者设计模式的实现

nextTick 解决的问题:让数据的变化到 DOM 的重新渲染是一个异步过程(不会改变数据就马上触发 watcher 的回调)，而是把这些 watcher 先添加到一个队列里，然后在 nextTick(flushSchedulerQueue)。
变成什么异步？
会优先使用 Promise 等 microtask，保证在同一个事件循环里面执行，这样页面只需要渲染一次。实在不行的话用 setTimeout 来兜底，虽然会造成二次渲染，但这也是最差的情况。

```JS
nextTick(flushSchedulerQueue)  // flushSchedulerQueue函数依次调用了wacther对象的run方法执行更新。即()=>watcher.run()。
...
setTimeout(flushCallbacks, 0)  // flushCallbacks调用回调函数
```

flushSchedulerQueue 函数中的 waiting

```JS
/** queueWatcher函数*/
let has = {};
let queue = [];
let waiting = false;

function queueWatcher (watcher: Watcher) {
  const id = watcher.id
  // 防止queue队列wachter对象重复
  if (has[id] == null) {
    has[id] = true
    queue.push(watcher)

    // 传递本次的更新任务
    if (!waiting) {
      waiting = true
      nextTick(flushSchedulerQueue)
    }
  }
}

/** flushSchedulerQueue函数 */
function flushSchedulerQueue () {
    let watcher, id;
    for (index = 0; index < queue.length; index++) {
        watcher = queue[index];
        id = watcher.id;
        has[id] = null;
        // 执行更新
        watcher.run();
    }
    // 更新完毕恢复标志位
    waiting = false;
}

waiting(pending)这个标记位代表我们是否已经向nextTick函数传递了更新任务，nextTick会在当前task结束后再去处理传入的回掉，只需要传递一次，更新完毕再重置这个标志位。
```

```JS
getData(res).then(()=>{
  this.xxx = res.data
  this.$nextTick(() => {
    // 这里我们可以获取变化后的 DOM
  })
})
```

总结:
数据发生变化的时候，触发 setter 逻辑，进而触发 **dep.notify()**，进而触发 watcher 的 **update** 方法，其内部主要执行了一个 **queueWatcher** 函数判重操作，相同的 watcher 对象只会被加入到队列一次。然后 **nextTick(flushSchedulerQueue)**，nextTick 里面回调调用**flushCallbacks**触发 watcher.run()，如果是渲染 wacher 结束.如果是用户 watcher，调用 this.cb()触发更新视图。`render 进而触发 getter，收集依赖。`

6. nextTick
   vue 异步更新，本质上是 js 事件机制的一种运用，优先考虑了具有高优先级的 microtask，为了兼容，又做了降级策略。(Promise=>MutationObserver=>setImmediate=>setTimeout)

   1. for 循环更新 count 数值，dom 会被更新 100 次吗？
      不会，因为 queueWatcher 函数做了过滤，相同的 watcher 对象不会被重复添加。

   2. nextTick 是如何做到监听 dom 更新完毕的？
      vue 用异步队列的方式来控制 DOM 更新和 nextTick 回调先后执行，保证了能在 dom 更新后在执行回调。

7. Vue.set
   当我们去给这个对象添加一个新的属性的时候，是不能够触发它的 setter 的
   Vue 为了解决这个问题，定义了一个全局 API Vue.set 方法

8. Vue 不能检测到以下变动的数组

```
1.当你利用索引直接设置一个项时，例如：vm.items[indexOfItem] = newValue
2.当你修改数组的长度时，例如：vm.items.length = newLength
```

9. 计算属性(computed watcher) VS 侦听属性(user watcher)
   **三种类型的 Watcher 对象**
   创建顺序：计算属性 Watcher、用户 Watcher (侦听器)、渲染 Watcher

## 计算属性 computed 是怎么实现的

computed 的依赖收集是借助 vue 的 watcher 来实现的，我们称之为 computed watcher,每一个计算属性会对应一个 computed watcher 对象,computed watcher 同时持有一个 **dep 实例**(每个响应式数据都有).
`vm._render`触发了计算属性的 getter，它会拿到计算属性对应的 watcher，然后执行 `watcher.depend()`,然后再执行 `watcher.evaluate()` 去求值.
判断 `this.dirty`，如果为 true 则通过 `this.get()` 求值，然后把 this.dirty 设置为 false。在求值过程中，会执行 `value = this.getter.call(vm, vm)`，这实际上就是执行了计算属性定义的 getter 函数
一旦我们对计算属性依赖的数据做修改，则会触发 setter 过程，通知所有订阅它变化的 watcher 更新，执行 `watcher.update()` 方法,通过 `this.get()` 求值。
**当计算属性最终计算的值发生变化**才会触发渲染 watcher 重新渲染

## 侦听属性 watch 是怎么实现的

watcher.run()，如果是渲染 wacher 结束.如果是用户 watcher，**调用 this.cb()触发更新视图**
设置了 deep 后会执行 traverse 函数，会有一定的性能开销，所以一定要根据应用场景权衡是否要开启这个配置.那么在执行了 traverse 后，我们再对 watch 的对象内部任何一个值做修改，也会调用 watcher 的回调函数了

```JS
get() {
  let value = this.getter.call(vm, vm)
  // ...
  if (this.deep) {
    traverse(value)
  }
}
```

10. 组件更新
    调用了 `vm._update`=>`vm.__patch__(prevVnode, vnode)`=>diff 算法
    原理图
    ![img](https://ustbhuangyi.github.io/vue-analysis/assets/reactive.png)
11. 编译
    模板到真实 DOM 渲染的过程，中间有一个环节是把**模板编译成 render 函数**，这个过程我们把它称作编译
    对编译过程的了解会让我们对 Vue 的指令、内置组件、slot 等有更好的理解
    编译的入口：

    ```JS
    compileToFunctions (
      template: string,
      options?: CompilerOptions,
      vm?: Component
    )
    ```

    1. 解析模板字符串生成 AST:正则解析字符串
       const ast = parse(template.trim(), options)
    2. 优化语法树:把一些 AST 节点优化成静态节点,这部分数据生成的 DOM 也不会变化，我们可以在 patch 的过程跳过对他们的比对
       optimize(ast, options)
       整个 optimize 的过程实际上就干 2 件事情，markStatic(root) **标记静态节点** ，markStaticRoots(root, false) **标记静态根**
    3. 生成代码:类似于 intermock 库,但是 intermock 处理函数写的太乱了
       const code = generate(ast, options)
       基本就是判断当前 AST 元素节点的属性执行不同的代码生成函数
       ```JS
        export function genElement (el: ASTElement, state: CodegenState): string {
        if (el.staticRoot && !el.staticProcessed) {
          return genStatic(el, state)
        } else if (el.once && !el.onceProcessed) {
          return genOnce(el, state)
        } else if (el.for && !el.forProcessed) {
          return genFor(el, state)
        } else if (el.if && !el.ifProcessed) {
          return genIf(el, state)
        } else if (el.tag === 'template' && !el.slotTarget) {
          return genChildren(el, state) || 'void 0'
        } else if (el.tag === 'slot') {
          return genSlot(el, state)
        }...
       ```

12. 编译入口逻辑之所以这么绕，是因为 Vue.js 在不同的平台下都会有编译的过程，因此编译过程中的依赖的配置 baseOptions 会有所不同。而编译过程会多次执行，但这同一个平台下每一次的编译过程配置又是相同的，为了**不让这些配置在每次编译过程都通过参数传入**，Vue.js 利用了函数**柯里化的技巧很好的实现了 baseOptions 的参数保留**。同样，Vue.js 也是**利用函数柯里化技巧把基础的编译过程函数抽出来**，通过 createCompilerCreator(baseCompile) 的方式把真正编译的过程和其它逻辑如对编译配置处理、缓存处理等剥离开，这样的设计还是非常巧妙的。

13. event
    编译的 parse 阶段对于属性的处理，在 processAttr 中
    Vue 支持 2 种事件类型，原生 DOM 事件和自定义事件，它们主要的区别在于添加和删除事件的方式不一样，并且**自定义事件的派发是往当前实例上派发**，但是可以利用在父组件环境定义回调函数来实现父子组件的通讯。另外要注意一点，只有组件节点才可以添加自定义事件，并且添加原生 DOM 事件需要使用 native 修饰符；而普通元素使用 .native 修饰符是没有作用的，也只能添加原生 DOM 事件。
14. v-model
    Vue 双向绑定的真正实现，但本质上就是一种语法糖，它即可以支持原生表单元素，也可以支持自定义组件。
    原生表单元素::value+@input 语法糖
    在组件的实现中，我们是可以配置子组件接收的 prop 名称，以及派发的事件名称
15. slot
    插槽分为普通插槽(children)和作用域插槽(一段函数)
    简单地说，两种插槽的目的都是让子组件 slot 占位符生成的内容由父组件来决定，但数据的作用域会根据它们 vnodes 渲染时机不同而不同。
    普通插槽是在**父组件编译和渲染阶段生成 vnodes**，所以数据的作用域是父组件实例，子组件渲染的时候直接拿到这些渲染好的 vnodes。而对于**作用域插槽，父组件在编译和渲染阶段并不会直接生成 vnodes，而是在父节点 vnode 的 data 中保留一个 scopedSlots 对象**，存储着不同名称的插槽以及它们对应的渲染函数，只有在编译和渲染子组件阶段才会执行这个渲染函数生成 vnodes，由于是在子组件环境执行的，所以对应的数据作用域是子组件实例。
16. keep-alive
17. transition
    对单一元素
18. transition-group
    对列表元素进行添加和删除，很好地帮助我们实现了列表的过渡效果
19. VuRouter
    Vue-Router 的能力十分强大，它支持 hash、history、abstract 3 种路由方式，提供了 <router-link> 和 <router-view> 2 种组件，还提供了简单的路由配置和一系列好用的 API。
    1. 路由注册 Vue.use(VueRouter)
       当用户执行 Vue.use(VueRouter) 的时候，实际上就是在执行 install 函数，为了确保 install 逻辑只执行一次，用了 install.installed 变量做已安装的标志位。另外用一个全局的 `_Vue` 来接收参数 Vue，因为作为 Vue 的插件对 Vue 对象是有依赖的，但又不能去单独去 import Vue，因为那样会增加包体积，所以就通过这种方式拿到 Vue 对象。
       Vue-Router 安装最重要的一步就是利用 Vue.mixin 去把 **beforeCreate 和 destroyed 钩子函数注入到每一个组件中**。
       接着给 Vue 原型上定义了 $router 和 $route 2 个属性的 get 方法，这就是为什么我们可以在组件实例上可以访问 this.$router 以及 this.$route，它们的作用之后介绍。
       接着又通过 Vue.component 方法定义了全局的 <router-link> 和 <router-view> 2 个组件，这也是为什么我们在写模板的时候可以使用这两个标签，它们的作用也是之后介绍。
    2. VueRouter 对象
    3. matcher
    4. 路径切换
20. Vuex
21. Vue.use 实现
22. Vue.mixin 实现
    ```JS
    export function initMixin (Vue: GlobalAPI) {
      Vue.mixin = function (mixin: Object) {
        this.options = mergeOptions(this.options, mixin)
        return this
      }
    }
    ```
