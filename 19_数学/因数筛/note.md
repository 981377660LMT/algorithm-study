这种题一般都是直接枚举，`且涉及到 gcd`
这类题的特点是 `nums[i]<=10^5` 筛法枚举因子是 nlogn

1. 遍历因子 range(1,MAX+1)
2. 遍历因子的倍数
3. 统计每个因子的倍数在原数组中出现了多少次

```Python
for factor in range(1, MAX + 1):
   for multi in range(factor, MAX + 1, factor):
       # 获取每个因子的信息
       multiCouner[factor] += counter[multi]
```

经典题
`1627. 带阈值的图连通性-枚举因子`
`1819. 序列中不同最大公约数的数目-遍历范围枚举`
`6015. 统计可以被 K 整除的下标对数目`

关注每个因子，而不是关注 pair

力扣考数论基本只考筛法
