1. python int 的上限是多少: 无上限,int 等于别的语言的 bigint
   python float 的上限是多少: 1.79e308 (这也是 double 的最大值)
2. python int 大数相乘的复杂度是多少(karatsuba 算法计算大数乘法)

   https://cloud.tencent.com/developer/article/1557959

   - 一般的长乘法复杂度为 `O(n^2)`
   - karatsuba 算法复杂度为 `O(n^lg3) 约为 O(n^1.585)`
   - **优化:fft/ntt 求解大数乘法 O(nlogn)**

3. golang 标准库 big 里的算法 (karatsuba 算法计算大数乘法)
4. javascript 里的 bigint

## 一些大数运算的速度印象

- 乘法取模超级快 (50ms/0ms)

  ```python
  a = int("2" * int(2e5))
  b = int("3" * int(2e5))
  print((a * b) % 1000000007)  # bad
  print((a % 1000000007 * b % 1000000007))  # good
  ```

- 乘法慢 (3s)

  ```python
  a = int("2" * int(2e5))
  b = int("3" * int(2e5))
  print(a * b)
  ```
