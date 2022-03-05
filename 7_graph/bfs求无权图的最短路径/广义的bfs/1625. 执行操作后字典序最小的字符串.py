# 2 <= s.length <= 100

# 累加：将  a 加到 s 中所有下标为奇数的元素上（下标从 0 开始）
# 轮转：将 s 向右轮转 b 位。

# 在本题的数据范围内，可以枚举所有可能的操作结果
# Traverse the whole graph via the two operations, and return the minimum string
from collections import deque


class Solution:
    def findLexSmallestString(self, s: str, a: int, b: int) -> str:
        n = len(s)
        res = s

        # bfs搜索
        visited = set([s])
        queue = deque([s])
        while queue:
            cur = queue.popleft()
            res = min(res, cur)

            #### 方式a
            nxt = list(cur)
            for i in range(1, n, 2):
                nxt[i] = str((int(nxt[i]) + a) % 10)
            nxt_a_str = ''.join(nxt)

            if nxt_a_str not in visited:
                visited.add(nxt_a_str)
                queue.append(nxt_a_str)

            #### 方式b
            nxt_b_str = cur[-b:] + cur[:-b]
            if nxt_b_str not in visited:
                visited.add(nxt_b_str)
                queue.append(nxt_b_str)

        return res


print(Solution().findLexSmallestString(s="5525", a=9, b=2))
# 输出："2050"
# 解释：执行操作如下：
# 初态："5525"
# 轮转："2555"
# 累加："2454"
# 累加："2353"
# 轮转："5323"
# 累加："5222"
# 累加："5121"
# 轮转："2151"
# 累加："2050"​​​​​​​​​​​​
# 无法获得字典序小于 "2050" 的字符串。
