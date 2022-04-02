from typing import List

# Remove every point (a, b) if there's another point (c, d) such that a ≤ c and b ≤ d.
#  The relative ordering of the original list should be maintained.


class Solution:
    def numberOfWeakCharacters(self, coordinates: List[List[int]]) -> List[List[int]]:
        pairs = sorted(
            [(i, x, y) for i, (x, y) in enumerate(coordinates)], key=lambda x: (x[1], x[2])
        )

        notOk = set()
        preMax = -0x7FFFFFFF

        # 这里从后往前遍历,因为要获取最大的防御值
        for i, _, num in reversed(pairs):
            if preMax >= num:
                notOk.add(i)
            preMax = max(preMax, num)

        return [coordinates[i] for i in range(len(coordinates)) if i not in notOk]


print(Solution().numberOfWeakCharacters(coordinates=[[1, 0], [1, 1]]))
# 输出：1
# 解释：第三个角色是弱角色，因为第二个角色的攻击和防御严格大于该角色。
