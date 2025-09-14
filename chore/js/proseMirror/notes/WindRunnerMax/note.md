https://www.zhihu.com/people/WindrunnerMax/posts
https://blog.touchczy.top/

- [命令模式](https://blog.touchczy.top/?p=/Patterns/%E5%91%BD%E4%BB%A4%E6%A8%A1%E5%BC%8F)
- 与其将 Virtual DOM 视为一种技术，不如说它是一种模式，人们提到它时经常是要表达不同的东西。
- Virtual DOM 是一种编程概念，在这个概念里，UI 以一种理想化的，或者说虚拟的表现形式被保存于内存中，并通过如 ReactDOM 等类库**使之与真实的 DOM 同步，这一过程叫做协调。**
- 维护一个使用 Js 对象表示的 Virtual DOM，与真实 DOM 一一对应。
  对前后两个 Virtual DOM 做 diff，生成变更 Mutation。
  把变更应用于真实 DOM，生成最新的真实 DOM。

- Vue 中$nextTick 方法将回调延迟到下次 DOM 更新循环之后执行，也就是在下次 DOM 更新循环结束之后执行延迟回调，在修改数据之后立即使用这个方法，能够获取更新后的 DOM。简单来说就是当数据更新时，在 DOM 中渲染完成后，执行回调函数。

- [低代码场景的状态管理方案](https://blog.touchczy.top/?p=/React/%E4%BD%8E%E4%BB%A3%E7%A0%81%E5%9C%BA%E6%99%AF%E7%9A%84%E7%8A%B6%E6%80%81%E7%AE%A1%E7%90%86%E6%96%B9%E6%A1%88)
  基于 Immer 以及 OT-JSON 实现原子化、可协同、高扩展的应用级状态管理方案。
  飞书的数据结构协同是使用 OT-JSON 来实现的，文本的协同则是借助了 EasySync 作为 OT-JSON 的子类型来实现的，以此来提供更高的扩展性。

  - transform 解决了 a 操作对 b 操作造成的影响，即维护因果关系。
  - 同步变更数据，异步渲染视图

- 无代码 NoCode 和低代码 LowCode 还是比较容易混淆的，在我的理解上，NoCode 强调自己编程给自己用，给用户的感觉是一个更强大的实用软件，是一个上层的应用，也就是说 `NoCode 需要面向非常固定的领域才能做到好用`。而对于 LowCode 而言，除了要考虑能用界面化的方式搭建流程，还要考虑在需要扩展的时候，把底层也暴露出来，`拥有更强的可定制化功能，也就是说相比 NoCode 可以不把使用场景限定得那么固定。`

- Vue 中的三种 Watcher 形态
  负责视图更新的 render watcher。
  执行计算属性更新的 computed watcher。
  用户注册的普通 watcher api。

  他们都是通过 class Watcher 类来实现的
