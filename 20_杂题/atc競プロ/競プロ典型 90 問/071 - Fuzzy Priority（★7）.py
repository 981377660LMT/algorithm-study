# 输出k种拓扑排序方案
# 不存在k种输出-1

# n,m<=1e5
# k<=10


# 拓扑排序+dfs全探索
# 不能bfs 时间空间都会爆炸
# !因为只要输出k种拓扑排序方案 所以用dfs回溯
from collections import defaultdict, deque
import sys
from typing import Deque, List

sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline
MOD = int(1e9 + 7)


n, m, k = map(int, input().split())
adjMap = defaultdict(set)
deg = defaultdict(int)
for _ in range(m):
    a, b = map(int, input().split())
    adjMap[a].add(b)
    deg[b] += 1


queue = deque([i for i in range(1, n + 1) if deg[i] == 0])  # !合法的拓扑排序顺序
res = []


def bt(path: List[int]) -> None:
    """因为要找到k个解马上退出 所以不能bfs 要dfs回溯来拓扑排序层序遍历"""
    if len(path) == n:
        res.append(path[:])
        if len(res) == k:
            for nums in res:
                print(*nums)
            exit(0)
        return  # 注意return

    # !剪枝
    if len(queue) == 0:
        print(-1)
        exit(0)

    len_ = len(queue)
    for _ in range(len_):  # !逐个元素作为当前队头看一遍 (类似rotate 取出 popleft 最后 append)
        cur = queue.popleft()
        path.append(cur)
        tmp = []

        for next in adjMap[cur]:
            deg[next] -= 1
            tmp.append(next)
            if deg[next] == 0:
                queue.append(next)
        bt(path)  # 这之前是bfs

        path.pop()
        for next in tmp:
            deg[next] += 1
            if deg[next] == 1:
                queue.pop()

        queue.append(cur)


bt([])
print(-1)

