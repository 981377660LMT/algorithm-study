# https://www.cnblogs.com/zengzk/p/16584644.html

# 对于字符串S，在其后面添加一个新的字符，
# 增加的收益只和新字符以及S最后两个字符有关。
# !n<=18278(26+26**2+26**3)
# !len(wordi)<=3
# 因为len(wordi)<=3,所以添加一个字符的状态转移可以全部处理出来
# 等价于求图中的最长路，起点为空字符串 (原图中最多有27*27个结点)

import sys
from collections import defaultdict, deque


sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


n = int(input())
adjMatrix = [[0] * 26 for _ in range(27 * 27)]  # 邻接矩阵
for _ in range(n):
    word, score = input().split()
    score = int(score)
    cur = ord(word[-1]) - 97
    if len(word) == 1:
        for pre in range(27 * 27):
            adjMatrix[pre][cur] += score
    elif len(word) == 2:
        for pre in range(ord(word[0]) - 97, 27 * 27, 27):
            adjMatrix[pre][cur] += score
    elif len(word) == 3:
        pre = 27 * (ord(word[0]) - 97) + ord(word[1]) - 97
        adjMatrix[pre][cur] += score

# !spfa判断正环+求最长路

start = 27 * 27 - 1  # 起点为两个空字符串$$
dist = defaultdict(lambda: -INF, {start: 0})
queue = deque([start])
inQueue = defaultdict(lambda: False)
count = defaultdict(int)  # 边数更新次数
while queue:
    cur = queue.popleft()
    inQueue[cur] = False
    for i in range(26):
        weight = adjMatrix[cur][i]
        cand = dist[cur] + weight
        next = 27 * (cur % 27) + i
        if cand > dist[next]:
            dist[next] = cand
            count[next] = count[cur] + 1
            if count[next] >= 27 * 27:  # 检测到正环
                print("Infinity")
                exit(0)
            if not inQueue[next]:
                inQueue[next] = True
                if queue and dist[next] < dist[queue[0]]:  # !酸辣粉优化
                    queue.appendleft(next)
                else:
                    queue.append(next)


res = -INF
for key in dist:
    if key == start:
        continue
    res = max(res, dist[key])
print(res)
