- 一维卷积的形式

`c[k] = ∑_i+j=k a[i]*b[j]`
入参是长为 n 的数组 A 和长为 m 的数组 B，返回值是长为 n+m-1 的数组 C

```py
朴素的 O(nm) 实现
def conv(a, b):
    n, m = len(a), len(b)
    c = [0] * (n + m - 1)
    for i in range(n):
        for j in range(m):
            c[i + j] += a[i] * b[j]
    return c
```

- 解决这样的问题 => 看似 O(n^2),通过 fft/ntt 卷积可以 O(nlogn)
  朴素的多项式乘法
  朴素的异或和
  **问题变形，能够把高复杂度部分变成卷积形式**
  https://atcoder.jp/contests/abc291/tasks/abc291_g
  a|b = a+b - (a&b), 当 a,b 都是 0/1 时，相当于 a|b = a+b - `a*b`
  每个位的错位异或和可以 reverse+卷积 求出

- 需要的卷积模板类型
  https://atcoder.jp/contests/practice2/tasks/practice2_f
  https://judge.yosupo.jp/

1. fft

   - 普通 fft (不带模数)
   - 任意模数 fft

2. ntt (带模)
   - 998244353 模数 ntt
   - 任意模数 ntt (跑得很慢)
     https://suisen-cp.github.io/cp-library-cpp/library/convolution/arbitrary_mod_convolution.hpp

- 一般的卷积理解
  https://suisen-cp.github.io/cp-library-cpp/library/convolution/convolution.hpp
