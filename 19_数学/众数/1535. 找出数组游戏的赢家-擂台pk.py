from typing import List

# 比较 arr[0] 与 arr[1] 的大小
# 较大的整数将会取得这一回合的胜利并保留在位置 0 ，
# 较小的整数移至数组的末尾。当一个整数赢得 k 个连续回合时，
# 游戏结束，该整数就是比赛的 赢家 。
# 返回赢得比赛的整数。
# 题目数据 保证 游戏存在赢家。

# 移到末尾是不必须的 模拟题


class Solution:
    def getWinner(self, arr: List[int], k: int) -> int:
        cur = arr[0]
        win = 0
        for enemy in arr[1:]:
            if enemy > cur:
                cur = enemy
                win = 0
            win += 1
            if win == k:
                break

        return cur


print(Solution().getWinner(arr=[2, 1, 3, 5, 4, 6, 7], k=2))
