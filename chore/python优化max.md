function lookup 和 function call 导致的性能问题
[max 函数用 if 代替](python%E7%9A%84max%E5%87%BD%E6%95%B0%E5%BE%88%E6%85%A2.py)
Python performance 推荐这篇文章：https://lwn.net/Articles/893686/

- 变量查找
  len()、int()和更多内置函数也是如此；每次使用它们都需要进行昂贵的查找。
- 动态属性查找
  即使在单个类中，解释器也`不知道该类的对象上存在哪些属性`(所以 slots 非常关键)。
  这意味着它需要对象的动态表示来跟踪属性。在 Python 中查找属性非常快，但仍然比静态属性列表要慢得多。
- 变量类型查找
  解释器不知道任何变量的类型。
  这意味着对变量的任何操作都需要查找它是什么类型，以便弄清楚如何对该类型执行操作。

- int 大小
  python 类似的多倍长整形的语言中
  运算图中如果超过 `2^63-1` 计算会慢很多
  `因此 INF 设置成 int(1e18) 比 int(1e20) 好`
  例如
  `a*b%MOD*c*d%MOD*e*f%MOD`只需要 0.9 秒
  但是`a*b*c*d*e*f%MOD` 需要 5 秒

```Python
# !3. 转移
for _ in range(1, k):
    # !这里把1<<63(9223372036854775808) 改小成 1<<63-1(9223372036854775807) 快了700ms pypy3 (碰到超过1<<63-1的数就会变慢)
    ndp = [int(1e18)] * (2 ** n)  # !比较合理的INF是int(1e18) (longlong 是9e18多一点)，比较慢的INF是int(1e20)
```

- 遍历保存中间结果和即使取余

```Python
# # !1.TLE
res = 0
for sub in combinations(nums, 5):
    mod_ = reduce(lambda x, y: x * y % p, sub, 1)
    res += mod_ % p == q
print(res)

# # !3. 2的基础上保存中间结果   489 ms ac
res = 0
for i in range(n):
    m1 = nums[i]
    for j in range(i + 1, n):
        m2 = m1 * nums[j] % p
        for k in range(j + 1, n):
            m3 = m2 * nums[k] % p
            for l in range(k + 1, n):
                m4 = m3 * nums[l] % p
                for m in range(l + 1, n):
                    m5 = m4 * nums[m] % p
                    res += m5 == q
print(res)


# 优化：
# 1. combinations会比直接遍历慢很多
# 2. 计算过程及时取余 不超过 1<<63-1 (这一点似乎只对pypy3有效)
# 3. 遍历时保存中间计算结果

```
