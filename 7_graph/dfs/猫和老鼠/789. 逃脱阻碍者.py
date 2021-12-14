from typing import List

# 每一回合，你和阻碍者们可以同时向东，西，南，北四个方向移动，
# 每次可以移动到距离原位置 1 个单位 的新位置。
# 当然，也可以选择 不动 。所有动作 同时 发生。


# 如果你可以在任何阻碍者抓住你 之前 到达目的地（阻碍者可以采取任意行动方式），则被视为逃脱成功。
# 如果你和阻碍者同时到达了一个位置（包括目的地）都不算是逃脱成功。

# 你从 [0, 0] 点开始出发


# 其实问的就是有没鬼能比你更快到达终点
class Solution:
    def escapeGhosts(self, ghosts: List[List[int]], target: List[int]) -> bool:
        # 启发式: 我们到达终点需要的步数至少是 abs(targetX - 0) + abs(targetY - 0)
        # 敌人可以任意移动，那么存在任意敌人在这或者这之前能到达终点，我们必然就不能走到终点
        # 以上条件可以想象为我们在某些半路会被抓，那么敌人必然可以以同样时间到达终点
        # 不会有傻子觉得自己走更多步数（绕路）就可以避开敌人，因为那只是给了对方更多时间到达终点
        def manhattanDistance(p1, p2):
            return abs(p2[0] - p1[0]) + abs(p2[1] - p1[1])

        m = manhattanDistance((0, 0), target)
        return all(manhattanDistance(g, target) > m for g in ghosts)


print(Solution().escapeGhosts(ghosts=[[1, 0], [0, 3]], target=[0, 1]))
# 输出：true
# 解释：你可以直接一步到达目的地 (0,1) ，在 (1, 0) 或者 (0, 3) 位置的阻碍者都不可能抓住你。

