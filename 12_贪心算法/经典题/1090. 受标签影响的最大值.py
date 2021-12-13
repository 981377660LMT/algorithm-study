from typing import List
from collections import Counter

# 我们从这些项中选出一个子集 S，这样一来：

# |S| <= num_wanted
# 对于任意的标签 L，子集 S 中标签为 L 的项的数目总满足 <= use_limit。
# 返回子集 S 的最大可能的 和。


class Solution:
    def largestValsFromLabels(
        self, values: List[int], labels: List[int], numWanted: int, useLimit: int
    ) -> int:
        counter = Counter()
        res = []
        for val, lab in sorted(zip(values, labels), reverse=True):
            if counter[lab] < useLimit:
                counter[lab] += 1
                res.append(val)
            if len(res) == numWanted:
                break
        return sum(res)


print(
    Solution().largestValsFromLabels(
        values=[5, 4, 3, 2, 1], labels=[1, 1, 2, 2, 3], numWanted=3, useLimit=1
    )
)
# 输出：9
# 解释：选出的子集是第一项，第三项和第五项。
