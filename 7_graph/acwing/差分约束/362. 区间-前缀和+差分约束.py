# 给定 n 个区间 [ai,bi] 和 n 个整数 ci。
# 你需要构造一个整数集合 Z，使得 ∀i∈[1,n]，Z 中满足 ai≤x≤bi 的整数 x 不少于 ci 个。
# 求这样的整数集合 Z 最少包含多少个数。
# 1≤n≤50000,
# 0≤ai,bi≤50000

# 前缀和Si表示[0,i]选择了多少个数 题目要求S50001的最小值  即求最长路


# 所有的限制要找全：
# 1. Si>=Si-1
# 2. Si-Si-1<=1
# 3. Sb-Sa-1>=c

from collections import defaultdict, deque
from typing import List, Tuple


n = int(input())
adjMap = defaultdict(set)
max_ = 50001
for _ in range(n):
    a, b, c = map(int, input().split())
    a, b = a + 1, b + 1
    adjMap[a - 1].add((b, c))
    max_ = max(max_, a, b)

# 前缀和满足的约束
for i in range(1, max_ + 1):
    adjMap[i - 1].add((i, 0))
    adjMap[i].add((i - 1, -1))


def spfa(n: int, adjMap: defaultdict) -> Tuple[bool, List[int]]:
    """spfa求单源最长路，顺便判断正环"""
    dist = [-int(1e20)] * (n + 1)
    # n+1个点，此时0是虚拟源点
    dist[0] = 0

    queue = deque([0])
    count = [0] * (n + 1)  # 边数
    isInqueue = [False] * (n + 1)  # 在队列里的点
    isInqueue[0] = True

    while queue:
        cur = queue.popleft()
        isInqueue[cur] = False  # 点从队列出来了

        # 更新过谁，就拿谁去更新别人
        for next, weight in adjMap[cur]:
            if dist[cur] + weight > dist[next]:
                dist[next] = dist[cur] + weight
                count[next] = count[cur] + 1
                if count[next] >= n + 1:  # n+1 个点
                    return False, []
                if not isInqueue[next]:
                    isInqueue[next] = True
                    # 酸辣粉优化：队列不为空，且当前元素距离小于队头，则加入队头，否则加入队尾
                    if queue and dist[next] > dist[queue[0]]:
                        queue.appendleft(next)
                    else:
                        queue.append(next)

    return True, dist


ok, dist = spfa(max_, adjMap)
if not ok:
    print(-1)
else:
    # 答案是前缀S50001这个
    print(dist[max_])
