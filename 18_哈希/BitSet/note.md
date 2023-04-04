1. 优化: 对二的整数幂取模(位运算时) 可以换成与运算

```ts
const id = index >> 5
const mask = index & 31
```
