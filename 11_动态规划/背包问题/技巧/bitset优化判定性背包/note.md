01 背包中改变一个元素 x 的正负，等价于将 dp 数组整体平移 x 个单位。

```go
dp.IOr(dp.Copy().Lsh(x))
```
