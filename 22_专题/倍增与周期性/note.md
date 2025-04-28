## 哈希表找周期

1. 保存状态
2. 线性转移
3. 寻找周期

## 倍增

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

---

例题：
https://leetcode.cn/problems/prison-cells-after-n-days/description/
https://leetcode.cn/problems/sentence-screen-fitting/description/
https://atcoder.jp/contests/abc167/tasks/abc167_d
https://atcoder.jp/contests/abc241/tasks/abc241_e
https://atcoder.jp/contests/typical90/tasks/typical90_bf
https://atcoder.jp/contests/abc258/tasks/abc258_e
