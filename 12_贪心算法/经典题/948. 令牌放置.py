from typing import List
from collections import deque

# 在使用任意数量的令牌后，返回我们可以得到的最大 分数 。
# 如果你至少有 token[i] 点 能量 ，可以将令牌 i 置为正面朝上，失去 token[i] 点 能量 ，并得到 1 分 。
# 如果我们至少有 1 分 ，可以将令牌 i 置为反面朝上，获得 token[i] 点 能量 ，并失去 1 分 。

# Sort tokens.
# Buy at the cheapest and sell at the most expensive.


class Solution:
    def bagOfTokensScore(self, tokens: List[int], power: int) -> int:
        queue = deque(sorted(tokens))
        score = 0

        while queue and power >= queue[0]:
            while queue and power >= queue[0]:
                power -= queue.popleft()
                score += 1
            if score > 0 and len(queue) > 1:
                power += queue.pop()
                score -= 1

        return score


print(Solution().bagOfTokensScore(tokens=[100, 200, 300, 400], power=200))
# 输出：2
# 解释：按下面顺序使用令牌可以得到 2 分：
# 1. 令牌 0 正面朝上，能量变为 100 ，分数变为 1
# 2. 令牌 3 正面朝下，能量变为 500 ，分数变为 0
# 3. 令牌 1 正面朝上，能量变为 300 ，分数变为 1
# 4. 令牌 2 正面朝上，能量变为 0 ，分数变为 2
