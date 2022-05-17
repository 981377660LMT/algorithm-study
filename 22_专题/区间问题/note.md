所谓区间问题，就是线段问题，让你合并所有线段、找出线段的交集等等。主要有两个技巧：

1. `排序`区间，`再扫描区间比较 preEnd 与 curStart、curEnd，更新 preEnd`，讨论相交 覆盖 相离
   `注意排序一般是 sort() 但有时候需要将长的区间排在前面`

```Python
   intervals.sort()
   preEnd = intervals[0][1]
   res = intervals[0][1] - intervals[0][0]
   for i in range(1, len(intervals)):
       curStart, curEnd = intervals[i]
       # 包含
       if curEnd <= preEnd:
           continue
       # 相交
       elif curStart <= preEnd < curEnd:
           res += curEnd - preEnd
       # 相离
       elif preEnd < curStart:
           res += curEnd - curStart
       preEnd = max(preEnd, curEnd)

   return res
```

2. 差分+扫描线 会议室系列
   [1943. 描述绘画结果-差分+扫描线](..%5C%E6%97%A5%E7%A8%8B%E5%AE%89%E6%8E%92-%E6%89%AB%E6%8F%8F%E7%BA%BF+%E5%B7%AE%E5%88%86%5C1943.%20%E6%8F%8F%E8%BF%B0%E7%BB%98%E7%94%BB%E7%BB%93%E6%9E%9C-%E5%B7%AE%E5%88%86+%E6%89%AB%E6%8F%8F%E7%BA%BF.py)
