# 1 <= s.length <= 105
from collections import Counter


class Solution:
    def minimumKeypresses(self, s: str) -> int:
        """最少的按键次数分配键盘
        
        贪心 出现最多的字符要用最少的按键
        """
        return sum((i // 9 + 1) * count for i, (_, count) in enumerate(Counter(s).most_common()))


print(Solution().minimumKeypresses("abcdefghijkl"))
