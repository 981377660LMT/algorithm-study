# 2_树上所有路径的位运算异或的和
# 解法：
# !由于任意路径异或和可以用从根节点出发的路径异或和表示
# 对每一位，统计从根节点出发的路径异或和在该位上的 0 的个数和 1 的个数，
# 只有当 0 与 1 异或时才对答案有贡献，所以贡献即为这两个个数之积


from typing import List, Tuple


def xorPathSum(n: int, edges: List[Tuple[int, int]], values: List[int]) -> int:
    adjList = [[] for _ in range(n)]
    for u, v in edges:
        adjList[u].append(v)
        adjList[v].append(u)

    maxLog = max(values).bit_length()
    res = 0
    bitCounter = [0] * (maxLog + 1)

    def dfs(cur: int, pre: int, xor: int) -> int:
        xor ^= values[cur]
        for i in range(maxLog + 1):
            bitCounter[i] += (xor >> i) & 1
        for next_ in adjList[cur]:
            if next_ != pre:
                dfs(next_, cur, xor)

    dfs(0, -1, 0)
    res = 0
    for i, c in enumerate(bitCounter):
        res += (1 << i) * c * (n - c)
    return res
