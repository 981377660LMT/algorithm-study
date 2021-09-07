```JS
// 最后使用总长度,加个 0 防止 Infinity
 return wall.length - Math.max(...Object.values(rowMap), 0)

console.log(Math.max())  // -Infinity
console.log(Math.min())  // Infinity
```
