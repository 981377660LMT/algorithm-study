js 的基本类型 7 个

```TS
type Primitive = string | number | boolean | bigint | symbol | null | undefined
// sun sb nb

```

typeof 结果 8 个

```TS
function undefined number bigint boolean object symbol string
fun bboss
```

defer 与 async 区别
不加会阻塞
defer 会按顺序异步，defer 标签会按照在文档中的声明顺序执行
async 不会按顺序异步，谁先来就谁先执行

闭包：词法作用域才是闭包产生的原因
**与 python 闭包的区别**
python 的 inner 函数读取 outer 函数里的变量需要加 nonlocal 关键词
否则会默认为 inner 局部变量
而 JS 不一样 是静态作用域:代码写好**作用域**就已经确定了

V8 借鉴了 JVM 的垃圾回收机制
