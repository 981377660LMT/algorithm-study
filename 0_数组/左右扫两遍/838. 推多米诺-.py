from collections import deque

# 我们要统计所有点最近的'L'和'R'，来判断最终该点的状态。
# 如果‘L’更近，最终为'L'；如果'R'更近，最终为'R'；如果一样近，最终为'.'
INF = int(1e20)


class Solution:
    def pushDominoes(self, dominoes: str) -> str:
        n = len(dominoes)
        # 到最近的L/R的距离
        leftR, rightL = [INF] * n, [INF] * n

        pre = -INF
        for i in range(n):
            if dominoes[i] == "L":
                pre = -INF
            elif dominoes[i] == "R":
                pre = i
            leftR[i] = i - pre

        pre = INF
        for i in range(n - 1, -1, -1):
            if dominoes[i] == "R":
                pre = INF
            elif dominoes[i] == "L":
                pre = i
            rightL[i] = pre - i

        res = []
        for d1, d2 in zip(leftR, rightL):
            if d1 == d2 or (d1 > n and d2 > n):
                res.append(".")
            elif d1 < d2:
                res.append("R")
            else:
                res.append("L")

        return "".join(res)


print(Solution().pushDominoes(".L.R...LR..L.."))
# 输出："LL.RR.LLRRLL.."
