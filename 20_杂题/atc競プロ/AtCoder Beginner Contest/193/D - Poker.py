# 两个人打扑克 问高桥获胜的概率

from collections import Counter
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


def poker(k: int, s1: str, s2: str) -> float:
    counter = Counter({i: k for i in range(1, 10)})
    for c in s1[:-1] + s2[:-1]:
        counter[int(c)] -= 1

    # 枚举最后一张牌
    count, win = 0, 0
    for i in range(1, 10):
        for j in range(1, 10):
            score1, score2 = calScore(s1[:-1] + str(i)), calScore(s2[:-1] + str(j))
            if i == j:
                count += counter[i] * (counter[i] - 1)
                if score1 > score2:
                    win += counter[i] * (counter[i] - 1)
            else:
                count += counter[i] * counter[j]
                if score1 > score2:
                    win += counter[i] * counter[j]
    return win / count


def calScore(digits: str) -> int:
    counter = [0] * 10
    for c in digits:
        counter[int(c)] += 1
    return sum(i * pow(10, counter[i]) for i in range(1, 10))


if __name__ == "__main__":
    k = int(input())
    s1 = input()
    s2 = input()
    print(poker(k, s1, s2))
