https://gpingfeng.github.io/ahooks-analysis/
https://github.dev/GpingFeng/hooks/tree/guangping/read-code/packages/hooks/src/useRequest

一些实用的 hooks.

## useRequest

- useRequest 是一个强大的异步数据管理的 Hooks，React 项目中的网络请求场景使用 useRequest 就够了。
  - 通过插件式组织代码，核心代码极其简单，并且可以很方便的扩展出更高级的功能
    - 入口 useRequest。它负责的是初始化处理数据以及将结果返回。
    - Fetch。是整个 useRequest 的核心代码，它处理了整个请求的生命周期。
    - plugin。在 Fetch 中，会通过插件化机制在不同的时机触发不同的插件方法，拓展 useRequest 的功能特性。
    - utils 和 types.ts。提供工具方法以及类型定义。
- 基本所有的插件功能都是在一个请求的一个或者多个阶段中实现的，也就是说我们只需要在请求的相应阶段，执行我们的插件的逻辑，就能执行和完成我们插件的功能
  Fetch 类的代码会变得非常的精简，只需要完成整体流程的功能，所有额外的功能（比如重试、轮询等等）都交给插件去实现。

  ```ts
  // 执行插件中的某个事件（event），rest 为参数传入
  runPluginHandler(event: keyof PluginReturn<TData, TParams>, ...rest: any[]) {
    // @ts-ignore
    const r = this.pluginImpls.map((i) => i[event]?.(...rest)).filter(Boolean);
    return Object.assign({}, ...r);
  }
  ```

  - 单一职责，一个 plugin 只做一件事
  - 最好的模块是深的：他们有很多功能隐藏在简单的接口后。深模块是好的抽象，因为它只把自己内部的一小部分复杂度暴露给了用户。

- currentCount 变量，控制取消请求。

```ts
// 取消当前正在进行的请求
cancel() {
  // 设置 + 1，在执行 runAsync 的时候，就会发现 currentCount !== this.count，从而达到取消请求的目的
  this.count += 1;
  this.setState({
    loading: false,
  });

  // 执行 plugin 中所有的 onCancel 方法
  this.runPluginHandler('onCancel');
}

// 假如不是同一个请求，则返回空的 promise
if (currentCount !== this.count) {
  // prevent run.then when request is canceled
  return new Promise(() => {});
}
```

- onBefore、onRequest、onCancel、onSuccess/onError/onFinally

## useAntdTable

基于 useRequest 实现，封装了常用的 Ant Design Form 与 Ant Design Table 联动逻辑，并且同时支持 antd v3 和 v4

## useInfiniteScroll

点击加载更多或者说下拉加载更加功能

## useVirtualList

提供虚拟化列表能力的 Hook，用于解决展示海量数据渲染时首屏渲染缓慢和滚动卡顿问题。

## use-immer

## useUrlState

工具库中假如某个工具函数/hook 依赖于一个开发者可能并不会使用的包，而且这个包的体积还比较大的时候，`可以将这个工具函数/hook 独立成一个 npm 包，开发者使用的时候才进行安装`。
另外这种可以考虑使用 monoRepo 的包管理方法，方便进行文档管理以及一些公共包管理等。

## useAsyncEffect

## effect function 应该返回一个销毁函数（effect：是指 return 返回的 cleanup 函数），如果 useEffect 第一个参数传入 async，返回值则变成了 Promise，会导致 react 在调用销毁函数的时候报错。

useEffect 怎么支持 async...await...

- 在内部使用

```ts
useEffect(() => {
  const asyncFun = async () => {
    setPass(await mockCheck())
  }
  asyncFun()
}, [])
```

- 实现方法

  1. 延迟清除(有人认为不对)

  ```tsx
  function useAsyncEffect(
    effect: () => AsyncGenerator<void, void, void> | Promise<void>,
    // 依赖项
    deps?: DependencyList
  ) {
    useEffect(() => {
      const e = effect()
      // 这个标识可以通过 yield 语句可以增加一些检查点
      // 如果发现当前 effect 已经被清理，会停止继续往下执行。
      let cancelled = false
      // 执行函数
      async function execute() {
        // 如果是 Generator 异步函数，则通过 next() 的方式全部执行
        if (isAsyncGenerator(e)) {
          while (true) {
            const result = await e.next()
            // Generate function 全部执行完成
            // 或者当前的 effect 已经被清理
            if (result.done || cancelled) {
              break
            }
          }
        } else {
          await e
        }
      }
      execute()
      return () => {
        // 当前 effect 已经被清理
        cancelled = true
      }
    }, deps)
  }
  ```

  2. 取消

  ```tsx
  function useAsyncEffect(
    effect: (isCanceled: () => boolean) => Promise<void>,
    dependencies?: any[]
  ) {
    return useEffect(() => {
      let canceled = false
      effect(() => canceled)
      return () => {
        canceled = true
      }
    }, dependencies)
  }

  useAsyncEffect(
    async isCanceled => {
      const result = await doSomeAsyncStuff(stuffId)
      if (!isCanceled()) {
        // TODO: Still OK to do some effect, useEffect hasn't been canceled yet.
      }
    },
    [stuffId]
  )
  ```

## useInterval 和 useTimeout

通过 useEffect 的返回清除机制，开发者不需要关注清除定时器的逻辑，避免内存泄露问题。这点是很多开发者会忽略的点。

## useRafInterval 和 useRafTimeout

在`页面不可见的时候，不执行定时器`，可以选择 useRafInterval 和 useRafTimeout，其内部是使用 requestAnimationFrame 进行实现。

## useExternal

动态注入 JS 或 CSS 资源，useExternal 可以保证资源全局唯一。
其实现原理创建 link 标签加载 CSS 资源或者 script 标签加载 JS 资源。通过 `document.createElement 返回 Element 对象，监听该对象获取加载状态。`

## useFullscreen

.request(element, options?)。使一个元素全屏显示。默认元素是 <html>
.exit()。退出全屏。
.toggle(element, options?)。假如目前是全屏，则退出，否则进入全屏。
.on(event, function)。添加一个监听器，用于当浏览器切换到全屏或切换出全屏或出现错误时。event 支持 'change' 或者 'error'。另外两种写法：.onchange(function) 和 .onerror(function)。
.isFullscreen。判断是否是全屏。
.isEnabled。判断当前环境是否支持全屏。

## useLongPress

其主要原理是判断当前是否支持 touch 事件，假如支持，则监听 touchstart 和 touchend 事件。假如不支持，则监听 mousedown、mouseup 和 mouseleave 事件。根据定时器设置标识，判断是否达到长按，触发回调，从而实现长按事件。

---

js-cookie、query-string、screenfull 这几个库
