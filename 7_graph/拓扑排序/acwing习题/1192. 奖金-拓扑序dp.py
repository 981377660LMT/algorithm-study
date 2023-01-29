# 每位参加会谈的代表提出了自己的意见：“我认为员工 a 的奖金应该比 b 高！”
# Mr.Z 决定要找出一种奖金方案，满足各位代表的意见，且同时使得总奖金数最少。
# 每位员工奖金最少为 100 元，且必须是整数。


from collections import deque


n, m = map(int, input().split())
adjList = [set() for _ in range(n)]
indegree = [0] * n
visitedPair = set()
for _ in range(m):
    next, cur = map(int, input().split())
    next, cur = next - 1, cur - 1
    if (cur, next) not in visitedPair:
        visitedPair.add((cur, next))
        adjList[cur].add(next)
        indegree[next] += 1

count = 0
queue = deque([i for i in range(n) if indegree[i] == 0])
dp = [100] * n
while queue:
    cur = queue.popleft()
    count += 1
    for next in adjList[cur]:
        indegree[next] -= 1
        dp[next] = max(dp[next], dp[cur] + 1)
        if indegree[next] == 0:
            queue.append(next)

if count != n:
    print('Poor Xed')
else:
    print(sum(dp))

