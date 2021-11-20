from typing import List
from collections import Counter
from math import gcd


class Solution:
    def hasGroupsSizeX(self, deck: List[int]) -> bool:
        return gcd(*Counter(deck).values()) >= 2


print(Solution().hasGroupsSizeX([1, 1, 2, 2, 2, 2]))
# 输出：true

# 解释：可行的分组是 [1,1]，[2,2]，[2,2]
# 你需要选定一个数字 X，使我们可以将整副牌按下述规则分成 1 组或更多组：

# 每组都有 X 张牌。
# 组内所有的牌上都写着相同的整数。
# 仅当你可选的 X >= 2 时返回 true。

