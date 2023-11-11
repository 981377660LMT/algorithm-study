# n个白球 n个道具
# 第i个道具只有在pre里的至少一个球为白色时才能被使用 使得第i个球变为黑色·
# 问能否所有球都变为黑色？
# 如果存在 输出道具的一种使用顺序，否则输出-1
# n<=1e5

# !倒着考虑 把所有黑球变白 把哪个球变白后 可以让哪个操作变有效
# !自环的点为起点 bfs检查是否能遍历完所有点

from collections import defaultdict, deque
import sys

sys.setrecursionlimit(int(1e6))
input = sys.stdin.readline
MOD = int(1e9 + 7)


n = int(input())
adjMap = defaultdict(set)
queue = deque()
visited = set()
for i in range(n):
    a, b = map(int, input().split())
    a, b = a - 1, b - 1
    adjMap[a].add(i)
    adjMap[b].add(i)
    if i == a or i == b:
        queue.append(i)
        visited.add(i)

res = []
while queue:
    cur = queue.popleft()
    res.append(cur)
    for next in adjMap[cur]:
        if next not in visited:
            queue.append(next)
            visited.add(next)

if len(visited) != n:
    print(-1)
    exit(0)

res.reverse()
for num in res:
    print(num + 1)
