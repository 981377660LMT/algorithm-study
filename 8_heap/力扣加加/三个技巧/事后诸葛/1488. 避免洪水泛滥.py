# /**
#  * @param {number[]} rains
#  * @return {number[]}
#  * @description 这道题没有使用到堆.事后诸葛亮这个技巧并不是堆特有的，实际上这就是一种普通的算法思想
#  * 所有湖泊一开始都是空的。当第 n 个湖泊下雨的时候，如果第 n 个湖泊是空的，那么它就会装满水，否则这个湖泊会发生洪水。你的目标是避免任意一个湖泊发生洪水。
#  * rains[i] > 0 表示第 i 天时，第 rains[i] 个湖泊会下雨。
#  * rains[i] === 0 表示第 i 天没有湖泊会下雨，你可以选择 一个 湖泊并 抽干 这个湖泊的水。
#  * 请返回一个数组 ans
#  * 如果 rains[i] > 0 ，那么ans[i] == -1 。
#    如果 rains[i] == 0 ，ans[i] 是你第 i 天选择抽干的湖泊。
#    如果没办法阻止洪水，请返回一个 空的数组 。
#  */
from typing import List
from collections import defaultdict, deque
from heapq import heappush, heappop

# 总结:对每个湖泊记录下雨天数；优先队列存离现在最近的下雨天；晴天直接抽干
# 贪心+优先队列, 从前向后遍历, 将`已经满了的湖的下一个下雨日期`加入优先队列中,
# 遇到0优先抽当前已满的湖中下次下雨日期距离现在最近的湖
# https://leetcode-cn.com/problems/avoid-flood-in-the-city/solution/tan-xin-you-xian-dui-lie-he-xin-si-lu-yi-ju-hua-by/
class Solution:
    def avoidFlood(self, rains: List[int]) -> List[int]:
        n = len(rains)
        res = [-1] * n
        pq = []
        full = set()
        lakeRainDays = defaultdict(deque)
        for i, rain in enumerate(rains):
            if rain > 0:
                lakeRainDays[rain].append(i)

        for i, rain in enumerate(rains):
            if rain > 0:
                if rain in full:
                    return []

                res[i] = -1
                full.add(rain)
                lakeRainDays[rain].popleft()
                # 如果这个湖接下来某天还要下雨, 将`下一个下雨日期`加入优先队列
                if lakeRainDays[rain]:
                    heappush(pq, lakeRainDays[rain][0])

            else:
                if not pq:
                    # 随便抽一个湖就行, 这里选1
                    res[i] = 1
                else:
                    toRemoveIndex = heappop(pq)
                    res[i] = rains[toRemoveIndex]
                    full.remove(rains[toRemoveIndex])

        return res


print(Solution().avoidFlood(rains=[1, 2, 0, 0, 2, 1]))
# 输入：rains = [1,2,0,0,2,1]
# 输出：[-1,-1,2,1,-1,-1]
# 解释：第一天后，装满水的湖泊包括 [1]
# 第二天后，装满水的湖泊包括 [1,2]
# 第三天后，我们抽干湖泊 2 。所以剩下装满水的湖泊包括 [1]
# 第四天后，我们抽干湖泊 1 。所以暂时没有装满水的湖泊了。
# 第五天后，装满水的湖泊包括 [2]。
# 第六天后，装满水的湖泊包括 [1,2]。
# 可以看出，这个方案下不会有洪水发生。同时， [-1,-1,1,2,-1,-1] 也是另一个可行的没有洪水的方案。

