# 循环队列排队买票，每个人一次只能买一张，买完就到队伍后面继续排
# 问每个人的等待时间
# 0 ≤ n ≤ 100,000


from collections import deque
from sortedcontainers import SortedList


class Solution:
    def solve1(self, tickets):
        # 暴力模拟：TLE
        queue = deque([(i, ticket) for i, ticket in enumerate(tickets)])
        res = [0] * len(tickets)
        time = 0

        while queue:
            i, cur = queue.popleft()
            time += 1
            if cur - 1 > 0:
                queue.append((i, cur - 1))
            else:
                res[i] = time

        return res

    def solve(self, tickets):
        # 有序集合解法 https://leetcode-solution-leetcode-pp.gitbook.io/leetcode-solution/hard/ticket-order
        # 对于一个人 p 来说，他需要等待的时间 t 为 ：
        # 比他票少的人的等待总时长 allFast + 比他票多的人的等待总时长（包括自己，截止到t）allSlow + 排在他前面且票不比他少的总人数 leftSlow。
        n = len(tickets)
        res = [0] * len(tickets)
        sortedList = SortedList()
        people = sorted((need, rawId) for rawId, need in enumerate(tickets))
        allFast = 0
        for index, (need, rawId) in enumerate(people):
            allSlow = (n - index) * (need - 1)
            leftSlow = rawId + 1 - sortedList.bisect_left(rawId)
            res[rawId] = allFast + allSlow + leftSlow
            sortedList.add(rawId)
            allFast += need
        return res


print(Solution().solve(tickets=[2, 1, 3]))
