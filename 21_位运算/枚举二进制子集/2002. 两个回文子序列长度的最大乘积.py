# 请你找到 s 中两个 不相交回文子序列 ，使得它们长度的 乘积最大 。
# 两个子序列在原字符串中如果没有任何相同下标的字符，则它们是 不相交 的。
# 2 <= s.length <= 12


# 直接枚举状态出所有的回文子序列(非空)即可
from typing import List, Tuple


class Solution:
    def maxProduct(self, s: str) -> int:
        n = len(s)

        store: List[Tuple[int, int]] = []
        for state in range(1, 1 << n):
            sb = []
            for i in range(n):
                if state & (1 << i):
                    sb.append(s[i])
            if sb == sb[::-1]:
                store.append((state, len(sb)))

        res = 0
        for i in range(len(store)):
            s1, l1 = store[i]
            for j in range(i + 1, len(store)):
                s2, l2 = store[j]
                if not s1 & s2:
                    res = max(res, l1 * l2)

        return res


print(Solution().maxProduct("leetcodecom"))
# 输出：9
# 解释：最优方案是选择 "ete" 作为第一个子序列，"cdc" 作为第二个子序列。
# 它们的乘积为 3 * 3 = 9 。
