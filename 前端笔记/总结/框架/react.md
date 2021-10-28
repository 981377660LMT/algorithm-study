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
   为什么合成事件:跨平台、挂到 document 减少消耗、频繁解绑
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
