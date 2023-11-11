# 有n张反着的牌，从上到下第i张牌有一个数字Pi;(1<Pi≤n)，互不相同。
# 你将做以下操作n 次:
# 每次从最上面取一张，假设为X，准备将它放到已经取出的牌堆中。
# 找到牌堆中第一个>=X的牌，将X放在这堆牌的顶部。
# 如果找不到就自己新立一堆。
# 特别的，当一个牌堆的牌数等于K时，这堆牌将会被清除。
# 问:数字为i的牌是在第几次操作后被清除的?如果没被清除，输出―1。(1≤K ≤n≤2·105) .


# !1. 有序容器来查找/维护 牌堆
# !2. "并查集" 或者 "pre/next 数组模拟链表" 记录每个牌堆的牌
from sortedcontainers import SortedList


import sys
import os
from typing import List, Tuple


sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    n, k = map(int, input().split())
    nums = list(map(int, input().split()))

    sl = SortedList()  # 牌堆
    pre: List[Tuple[int, int]] = [None] * (n + 1)  # type: ignore 每张牌的(前驱，所在牌堆牌数)
    res = [-1] * (n + 1)

    for i, num in enumerate(nums, start=1):
        pos = sl.bisect_left(num)
        if pos == len(sl):
            sl.add(num)
            pre[num] = (-1, 1)
        else:
            toRemove = sl.pop(pos)
            sl.add(num)
            pre[num] = (toRemove, pre[toRemove][1] + 1)

        if pre[num][1] == k:
            sl.pop(pos)
            cur = num
            res[num] = i
            while pre[cur][0] != -1:
                cur = pre[cur][0]
                res[cur] = i

    print(*res[1:], sep="\n")


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
