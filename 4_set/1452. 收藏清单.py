from typing import List

# 请找出不是其他任何人收藏的公司清单的子集的收藏清单，并返回该清单下标。下标需要按升序排列。
class Solution:
    def peopleIndexes(self, favoriteCompanies: List[List[str]]) -> List[int]:
        return [
            i
            for i, cur in enumerate(favoriteCompanies)
            if not any(set(other) > set(cur) for other in favoriteCompanies)
        ]


print(
    Solution().peopleIndexes(
        favoriteCompanies=[
            ["leetcode", "google", "facebook"],
            ["google", "microsoft"],
            ["google", "facebook"],
            ["google"],
            ["amazon"],
        ]
    )
)
# 输出：[0,1,4]
