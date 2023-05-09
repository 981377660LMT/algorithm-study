from typing import Deque, List, Tuple, Optional
from collections import defaultdict, Counter, deque
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 在永恒之森中，存在着一本生物进化录，以 一个树形结构 记载了所有生物的演化过程。经过观察并整理了各节点间的关系，parents[i] 表示编号 i 节点的父节点编号(根节点的父节点为 -1)。

# 为了探索和记录其中的演化规律，队伍中的炼金术师提出了一种方法，可以以字符串的形式将其复刻下来，规则如下：

# 初始只有一个根节点，表示演化的起点，依次记录 01 字符串中的字符，
# 如果记录 0，则在当前节点下添加一个子节点，并将指针指向新添加的子节点；
# 如果记录 1，则将指针回退到当前节点的父节点处。
# 现在需要应用上述的记录方法，复刻下它的演化过程。请返回能够复刻演化过程的字符串中， 字典序最小 的 01 字符串。

# 注意：

# 节点指针最终可以停在任何节点上，不一定要回到根节点。


# !0:父结点，1:子结点
# 字典序最小是什么意思
# 每个结点处dfs沿着最长链走
class Solution:
    def evolutionaryRecord(self, parents: List[int]) -> str:
        def dfs(u: int, p: int, d: int) -> Deque[int]:
            V = []
            for v in adjList[u]:
                if v != p:
                    V.append(dfs(v, u, d + 1))
            res = deque([d])
            V.sort(reverse=True)
            for A in V:
                if len(res) > len(A):
                    while A:
                        res.append(A.popleft())
                else:
                    res, A = A, res
                    while A:
                        res.appendleft(A.pop())
                res.append(d)
            return res

        n = len(parents)
        root = -1
        adjList = [[] for _ in range(n)]
        for i, p in enumerate(parents):
            if p == -1:
                root = i
            else:
                adjList[p].append(i)

        path = list(dfs(root, -1, 0))
        res = []
        for pre, cur in zip(path, path[1:]):
            if pre < cur:
                res.append("0")
            else:
                res.append("1")
        while res and res[-1] == "1":
            res.pop()
        return "".join(res)


# parents = [-1,0,0,2]
# parents = [-1,0,0,1,2,2]

parents = [-1, 0, 0, 2]
print(Solution().evolutionaryRecord(parents))
print(Solution().evolutionaryRecord([-1, 0, 0, 1, 2, 2]))
