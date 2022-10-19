from typing import List
from collections import defaultdict, deque

# meetings[i] = [xi, yi, timei] 表示专家 xi 和专家 yi 在时间 timei 要开一场会
# 专家 0 有一个 秘密 ，最初，他在时间 0 将这个秘密分享给了专家 firstPerson 。
# 每次会议，如果专家 xi 在时间 timei 时知晓这个秘密，那么他将会与专家 yi 分享这个秘密，反之亦然。
# 在所有会议都结束之后，返回所有知晓这个秘密的专家列表

# 总结：
# !分组bfs/dfs 把开会时间相同的分到一组，按时间顺序bfs/dfs


class Solution:
    def findAllPeople(self, n: int, meetings: List[List[int]], firstPerson: int) -> List[int]:
        res = set([0, firstPerson])
        group = defaultdict(list)
        for u, v, time in meetings:
            group[time].append((u, v))

        for time in sorted(group):
            edges = group[time]
            adjMap = defaultdict(set)
            queue = deque()
            for u, v in edges:
                adjMap[u].add(v)
                adjMap[v].add(u)
                if u in res:
                    queue.append(v)
                if v in res:
                    queue.append(u)

            while queue:
                cur = queue.popleft()
                res.add(cur)
                for next in adjMap[cur]:
                    if next not in res:
                        queue.append(next)

        return list(res)


print(Solution().findAllPeople(n=6, meetings=[[1, 2, 5], [2, 3, 8], [1, 5, 10]], firstPerson=1))
# 输出：[0,1,2,3,5]
# 解释：
# 时间 0 ，专家 0 将秘密与专家 1 共享。
# 时间 5 ，专家 1 将秘密与专家 2 共享。
# 时间 8 ，专家 2 将秘密与专家 3 共享。
# 时间 10 ，专家 1 将秘密与专家 5 共享。
# 因此，在所有会议结束后，专家 0、1、2、3 和 5 都将知晓这个秘密。
