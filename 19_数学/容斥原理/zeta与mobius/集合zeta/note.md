1. 卷积

- 子集卷积
  两个基础:zeta/mobius `O(n*2^n)`
  四个应用:
  xor/and/or 卷积 `O(n*2^n)`
  子集卷积 `O(n^2*2^n)`
  https://qiita.com/convexineq/items/afc84dfb9ee4ec4a67d5
  https://nyaannyaan.github.io/library/set-function/subset-convolution.hpp

  > 我的理解:
  > zeta 变换 -> 高维前缀和
  > mobius 变换 -> 高维差分(高维前缀和的逆运算)
  > 快速 zeta 变换 -> 高速 d 维前缀和,1 维前缀和在 d 个维度上分别进行
  > 快速 mobius 变换 -> 高速 d 维差分,1 维差分在 d 个维度上分别进行

  2 次元累積和は 1 次元累積和をそれぞれの次元について行えばよい。

  ```python
      #ヨコに累積和
  for i in range(n):
      for j in range(n):
          if i: a[i][j] += a[i-1][j]
  #タテに累積和
  for i in range(n):
      for j in range(n):
          if j: a[i][j] += a[i][j-1]


  for i in range(n-1,-1,-1):
      for j in range(n-1,-1,-1):
          if i: a[i][j] -= a[i-1][j]
  for i in range(n-1,-1,-1):
      for j in range(n-1,-1,-1):
          if j: a[i][j] -= a[i][j-1]
  ```

  推广到 d 维`O(d*2^d)`:
  https://oi-wiki.org/basic/prefix-sum/#%E5%9F%BA%E4%BA%8E-dp-%E8%AE%A1%E7%AE%97%E9%AB%98%E7%BB%B4%E5%89%8D%E7%BC%80%E5%92%8C

  ```python
  for i in range(d):
     for j in range(1<<d):  # 不同维度的组合
       if (j&(1<<i)): f[j] += f[j^(1<<i)] 将f进行zeta变换

  for i in range(d):
     or j in range(1<<d):
       if (j&(1<<i)): g[j] -= g[j^(1<<i)] 将g进行mobius变换
  ```

  如果要计算 `h[x] = ∑f[i]*g[j] (0<=i,j<n,max(i,j)=x)`呢?
  将 f 和 g 分别进行快速 zeta 变换,然后对应元素相乘,再进行快速 mobius 变换即可.

  推广到 gcd:

  ```python
  #注意：あらかじめ素数のリスト primes を作成する
  def zeta_divisor(a,primes): #aを約数ゼータ
      n = len(a)-1
      for p in primes:
          for i in range(n//p,0,-1):
              a[i] += a[p*i]
  def mobius_divisor(a,primes): #aを約数メビウス
      n = len(a)
      for p in primes:
          for i in range(1,n):
              if i*p >= n: break
              a[i] -= a[p*i]
  ```

- 正常的卷积
