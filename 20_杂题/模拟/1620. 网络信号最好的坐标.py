from typing import List
from math import floor

# towers[i] = [xi, yi, qi] 表示第 i 个网络信号塔的坐标是 (xi, yi) 且信号强度参数为 qi
# 整数 radius 表示一个塔 能到达 的 最远距离
# 如果第 i 个塔能到达 (x, y) ，那么该塔在此处的信号为 ⌊qi / (1 + d)⌋ ，
# 其中 d 是塔跟此坐标的距离。一个坐标的 网络信号 是所有 能到达 该坐标的塔的信号强度之和。

# 请你返回 网络信号 最大的整数坐标点。如果有多个坐标网络信号一样大，请你返回字典序最小的一个坐标。
# 0 <= xi, yi, qi <= 50
class Solution:
    def bestCoordinate(self, towers: List[List[int]], radius: int) -> List[int]:
        max_signal = 0
        res = [0, 0]
        ################ 遍历所有网格点
        for x in range(51):
            for y in range(51):
                cur = 0
                for [xi, yi, qi] in towers:
                    d = ((x - xi) ** 2 + (y - yi) ** 2) ** 0.5
                    ## 距离在范围内
                    if d <= radius:
                        cur += floor(qi / (1 + d))
                #### 更新res
                if cur > max_signal:
                    max_signal = cur
                    res = [x, y]
        return res


print(Solution().bestCoordinate(towers=[[1, 2, 5], [2, 1, 7], [3, 1, 9]], radius=2))
# 输出：[2,1]
# 解释：
# 坐标 (2, 1) 信号强度之和为 13
# - 塔 (2, 1) 强度参数为 7 ，在该点强度为 ⌊7 / (1 + sqrt(0)⌋ = ⌊7⌋ = 7
# - 塔 (1, 2) 强度参数为 5 ，在该点强度为 ⌊5 / (1 + sqrt(2)⌋ = ⌊2.07⌋ = 2
# - 塔 (3, 1) 强度参数为 9 ，在该点强度为 ⌊9 / (1 + sqrt(1)⌋ = ⌊4.5⌋ = 4
# 没有别的坐标有更大的信号强度。
