from typing import List
from bisect import bisect_left


# 文明等级（C），资源储备（R）以及人口数量（H）。在游戏开始时（第 0 天），三种属性的值均为 0。
# [[1,2,1],[3,4,2]] 表示第一天三种属性分别增加 1,2,1 而第二天分别增加 3,4,2。
# 对于某个剧情的触发条件 c[i], r[i], h[i]，如果当前 C >= c[i] 且 R >= r[i] 且 H >= h[i] ，则剧情会被触发。
# 根据所给信息，请计算每个剧情的触发时间，并以一个数组返回。

# 1 <= increase.length <= 10000
# 1。前缀和数组
# 2. 总结:对每个条件requirement，寻找符合条件的最左


class Solution:
    def getTriggerTime(self, increase: List[List[int]], requirements: List[List[int]]) -> List[int]:
        a, b, c = [0], [0], [0]
        for da, db, dc in increase:
            a.append(a[-1] + da)
            b.append(b[-1] + db)
            c.append(c[-1] + dc)

        res = []
        for x, y, z in requirements:
            la = bisect_left(a, x)
            lb = bisect_left(b, y)
            lc = bisect_left(c, z)
            upper = max(la, lb, lc)
            if upper <= len(increase):
                res.append(upper)
            else:
                res.append(-1)

        return res


print(
    Solution().getTriggerTime(
        increase=[[2, 8, 4], [2, 5, 0], [10, 9, 8]],
        requirements=[[2, 11, 3], [15, 10, 7], [9, 17, 12], [8, 1, 14]],
    )
)
# 输出: [2,-1,3,-1]
# 解释：
# 初始时，C = 0，R = 0，H = 0
# 第 1 天，C = 2，R = 8，H = 4
# 第 2 天，C = 4，R = 13，H = 4，此时触发剧情 0
# 第 3 天，C = 14，R = 22，H = 12，此时触发剧情 2
# 剧情 1 和 3 无法触发。

