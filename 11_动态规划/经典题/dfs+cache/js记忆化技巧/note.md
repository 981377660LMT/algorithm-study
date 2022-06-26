1. 将状态的 key 哈希成一个数

```JS
const hash = state1 * 7 + state2
if (cache[index][hash] !== -1) return cache[index][hash]
```

2. 使用静态数组

```JS
const cache = Array.from({ length: n }, () => new Int32Array(100).fill(-1))
```

3. 使用字符串 join 模拟很长的元组 (`当然也可以写一个计算 key 的哈希函数`)
4. 数组存的性能比 Map 存字符串好很多
5. 注意 js 用 map 存不需要 clear
