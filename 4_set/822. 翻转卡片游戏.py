from typing import List

# 哪个数是这些想要的数字中最小的数（找到这些数中的最小值）呢？如果没有一个数字符合要求的，输出 0。
# 如果选中的那张卡片`背面的数字 X 与任意一张卡片的正面的数字都不同`，那么这个数字是我们想要的数字。

# 首先
# 1. 把正反相同的卡的`值`排除掉
# 2. 如果背面最小为2 但是另一张牌正面为2 那我们就把那张牌翻成背面 背面肯定不为2(正反相等的卡的值已经排除了)
class Solution:
    def flipgame(self, fronts: List[int], backs: List[int]) -> int:
        inValidRes = set(f for f, b in zip(fronts, backs) if f == b)
        return min([val for val in fronts + backs if val not in inValidRes], default=0)


print(Solution().flipgame(fronts=[1, 2, 4, 4, 7], backs=[1, 3, 4, 1, 3]))
# 输出：2
# 解释：假设我们翻转第二张卡片，那么在正面的数变成了 [1,3,4,4,7] ， 背面的数变成了 [1,2,4,1,3]。
# 接着我们选择第二张卡片，因为现在该卡片的背面的数是 2，2 与任意卡片上正面的数都不同，所以 2 就是我们想要的数字。
