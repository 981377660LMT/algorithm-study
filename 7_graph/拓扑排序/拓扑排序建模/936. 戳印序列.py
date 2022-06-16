from typing import List
from collections import deque

# 如果可以印出序列，那么返回一个数组，该数组由每个回合中被印下的最左边字母的索引组成
# https://leetcode-cn.com/problems/stamping-the-sequence/solution/tuo-bu-pai-xu-by-kevinchen147-09ad/
# 这题很像
# 1591. 奇怪的打印机 II.py

# 拓扑排序：
# 倒着看，可以想象该窗口中的字符全部替换为 '*' ，表示可以匹配任意字符。
# !入度为窗口内与印章不同的字符数
# 边为字符 => 窗口

# 邻接表:字符index对应窗口


class Solution:
    def movesToStamp(self, stamp: str, target: str) -> List[int]:
        m, n = len(stamp), len(target)
        indeg = [m] * (n - m + 1)
        adjList = [[] for _ in range(n)]
        queue = deque()
        visited = [False] * n

        # 窗口id为区间左侧索引
        for i in range(n - m + 1):
            for j in range(m):
                if target[i + j] == stamp[j]:
                    indeg[i] -= 1
                    if indeg[i] == 0:
                        queue.append(i)
                else:
                    adjList[i + j].append(i)

        res = []
        while queue:
            cur = queue.popleft()
            res.append(cur)
            for i in range(m):
                if visited[cur + i]:
                    continue
                visited[cur + i] = True
                for next in adjList[cur + i]:
                    indeg[next] -= 1
                    if indeg[next] == 0:
                        queue.append(next)

        return res[::-1] if len(res) == n - m + 1 else []


print(Solution().movesToStamp(stamp="abc", target="ababc"))
# 输出：[0,2]
# （[1,0,2] 以及其他一些可能的结果也将作为答案被接受）

