from itertools import accumulate
from operator import xor
from typing import List
from collections import Counter

# 待检子串都可以表示为 queries[i] = [left, right, k]。
# 我们可以 `重新排列` 子串 s[left], ..., s[right]，并从中选择 最多 k 项替换成任何小写英文字母。
# 子串可以变成回文形式的字符串，那么检测结果为 true，
# 1 <= s.length, queries.length <= 10^5

# 1371. 每个元音包含偶数次的最长子字符串

# 规则：偶数个可以忽略，奇数要//2
# 'aaaac' => 'aacaa'
# 'abcd' => 'abba' (4 // 2 = 2)
# 'abcde' => 'abcba' (5 // 2 = 2)


class Solution:
    # 超时，因为我们重复统计Counter(str),没做预处理
    def canMakePaliQueries2(self, s: str, queries: List[List[int]]) -> List[bool]:
        res = []
        for l, r, k in queries:
            str = s[l : r + 1]
            counter = Counter(str)
            oddCount = 0
            for _, count in counter.items():
                oddCount += count & 1
            needReplace = oddCount >> 1
            res.append(needReplace <= k)
        return res

    # prefix预处理前缀 一个
    def canMakePaliQueries(self, s: str, queries: List[List[int]]) -> List[bool]:
        states = [1 << (ord(char) - ord('a')) for char in s]
        preXor = list(accumulate(states, initial=0, func=xor))  # 因为只关心奇偶，所以可以用xor
        return [((preXor[l] ^ preXor[r + 1]).bit_count() // 2) <= k for l, r, k in queries]


print(
    Solution().canMakePaliQueries(
        s="abcda", queries=[[3, 3, 0], [1, 2, 0], [0, 3, 1], [0, 3, 2], [0, 4, 1]]
    )
)
# 输出：[true,false,false,true,true]
# 解释：
# queries[0] : 子串 = "d"，回文。
# queries[1] : 子串 = "bc"，不是回文。
# queries[2] : 子串 = "abcd"，只替换 1 个字符是变不成回文串的。
# queries[3] : 子串 = "abcd"，可以变成回文的 "abba"。 也可以变成 "baab"，先重新排序变成 "bacd"，然后把 "cd" 替换为 "ab"。
# queries[4] : 子串 = "abcda"，可以变成回文的 "abcba"。

