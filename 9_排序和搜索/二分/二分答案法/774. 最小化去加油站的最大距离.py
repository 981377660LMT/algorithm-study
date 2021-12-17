from typing import List

# 在数轴上增设 k 个加油站，新增加油站可以位于 水平数轴 上的任意位置，而不必放在整数位置上。
# 请你最小化 增设 k 个新加油站后，相邻 两个加油站间的最大距离。
# 10 <= stations.length <= 2000
# 0 <= stations[i] <= 108
# stations 按 严格递增 顺序排列

# 二分
class Solution:
    def minmaxGasDist(self, stations: List[int], k: int) -> float:
        eps = 1e-8

        n = len(stations)
        dist = []
        for i in range(n - 1):
            dist.append(stations[i + 1] - stations[i])

        def check(mid: float) -> bool:
            """"两个加油站间最大距离为mid时需要的新增站数<=k"""

            need = 0
            for d in dist:
                need += int(d / mid)
            return need <= k

        # 与实际答案误差在 10-6 范围内的答案将被视作正确答案。
        l = 0
        r = 10 ** 8
        while l <= r:
            mid = (l + r) / 2
            if check(mid):
                r = mid - eps
            else:
                l = mid + eps
        return l


print(Solution().minmaxGasDist(stations=[23, 24, 36, 39, 46, 56, 57, 65, 84, 98], k=1))
# 输出：14.00000
