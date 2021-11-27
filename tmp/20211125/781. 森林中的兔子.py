from typing import List
from collections import Counter
from math import ceil

# 其中一些兔子（可能是全部）告诉你还有多少其他的兔子和自己有相同的颜色。我们将这些回答放在 answers 数组里。
# 返回森林中兔子的最少数量。
# 当某个兔子回答 x 的时候，那么数组中最多允许 x+1 个同花色的兔子🐰同时回答 x。
# 否则他们花色不全相同


class Solution:
    def numRabbits(self, answers: List[int]) -> int:
        counter = Counter(answers)
        res = 0
        for same_color, freq in counter.items():
            if freq > same_color + 1:
                res += (same_color + 1) * ceil(freq / (same_color + 1))
            else:
                res += same_color + 1
        return res


print(Solution().numRabbits([1, 1, 2]))
# 两只回答了 "1" 的兔子可能有相同的颜色，设为红色。
# 之后回答了 "2" 的兔子不会是红色，否则他们的回答会相互矛盾。
# 设回答了 "2" 的兔子为蓝色。
# 此外，森林中还应有另外 2 只蓝色兔子的回答没有包含在数组中。
# 因此森林中兔子的最少数量是 5: 3 只回答的和 2 只没有回答的。


# 比如有一个红色的兔子回答了 2，那么数组中最多有 3 个红色的兔子。
# 如果数组是 [2,2,2] ，那么至少有一种颜色的兔子。
# 如果数组是 [2,2,2,2] ，那么说明至少有两种颜色的兔子，比如说前 3 个兔子构成一种颜色；那么最后一个兔子说的必须是其他颜色。
# 如果数组是 [2,2,2,2,2,2] ，那么说明至少有两种颜色的兔子，比如说前 3 个兔子构成一种颜色；那么后 3 个兔子说的必须是其他颜色。

