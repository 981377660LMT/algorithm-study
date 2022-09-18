from typing import List
from sortedcontainers import SortedList

# 第 j 个查询的答案是满足如下条件的最近的房间 id ：

# 房间的面积 至少 为 minSizej ，且
# abs(id - preferredj) 的值 最小 ，其中 abs(x) 是 x 的绝对值。

# 将区间和查询分别排序，然后离线处理
# `将查询从大到小排序`，然后把所有符合条件的房间放到集合里面，在集合里面二分房间号，找到最接近pre的。
class Solution:
    def closestRoom(self, rooms: List[List[int]], queries: List[List[int]]) -> List[int]:
        # 存放房间id
        availRooms = SortedList()
        queries = sorted(
            [[size, prefer, i] for i, (prefer, size) in enumerate(queries)], reverse=True
        )
        rooms = sorted([[size, id] for id, size in rooms], reverse=True)

        m, n = len(rooms), len(queries)
        ri, qi = 0, 0
        res = [-1] * n

        for qi in range(n):
            while ri < m and rooms[ri][0] >= queries[qi][0]:
                availRooms.add(rooms[ri][1])
                ri += 1

            if availRooms:
                _, prefer, index = queries[qi]
                pos = availRooms.bisect_right(prefer) - 1

                # 直接调最右二分，然后看i和i-1，减少了思考难度
                cands = []
                if pos >= 0:
                    cands.append(availRooms[pos])
                if pos + 1 < len(availRooms):
                    cands.append(availRooms[pos + 1])
                res[index] = min(cands, key=lambda x: abs(x - prefer))

        return res


print(Solution().closestRoom(rooms=[[2, 2], [1, 2], [3, 2]], queries=[[3, 1], [3, 3], [5, 2]]))
# 输出：[3,-1,3]
# 解释：查询的答案如下：
# 查询 [3,1] ：房间 3 的面积为 2 ，大于等于 1 ，且号码是最接近 3 的，为 abs(3 - 3) = 0 ，所以答案为 3 。
# 查询 [3,3] ：没有房间的面积至少为 3 ，所以答案为 -1 。
# 查询 [5,2] ：房间 3 的面积为 2 ，大于等于 2 ，且号码是最接近 5 的，为 abs(3 - 5) = 2 ，所以答案为 3 。
