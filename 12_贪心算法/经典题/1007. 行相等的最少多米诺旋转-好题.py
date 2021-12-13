from typing import List

# A[i] 和 B[i] 分别代表第 i 个多米诺骨牌的上半部分和下半部分。
# （一个多米诺是两个从 1 到 6 的数字同列平铺形成的 —— 该平铺的每一半上都有一个数字。）


# 我们可以旋转第 i 张多米诺，使得 A[i] 和 B[i] 的值交换。
# 返回能使 A 中所有值或者 B 中所有值都相同的最小旋转次数。
# 如果无法做到，返回 -1.

# 2 <= A.length == B.length <= 20000

# 暴力
class Solution:
    def minDominoRotations(self, tops: List[int], bottoms: List[int]) -> int:
        for sameValue in (tops[0], bottoms[0]):
            if all(sameValue in pair for pair in zip(tops, bottoms)):
                return len(tops) - max(tops.count(sameValue), bottoms.count(sameValue))
        return -1


print(Solution().minDominoRotations(tops=[2, 1, 2, 4, 2, 2], bottoms=[5, 2, 6, 2, 3, 2]))
# 输出：2
# 解释：
# 图一表示：在我们旋转之前， A 和 B 给出的多米诺牌。
# 如果我们旋转第二个和第四个多米诺骨牌，我们可以使上面一行中的每个值都等于 2，如图二所示。
