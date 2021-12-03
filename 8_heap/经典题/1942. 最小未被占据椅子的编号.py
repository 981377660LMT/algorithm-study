from typing import List
from heapq import heappush, heappop

# 有 n 个朋友在举办一个派对，这些朋友从 0 到 n - 1 编号。派对里有 无数 张椅子，编号为 0 到 infinity
# 当一个朋友到达派对时，他会占据 编号最小 且未被占据的椅子。


# 总结:排序,模拟sweeping
# 保持pq队头是可供使用的编号最小的椅子
class Solution:
    def smallestChair(self, times: List[List[int]], targetFriend: int) -> int:
        action = []
        for i, (arrival, leaving) in enumerate(times):
            action.append((arrival, True, i))
            action.append((leaving, False, i))

        id = 0
        # 可供使用的椅子
        pq = []
        seatByPerson = dict()
        for _, isArrival, person in sorted(action):
            if isArrival:
                # 取一个椅子
                if pq:
                    seat = heappop(pq)
                else:
                    seat = id
                    id += 1
                if person == targetFriend:
                    return seat
                seatByPerson[person] = seat
            else:
                # 放回椅子
                heappush(pq, seatByPerson[person])


print(Solution().smallestChair(times=[[1, 4], [2, 3], [4, 6]], targetFriend=1))
# 输出：1
# 解释：
# - 朋友 0 时刻 1 到达，占据椅子 0 。
# - 朋友 1 时刻 2 到达，占据椅子 1 。
# - 朋友 1 时刻 3 离开，椅子 1 变成未占据。
# - 朋友 0 时刻 4 离开，椅子 0 变成未占据。
# - 朋友 2 时刻 4 到达，占据椅子 0 。
# 朋友 1 占据椅子 1 ，所以返回 1 。
