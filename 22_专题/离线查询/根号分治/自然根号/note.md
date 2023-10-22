# 自然根号

https://ddosvoid.github.io/2020/10/18/%E6%B5%85%E8%B0%88%E6%A0%B9%E5%8F%B7%E7%AE%97%E6%B3%95/

1. **有若干非负数的和为 n，则不同的数最多有 O(sqrt(n))个**

2. 二元组自然根号：给定一个和为 n 的数组 `A[i]` 和 O(n) 个二元组`(i,j)`，每个二元组的代价为 `O(min(A[i],A[j]))`，则总代价为 `O(n*sqrt(n))`

   证明：
   令 min(A[i],A[j]) = k
   如果 k>sqrt(n)，则最多有 n/k 个这样的二元组，这一部分二元组总代价为 O(n\*sqrt(n))；
   如果 k<=sqrt(n), 则总代价为 O(n\*sqrt(n))。

- https://ddosvoid.github.io/2021/02/04/Luogu-P5901-IOI2009-regions/
