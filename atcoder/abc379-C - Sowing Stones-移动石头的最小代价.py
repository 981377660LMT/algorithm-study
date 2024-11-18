# abc379-C - Sowing Stones-移动石头的最小代价
# https://atcoder.jp/contests/abc379/tasks/abc379_c
#
# n个格子，初始有m个格子有棋子数ai。
# 现可进行操作，如果第i个格子有棋子，可以将其一个棋子移动到第i+1个格子。
# 问最少进行的操作数，使得每个格子都恰好有一个棋子。
# 如果无法实现，输出-1。


from typing import List


def sowingStones(n: int, pos: List[int], stones: List[int]) -> int:
    if sum(stones) != n:
        return -1
    pairs = [(p, v) for p, v in zip(pos, stones)]
    pairs.sort()
    curSum1, curSum2 = 0, 0
    for p, v in pairs:
        if curSum1 < p - 1:
            return -1
        curSum1 += v
        curSum2 += p * v
    return n * (n + 1) // 2 - curSum2


if __name__ == "__main__":
    N, M = map(int, input().split())
    X = list(map(int, input().split()))
    A = list(map(int, input().split()))
    print(sowingStones(N, X, A))
