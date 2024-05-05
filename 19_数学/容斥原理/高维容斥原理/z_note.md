容斥原理 `用于求集合的并`

```Python
# 枚举容斥原理里的交叉项，选哪几个集合
for state in range(1, 1 << n):
    count = 0  # 用于统计奇数/偶数个集合 +1 还是 -1 的贡献
    for i in range(n):
        if (state >> i) & 1:
            count += 1
            ...
    tmp = 选这些集合时的方案数
    res += tmp * (1 if count & 1 else -1)
    res %= MOD
```

---

https://compro.tsutaj.com//archive/181015_incexc.pdf

---

https://github.com/EndlessCheng/codeforces-go/blob/86d1fd150c7b53861a52fa81ce456666bf547691/copypasta/math_comb.go#L376
