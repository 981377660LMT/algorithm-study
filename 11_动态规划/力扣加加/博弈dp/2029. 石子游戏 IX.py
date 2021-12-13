from typing import List

# Alice 和 Bob 轮流进行自己的回合，Alice 先手。每一回合，玩家需要从 stones 中移除任一石子。
# 如果玩家移除石子后，导致 所有已移除石子 的价值 总和 可以被 3 整除，那么该玩家就 输掉游戏 。
# 如果不满足上一条，且移除后没有任何剩余的石子，那么 Bob 将会直接获胜

# 这题没什么意义
# https://leetcode-cn.com/problems/stone-game-ix/solution/pythonpan-duan-aliceneng-fou-huo-sheng-b-unyt/
class Solution:
    def stoneGameIX(self, stones: List[int]) -> bool:
        cnt = [0] * 3
        for v in stones:
            cnt[v % 3] += 1
        if cnt[0] % 2 == 0 and cnt[1] * cnt[2] > 0:
            return True
        if cnt[0] % 2 == 1 and abs(cnt[1] - cnt[2]) >= 3:
            return True
        return False


print(Solution().stoneGameIX([2, 1]))
# 输出：true
# 解释：游戏进行如下：
# - 回合 1：Alice 可以移除任意一个石子。
# - 回合 2：Bob 移除剩下的石子。
# 已移除的石子的值总和为 1 + 2 = 3 且可以被 3 整除。因此，Bob 输，Alice 获胜。
