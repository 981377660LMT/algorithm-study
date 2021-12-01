from typing import List

# 3 <= stones.length <= 10^4
# stones[i] 的值各不相同。


# 你可以将一颗端点石子拿起并移动到一个未占用的位置，使得该石子`不再是一颗端点石子`。
# 当你无法进行任何移动时，即，这些石子的位置连续时，游戏结束。
# 要使游戏结束，你可以执行的最小和最大移动次数分别是多少？
class Solution:
    def numMovesStonesII(self, stones: List[int]) -> List[int]:
        ...


print(Solution().numMovesStonesII([6, 5, 4, 3, 10]))
# 输出：[2,3]
# 解释：
# 我们可以移动 3 -> 8，接着是 10 -> 7，游戏结束。
# 或者，我们可以移动 3 -> 7, 4 -> 8, 5 -> 9，游戏结束。
# 注意，我们无法进行 10 -> 2 这样的移动来结束游戏，因为这是不合要求的移动。


# 直接放弃此题 没什么意义
