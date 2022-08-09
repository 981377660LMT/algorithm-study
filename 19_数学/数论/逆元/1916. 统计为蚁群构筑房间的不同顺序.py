# https://leetcode-cn.com/problems/count-ways-to-build-rooms-in-an-ant-colony/solution/python-shu-zhuang-shu-zu-qiu-tuo-bu-fang-5e8f/
# https://leetcode.com/problems/count-ways-to-build-rooms-in-an-ant-colony/discuss/1299540/PythonC%2B%2B-clean-DFS-solution-with-explanation
from typing import List, Tuple


# 2 <= n <= 105
# 在完成所有房间的构筑之后，从房间 0 可以访问到每个房间。
# !prevRoom[i] 表示在构筑房间 i 之前，你必须先构筑房间 prevRoom[i] (即连接过程中边不能断开)
# 每个房间只能有一个 prevRoom

MOD = int(1e9 + 7)


fac = [1]
ifac = [1]
for i in range(1, int(1e5 + 10)):
    fac.append((fac[-1] * i) % MOD)
    ifac.append((ifac[-1] * pow(i, MOD - 2, MOD)) % MOD)


def C(n: int, k: int) -> int:
    if n < 0 or k < 0 or n < k:
        return 0
    return ((fac[n] * ifac[k]) % MOD * ifac[n - k]) % MOD


# 蚂蚁
class Solution:
    def waysToBuildRooms(self, prevRoom: List[int]) -> int:
        n = len(prevRoom)
        adjList = [[] for _ in range(n)]
        for cur, pre in enumerate(prevRoom):
            if pre == -1:
                continue
            adjList[cur].append(pre)
            adjList[pre].append(cur)

        # !返回:元素个数,排序方案数
        # 计算组合两个数组并保持其原始顺序的方法的数量
        # 假设这两个数组的长度分别是 l 和 r，那么答案是 C(l + r，l)
        def dfs(cur: int, pre: int) -> Tuple[int, int]:
            nodeCount, res = 0, 1
            for next in adjList[cur]:
                if next == pre:
                    continue
                subCount, nextRes = dfs(next, cur)
                nodeCount += subCount

                # 子树1排序数*子树2排序数*组内保持顺序合并数组的方式
                res = (res * nextRes * C(nodeCount, subCount)) % MOD
            return (nodeCount + 1, res)

        return dfs(0, -1)[1]


print(Solution().waysToBuildRooms([-1, 0, 0, 1, 2]))

# 解释：
# 有 6 种不同顺序：
# 0 → 1 → 3 → 2 → 4
# 0 → 2 → 4 → 1 → 3
# 0 → 1 → 2 → 3 → 4
# 0 → 1 → 2 → 4 → 3
# 0 → 2 → 1 → 3 → 4
# 0 → 2 → 1 → 4 → 3
