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

n 的**约数个数(因数个数)**是啥数量级:在 int 范围内近似为`立方根`
n 的**素因子个数(质因子个数)**是啥数量级:`log(n)`，因为每个至少为 2

---

`a*b<=c 等价于 a<=c//b (a,c为整数,b为正整数)`

---

- 等比数列之和：
  `Sn = a1*(q^n-1)/(q-1), q!=1`

- 等差乘以等比数列之和：
  `∑i*q^(i-1) = n*q^n - (q^n-1)/(q-1)`

- 若干无穷级数之和的公式
  `∑^∞ r^i = 1/(1-r)`
  `∑^∞ i*r^i = r/(1-r)^2`

- 等幂和
  `1^2 + ... + n^2 = n*(n+1)*(2*n+1)/6`
  `1^3 + ... + n^3 = [n(n+1)/2]^2`

- 长为 n 的数组的所有子数组长度之和
  `n*(n+1)*(n+2)/6`

- 长为 n 的数组的所有子数组的「长度/2 下取整」之和
  n 为偶数时：`m=n/2, m*(m+1)*(4*m-1)/6` https://oeis.org/A002412
  n 为奇数时：`m=n/2, m*(m+1)*(4*m+5)/6` https://oeis.org/A016061
  综合：`m*(m+1)*(m*4+n%2*6-1)/6`
- `floor(log2(i)) == (bits.Len(i)-1)`

- a 中任意两数互质 <=> 每个质数至多整除一个 a[i]
- 两种硬币面额为 a 和 b，互质，数量无限，所不能凑出的数值的最大值为 `a*b-a-b`
- 对于一任意非负序列，前 i 个数的 GCD 是非增序列，且至多有 O(logMax) 个不同值
