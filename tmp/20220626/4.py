from itertools import combinations
from typing import List, Tuple
from collections import defaultdict, Counter, deque

MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def minimumScore(self, nums: List[int], edges: List[List[int]]) -> int:
        """删除树中两条 不同 的边以形成三个连通组件
        
        分别获取三个组件 每个 组件中所有节点值的异或值。
        最大 异或值和 最小 异或值的 差值 就是这一种删除边方案的分数。

        # 3 <= n <= 1000
        枚举要删除的两条边(的顶点)
        
        parent写错了 不能只记录一个父节点

        判断祖先可以用 dfs 序
        """

        def dfs(cur: int, pre: int, rootXorVal: int) -> None:
            subxor[cur] = nums[cur]
            rootXor[cur] = rootXorVal
            for next in adjMap[cur]:
                if next == pre:
                    continue
                dfs(next, cur, rootXorVal ^ nums[next])
                parents[next].add(cur)  # !parent写错了 这里只记录了父结点 忘记记录祖先节点了🤣
                subxor[cur] ^= subxor[next]

        n = len(nums)
        adjMap = defaultdict(set)
        for u, v in edges:
            adjMap[u].add(v)
            adjMap[v].add(u)

        subxor = [0] * n
        parents = [set()] * n
        rootXor = [0] * n  # 到根节点的异或值
        dfs(0, -1, nums[0])

        allXor = subxor[0]
        res = int(1e20)
        # 枚举中间点

        "枚举两个低的点 p1 可能是 p2 的父节点"
        for p1, p2 in combinations(range(n), 2):
            # if p1 == 0 or p2 == 0:
            #     continue
            isP1Parent = p1 in parents[p2]
            isP2Parent = p2 in parents[p1]
            isParent = isP1Parent or isP2Parent
            if isP2Parent:
                p1, p2 = p2, p1

            if not isParent:
                xor2 = subxor[p2]
                xor1 = subxor[p1]
                xor3 = allXor ^ xor1 ^ xor2
                xor1, xor2, xor3 = sorted([xor1, xor2, xor3])
                cand = xor3 - xor1
                if cand < res:
                    # print(p1, p2, xor1, xor2, xor3, 999, isParent, subxor[p1])
                    res = cand
            else:
                xor2 = subxor[p2]
                xor1 = rootXor[p1] ^ rootXor[p2] ^ nums[p1] ^ nums[p2]
                xor3 = allXor ^ xor1
                xor1, xor2, xor3 = sorted([xor1, xor2, xor3])
                cand = xor3 - xor1
                if cand < res:
                    # print(p1, p2, xor1, xor2, xor3, 1000, isParent, subxor[p1])
                    res = cand

        return res


print(Solution().minimumScore(nums=[1, 5, 5, 4, 11], edges=[[0, 1], [1, 2], [1, 3], [3, 4]]))  # 9
# print(
#     Solution().minimumScore(nums=[5, 5, 2, 4, 4, 2], edges=[[0, 1], [1, 2], [5, 2], [4, 3], [1, 3]])
# )

print(Solution().minimumScore(nums=[29, 29, 23, 32, 17], edges=[[3, 1], [2, 3], [4, 1], [0, 4]]))
# 15
print(
    Solution().minimumScore(
        nums=[28, 24, 29, 16, 31, 31, 17, 18],
        edges=[[0, 1], [6, 0], [6, 5], [6, 7], [3, 0], [2, 1], [2, 4]],
    )
)
# 8
