# 2387. 缓存交换( 最优页面置换算法)
# https://www.acwing.com/problem/content/2389/
# 在现代计算机中，往往采用 LRU(最近最少使用)的算法来进行 Cache 调度
# 但这并不是最优的算法
# 对于一个固定容量的空 Cache 和连续的若干主存访问请求，
# 求出如何在每次 Cache 缺失时换出正确的主存单元，以达到最少的 Cache 缺失次数
# !输出 Cache 缺失次数的最小值
#
# 若不满，则直接放。
# !若满，则弹出当前 Cache 中距离下一次出现位置最远的元素
# 使用堆来维护即可


from heapq import heappop, heappush
from typing import List

INF = int(1e18)


def cacheExchange(nums: List[int], cap: int) -> int:
    pool = dict()
    for i, v in enumerate(nums):
        nums[i] = pool.setdefault(v, len(pool))

    n, m = len(nums), len(pool)
    nexts, valueNexts = [INF] * n, [INF] * m
    for i, v in reversed(list(enumerate(nums))):
        tmp = valueNexts[v]
        nexts[i] = tmp if tmp != INF else n
        valueNexts[v] = i

    inCache = [False] * m
    pq = []
    res = 0
    used = 0
    for i, v in enumerate(nums):
        if not inCache[v]:
            if used < cap:
                used += 1
            else:
                _, index = heappop(pq)
                inCache[nums[index]] = False
            res += 1
            inCache[v] = True
        heappush(pq, (-nexts[i], i))
    return res


if __name__ == "__main__":
    import sys

    input = lambda: sys.stdin.readline().rstrip("\r\n")
    n, size = map(int, input().split())
    nums = list(map(int, input().split()))
    print(cacheExchange(nums, size))
