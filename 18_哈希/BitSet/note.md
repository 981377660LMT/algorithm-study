1. 优化: 对二的整数幂取模(位运算时) 可以换成与运算

```ts
const id = index >> 5
const mask = index & 31
```

---

当 Set 存储的值为 0-1e9 时,可以将 Set 换成 BitSet,这样可以节省空间,加快查询速度
