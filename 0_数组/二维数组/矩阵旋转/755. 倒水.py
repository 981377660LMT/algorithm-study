from typing import List

# 有点像1861. 旋转盒子.py
# 液滴优先向左移动到低位
#  一滴一滴地找，先左后右，找不到就原地+1
class Solution:
    def pourWater(self, heights: List[int], volume: int, k: int) -> List[int]:
        n = len(heights)
        for _ in range(volume):
            for d in (-1, 1):
                fillPos = k
                index = k
                while 0 <= index + d < n and heights[index + d] <= heights[index]:
                    if heights[index + d] < heights[index]:
                        fillPos = index + d
                    index += d
                if fillPos != k:
                    heights[fillPos] += 1
                    break
            else:
                heights[k] += 1

        return heights


# 在 V 个单位的水落在索引 K 处以后，每个索引位置有多少水？
print(Solution().pourWater(heights=[2, 1, 1, 2, 1, 2, 2], volume=4, k=3))
# 输出：[2,2,2,3,2,2,2]
