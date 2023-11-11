"""
n个操作 [t,y]
如果 t==1 那么就将当前的和替换为y
如果 t==2 那么就将当前的和加上y
有k次跳过操作的机会
最后和的最大值是多少
n<=2e5

nlogn
1.倒着枚举最后的覆盖位置 (因为覆盖某个位置的话前面就没有意义了 所以要倒着枚举)
2.大根堆反悔存最差的k个数
"""

from heapq import heappop, heappush
import sys
import os

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    n, k = map(int, input().split())
    opts = [(1, 0)]  # !开头算一次用0覆盖 哨兵
    for _ in range(n):
        t, y = map(int, input().split())
        opts.append((t, y))

    pq = []  # 大根堆 容量为k 可以存储k个最差的数不选
    res, curSum = -int(1e18), 0  # curSum 表示选取了后面的哪些yi
    for t, addOrTarget in opts[::-1]:
        if t == 1:
            res = max(res, curSum + addOrTarget)
            k -= 1  # 本身取消覆盖就用去了一次
            if k < 0:
                break
        else:
            if addOrTarget >= 0:
                curSum += addOrTarget
            else:
                heappush(pq, -addOrTarget)

        while len(pq) > k:
            preMax = -heappop(pq)
            curSum += preMax

    print(res)


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
