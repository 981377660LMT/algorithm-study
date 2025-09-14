TODO

## 如何排查一个函数为什么被调用到了？

在排查问题的时候有一个常见的 case，怀疑是走到了一个有问题的函数导致的，但是不确定为啥这个函数被调用了，很多时候也不太方便去修改源码，添加 debugger，重新部署去调试，其实 chrome 有一个更简单的方案：利用在 chrome devtool 环境下特有的一个 ”debug“ 全局函数:`typescriptfunction foo() {}debug(foo);// 此后任意对 foo 函数的调用都会断点断住// native code 同样可以断点断住debug(window.alert);window.alert('xx');`> 注意在源码里面书写 debug 是找不到这个函数的，只要在 chrome dev tool 的环境里面才会注入这个函数

---

这是一个在 Chrome 浏览器开发者工具（DevTools）中进行 JavaScript 调试的技巧。
它介绍了 `debug()` 这个在开发者工具控制台（Console）中内置的特殊函数。

**作用：**
当你想要知道某个函数是在何时、被哪个调用链所调用时，可以在控制台里执行 `debug(函数名)`。这样，当该函数下一次被执行时，浏览器会自动在该函数的第一行暂停（断点），让你能够检查当前的调用堆栈（Call Stack）和作用域中的变量。

**优点：**
这种方法非常方便，因为它不需要你修改源代码、添加 `debugger` 语句，然后重新部署或刷新页面。你可以直接在正在运行的页面上对任何函数设置断点，包括浏览器内置的函数（如 `window.alert`）。

与之对应的，还有一个 `undebug(函数名)` 函数，用于取消对特定函数的断点。

---

除了 `debug()` 之外，Chrome 开发者工具还提供了一些其他非常有用的调试技巧：

### 1. 监控函数调用: `monitor()`

如果你不想中断执行，只是想知道函数是否被调用以及传入了什么参数，可以使用 `monitor()`。

```javascript
function greet(name) {
  console.log('Hello, ' + name)
}

// 在控制台执行
monitor(greet)

// 之后每次调用 greet，控制台都会打印出函数名和传入的参数
greet('world')
// 输出: function greet called with arguments: "world"
// 输出: Hello, world

// 使用 unmonitor() 来停止监控
unmonitor(greet)
```

### 2. 监控事件: `monitorEvents()`

当你想知道一个 DOM 元素上触发了哪些事件时，这个函数非常有用。

```javascript
// 在 Elements 面板选中一个按钮，然后在控制台执行：
// $0 代表当前选中的 DOM 元素
monitorEvents($0, 'click')

// 当你点击该按钮时，控制台会打印出 click 事件对象。
// 你也可以监控多种事件
monitorEvents(window, ['resize', 'scroll'])

// 使用 unmonitorEvents() 停止监控
unmonitorEvents($0)
```

### 3. 快速引用 DOM 元素: `$0` - `$4`

在控制台中，`$0`、`$1`、`$2`、`$3` 和 `$4` 是对你在 **Elements** 面板中最近选择的五个 DOM 元素的引用。`$0` 是最近一次选择的元素。

这让你不必为了获取一个元素的引用而费力地编写 `document.querySelector`。

### 4. 复制对象到剪贴板: `copy()`

当你在控制台中得到一个复杂的 JSON 对象或数组时，可以用 `copy()` 函数将其内容快速复制到剪贴板。

```javascript
const myData = { user: 'test', id: 123, settings: { theme: 'dark' } }

// 在控制台执行
copy(myData)

// 现在你可以把这个对象的 JSON 字符串粘贴到任何地方。
```

### 5. 条件断点 (Conditional Breakpoints)

这是一个基于 UI 的功能。在 **Sources** 面板中，你可以右键点击行号，选择 "Add conditional breakpoint..."。然后输入一个表达式，只有当该表达式结果为 `true` 时，断点才会触发。这对于调试循环或者只在特定条件下才出现问题的场景非常有用。
