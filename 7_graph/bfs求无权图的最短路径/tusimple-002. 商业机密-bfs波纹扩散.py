# 员工的索引是从 1 到 N 递增的整数，假设每条边的长度都是 1 。

from collections import defaultdict, deque


# 员工数量，关系数量，参加秘密会议的人数，以及任务描述中的 k。
n, m, q, k = map(int, input().split())
joins = list(map(int, input().split()))
# 从会议结束后的第一天开始，
# 所有知道秘密的人都会和距离（在上述关系图中的距离）自己不超过 k 的所有人共享这个秘密，
# 而第二天，
# 所有知道秘密的人都会和距离自己不超过 2k 的所有人共享这个秘密。也就是说，在会议结束后的第 x 天，
# 所有知道秘密的人都开始与距自己最多 的所有人共享该秘密。
# 现在 TuTu 被委托计算每个人都会在什么时候知道这个秘密。

res = [-1] * (n + 1)
adjMap = defaultdict(set)
for _ in range(m):
    u, v = map(int, input().split())
    adjMap[u].add(v)
    adjMap[v].add(u)

for i in joins:
    res[i] = 0

count = 0
day = 1
queue = deque(joins)
while count < n:
    for _ in range(day * k):  # 这一次看几层
        len_ = len(queue)
        for _ in range(len_):  # 每一层
            cur = queue.popleft()
            count += 1
            for next in adjMap[cur]:
                if res[next] == -1:
                    res[next] = day  # 记录第一次知道秘密的时间
                    queue.append(next)
    day += 1


for num in res[1:]:
    print(num, end=' ')


# bfs波纹扩散
