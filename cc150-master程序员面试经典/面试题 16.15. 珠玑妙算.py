from collections import Counter
from typing import List


# 要是猜对某个槽的颜色，则算一次“猜中”；要是只猜对颜色但槽位猜错了，则算一次“伪猜中”。注意，“猜中”不能算入“伪猜中”
# 伪猜中 = 总次数 - 猜中次数
class Solution:
    def masterMind(self, solution: str, guess: str) -> List[int]:
        a = sum(i == j for i, j in zip(solution, guess))
        b = sum((Counter(solution) & Counter(guess)).values())
        # print(Counter(solution))
        # print((Counter(solution).values()))
        # print(Counter(guess))
        # print((Counter(solution) & Counter(guess)))
        return [a, b - a]


print(Solution().masterMind("RGBYYR", "GGRRRR"))

