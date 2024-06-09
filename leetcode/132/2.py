from typing import List, Tuple, Optional
from collections import defaultdict, Counter, deque
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 有 n 位玩家在进行比赛，玩家编号依次为 0 到 n - 1 。

# 给你一个长度为 n 的整数数组 skills 和一个 正 整数 k ，其中 skills[i] 是第 i 位玩家的技能等级。skills 中所有整数 互不相同 。

# 所有玩家从编号 0 到 n - 1 排成一列。

# 比赛进行方式如下：

# 队列中最前面两名玩家进行一场比赛，技能等级 更高 的玩家胜出。
# 比赛后，获胜者保持在队列的开头，而失败者排到队列的末尾。
# 这个比赛的赢家是 第一位连续 赢下 k 场比赛的玩家。


# 请你返回这个比赛的赢家编号。
class Solution:
    def findWinningPlayer(self, skills: List[int], k: int) -> int:
        n = len(skills)
        win = [0] * n
        queue = deque(range(n))
        for _ in range(2 * n):
            i, j = queue[0], queue[1]
            if skills[i] > skills[j]:
                win[i] += 1
                if win[i] == k:
                    return i
                queue.popleft()
                queue.popleft()
                queue.appendleft(i)
                queue.append(j)
            else:
                win[j] += 1
                if win[j] == k:
                    return j
                queue.append(queue.popleft())

        max_ = max(win)
        return win.index(max_)


# [8,9,7,19,11]
# 3

print(Solution().findWinningPlayer(skills=[8, 9, 7, 19, 11], k=3))
