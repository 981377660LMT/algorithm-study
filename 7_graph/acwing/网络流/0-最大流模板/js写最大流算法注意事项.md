**为了保证速度需要注意的事项**

1. `不要用解构/es6特性 解构/es6特性 会慢很多 尽量用最简单的索引读取`

例如:[Dinic 最大流](atc-useDinic.ts)中的

```JS
const next = edges[i][0]
const capacity = edges[i][1]
改为
const [next, capacity] = edges[i]

从 4800ms 变为了 6900ms
```

2. 循环尽量用最普通的 `for 循环` ，不要用 迭代器和 forEach 循环
3. `INF 一般取 2e15` ，不要用 Infinity 相加会导致 NaN
4. 尽量用`静态数组`
   布尔值:Uint8Array
   数字:Uint32Array、Int32Array(需要-1 初始化时)
5. 如果可以用多个条件来约束，就不要 for 循环

6. 判断语句中把`运算量小的放在前面`，把运算量大的放在后面
   例如
   ```JS
   if (capacity > 0 && levels[next] === -1) {
     // ...
   }
   ```
