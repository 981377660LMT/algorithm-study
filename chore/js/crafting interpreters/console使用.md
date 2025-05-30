# JavaScript Console 对象详解

`console`对象是 JavaScript 中用于调试的重要工具，提供了多种输出和调试方法。下面详细讲解各种方法及其使用案例。

## 基础输出方法

### 1. console.log()

最常用的输出方法，可接受多个参数。

```javascript
console.log('Hello World') // 输出字符串
console.log('数值:', 42) // 输出多个值
console.log('对象:', { name: '张三', age: 25 }) // 输出对象
console.log('%c彩色文本', 'color:red; font-size:20px') // 样式文本
```

### 2. console.info()

与 log 类似，但在某些浏览器中会显示信息图标。

```javascript
console.info('这是一条信息')
```

### 3. console.warn()

输出警告信息，通常显示黄色背景。

```javascript
console.warn('警告：此功能即将废弃')
```

### 4. console.error()

输出错误信息，通常显示红色背景。

```javascript
console.error('发生错误！', new Error('出错了'))
```

### 5. console.debug()

输出调试信息，需要开启调试级别才能看到。

```javascript
console.debug('调试信息')
```

## 格式化与表格

### 6. console.table()

以表格形式显示数据，适合展示数组和对象。

```javascript
const users = [
  { name: '张三', age: 25, role: '开发' },
  { name: '李四', age: 30, role: '设计' },
  { name: '王五', age: 28, role: '测试' }
]
console.table(users)
console.table(users, ['name', 'role']) // 只显示特定列
```

### 7. console.dir()

显示对象的所有属性，适合查看 DOM 对象。

```javascript
const obj = { a: 1, b: { c: 2, d: { e: 3 } } }
console.dir(obj, { depth: null }) // 展开所有层级
```

### 8. console.dirxml()

以 XML/HTML 形式显示元素。

```javascript
console.dirxml(document.body)
```

## 分组与缩进

### 9. console.group() / console.groupCollapsed() / console.groupEnd()

创建分组输出，便于组织信息。
group/groupCollapsed() 与 groupEnd() 需要成对使用。

```javascript
console.group('用户信息')
console.log('姓名: 张三')
console.log('年龄: 25')
console.group('详细信息')
console.log('职位: 开发工程师')
console.log('部门: 技术部')
console.groupEnd()
console.groupEnd()

// 折叠分组
console.groupCollapsed('系统信息') // 初始折叠显示
console.log('操作系统: Windows')
console.log('浏览器: Chrome')
console.groupEnd()
```

## 计时功能

### 10. console.time() / console.timeEnd() / console.timeLog()

测量代码执行时间。

```javascript
console.time('排序操作')
const arr = Array(10000)
  .fill()
  .map(() => Math.random())
arr.sort()
console.timeLog('排序操作') // 打印当前计时
// 执行更多操作...
console.timeEnd('排序操作') // 结束计时并打印总时间
```

## 计数功能

### 11. console.count() / console.countReset()

计算调用次数。
eg：检查死循环.

```javascript
function doSomething() {
  console.count('doSomething被调用')
}

doSomething() // doSomething被调用: 1
doSomething() // doSomething被调用: 2
doSomething() // doSomething被调用: 3

console.countReset('doSomething被调用') // 重置计数
doSomething() // doSomething被调用: 1
```

## 断言与跟踪

### 12. console.assert()

条件为 false 时才输出信息。

```javascript
const x = 5
console.assert(x > 10, 'x不大于10', x) // 条件为false，输出断言失败信息
console.assert(x > 3, 'x不大于3') // 条件为true，不输出
```

### 13. console.trace()

打印调用堆栈跟踪信息。

```javascript
function firstFunction() {
  secondFunction()
}

function secondFunction() {
  thirdFunction()
}

function thirdFunction() {
  console.trace('跟踪调用堆栈')
}

firstFunction() // 显示完整的调用路径
```

## 其他功能

### 14. console.clear()

清空控制台。

```javascript
console.clear()
```

### 15. console.memory

查看内存使用情况（属性，不是方法）。

```javascript
console.log(console.memory)

// MemoryInfo {totalJSHeapSize: 11900000, usedJSHeapSize: 10000000, jsHeapSizeLimit: 3760000000}
```

### 16. 性能分析（部分浏览器支持）

```javascript
console.profile('性能分析')
// 执行需要分析的代码
console.profileEnd('性能分析')
```

## 实用技巧

### 对象解构输出

```javascript
const user = { name: '张三', age: 25, email: 'zhangsan@example.com' }
console.log({ user }) // 自动展示变量名和值
```

### 条件断点调试

```javascript
let sum = 0
for (let i = 1; i <= 1000; i++) {
  sum += i
  // 当i为500时在控制台输出，类似条件断点
  i === 500 && console.log(`中间值: i=${i}, sum=${sum}`)
}
```

大多数方法在 Node.js 和浏览器环境中都可用，但某些方法可能在特定环境中有差异或不可用。这些方法主要用于开发调试，生产环境中应适当清理 console 语句。

---

## JavaScript Console 实用技巧

1. **对象解构输出**

   - 直接输出对象变量名和值，便于调试。

   ```javascript
   const user = { name: '张三', age: 25 }
   console.log({ user }) // 输出: { user: { name: '张三', age: 25 } }
   ```

2. **条件断点调试**

   - 利用短路运算在特定条件下输出，模拟条件断点。

   ```javascript
   for (let i = 0; i < 100; i++) {
     i % 10 === 0 && console.log('i是10的倍数:', i)
   }
   ```

3. **彩色和样式化输出**

   - 使用 `%c` 和 CSS 样式美化输出内容。

   ```javascript
   console.log('%c重要信息', 'color: red; font-weight: bold; font-size: 16px')
   ```

4. **表格展示数据**

   - 用 `console.table()` 直观展示数组或对象。

   ```javascript
   const arr = [
     { a: 1, b: 2 },
     { a: 3, b: 4 }
   ]
   console.table(arr)
   ```

5. **分组输出**

   - 用 `console.group()` 和 `console.groupEnd()` 组织输出，便于阅读嵌套信息。

   ```javascript
   console.group('用户信息')
   console.log('姓名: 张三')
   console.groupEnd()
   ```

6. **计时与性能分析**

   - 用 `console.time()` 和 `console.timeEnd()` 测量代码块耗时。

   ```javascript
   console.time('loop')
   for (let i = 0; i < 100000; i++) {}
   console.timeEnd('loop')
   ```

7. **计数调用次数**

   - 用 `console.count()` 统计某段代码被执行的次数。

   ```javascript
   function foo() {
     console.count('foo被调用')
   }
   foo()
   foo()
   ```

8. **断言输出**

   - 用 `console.assert()` 只在条件为假时输出错误信息。

   ```javascript
   console.assert(1 > 2, '断言失败：1不大于2')
   ```

9. **堆栈跟踪**

   - 用 `console.trace()` 输出当前调用堆栈，定位问题来源。

   ```javascript
   function a() {
     b()
   }
   function b() {
     console.trace('跟踪')
   }
   a()
   ```

10. **清空控制台**
    - 用 `console.clear()` 快速清屏，便于查看新输出。

> 这些技巧可以极大提升调试效率，建议结合实际开发场景灵活使用。
