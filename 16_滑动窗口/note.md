滑动窗口的滑动需要具**有单向性**:
如果区间移动`对某个性质的影响是具有单调性的` 那都可以考虑滑窗
例如移动右使得某个量一直增，移动左使得某个量一直减
思路都是用 map 记录数字对应的索引/出现次数

1. 固定窗口大小 只需要一个 for 循环和两个时间点
   `1852. 每个子数组的数字种类数.py`
   `2107_分k颗糖果后的不同风味-固定大小滑窗.py`

   ```Python
        res = 0
        for right, cur in enumerate(candies):
            counter[cur] -= 1
            if not counter[cur]:
                del counter[cur]
            # 更新左端点
            if right >= k:
                counter[candies[right - k]] += 1
            # 更新答案
            if right >= k - 1:
                res = max(res, len(counter))
        return res
   ```

2. 移动窗口大小/固定端点找边界(一般是固定右、找左) => 不符合条件就收缩
   注意到窗口移动的单调性 直到满足某个条件才停止继续寻找边界

   - 内部窗口,枚举右端点 right,找到满足条件的左端点 left
     [固定左端点寻找右边界](%E5%9B%BA%E5%AE%9A%E7%AB%AF%E7%82%B9%E6%89%BE%E5%8F%A6%E4%B8%80%E4%B8%AA%E7%AB%AF%E7%82%B9%E8%BE%B9%E7%95%8C/E%20-%20At%20Least%20One.py)

   ```Python
    for right in range(n):
        curSum += nums[right]  # add 逻辑
        while left <= right and ...:  # 不要忘记 left <= right
            curSum -= nums[left]  # remove 逻辑
            left += 1
            ...
        res = max(res, right - left + 1)  # !`合法`的时候更新答案
    return res
   ```

   - 外部窗口(前缀+后缀)，对每一个可能的起点 left，求出第一个符合题意的右边界 right
     这种情形也可以用二分写，更加简洁

- 双指针包括同向/反向，滑窗是双指针的一种(同向双指针)

---

最优秀的限流算法：滑动窗口
https://blog.herbert.top/2020/10/02/network_flow_control/
避免了漏桶的瞬时大流量问题，以及计数器实现的突刺现象
