# 送快递
# F - Shipping
# https://atcoder.jp/contests/abc374/tasks/abc374_f
# 给定N个货物，每个货物下单时间为Ti.
# 快递一次最多可以运送K个货物.
# 每次出发送快递后，需要X天才能再次发送快递.
# !记不满度为所有货物的送达时间减去下单时间之和.
# 求最小不满度.
# K<=N<=100
#
# !dfs(i,time) 表示当前送完了前i个货物，此时时间为time 状态下的最小不满度.
# O(N^2*K)


from itertools import accumulate
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")

from functools import lru_cache

INF = int(4e18)


def min2(a: int, b: int) -> int:
    return a if a < b else b


if __name__ == "__main__":
    N, K, X = map(int, input().split())
    T = list(map(int, input().split()))

    T.sort()
    preSum = [0] + list(accumulate(T))

    @lru_cache(None)
    def dfs(index: int, time: int) -> int:
        if index == N:
            return 0

        res = INF
        for j in range(K):  # 送(j+1)个货物
            if index + j >= N:
                break
            if time >= T[index + j]:  # !送货[index, index+j+1)
                curCost = time * (j + 1) - (preSum[index + j + 1] - preSum[index])
                res = min2(res, dfs(index + j + 1, time + X) + curCost)
            else:
                res = min2(res, dfs(index, T[index + j]))  # !到下一个送货时间
        return res

    res = dfs(0, 0)
    dfs.cache_clear()
    print(res)
