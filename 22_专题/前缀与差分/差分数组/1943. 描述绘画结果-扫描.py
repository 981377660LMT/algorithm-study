from typing import List
from collections import defaultdict

# 刷油漆问题
# 线段间重叠部分的颜色会被 混合 。如果有两种或者更多颜色混合时，它们会形成一种新的颜色，用一个 集合 表示这个混合颜色。
# 比方说，如果颜色 2 ，4 和 6 被混合，那么结果颜色为 {2,4,6} 。
# 1 <= segments.length <= 2 * 104
# 请你返回二维数组 painting ，它表示最终绘画的结果

# 会议室问题

# 扫描线：记录pre,扫描cur
class Solution:
    def splitPainting(self, segments: List[List[int]]) -> List[List[int]]:
        dic = defaultdict(int)
        for left, right, delta in segments:
            dic[left] += delta
            dic[right] -= delta

        res = []
        pre = -1
        for cur in sorted(dic):
            if dic[pre] > 0:
                res.append([pre, cur, dic[pre]])
            dic[cur] += dic[pre]
            pre = cur
        return res


print(Solution().splitPainting(segments=[[1, 4, 5], [4, 7, 7], [1, 7, 9]]))
# 输出：[[1,4,14],[4,7,16]]
# 解释：绘画借故偶可以表示为：
# - [1,4) 颜色为 {5,9} （和为 14），分别来自第一和第二个线段。
# - [4,7) 颜色为 {7,9} （和为 16），分别来自第二和第三个线段。

