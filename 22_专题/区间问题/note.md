所谓区间问题，就是线段问题，让你合并所有线段、找出线段的交集等等。主要有两个技巧：

1. `排序`区间，`再扫描区间比较 preEnd 与 curStart，更新 preEnd`。常见的排序方法就是按照区间起点排序，或者先按照起点升序排序，若起点相同，则按照终点降序排序。当然，如果你非要按照终点排序，无非对称操作，本质都是一样的。
2. 画图。就是说不要偷懒，勤动手，两个区间的相对位置到底有几种可能(3 种)，不同的相对位置我们的代码应该怎么去处理。
   相交 覆盖 相离

3. `sortedDict+上车下车扫描` 会议室系列

```JS
预处理
intervals.sort((a, b) => a[0] - b[0] || b[1] - a[1])  => 起点相同时，区间短的排在前面
记录preEnd
区间相交只需比较 preEnd 与此区间的start 大小
更新preEnd=Math.max(preEnd,end)
`82. 寻找合适开会的时间`
```

模板

```Python

        res = []
        preEnd = intervals[0][1]
        for curStart, curEnd in intervals:
            if preEnd < curStart:
                res.append(Interval(preEnd, curStart))
            preEnd = max(preEnd, curEnd)

        return res
```
