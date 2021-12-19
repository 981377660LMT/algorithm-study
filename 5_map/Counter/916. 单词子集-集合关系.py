from typing import List
from collections import Counter


class Solution:
    def wordSubsets(self, words1: List[str], words2: List[str]) -> List[str]:
        need = Counter()
        for word in words2:
            need |= Counter(word)
        # 子集关系：交小并大
        return [w for w in words1 if need & Counter(w) == need]


print(
    Solution().wordSubsets(
        words1=["amazon", "apple", "facebook", "google", "leetcode"], words2=["e", "oo"]
    )
)
