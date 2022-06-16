# https://leetcode-cn.com/problems/count-ways-to-build-rooms-in-an-ant-colony/solution/python-shu-zhuang-shu-zu-qiu-tuo-bu-fang-5e8f/
# https://leetcode.com/problems/count-ways-to-build-rooms-in-an-ant-colony/discuss/1299540/PythonC%2B%2B-clean-DFS-solution-with-explanation
from typing import List, Tuple
from math import comb
from collections import defaultdict

# 2 <= n <= 105
# 在完成所有房间的构筑之后，从房间 0 可以访问到每个房间。
# prevRoom[i] 表示在构筑房间 i 之前，你必须先构筑房间 prevRoom[i]
# 每个房间只能有一个 prevRoom

MOD = 10 ** 9 + 7


class Solution:
    def waysToBuildRooms(self, prevRoom: List[int]) -> int:
        adjmap = defaultdict(list)
        for i, num in enumerate(prevRoom):
            adjmap[num].append(i)

        # 返回:元素个数,排序方案数
        # 计算组合两个数组并保持其原始顺序的方法的数量
        # 假设这两个数组的长度分别是 l 和 r，那么答案是 math.com b (l + r，l)
        def dfs(cur: int) -> Tuple[int, int]:
            nodeCount, res = 0, 1
            for next in adjmap[cur]:
                subCount, nextRes = dfs(next)
                nodeCount += subCount

                # 子树1排序数*子树2排序数*组内保持顺序合并数组的方式
                res = (res * nextRes * comb(nodeCount, subCount)) % MOD
            return (nodeCount + 1, res)

        return dfs(0)[1]


print(Solution().waysToBuildRooms([-1, 0, 0, 1, 2]))

# 解释：
# 有 6 种不同顺序：
# 0 → 1 → 3 → 2 → 4
# 0 → 2 → 4 → 1 → 3
# 0 → 1 → 2 → 3 → 4
# 0 → 1 → 2 → 4 → 3
# 0 → 2 → 1 → 3 → 4
# 0 → 2 → 1 → 4 → 3

