倍增

1. 初始化 选定
   `MAXJ=floor(log2(k)) + 1`
2. 倍增
   `dp[j + 1][i] = dp[j]dp[j][i]]`
3. 二进制分解 k

```Python
res = 0
for bit in range(maxJ + 1):
    if (k >> bit) & 1:
        res = dp[bit][res]

```

或者

```JS
let bit = 0
while (k) {
     if (k & 1) res = this._dp[bit][res]
     bit++
     k >>>= 1
}
```
