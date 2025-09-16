`如果要查询某个状态的所有子集的贡献，就要用 sosdp(高维前缀和)查询`

```python
zeta变换(高维前缀和)
for i in range(d):
    for j in range(1<<d):  # 不同维度的组合
      if (j&(1<<i)): f[j] += f[j^(1<<i)]  # f里保存了每个子集的信息

f代表着2^d个子集的前缀和,这里的子集一般用于描述值域(二进制每个位取还是不取)
f[0b110011] 记录着 0b010011/0b100011/0b110001/0b110010 这四个子集的贡献和

```

- https://blog.csdn.net/weixin_38686780/article/details/100109753
- https://issue-is-vegetable.blog.csdn.net/article/details/112488108?spm=1001.2101.3001.6661.1&utm_medium=distribute.pc_relevant_t0.none-task-blog-2%7Edefault%7ECTRLIST%7ERate-1.pc_relevant_antiscanv2&depth_1-utm_source=distribute.pc_relevant_t0.none-task-blog-2%7Edefault%7ECTRLIST%7ERate-1.pc_relevant_antiscanv2&utm_relevant_index=1
- https://blog.csdn.net/weixin_30249203/article/details/97527552?spm=1001.2101.3001.6661.1&utm_medium=distribute.pc_relevant_t0.none-task-blog-2%7Edefault%7ECTRLIST%7ERate-1.pc_relevant_paycolumn_v3&depth_1-utm_source=distribute.pc_relevant_t0.none-task-blog-2%7Edefault%7ECTRLIST%7ERate-1.pc_relevant_paycolumn_v3&utm_relevant_index=1
  ![](image/note/1651762471154.png)
  全称是 Sum over Subsets dynamic programmingSum over Subsets dynamic programming, 意思就是子集和 dp

```py
# memory optimized, super easy to code.
# dp[state]记录所有子集的和(高维前缀和) 而不是这一个状态的和

# n = ceil(log2(1e5))
# upper = 1 << n
# preSum = [0] * upper

for i in range(upper):
    sosdp[i] = ...  # 初始化每个元素的贡献
for i in range(n):
    for state in range(upper):
        if (state >> i) & 1:
            sosdp[state] += sosdp[state ^ (1 << i)]  # 贡献求子集前缀和
            # sosdp[state ^ (1 << i)] += sosdp[state] # 贡献求超集前缀和
# print(preSum[5])

##############################################################
@lru_cache(None)
def sosdp(state: int) -> int:
    """state真子集的贡献"""
    if state == 0:
        return 1
    res = 0
    for i in range(20):
        if (state >> i) & 1:
            res += sosdp(state ^ (1 << i))
    return res
```

将时间复杂度从枚举子集的子集的`O(3^n)`优化到了 `O(n*2^n)`
可以`O(n*2^n)`算出 n 位每个 mask 值所包含子集的二进制码下标的贡献
比如 f[5]=a[0]+a[1]+a[4]+a[5]这种的，101 包含了 000,001,100,101
以 5(101)为例，状压记录的是(100, 1)这两个状态，而 sosdp 记录的是(101, 100, 1, 0)这四个状态，也可以说是他的子集。即 sosdp[mask]存的是 所有的 a[i]，其中 i&mask == i
核心思想就是从低位枚举到高位，

f[mask][i]表示 mask 码低 i 位子集的贡献
如果 mask 的第 i 位是 1，那么 f[mask][i]=f[mask][i-1]+f[mask^(1<<i)][i-1]，很巧妙地以第 i 位的 01 作为子集划分
如果 mask 的第 i 位是 0，那么 f[mask][i]=f[mask][i-1]

- 二维前缀和与高维前缀和

```C++
这是用容斥来做的
for(int i=1;i<=n;i++)
{
	for(int j=1;j<=m;j++) f[i][j]=f[i-1][j]+f[i][j-1]-f[i-1][j-1]+a[i][j];
}

但我们还可以依次对行和列进行前缀和，即：
for(int i=1;i<=n;i++) for(int j=1;j<=m;j++) f[i][j]=f[i-1][j]+a[i][j];
for(int i=1;i<=n;i++) for(int j=1;j<=m;j++) f[i][j]=f[i][j-1]+a[i][j];
看起来好像带常数，但这是可以仅通过加入若干层枚举轻松拓展到更高维的，而容斥每加一维都并不好推。
因此，不妨以这个思路，将集合的元素个数用维来代替，那就每次处理某一维的前缀和即可。
```

---

https://maspypy.github.io/library/setfunc/zeta.hpp

Zeta 变换(SOSDp)，也被称为子集和变换或者幂集变换，是一种在组合数学和计算机科学中常见的运算。它主要用于处理和优化与集合相关的问题。
在 Zeta 变换中，我们有一个函数 f，它定义在大小为 n 的集合的所有子集上。Zeta 变换就是将这个函数转换为另一个函数 g，其中 g 的每个值都是 f 在某个子集及其所有子集上的值的总和。
具体来说，对于一个集合 S，我们定义 Zeta 变换后的函数 g(S)为 f 在 S 的所有子集上的值的总和，即：
g(S) = Σ f(T)，其中 T 是 S 的一个子集。
Zeta 变换在很多场合都有应用，比如在计算组合数、优化动态规划算法等方面。

**这类题的特点是 nums[i]<=1e6，与按位与有关**

---

| 函数名           | 作用                            | 别名/概念              |
| ---------------- | ------------------------------- | ---------------------- |
| `SubsetZeta`     | 计算所有**子集**的和            | Sum over Subsets (SOS) |
| `SubsetMobius`   | `SubsetZeta` 的逆变换（容斥）   | 子集莫比乌斯反演       |
| `SuperSetZeta`   | 计算所有**超集**的和            | Sum over Supersets     |
| `SupersetMobius` | `SuperSetZeta` 的逆变换（容斥） | 超集莫比乌斯反演       |
