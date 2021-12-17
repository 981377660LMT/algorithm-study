from typing import List


class Solution:
    def wateringPlants(self, plants: List[int], capacity: int) -> int:
        ...


print(Solution().wateringPlants(plants=[2, 2, 3, 3], capacity=5))
# 输出：14
# 解释：从河边开始，此时水罐是装满的：
# - 走到植物 0 (1 步) ，浇水。水罐中还有 3 单位的水。
# - 走到植物 1 (1 步) ，浇水。水罐中还有 1 单位的水。
# - 由于不能完全浇灌植物 2 ，回到河边取水 (2 步)。
# - 走到植物 2 (3 步) ，浇水。水罐中还有 2 单位的水。
# - 由于不能完全浇灌植物 3 ，回到河边取水 (3 步)。
# - 走到植物 3 (4 步) ，浇水。
# 需要的步数是 = 1 + 1 + 2 + 3 + 3 + 4 = 14 。
