from typing import List

# 你的目标是找到松鼠收集所有坚果的最小路程，
# 且坚果是一颗接一颗地被放在树下。
# 松鼠一次最多只能携带一颗坚果

# 先假设松鼠一开始在就在树的位置，最小路程就是曼哈顿距离之和。

# 然后要将第一步树到某坚果的路程替换为松鼠到某坚果的路程，找替换代价最小的即可。
class Solution:
    def minDistance(
        self, height: int, width: int, tree: List[int], squirrel: List[int], nuts: List[List[int]]
    ) -> int:
        def cal(p1, p2):
            return abs(p1[0] - p2[0]) + abs(p1[1] - p2[1])

        res = 2 * sum(cal(tree, p) for p in nuts)
        return res + min(cal(squirrel, p) - cal(tree, p) for p in nuts)
