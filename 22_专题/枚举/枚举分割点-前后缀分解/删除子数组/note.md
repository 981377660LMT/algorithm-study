1. 预处理 pre suf 数组，维护`每个数作为开头/结尾`的某种性质
2. 枚举分割点、双指针(滑窗)求出答案

   - 一般是**枚举前缀长度**，对 pre 数组里的每一个下标，
     在 suf 数组里找到第一个满足条件的下标，求出 first 数组
     可以二分或者双指针求，二分一般更简单
   - 有的题注意特判没有前缀或后缀(长度为 0)的情况

   ```py
    # 双指针(外部滑窗)
    # 这里left的含义是前缀的长度
   for left, v in enumerate(pre):
       while right < len(suf) and suf[right] + v > n:
           right += 1
       if right < len(suf) and v + suf[right] <= n:
           first[left] = right
   ```

   ```py
    # 二分
    for i, v in enumerate(pre):
       left, right = 0, m
       ok = False
       while left <= right:
           mid = (left + right) // 2
           if suf[mid] + v <= n:
               right = mid - 1
               ok = True
           else:
               left = mid + 1
       if ok:
           first[i] = left
   ```
