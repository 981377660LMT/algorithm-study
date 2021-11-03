函数参数中会反复出现<S\>、<T\>、<P\>、<I\>、<I\>，这些大写字母，react 约定他们对应的单词如下：  
state -> S -> 约定表示某种“数据”  
type -> T -> 约定表示某种“类型”  
props -> P -> 约定表示“属性传值对应的 props”  
initial -> I -> 约定表示某个“初始值”

1. 为什么要用 Hooks？
   1. 类组件的一些缺点：
      缺点一：复杂且不容易理解的“this”
      缺点二：组件数据状态逻辑不能重用、组件之间传值过程复杂
      子组件为了达到修改父组件中的数据状态，通常采用“高阶组件(HOC)”或“父组件暴露修改函数给子组件(render props)”这 2 种方式。 这 2 种方式都会让组件变得复杂且降低可复用性。
      缺点三：复杂场景下代码难以组织在一起
2. “勾住”？是什么意思
   Hook 本身单词意思是“钩子”，作用就是“勾住某些生命周期函数或某些数据状态，并进行某些关联触发调用”。
   “如何勾住”？ 在 React 底层代码中，是通过自定义 dispatcher，采用“发布订阅模式”实现的。
3. useState 是来解决类组件什么问题的？
   答：useState 能够解决类组件 所有自定义变量只能存储在 this.state 的问题。

   **关于闭包**

   ```JS
   for(let i=0; i<3; i++){
     setCount(count+1);
   }
   ```

   无论 for 循环执行几次，最终实际结果都将是 +1(闭包)。
   把代码修改为：

   ```JS
   for(let i=0; i<3; i++){
     setCount(prevData => {return prevData+1});
     //可以简化为 setCount(prevData => prevData+1);
   }
   ```

   就可以+3 了

   **数据类型为 Objcet 的修改方法**
   注意 state 必须是**不可变**数据
   正确的做法：
   我们需要先将 person 拷贝一份，修改之后再进行赋值。

   ```JS
   const [person, setPerson] = useState({name:'puxiao',age:34});
   setPerson({...person,age:18});
   ```

   ```JSX
   import React, { useState } from 'react';

   function Component() {

     const [person, setPerson] = useState({name:'puxiao',age:34});

     const nameChangeHandler = (eve) => {
       setPerson({...person,name:eve.target.value});
     }

     const ageChangeHandler = (eve) => {
       setPerson({...person,age:eve.target.value});
     }

     return <div>
       <input type='text' value={person.name} onChange={nameChangeHandler} />
       <input type='number' value={person.age} onChange={ageChangeHandler} />
       {JSON.stringify(person)}
     </div>
   }
   export default Component;

   ```

   **性能优化**
   通过 setXxx 设置新值，但是如果新值和当前值完全一样，那么会引发 React 重新渲染吗？
   通过 React 官方文档可以知道，当使用 setXxx 赋值时，Hook 会使用 Object.is()来对比当前值和新值，结果为 true 则不渲染，结果为 flase 就会重新渲染。
   如果可以，最好避免使用复杂类型的值

4. useEffect 概念解释
   都能勾住哪些生命周期函数？
   答：componentDidMount(组件被挂载完成后)、componentDidUpdate(组件重新渲染完成后)、componentWillUnmount(组件即将被卸载前)
5. useContext 概念解释
   他的作用是“勾住”获取由 React.createContext()创建、<XxxContext.Provider>添加设置的共享数据 value 值。useContext 可以替代<XxxContext.Consumer>标签，简化获取共享数据的代码。
6. useCallback/useMemo 解决传入 props 不同导致子组件重新渲染的问题
   注意此时子组件要 React.memo 进行 shallowDiff 否则 useCallback/useMemo 加了也白加
7. useRef 基础用法
   他的作用是“勾住”**某些组件挂载完成或重新渲染完成后才拥有的某些对象，并返回该对象的引用**。**该引用在组件整个生命周期中都固定不变，该引用并不会随着组件重新渲染而失效**。

   “某些组件挂载完成或重新渲染完成后才拥有的某些对象”：是什么?
   这句话中的“某些对象”主要分为 3 种：JSX 组件转换后对应的真实 **DOM** 对象、在 **useEffect 中创建的变量(例如定时器)**、子组件内自定义的函数(方法)。
   第 1：JSX 组件转换后对应的真实 DOM 对象：
   举例：假设在 JSX 中，有一个输入框<input type='text' />，这个标签最终将编译转换成真正的 html 标签中的<input type='text'/>。
   你应该知道以下几点：
   1、JSX 中小写开头的组件看似和原生 html 标签相似，但是并不是真的原生标签，依然是 react 内置组件。
   2、什么时候转换？ 虚拟 DOM 转化为真实 DOM
   3、什么时候可访问？组件挂载完成或重新渲染完成后

   注意注意注意!!!!!!!!!!!!!
   useRef **只适合“勾住”小写开头的类似原生标签的组件**。如果是自定义的 react 组件(自定义的组件必须大写字母开头)，那么是无法使用 useRef 的。
   所以才有`React.forwardRef`来传递 ref 来获取子组件里的原生标签

   第 2：在 useEffect 中创建的变量：

   ```JS
   useEffect(() => {
       let timer = setInterval(() => {
           setCount(prevData => prevData +1);
       }, 1000);
       return () => {
           clearInterval(timer);
       }
   },[]);
    上述代码中，请注意这个timer是在useEffect中才定义的。
    思考：useEffect 以外的地方，该如何获取这个 timer 的引用？
    答：用useRef
   ```

   第 3：子组件内自定义的函数(方法)
   需要结合 `useImperativeHandle` 才可以实现
   imperative:重要的,必要的,命令式

   “并返回该对象的引用”：
   上面的前 2 种情况，都提到用 useRef 来获取对象的引用。具体如何获取，稍后在 useRef 用法中会有演示。

   “该引用在组件整个生命周期中都固定不变”：
   假设通过 useRef 获得了该对象的引用，那么当 react 组件重新渲染后，如何保证该引用不丢失？
   答：react 在底层帮我们做了这个工作，我们只需要相信之前的引用可以继续找到目标对象即可。

8. useRef 是来解决什么问题的？
   答：useRef 可以“获取某些组件挂载完成或重新渲染完成后才拥有的某些对象”的引用，且保证该引用在组件整个生命周期内固定不变，都能准确找到我们要找的对象。
   具体已经在 useRef 中做了详细阐述，这里不再重复。
   补充说明：
   1、useRef 是针对函数组件的，如果是类组件则使用 React.createRef()。
   2、React.createRef()也可以在函数组件中使用。
   只不过 React.createRef 创建的引用不能保证每次重新渲染后引用固定不变。如果你只是使用 React.createRef“勾住”JSX 组件转换后对应的真实 DOM 对象是没问题的，但是如果想“勾住”在 useEffect 中创建的变量，那是做不到的。

   2 者都想可以“勾住”，只能使用 useRef。

9. useRef 的重载

```TS
function useRef<T>(initialValue: T): MutableRefObject<T>;  // 取值
function useRef<T>(initialValue: T|null): RefObject<T>;  // 取dom
function useRef<T = undefined>(): MutableRefObject<T | undefined>; // 懒人专用取值，不指定初始值

interface MutableRefObject<T> {
    current: T;  // 值
}
interface RefObject<T> {
    readonly current: T | null;  // dom
}

const timerRef = useRef<NodeJS.Timer>()  // 取值  NodeJS.Timer|undefined
const buttonRef = useRef<HTMLButtonElement>(null) // 取dom  HTMLButtonElement|null
```

结论：
1、如果需要对渲染后的 DOM 节点进行操作，必须使用 useRef。
2、如果需要对渲染后才会存在的变量对象进行某些操作，建议使用 useRef。
第 3 遍强调：useRef 只适合“勾住”小写开头的类似原生标签的组件。如果是自定义的 react 组件(自定义的组件必须大写字母开头)，那么是无法使用 useRef 的。

**在 TypeScript 中使用 useRef 创建计时器注意事项**：

```JS
timerRef.current = setInterval(() => {
        setCount((prevData) => { return prevData +1});
    }, 1000);
如果是在 TS 语法下，上面的代码会报错误：
不能将类型“Timeout”分配给类型“number”。
造成这个错误提示的原因是：

1. TypeScript 是运行在 Nodejs 环境下的，TS 编译之后的代码是运行在浏览器环境下的。
2. Nodejs 和浏览器中的 window 他们各自实现了一套自己的 setInterval
3. 原来代码 timerRef.current = setInterval( ... ) 中 setInterval 会被 TS 认为是 Nodejs 定义的 setInterval，而 Nodejs 中 setInterval 返回的类型就是 NodeJS.Timeout。
4. 所以，我们需要将上述代码修改为：timerRef.current = window.setInterval( ... )，明确我们调用的是 window.setInterval，而不是 Nodejs 的 setInterval。
```

10. useImperativeHandle 基础用法
    他的作用是“勾住”子组件中某些函数(方法)供父组件调用。

    - 如果子组件想调用父组件内函数，该怎么办？
      react 属于单向数据流，父组件可以通过属性传值，将父组件内的函数(方法)传递给子组件
    - **如果父组件想调用子组件中自定义的方法**，该怎么办？
      答：使用 useImperativeHandle()。

      **useImperativeHandle 是来解决什么问题的？**
      答：useImperativeHandle 可以让父组件获取并执行子组件内某些自定义函数(方法)。本质上其实是子组件将自己内部的函数(方法)通过 useImperativeHandle 添加到父组件中 useRef 定义的对象中。

      补充说明：
      1、useRef 创建引用变量
      2、React.forwardRef 将引用变量传递给子组件
      3、useImperativeHandle 将子组件内定义的函数作为属性，添加到父组件中的 ref 对象上。

      因此，如果想使用 useImperativeHandle，那么还要结合 useRef、React.forwardRef 一起使用。

```TS
function useImperativeHandle<T, R extends T>(ref: Ref<T>|undefined, init: () => R, deps?: DependencyList): void;
useImperativeHandle(ref,create,[deps])函数前2个参数为必填项，第3个参数为可选项。
第1个参数为父组件通过useRef定义的引用变量；
第2个参数为子组件要附加给ref的对象，该对象中的属性即子组件想要暴露给父组件的函数(方法)；
第3个参数为可选参数，为函数的依赖变量。凡是函数中使用到的数据变量都需要放入deps中，如果处理函数没有任何依赖变量，可以忽略第3个参数。
```

11. useLayoutEffect 基础用法
    你可以把 useLayoutEffect 等同于 componentDidMount、componentDidUpdate，因为他们调用阶段是相同的。而 useEffect 是在 componentDidMount、componentDidUpdate 调用之后才会触发的。
    也就是说，当组件所有 DOM 都渲染完成后，同步调用 useLayoutEffect，然后再调用 useEffect。
    useLayoutEffect 永远要比 useEffect 先触发完成。
    **那通常在 useLayoutEffect 阶段我们可以做什么呢？**
    答：在触发 useLayoutEffect 阶段时，页面全部 DOM 已经渲染完成，此时可以获取当前页面所有信息，包括页面显示布局等，你可以根据需求修改调整页面。
    在 react 官方文档中，明确表示只有在 useEffect 不能满足你组件需求的情况下，才应该考虑使用 useLayoutEffect。 官方推荐优先使用 useEffect。
    请注意：**如果是服务器渲染，无论 useEffect 还是 useLayoutEffect 都无法在 JS 代码加载完成之前执行，**因此都会收到错误警告。 服务器渲染时若想使用 useEffect，解决方案不在本章中讨论。
    **useLayoutEffect 是来解决什么问题的？**
    **答：useLayoutEffect 的作用是“当页面挂载或渲染完成时，再给你一次机会对页面进行修改”。**
12. 像 useState、useEffect、useContext、useReducer、useCallback、useMemo、useRef、useImperativeHandle、useLayoutEffect、useDebugValue 这 10 个 hook 是 react 默认自带的 hook，而所谓自定义 hook 就是由我们自己编写的 hook。
    **自定义 hook 利用闭包实现了一个`类`**
13. react hook 回顾

    定义变量
    useState()：定义普通变量
    useReducer()：定义有不同类型、参数的变量

    组件传值
    useContext()：定义和接收具有全局性质的属性传值对象，必须配合 React.createContext()使用

    对象引用
    useRef()：获取渲染后的 DOM 元素对象，可调用该对象原生 html 的方法，可能需要配合 React.forwardRef()使用
    useImperativeHandle()：获取和调用渲染后的 DOM 元素对象拥有的自定义方法，必须配合 React.forwardRef()使用

    生命周期
    useEffect()：挂载或渲染完成后、即将被卸载前，调度
    useLayoutEffect()：挂载或渲染完成后，同步调度

    性能优化
    useCallback()：获取某处理函数的引用，必须配合 React.memo()使用
    useMemo()：获取某处理函数返回值的副本

    代码调试
    useDebugValue()：对 react 开发调试工具中的自定义 hook，增加额外显示信息

    自定义 hook
    useCustomHook()：将 hook 相关逻辑代码从组件中抽离，提高 hook 代码可复用性

14. Hook 是 React 团队在大量实践后的产物, 更优雅的代替 class, 且性能更高. 故从开发使用者的角度来讲, 应该拥抱 Hook 所带来的便利.
