# P1409 骰子-排队、移到队尾的概率
# https://www.luogu.com.cn/problem/P1409
# n 个人排成一排，你排在第 m 个。
# 每轮队首的人投一次骰子。
# 若掷到 1，则队首的人获胜。若掷到 2,4,6，则队首的人排到队尾。若掷到3,5，则队首的人出队。
# 若队列中仅剩一人，则该人获胜，求你获胜的概率。
# n,m <=1000

from typing import List


def solve(n: int, win: float, shift: float, out: float) -> List[float]:
    dp = [0.0] * (n + 1)
    a, tmp = [0.0] * (n + 1), [0.0] * (n + 1)
    dp[1] = 1
    for i in range(2, n + 1):
        a[1] = win
        for j in range(2, i + 1):
            a[j] = out * dp[j - 1]
        s, t = 0, 1
        for j in range(1, i + 1):
            s *= shift
            t *= shift
            s += a[j]
        tmp[i] = s / (1 - t)
        tmp[0] = tmp[i]
        for j in range(1, i):
            tmp[j] = a[j] + tmp[j - 1] * shift
        for j in range(1, i + 1):
            dp[j] = tmp[j]
    return dp


if __name__ == "__main__":
    n, m = map(int, input().split())
    dp = solve(n, 1 / 6, 1 / 2, 1 / 3)
    # 精确到小数点后 9 位。
    print(f"{dp[m]:.9f}")
