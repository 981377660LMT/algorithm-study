from typing import List
from collections import defaultdict

# 2015. 每段建筑物的平均高度
# 将重叠区间内的所有大楼高度求平均值，返回这些区域可进行分段的大楼跨度及其高度平均值数组
# 且遇到楼层平均高度相同的相邻区间时要合并
# 1943. 描述绘画结果-扫描.py
# 1. 使用map记录每个分界处的 差分信息
# 2. 排序，遍历每个分界点，更新 pre/preSum/preCount


class Solution:
    def averageHeightOfBuildings(self, buildings: List[List[int]]) -> List[List[int]]:
        # 总高度，个数
        diff = defaultdict(lambda: [0, 0])
        for start, end, delta in buildings:
            diff[start][0] += delta
            diff[end][0] -= delta
            diff[start][1] += 1
            diff[end][1] -= 1

        res = []
        # 区间起点，高度累加，个数累加
        curPos, curSum, curCount = 0, 0, 0
        for key in sorted(diff):
            delta, deltaCount = diff[key]
            if curSum > 0:
                cand = [curPos, key, curSum // curCount]

                # 区间合并
                if res and res[-1][1] == curPos and res[-1][2] == cand[2]:
                    res[-1][1] = key
                else:
                    res.append(cand)

            curPos = key
            curSum += delta
            curCount += deltaCount

        return res


print(Solution().averageHeightOfBuildings(buildings=[[1, 4, 2], [3, 9, 4]]))
# Output: [[1,3,2],[3,4,3],[4,9,4]]
# Explanation:
# From 1 to 3, there is only the first building with an average height of 2 / 1 = 2.
# From 3 to 4, both the first and the second building are there with an average height of (2+4) / 2 = 3.
# From 4 to 9, there is only the second building with an average height of 4 / 1 = 4.
