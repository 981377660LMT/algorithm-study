1. **母题：有序数组中差绝对值之和(所有点的曼哈顿距离之和)**

   1. **前缀和/后缀和**
      nums 是 `非递减有序`整数数组
      公式 (ni-n0) + (ni-n1) + (ni-n2) +...+(ni-ni) +(ni+1-ni) +(ni+2-ni)
      第 i 项等于
      `ni*i-preSum[i]`+`preSum[n]-preSum[i]-ni*(n-i)`
      **前面有 i+1 项 后面有 n-i 项**
   2. 三分法求严格凸函数极值
   3. 维护左边和右边的 cost，左右挪动更新总 cost:**不要这样做,容易写错索引(没有枚举到最后一个元素)，建议用前缀+后缀取代**

2. 枚举绝对值展开后的符号-O(n)
   `1131. 绝对值表达式的最大值`
   `1330. 翻转子数组得到最大的数组值`

---

`max(a,b,c) - min(a,b,c)` = `(abs(a-b) + abs(b-c) + abs(c-a)) //2`
