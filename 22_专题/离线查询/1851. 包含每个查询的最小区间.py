from typing import List
from heapq import heappop, heappush

# 区间的 长度 定义为区间中包含的整数数目，更正式地表达是 righti - lefti + 1 。
# 第 j 个查询的答案是满足 lefti <= queries[j] <= righti 的 长度最小区间 i 的长度 。如果不存在这样的区间，那么答案是 -1 。
# 以数组形式返回对应查询的所有答案。


# 1353. 最多可以参加的会议数目.py
# 1 <= intervals.length <= 105

# 总结:
# 1.开会题模板
# 2.离线查询先排序
class Solution:
    def minInterval(self, intervals: List[List[int]], queries: List[int]) -> List[int]:
        # 离线查询预处理
        intervals.sort()
        queriesWithIndedx = sorted([(query, index) for index, query in enumerate(queries)])

        ei = 0
        res = [-1] * len(queries)
        pq = []

        # 遍历intervals左区间的位置
        for qv, qi in queriesWithIndedx:
            # 将所有起始位置小于等于查询位置的区间intervals[i]添加到优先队列中
            while ei < len(intervals) and intervals[ei][0] <= qv:
                start, end = intervals[ei]
                heappush(pq, (end - start + 1, end))
                ei += 1

            # 将队列中不能覆盖要查询点的区间移除队列
            while pq and pq[0][1] < qv:
                heappop(pq)

            # 如果队列不为空，则代表队首区间是要查询的点的最短区间
            if pq:
                length, _ = pq[0]
                res[qi] = length

        return res


print(Solution().minInterval(intervals=[[1, 4], [2, 4], [3, 6], [4, 4]], queries=[2, 3, 4, 5]))
# 输出：[3,3,1,4]
# 解释：查询处理如下：
# - Query = 2 ：区间 [2,4] 是包含 2 的最小区间，答案为 4 - 2 + 1 = 3 。
# - Query = 3 ：区间 [2,4] 是包含 3 的最小区间，答案为 4 - 2 + 1 = 3 。
# - Query = 4 ：区间 [4,4] 是包含 4 的最小区间，答案为 4 - 4 + 1 = 1 。
# - Query = 5 ：区间 [3,6] 是包含 5 的最小区间，答案为 6 - 3 + 1 = 4 。
