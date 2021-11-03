https://zhuanlan.zhihu.com/p/80045489 源码流程

1. 在标签使用箭头函数的问题

```JS
class LoggingButton extends React.Component {
  handleClick() {
    console.log('this is:', this);
  }

  render() {
    // 此语法确保 `handleClick` 内的 `this` 已被绑定。
    return (
      <button onClick={() => this.handleClick()}>
        Click me
      </button>
    );
  }
}

此语法问题在于每次渲染 LoggingButton 时都会创建不同的回调函数。在大多数情况下，这没什么问题，`但如果该回调函数作为 prop 传入子组件时，这些组件可能会进行额外的重新渲染`。我们通常建议在构造器中绑定或使用 class fields 语法来避免这类性能问题。
```

2. && 语法需要注意的点
   数字 0 依旧会被渲染

```JS
数字0，仍然会被 React 渲染。例如，以下代码并不会像你预期那样工作，因为当 props.messages 是空数组时，0 仍然会被渲染
<div>
  {props.messages.length &&
    <MessageList messages={props.messages} />
  }
</div>
解决方式：确保 && 之前的表达式总是布尔值：
<div>
  {props.messages.length > 0 &&
    <MessageList messages={props.messages} />
  }
</div>
```

3. 受控和非受控组件
   我们对某个组件状态的掌控，它的值是否只能由用户设置，而不能通过代码控制
   **受控组件的定义**
   在 HTML 的表单元素中，它们通常自己维护一套 state，并随着用户的输入自己进行 UI 上的更新，这种行为是不被我们程序所管控的。而如果**将 React 里的 state 属性和表单元素的值建立依赖关系，再通过 onChange 事件与 setState()结合更新 state 属性**(v-model)，就能达到控制用户输入过程中表单发生的操作。被 React 以这种**方式控制取值的表单输入元素就叫做受控组件**。
   对于 **select 表单元素**来说，React 中将其转化为受控组件可能和原生 HTML 中有一些区别。
   使用 React 受控组件来写的话就不用那么麻烦了，因为它允许在根 select 标签上使用 value 属性，去控制选中了哪个。这样的话，对于我们也更加便捷，在用户每次重选之后我们只需要在根标签中更新它，就像是这个案例
   对于受控组件，我们需要为每个状态更新(例如 this.state.username)编写一个事件处理程序(例如 this.setState({ username: e.target.value }))。

```JS
class SelectComponent extends React.Component {
  constructor(props) {
    super(props);
    this.state = { value: 'cute' };
  }
  handleChange(event) {
    this.setState({value: event.target.value});
  }
  handleSubmit(event) {
    alert('你今日相亲对象的类型是: ' + this.state.value);
    event.preventDefault();
  }
  render() {
    return (
      <form onSubmit={(e) => this.handleSubmit(e)}>
        <label>
          你今日相亲对象的类型是:
          <select value={this.state.value} onChange={(e) => this.handleChange(e)}>
            <option value="sunshine">阳光</option>
            <option value="handsome">帅气</option>
            <option value="cute">可爱</option>
            <option value="reserved">高冷</option>
          </select>
        </label>
        <input type="submit" value="提交" />
      </form>
    );
  }
}
export default SelectComponent;

```

表单
https://github.com/LinDaiDai/niubility-coding-js/blob/master/%E6%A1%86%E6%9E%B6-%E5%BA%93/React/%E5%8F%97%E6%8E%A7%E5%92%8C%E9%9D%9E%E5%8F%97%E6%8E%A7%E7%BB%84%E4%BB%B6%E7%9C%9F%E7%9A%84%E9%82%A3%E4%B9%88%E9%9A%BE%E7%90%86%E8%A7%A3%E5%90%97.md

**非受控组件**
那么还有一种场景是：我们仅仅是想要获取某个表单元素的值，而不关心它是如何改变的。对于这种场景，我们有什么应对的方法吗 🤔️？

```JS
import React, { Component } from 'react';

export class UnControll extends Component {
  constructor (props) {
    super(props);
    this.inputRef = React.createRef();
  }
  handleSubmit = (e) => {
    console.log('我们可以获得input内的值为', this.inputRef.current.value);
    e.preventDefault();
  }
  render () {
    return (
      <form onSubmit={e => this.handleSubmit(e)}>
        <input defaultValue="lindaidai" ref={this.inputRef} />
        <input type="submit" value="提交" />
      </form>
    )
  }
}

```

特殊的文件 file 标签
另外在 input 中还有一个比较特殊的情况，那就是 file 类型的表单控件。
**对于 file 类型的表单控件它始终是一个不受控制的组件，因为它的值只能由用户设置**，而不是以编程方式设置。

```JS
import React, { Component } from 'react';

export default class UnControll extends Component {
  constructor (props) {
    super(props);
    this.state = {
      files: []
    }
  }
  handleSubmit = (e) => {
    e.preventDefault();
  }
  handleFile = (e) => {
    console.log(e.target.files);
    const files = [...e.target.files];
    console.log(files);
    this.setState({
      files
    })
  }
  render () {
    return (
      <form onSubmit={e => this.handleSubmit(e)}>
        <input type="file" value={this.state.files} onChange={(e) => this.handleFile(e)} />
        <input type="submit" value="提交" />
      </form>
    )
  }
}

```

在选择了文件之后，我试图用 setState 来更新，结果却报错了
`Failed to set the 'value' property on 'HTMLInputElement'`
所以我们**应当使用非受控组件的方式**来获取它的值，可以这样写：

```JS
import React, { Component } from 'react';

export default class FileComponent extends Component {
  constructor (props) {
    super(props);
    this.fileRef = React.createRef();
  }
  handleSubmit = (e) => {
    console.log('我们可以获得file的值为', this.fileRef.current.files);
    e.preventDefault();
  }
  render () {
    return (
      <form onSubmit={e => this.handleSubmit(e)}>
        <input type="file" ref={this.fileRef} />
        <input type="submit" value="提交" />
      </form>
    )
  }
}

```

这里获取到的 files 是一个数组哈，当然，如果你没有开启多选的话，这个数组的长度始终是 1，开启多选也非常简单，只需要添加 multiple 属性即可：

```JS
<input type="file" multiple ref={this.fileRef} />

```

**实际的应用场景**
，绝大部分时候推荐使用受控组件来实现表单，因为在受控组件中，表单数据由 React 组件负责处理；当然如果选择受受控组件的话，表单数据就由 DOM 本身处理。

4. hooks 的出现就是为了取代 HOC(高阶组件就是一个没有副作用的纯函数。) HOC 逻辑复杂 少用 HOC
   这也论证了`多用组合,少用继承`的设计原则
   发展 Mixin=>HOC=>Hooks

5. 事件机制
   事件绑定到 root 事件冒泡实例化生成统一的 SyntheticEvent
   再 dispatchEvent event 对象交由对应的 handler 处理
   注意: react 的事件体系, **不是全部都通过事件委托来实现的**. 有一些特殊情况, 是直接绑定到对应 DOM 元素上的(如:scroll, load),
   **为什么合成事件**:跨平台、挂到 document 减少消耗、频繁解绑
6. batchUpdate 机制(基于 transaction 机制,影响 setState 的 `异步`)
   setState 机制，batchUpdate 机制，transaction 机制

   1. newState 存入 pending 队列
   2. 调用函数前 isBatchingUpdates=true,调用函数后 isBatchingUpdates=false（transaction 机制;**类似于 python 的 上下文管理器 contextManager**）。执行 setState 时 判断是否处于 patchUpdate 是则保存组件到 dirtyComponents 不是则遍历 dirtyComponents 调用 updateComponent

```Python
这段代码的作用是任何对列表的修改只有当所有代码运行完成并且不出现异常的情况下才会生效
from contextlib import contextmanager

@contextmanager
def list_transaction(orig_list):
    working = list(orig_list)
    yield working
    orig_list[:] = working

>>> items = [1, 2, 3]
>>> with list_transaction(items) as working:
...     working.append(4)
...     working.append(5)
...
>>> items
[1, 2, 3, 4, 5]
>>> with list_transaction(items) as working:
...     working.append(6)
...     working.append(7)
...     raise RuntimeError('oops')
...
Traceback (most recent call last):
    File "<stdin>", line 4, in <module>
RuntimeError: oops
>>> items
[1, 2, 3, 4, 5]
>>>
```

7. Fiber 如何优化性能
8. react 渲染和更新的过程
   jsx 渲染、setState 更新页面(React 默认全部重新渲染)
9. redux 精通

10. react 事件和 DOM 事件区别
11. react16 所有事件挂载到 document react17 绑定到 root **有利于多个 react 版本共存**
12. event 是 SyntheticEvent ，模拟出来 DOM 事件所有能力
    event.nativeEvent 是原生事件对象
13. dispatchEvent 机制
14. React 性能优化
    1. Key
    2. 销毁
    3. 异步组件
    4. pure/memo
15. jsx 本质
    JSX 是 ECMAScript 一个类似 XML 的语法扩展。基本上，它只是为 React.createElement() 函数提供语法糖
    createElement 返回 vNode

    ```JS
        // DOM Elements
        function createElement<P extends DOMAttributes<T>, T extends Element>(
        type: string,
        props?: ClassAttributes<T> & P | null,
        ...children: ReactNode[]): DOMElement<P, T>;

        // Custom components
        function createElement<P extends {}>(
        type: FunctionComponent<P> | ComponentClass<P> | string,
        props?: Attributes & P | null,
        ...children: ReactNode[]): ReactElement<P>;
    ```

16. 组件通信
    1. props 传数据/传函数
17. setState 为何使用不可变值

```JS

// 不可变值（函数式编程，纯函数） - 数组
const list5Copy = this.state.list5.slice()
list5Copy.splice(2, 0, 'a') // 中间插入/删除
this.setState({
    list1: this.state.list1.concat(100), // 追加
    list2: [...this.state.list2, 100], // 追加
    list3: this.state.list3.slice(0, 3), // 截取
    list4: this.state.list4.filter(item => item > 100), // 筛选
    list5: list5Copy // 其他操作
})
// 注意，不能直接对 this.state.list 进行 push pop splice 等，这样违反不可变值

// 不可变值 - 对象
this.setState({
    obj1: Object.assign({}, this.state.obj1, {a: 100}),
    obj2: {...this.state.obj2, a: 100}
})
// 注意，不能直接对 this.state.obj 进行属性设置，这样违反不可变值

```

14. setState 批量更新合并

```JS
    第四，state 异步更新的话，更新前会被合并 ----------------------------

    // 传入对象，会被合并（类似 Object.assign ）。执行结果只一次 +1
    this.setState({
        count: this.state.count + 1
    })
    this.setState({
        count: this.state.count + 1
    })
    this.setState({
        count: this.state.count + 1
    })

    // 传入函数，不会被合并。执行结果是 +3
    this.setState((prevState, props) => {
        return {
            count: prevState.count + 1
        }
    })
    this.setState((prevState, props) => {
        return {
            count: prevState.count + 1
        }
    })
    this.setState((prevState, props) => {
        return {
            count: prevState.count + 1
        }
    })
```

15. setState 同步还是异步:无所谓，看是否命中 batchUpdate 机制(生命周期，react 中注册的事件即 React 可以管理的入口;定时器，自定义 DOM 事件 React 管不到的入口)

```JS

    第三，setState 可能是异步更新（有可能是同步更新） ----------------------------

    this.setState({
        count: this.state.count + 1
    }, () => {
        // 联想 Vue $nextTick - DOM
        console.log('count by callback', this.state.count) // 回调函数中可以拿到最新的 state
    })
    console.log('count', this.state.count) // 异步的，拿不到最新值

    // setTimeout 中 setState 是同步的
    setTimeout(() => {
        this.setState({
            count: this.state.count + 1
        })
        console.log('count in setTimeout', this.state.count)
    }, 0)

    自己定义的 DOM 事件，setState 是同步的。再 componentDidMount 中
```

16. React.StrictMode 带来的问题
17. ReactRom.createPortal
    使用场景:fixed 需要放在 body 第一层级
18. 异步组件:
    React.lazy(()=>import())
    React.Suspense
19. SCU 默认返回什么 (true,可以渲染)
    React 默认就是全部重新渲染
20. SCU 要配合不可变值

```JS
  onSubmitTitle = (title) => {
        // 正确的用法
        this.setState({
            list: this.state.list.concat({
                id: `id-${Date.now()}`,
                title
            })
        })

        // // 为了演示 SCU ，故意写的错误用法
        // this.state.list.push({
        //     id: `id-${Date.now()}`,
        //     title
        // })
        // this.setState({
        //     list: this.state.list
        // })
    }
```

21. immutabel.js  
    彻底的不可变值
    ```JS
    const arr = [1, 2, 3]
    arr.push(4) // 被修改
    const arr1 = a.concat(4) // 重新生成 arr1 ，但 arr 是一直不变的
    ```
22. 什么是 renderProps
    类似于 vue 里的作用域插槽 传下去是一个函数 可以获取子组件里的 props 或 state
23. ErrorBounday 的缺点
    React 16 提供了一个内置函数 componentDidCatch，使用它可以非常简单的获取到 react 下的错误信息

    但是它无法捕获

    1. 事件处理器
    2. 异步代码
    3. 服务端的渲染代码
    4. 在 error boundaries 区域内的错误

24. 当地址栏改变 url，组件的更新渲染都经历了什么？
    拿 history 模式做参考。当 url 改变，首先触发 histoy，调用事件监听 popstate 事件， 触发回调函数 handlePopState，触发 history 下面的 setstate 方法，产生新的 location 对象，然后通知 Router 组件更新 location 并通过 **context 上下文**传递，switch 通过传递的更新流，**匹配出符合的 Route 组件渲染**，最后有 Route 组件取出 context 内容，传递给渲染页面，渲染更新。
25. 当我们调用 history.push 方法，切换路由，组件的更新渲染又都经历了什么呢？
    我们还是拿 history 模式作为参考，当我们调用 history.push 方法，首先调用 history 的 push 方法，通过 history.pushState 来改变当前 url，接下来触发 history 下面的 setState 方法，接下来的步骤就和上面一模一样了，这里就不一一说了。
26. React 如何区分 Class 和 Function
    1. 检查原型链上的 render 方法
    2. React 为基类增加了一个特别的标记
27. 为什么 React 使用 className 而不是 class 属性?
    **class 是 JavaScript 中的关键字**，而 JSX 是 JavaScript 的扩展。这就是为什么 React 使用 className 而不是 class 的主要原因。传递一个字符串作为 className 属性。
    在实际项目中，我们经常使用`classnames`库来方便我们操作 className。
28. 在 React v16 中的错误边界是什么?
    错误边界是在其子组件树中的任何位置捕获 JavaScript 错误、记录这些错误并显示回退 UI 而不是崩溃的组件树的组件。
    如果一个类组件定义了一个名为 **componentDidCatch(error, info)** 或 **static getDerivedStateFromError()** 新的生命周期方法，则该类组件将成为错误边界：
29. 为什么我们需要将函数传递给 setState() 方法?

```JS
假设初始计数值为零。在连续三次增加操作之后，该值将只增加一个。
// assuming this.state.count === 0
this.setState({ count: this.state.count + 1 })
this.setState({ count: this.state.count + 1 })
this.setState({ count: this.state.count + 1 })
// this.state.count === 1, not 3

如果将函数传递给 setState()，则 count 将正确递增。
this.setState((prevState, props) => ({
  count: prevState.count + props.increment
}))
// this.state.count === 3 as expected
```

30. 是否可以在不调用 setState 方法的情况下，强制组件重新渲染?
    默认情况下，当组件的状态或属性改变时，组件将重新渲染。如果你的 render() 方法依赖于其他数据，你可以通过调用 forceUpdate() 来告诉 React，当前组件需要重新渲染。

```JS
component.forceUpdate(callback)
```

31. 如何有条件地应用样式类?
    模板字符串

```JSX
<div className={`btn-panel ${this.props.visible ? 'show' : 'hidden'}`}>
```

32. 如何使用 React label 元素?
    因为 for 是 JavaScript 的保留字，请使用 htmlFor 来替代。

```JSX
<label htmlFor={'user'}>{'User'}</label>
<input type={'text'} id={'user'} />
```

33. 在 React 状态中删除数组元素的推荐方法是什么?
    Array.prototype.filter() `不可变`
34. 如何用 React 漂亮地显示 JSON?
    我们可以使用 <pre> 标签，以便保留 JSON.stringify() 的格式：

```JSx
const data = { name: 'John', age: 42 }

class User extends React.Component {
  render() {
    return (
      <pre>
        {JSON.stringify(data, null, 2)}
      </pre>
    )
  }
}

React.render(<User />, document.getElementById('container'))
```

35. **为什么你不能更新 React 中的 props**
    React 的哲学是 props 应该是 **immutable** 和 **top-down**。这意味着父级可以向子级发送任何属性值，但子级不能修改接收到的属性。
36. 为什么 React 组件名称必须以大写字母开头?
    在 JSX 中，小写标签被认为是 HTML 标签。但是，含有 . 的大写和小写标签名却不是。

```JS
<component /> 将被转换为 React.createElement('component') (i.e, HTML 标签)
<obj.component /> 将被转换为 React.createElement(obj.component)
<Component /> 将被转换为 React.createElement(Component)
```

37. Redux 的核心原则是什么
    Redux 遵循三个基本原则：

    1. 单一数据来源： 整个应用程序的状态存储在单个对象树中。单状态树可以更容易地跟踪随时间的变化并调试或检查应用程序。
    2. 状态是只读的： 改变状态的唯一方法是发出一个动作，一个描述发生的事情的对象。这可以确保视图和网络请求都不会直接写入状态。
    3. 使用纯函数进行更改： 要指定状态树如何通过操作进行转换，您可以编写 reducers。Reducers 只是纯函数，它将先前的状态和操作作为参数，并返回下一个状态。
       我可以在 reducer 中触发一个 Action 吗?**-不可以**

38. 为什么 Redux 状态函数称为 reducers ?
    Reducers 总是返回状态的累积（基于所有先前状态和当前 Action）
39. 在 React 中 registerServiceWorker 的用途是什么?
    默认情况下，React 会为你创建一个没有任何配置的 service worker。Service worker 是一个 Web API，它帮助你缓存资源和其他文件，以便当用户离线或在弱网络时，他/她仍然可以在屏幕上看到结果，因此，它可以帮助你建立更好的用户体验，这是你目前应该了解的关于 Service worker 的内容。

```JS
   import React from 'react';
   import ReactDOM from 'react-dom';
   import App from './App';
   import registerServiceWorker from './registerServiceWorker';

   ReactDOM.render(<App />, document.getElementById('root'));
   registerServiceWorker();

```

40. 如何确保钩子遵循正确的使用规则?

```
npm install eslint-plugin-react-hooks@next
```

```JS
// Your ESLint configuration
{
  "plugins": [
    // ...
    "react-hooks"
  ],
  "rules": {
    // ...
    "react-hooks/rules-of-hooks": "error"
  }
}

```

41. React 16 中未捕获的错误的行为是什么?
    在 React 16 中，未被任何错误边界捕获的错误将导致整个 React 组件树的卸载。这一决定背后的原因是，与其显示已损坏的界面，不如完全移除它。例如，**对于支付应用程序来说，显示错误的金额比什么都不提供更糟糕**。
42. render 方法可能返回的类型是什么?
43. 什么是基于路由的代码拆分?

```JS
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';
import React, { Suspense, lazy } from 'react';

const Home = lazy(() => import('./routes/Home'));
const About = lazy(() => import('./routes/About'));

const App = () => (
  <Router>
    <Suspense fallback={<div>Loading...</div>}>
      <Switch>
        <Route exact path="/" component={Home}/>
        <Route path="/about" component={About}/>
      </Switch>
    </Suspense>
  </Router>
);

```

44. portals 的典型使用场景是什么?
    对话框、全局消息通知、悬停卡和工具提示。
45. 深入理解 JSX
    JSX 在编译时会被 Babel 编译为 React.createElement 方法。
    这也是为什么在每个使用 JSX 的 JS 文件中，你必须显式的声明
    import React from 'react';
    否则在运行时该模块内就会报未定义变量 React 的错误。
    **在 React17 中，已经不需要显式导入 React 了**
    老的：

    ```JS
    import React from 'react';

    function App() {
      return <h1>Hello World</h1>;
    }

    import React from 'react';

    function App() {
      return React.createElement('h1', null, 'Hello world');
    }
    ```

    新的

    ```JS
    function App() {
      return <h1>Hello World</h1>;
    }

    // 新的 JSX 转换不会将 JSX 转换为 React.createElement，而是自动从 React 的 package 中引入新的入口函数并调用。
    // 由编译器引入（禁止自己引入！）
    import {jsx as _jsx} from 'react/jsx-runtime';

    function App() {
      return _jsx('h1', { children: 'Hello world' });
    }
    ```

46. API React.isValidElement
    **具有$$typeof 的对象**

```JS
export function isValidElement(object) {
  return (
    typeof object === 'object' &&
    object !== null &&
    object.$$typeof === REACT_ELEMENT_TYPE
  );
}
JSX在运行时的返回结果（即React.createElement()的返回值）都是React Element。
```

47. 如何区分函数组件和类组件?
    React 通过 ClassComponent 实例原型上的 isReactComponent 变量判断是否是 ClassComponent。

```JS
ClassComponent.prototype.isReactComponent = {};
```

48. JSX 与 Fiber 节点的关系
    JSX 是一种描述当前组件内容的数据结构，他不包含组件 schedule、reconcile、render 所需的相关信息。
    比如如下信息就不包括在 JSX 中：

        - 组件在更新中的优先级
        - 组件的 state
        - 组件被打上的用于 Renderer 的标记

    在组件 mount 时，Reconciler 根据 JSX 描述的组件内容生成组件对应的 Fiber 节点。
    在 update 时，Reconciler 将 JSX 与 Fiber 节点保存的数据对比，生成组件对应的 Fiber 节点，并根据对比结果为 Fiber 节点打上标记。

49. React 的模式?

当前 React 共有三种模式：

- legacy，这是当前 React 使用的方式。当前没有计划删除本模式，但是这个模式可能不支持一些新功能。
  `ReactDOM.render(<App />, rootNode)`
- blocking，开启部分 concurrent 模式特性的中间模式。目前正在实验中。作为迁移到 concurrent 模式的第一个步骤。
  `ReactDOM.createBlockingRoot(rootNode).render(<App />)`
- concurrent，面向未来的开发模式。我们之前讲的任务中断/任务优先级都是针对 concurrent 模式。
  `ReactDOM.createRoot(rootNode).render(<App />)`

模式的变化影响整个应用的工作方式，所以无法只针对某个组件开启不同模式。

50. createElement 做了什么?(jsx->babel 转为 createElement->ReactElement)
    分离 props 与特殊属性

```JS
config 包括 props/key/ref 等
function createElement(type, config, children) {
  ...

  return ReactElement(
    type,
    key,
    ref,
    self,
    source,
    ReactCurrentOwner.current,
    props,
  );
}

const ReactElement = function(type, key, ref, self, source, owner, props) {
  const element = {
    // This tag allows us to uniquely identify this as a React Element
    $$typeof: REACT_ELEMENT_TYPE,

    // Built-in properties that belong on the element
    type: type,
    key: key,
    ref: ref,
    props: props,

    // Record the component responsible for creating this element.
    _owner: owner,
  };

  ...

   return element;
}

```

51. react 怎么检测开发者使用了错误的 props:
    `createElement`中 使用 Object.defineProperty 定义 key 和 ref 的属性描述符，dev 环境警告

```JS

function defineRefPropWarningGetter(props, displayName) {
  const warnAboutAccessingRef = function() {
    if (__DEV__) {
      if (!specialPropRefWarningShown) {
        specialPropRefWarningShown = true;
        console.error(
          '%s: `ref` is not a prop. Trying to access it will result ' +
            'in `undefined` being returned. If you need to access the same ' +
            'value within the child component, you should pass it as a different ' +
            'prop. (https://reactjs.org/link/special-props)',
          displayName,
        );
      }
    }
  };
  warnAboutAccessingRef.isReactWarning = true;
  Object.defineProperty(props, 'ref', {
    get: warnAboutAccessingRef,
    configurable: true,
  });
}
```

52. render 方法?

```JS

export function render(
  element: React$Element<any>,
  container: Container,
  callback: ?Function,
) {
  if (__DEV__) {
    console.error(
      'ReactDOM.render is no longer supported in React 18. Use createRoot ' +
        'instead. Until you switch to the new API, your app will behave as ' +
        "if it's running React 17. Learn " +
        'more: https://reactjs.org/link/switch-to-createroot',
    );
  }

  if (!isValidContainerLegacy(container)) {
    throw new Error('Target container is not a DOM element.');
  }

  if (__DEV__) {
    const isModernRoot =
      isContainerMarkedAsRoot(container) &&
      container._reactRootContainer === undefined;
    if (isModernRoot) {
      console.error(
        'You are calling ReactDOM.render() on a container that was previously ' +
          'passed to ReactDOM.createRoot(). This is not supported. ' +
          'Did you mean to call root.render(element)?',
      );
    }
  }
  return legacyRenderSubtreeIntoContainer(
    null,
    element,
    container,
    false,
    callback,
  );
}

```

53. legacyRenderSubtreeIntoContainer 干了什么？
    - 初始化 fiber 数据结构:构建 fiberRoot 与 rootFiber
    - 将第三个参数 callback this 绑定到 app 实例
    - 调用 updateContainer

```JS

// 初始化fiber数据结构:构建fiberRoot与rootFiber
function legacyRenderSubtreeIntoContainer(
  parentComponent: ?React$Component<any, any>,
  children: ReactNodeList,
  container: Container,
  forceHydrate: boolean,
  callback: ?Function,
) {
  if (__DEV__) {
    topLevelUpdateWarnings(container);
    warnOnInvalidCallback(callback === undefined ? null : callback, 'render');
  }

  // 构建fiberRoot与rootFiber
  // 通过container._reactRootContainer有无来区分mount还是update
  let root = container._reactRootContainer;
  let fiberRoot: FiberRoot;
  if (!root) {
    // Initial mount
    root = container._reactRootContainer = legacyCreateRootFromDOMContainer(
      container,
      forceHydrate,
    );
    fiberRoot = root;

    // 将callback里的this绑定到实例
    if (typeof callback === 'function') {
      const originalCallback = callback;
      callback = function() {
        const instance = getPublicRootInstance(fiberRoot);
        originalCallback.call(instance);
      };
    }
    // Initial mount should not be batched.
    // 初始化渲染不能批量更新，因为批量更新是异步的可以被打断，而初始化要尽量快
    flushSync(() => {
      updateContainer(children, fiberRoot, parentComponent, callback);
    });
  } else {
    fiberRoot = root;
    if (typeof callback === 'function') {
      const originalCallback = callback;
      callback = function() {
        const instance = getPublicRootInstance(fiberRoot);
        originalCallback.call(instance);
      };
    }
    // Update
    updateContainer(children, fiberRoot, parentComponent, callback);
  }

  // render方法第一个参数的DOM对象作为render方法的返回值
  // 渲染谁就返回谁的DOM对象
  return getPublicRootInstance(fiberRoot);
}
```

54.updateContainer 作用?

```JS

// 创建任务放到任务队列
export function updateContainer(
  element: ReactNodeList,
  container: OpaqueRoot,
  parentComponent: ?React$Component<any, any>,
  callback: ?Function,
): Lane {
  if (__DEV__) {
    onScheduleRoot(container, element);
  }

  // 计算任务过期时间
  // 这个过期时间和 requestIdleCallback 里的 options 的 timeout 作用一样
  // 防止任务因为优先级的原因一直被打断而无法执行
  // 到了过期时间就会强制执行该任务
  // 同步任务被设置成了1073741823 表示同步任务
  const current = container.current;
  const eventTime = requestEventTime();
  const lane = requestUpdateLane(current);

  if (enableSchedulingProfiler) {
    markRenderScheduled(lane);
  }

  const context = getContextForSubtree(parentComponent);
  if (container.context === null) {
    container.context = context;
  } else {
    container.pendingContext = context;
  }

  if (__DEV__) {
    if (
      ReactCurrentFiberIsRendering &&
      ReactCurrentFiberCurrent !== null &&
      !didWarnAboutNestedUpdates
    ) {
      didWarnAboutNestedUpdates = true;
      console.error(
        'Render methods should be a pure function of props and state; ' +
          'triggering nested component updates from render is not allowed. ' +
          'If necessary, trigger nested updates in componentDidUpdate.\n\n' +
          'Check the render method of %s.',
        getComponentNameFromFiber(ReactCurrentFiberCurrent) || 'Unknown',
      );
    }
  }

  // 创建任务
  const update = createUpdate(eventTime, lane);
  // Caution: React DevTools currently depends on this property
  // being called "element".
  update.payload = {element};

  callback = callback === undefined ? null : callback;
  if (callback !== null) {
    if (__DEV__) {
      if (typeof callback !== 'function') {
        console.error(
          'render(...): Expected the last optional `callback` argument to be a ' +
            'function. Instead received: %s.',
          callback,
        );
      }
    }
    update.callback = callback;
  }

  // 加入Fiber的更新队列 this.setState调用了这个方法
  enqueueUpdate(current, update, lane);
  const root = scheduleUpdateOnFiber(current, lane, eventTime);
  if (root !== null) {
    entangleTransitions(root, current, lane);
  }

  return lane;
}

```

55. scheduleUpdateOnFiber 作用？
    首先判断是否是无限循环，如果是则报错
    接着使用过期时间判断是否是同步任务
    经过一系列判断，最终调用同步任务入口的方法 `performSyncWorkOnRoot` 。

```JS
/**
 * 判断任务是否为同步 调用同步任务入口
 */
export function scheduleUpdateOnFiber(
  fiber: Fiber,
  expirationTime: ExpirationTime,
) {
  /**
   * fiber: 初始化渲染时为 rootFiber, 即 <div id="root"></div> 对应的 Fiber 对象
   * expirationTime: 任务过期时间 => 同步任务固定为 1073741823
   */
  /**
   * 判断是否是无限循环的 update 如果是就报错
   * 在 componentWillUpdate 或者 componentDidUpdate 生命周期函数中重复调用
   * setState 方法时, 可能会发生这种情况, React 限制了嵌套更新的数量以防止无限循环
   * 限制的嵌套更新数量为 50, 可通过 NESTED_UPDATE_LIMIT 全局变量获取
   */
  checkForNestedUpdates();
  // 开发环境下执行的代码 忽略
  warnAboutRenderPhaseUpdatesInDEV(fiber);
  // 遍历更新子节点的过期时间 返回 FiberRoot
  const root = markUpdateTimeFromFiberToRoot(fiber, expirationTime);
  if (root === null) {
    // 开发环境下执行 忽略
    warnAboutUpdateOnUnmountedFiberInDEV(fiber);
    return;
  }
  // 判断是否有高优先级任务打断当前正在执行的任务
  // 初始渲染时内部判断条件不成立 内部代码没有得到执行
  checkForInterruption(fiber, expirationTime);

  // 报告调度更新, 实际什么也没做，忽略
  recordScheduleUpdate();

  // 获取当前调度任务的优先级 数值类型 90-99 数值越大 优先级越高
  // 初始渲染时优先级为 97 表示普通优先级任务。
  // 这个变量在初始渲染时并没有用到，忽略
  const priorityLevel = getCurrentPriorityLevel();
  // 判断任务是否是同步任务 Sync的值为: 1073741823
  if (expirationTime === Sync) {
    if (
      // 检查是否处于非批量更新模式
      (executionContext & LegacyUnbatchedContext) !== NoContext &&
      // 检查是否没有处于正在进行渲染的任务
      (executionContext & (RenderContext | CommitContext)) === NoContext
    ) {
      // 在根上注册待处理的交互, 以避免丢失跟踪的交互数据
      // 初始渲染时内部条件判断不成立, 内部代码没有得到执行
      schedulePendingInteractions(root, expirationTime);
      // 同步任务入口点
      performSyncWorkOnRoot(root);
    } else {
      ensureRootIsScheduled(root);
      schedulePendingInteractions(root, expirationTime);
      if (executionContext === NoContext) {
        // Flush the synchronous work now, unless we're already working or inside
        // a batch. This is intentionally inside scheduleUpdateOnFiber instead of
        // scheduleCallbackForFiber to preserve the ability to schedule a callback
        // without immediately flushing it. We only do this for user-initiated
        // updates, to preserve historical behavior of legacy mode.
        flushSyncCallbackQueue();
      }
    }
  } else {
    ensureRootIsScheduled(root);
    schedulePendingInteractions(root, expirationTime);
  }
  // 初始渲染不执行
  if (
    (executionContext & DiscreteEventContext) !== NoContext &&
    // Only updates at user-blocking priority or greater are considered
    // discrete, even inside a discrete event.
    (priorityLevel === UserBlockingPriority ||
      priorityLevel === ImmediatePriority)
  ) {
    // This is the result of a discrete event. Track the lowest priority
    // discrete update per root so we can flush them early, if needed.
    if (rootsWithPendingDiscreteUpdates === null) {
      rootsWithPendingDiscreteUpdates = new Map([[root, expirationTime]]);
    } else {
      const lastDiscreteTime = rootsWithPendingDiscreteUpdates.get(root);
      if (lastDiscreteTime === undefined || lastDiscreteTime > expirationTime) {
        rootsWithPendingDiscreteUpdates.set(root, expirationTime);
      }
    }
  }
}

```

56. commit 阶段？
