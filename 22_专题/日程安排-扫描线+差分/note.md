对于区间判断是否重叠，我们可以反向判断，也可以正向判断。 暴力的方法是每次对所有的课程`先排序进行判断是否重叠`，这种解法可以 AC。我们也可以进一步优化，

使用**Count-Map 方法来通用解决所有的问题**，不仅可以完美解决这三道题，还可以扩展到《会议室》系列的两道题。(SortedMap+上车下车)
`759. 员工空闲时间.py`

离线区间重叠问题解法

1. 一类是使用 sortedDict + 上车下车 解法
   会议室系列+`759. 员工空闲时间.py`
2. 另一类是 排序遍历 + 比较 `preEnd` 与 curStart/curEnd 解法
   合并区间系列+`82. 寻找合适开会的时间.ts`

3. 扫描线 delta 问题
   `1943. 描述绘画结果-扫描.py`
   `2015_每个线段的平均高度.py`
   deltaDict 记录这一系列的 add 操作即可

```Python
deltaDict = defaultdict(int)
for start, end, delta in segments:
    deltaDict[start] += delta
    deltaDict[end] -= delta

res = []

# 区间起点,累加和
pre, preSum = 0, 0
for cur in sorted(deltaDict):
    delta = deltaDict[cur]
    # 有效的区间
    if preSum > 0:
        res.append([pre, cur, preSum])
    pre = cur
    preSum += delta
return res
```
