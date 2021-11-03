React LOGO 的图案是代表原子（atom）的符号。世间万物由原子组成，原子的类型与属性决定了事物的外观与表现。

同样，在 React 中，我们可以将 UI 拆分为很多独立的单元，每个单元被称为 Component。这些 Component 的属性与类型决定了 UI 的外观与表现。

讽刺的是，原子在希腊语中的意思为不可分割的（indivisible），但随后科学家在原子中发现了更小的粒子 —— 电子（electron）。电子可以很好的解释原子是如何工作的。

在 React 中，我们可以说 ClassComponent 是一类原子。

但对于 Hooks 来说，与其说是一类原子，不如说他是更贴近事物运行规律的电子。

我们知道，React 的架构遵循 schedule - render - commit 的运行流程，这个流程是 React 世界最底层的运行规律。

ClassComponent 作为 React 世界的原子，他的生命周期（componentWillXXX/componentDidXXX）是为了介入 React 的运行流程而实现的更上层抽象，这么做是为了方便框架使用者更容易上手。

相比于 ClassComponent 的更上层抽象，Hooks 则更贴近 React 内部运行的各种概念（state | context | life-cycle）。
作为使用 React 技术栈的开发者，当我们初次学习 Hooks 时，不管是官方文档还是身边有经验的同事，总会拿 ClassComponent 的生命周期来类比 Hooks API 的执行时机。

这固然是很好的上手方式，但是当我们熟练运用 Hooks 时，就会发现，这两者的概念有很多割裂感，并不是同一抽象层次可以互相替代的概念。

比如：替代 componentWillReceiveProps 的 Hooks 是什么呢？
可能有些同学会回答，是 useEffect：

```JS
useEffect( () => {
console.log('something updated');
}, [props.something])
```

但是 componentWillReceiveProps 是在 render 阶段执行，而 useEffect 是在 commit 阶段完成渲染后异步执行。

这篇文章可以帮你更好理解 componentWillReceiveProps：深入源码剖析 componentWillXXX 为什么 UNSAFE(opens new window)

所以，从源码运行规律的角度看待 Hooks，可能是更好的角度。这也是为什么上文说 **Hooks 是 React 世界的电子而不是原子**的原因。

Concurrent Mode 是 React 未来的发展方向，而 Hooks 是能够最大限度发挥 Concurrent Mode 潜力的 Component 构建方式。

正如 Dan 在 React Conf 2018 演讲结尾所说：你可以从 React 的 LOGO 中看到这些围绕着核心的电子飞行轨道，Hooks 可能一直就在其中。
