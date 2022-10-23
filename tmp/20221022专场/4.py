from functools import lru_cache
from itertools import accumulate
from typing import List, Tuple, Optional
from collections import defaultdict, Counter, deque
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 某电商平台举办了一个用户抽奖活动，奖池中共有若干个礼包，每个礼包中包含一些礼物。 present[i][j] 表示第 i 个礼包第 j 件礼（下标从 0 开始）物的价值。抽奖规则如下：

# 每个礼包中的礼物摆放是有顺序的，你必须从第 0 件礼物开始打开；
# 对于同一个礼包中的礼物，必须在打开该礼包的第 i 个礼物之后，才能打开第 i+1 个礼物；
# 每个礼物包中的礼物价值 非严格递增。
# 参加活动的用户总共可以打开礼物 limit 次，请返回用户能够获得的 最大 礼物价值总和。
# 1 <= present.length <= 2000
# 1 <= present[i].length <= 1000
# 1 <= present[i][j] <= present[i][j+1] <= 10^5
# 1 <= limit <= 1000


class Solution:
    def brilliantSurprise(self, present: List[List[int]], limit: int) -> int:
        # n, cap = map(int, input().split())
        n = len(present)
        goods = []
        for i in range(n):
            cur = [(0, 0)]
            for j in range(len(present[i])):
                preSum, preVal = cur[-1]
                cur.append((preSum + present[i][j], preVal + 1))
            goods.append(cur)

        dp = [0] * (limit + 1)
        queue = deque()  # 单减的单调队列
        for _ in range(n):
            cost, score, count = map(int, input().split())
            for j in range(cost):
                queue.clear()  # 注意每次新的循环都需要初始化队列
                remain = (limit - j) // cost  # 剩余的空间最多还能放几个当前物品
                for k in range(remain + 1):
                    val = dp[k * cost + j] - k * score
                    while queue and val >= queue[-1][1]:
                        queue.pop()
                    queue.append((k, val))
                    while queue and queue[0][0] < k - count:  # 存放的个数不能超出物品数量，否则弹出
                        queue.popleft()
                    dp[k * cost + j] = queue[0][1] + k * score

        return dp[limit]
