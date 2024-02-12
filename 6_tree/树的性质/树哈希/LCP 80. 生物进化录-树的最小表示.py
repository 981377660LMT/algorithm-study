# LCP 80. 生物进化录-字典序最小欧拉路径(树的最小表示法/树的括号序列)
# 把所有子树的结果排序，然后拼接起来，在开头加上 0，末尾加上 1，然后返回
# https://leetcode.cn/problems/qoQAMX/solution/di-gui-pai-xu-by-endlesscheng-hnjf/
# https://leetcode.cn/problems/qoQAMX/solution/bao-li-pai-xu-de-fu-za-du-fen-xi-by-hqzt-clpn/

from collections import deque
from typing import Deque, List


class Solution:
    def evolutionaryRecord(self, parents: List[int]) -> str:
        """O(n^2) 拼接字符串是瓶颈."""
        n = len(parents)
        adjList = [[] for _ in range(n)]
        for i, p in enumerate(parents):
            if p != -1:
                adjList[p].append(i)

        def dfs(cur: int, pre: int) -> str:
            sub = []
            for next in adjList[cur]:
                if next != pre:
                    sub.append(dfs(next, cur))
            sub.sort()
            return "0" + "".join(sub) + "1"

        res = dfs(0, -1)
        return res[1:].rstrip("1")  # 去掉根节点以及返回根节点的路径

    def evolutionaryRecord2(self, parents: List[int]) -> str:
        """O(nlogn) 双端队列优化字符串的拼接.

        可以把字符串改成双端队列,拼接的时候用启发式合并/启发式分裂,也是O(nlogn)的
        """
        n = len(parents)
        adjList = [[] for _ in range(n)]
        for i, p in enumerate(parents):
            if p != -1:
                adjList[p].append(i)

        def dfs(cur: int, pre: int) -> Deque[int]:
            sub = sorted([dfs(next, cur) for next in adjList[cur] if next != pre])
            res = deque()
            for d in sub:
                if len(res) > len(d):
                    while d:
                        res.append(d.popleft())
                else:
                    res, d = d, res
                    while d:
                        res.appendleft(d.pop())
            res.appendleft(0)
            res.append(1)
            return res

        res = dfs(0, -1)
        while res and res[-1] == 1:  # 去掉返回根节点的路径
            res.pop()
        res.popleft()  # 去掉根节点0
        return "".join(map(str, res))


if __name__ == "__main__":
    parents = [3, 3, -1, 2, 2]
    print(Solution().evolutionaryRecord(parents))
