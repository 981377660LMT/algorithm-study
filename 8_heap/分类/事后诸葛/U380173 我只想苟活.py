# https://www.luogu.com.cn/problem/U380173

# 魔塔问题。
# 给定q个查询，每个查询给定一个初始hp。
# 有n个物品，每个物品会加血或者扣血.
# 你可以跳过任意物品.
# 问：在生命值始终不低于0的情况下，最多可以获得多少个物品。

# 解:
# 1.每个正数必选，且每个正数只对右侧的负数有影响.为了消除后效性，需要倒序遍历.
# 2.倒序遍历维护一个大根堆，堆中的数可以看成无法选择的物品.每次遇到负数就加入，
#   遇到正数就用于抵消堆顶的负数(此时不用，更待何时)，同时计数+1.
# 3.对堆中的数据维护一个前缀和，二分寻找答案边界.

# !注意到可以在线查询.


from typing import List
from bisect import bisect_right
from heapq import heappop, heappush


def solve(points: List[int], initHps: List[int]) -> List[int]:
    selected = 0
    pq = []
    for v in reversed(points):
        if v > 0:
            selected += 1
            while pq and v:
                if v >= pq[0]:
                    v -= heappop(pq)
                    selected += 1
                else:
                    pq[0] -= v
                    v = 0
        elif v == 0:
            selected += 1
        else:
            heappush(pq, -v)

    preSum = [0] * (len(pq) + 1)
    for i in range(len(pq)):
        preSum[i + 1] = preSum[i] + heappop(pq)
    return [selected + bisect_right(preSum, hp) - 1 for hp in initHps]


if __name__ == "__main__":
    import sys

    input = sys.stdin.readline

    T = int(input())
    for _ in range(T):
        n, q = map(int, input().split())
        points = list(map(int, input().split()))
        initHps = [int(input()) for _ in range(q)]
        res = solve(points, initHps)
        print(*res, sep="\n")

    def bf(points: List[int], initHps: List[int]) -> List[int]:
        res = []
        for initHp in initHps:
            cur = 0
            for state in range(1 << len(points)):
                hp = initHp
                todo = []
                for i in range(len(points)):
                    if state & (1 << i):
                        todo.append(points[i])
                for v in todo:
                    hp += v
                    if hp < 0:
                        break
                else:
                    cur = max(cur, len(todo))
            res.append(cur)
        return res

    def check() -> None:
        from random import randint

        for _ in range(10):
            for n in range(4):
                for q in range(2):
                    points = [randint(-5, 5) for _ in range(n)]
                    initHps = [randint(0, 10) for _ in range(q)]
                    res1 = solve(points, initHps)
                    res2 = bf(points, initHps)
                    if res1 != res2:
                        print(points, initHps)
                        print(res1, res2)
                        return

    check()

    # print(solve([-6, 6], [3]))
