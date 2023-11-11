"""华容道 每次状态转移找到哪个位置没有棋子"""
# 题意: 8-puzzle 问题的变种，
# 一共9个格子,有8个棋子和1个空格
# !求将 0 到 7 编号的棋子分别放到 0 到 7 编号的节点上的最小步数，
# 同一时间每个节点最多只能放 1 个棋子，节点之间存在转移边
from collections import deque
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


if __name__ == "__main__":
    m = int(input())
    adjList = [[] for _ in range(9)]  # 9个格子
    for _ in range(m):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        adjList[u].append(v)
        adjList[v].append(u)
    pos = tuple([int(x) - 1 for x in input().split()])  # !初始时0-7号棋子所在的位置
    state = [-1] * 9  # !初始时每个格子上面的棋子编号 -1表示空格
    for i, p in enumerate(pos):
        state[p] = i

    visited = set([tuple(state)])
    queue = deque([(tuple(state), 0)])
    target = tuple(range(8)) + (-1,)
    while queue:
        curState, curStep = queue.popleft()
        if curState == target:
            print(curStep)
            exit(0)

        empty = curState.index(-1)  # !找到空格的位置
        curList = list(curState)
        for next in adjList[empty]:
            curList[empty], curList[next] = curList[next], curList[empty]
            nextState = tuple(curList)
            if nextState not in visited:
                visited.add(nextState)
                queue.append((nextState, curStep + 1))
            curList[empty], curList[next] = curList[next], curList[empty]

    print(-1)
