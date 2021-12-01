from typing import List
from heapq import heappop, heappush

# 在第 i 天，树上会长出 apples[i] 个苹果，这些苹果将会在 days[i] 天后（也就是说，第 i + days[i] 天时）腐烂
# 你打算每天 最多 吃一个苹果来保证营养均衡。注意，你可以在这 n 天之后继续吃苹果。
# 1353. 最多可以参加的会议数目.py
class Solution:
    def eatenApples(self, apples: List[int], days: List[int]) -> int:
        res = 0
        pq = []
        index = 0
        while pq or index < len(apples):
            if index < len(apples):
                heappush(pq, [index + days[index], apples[index]])
            while pq and (pq[0][0] <= index or pq[0][1] == 0):
                heappop(pq)
            if pq:
                res += 1
                pq[0][1] -= 1
            index += 1
        return res


print(Solution().eatenApples(apples=[1, 2, 3, 5, 2], days=[3, 2, 1, 4, 2]))
# 输出：7
# 解释：你可以吃掉 7 个苹果：
# - 第一天，你吃掉第一天长出来的苹果。
# - 第二天，你吃掉一个第二天长出来的苹果。
# - 第三天，你吃掉一个第二天长出来的苹果。过了这一天，第三天长出来的苹果就已经腐烂了。
# - 第四天到第七天，你吃的都是第四天长出来的苹果。

