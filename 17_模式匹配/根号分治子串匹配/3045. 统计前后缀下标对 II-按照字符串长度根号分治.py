# 3045. 统计前后缀下标对 II
# https://leetcode.cn/problems/count-prefix-and-suffix-pairs-ii/solutions/2644232/gen-hao-bao-li-by-hqztrue-2d0s/
# 当 str1 同时是 str2 的前缀和后缀时，isPrefixAndSuffix(str1, str2) 返回 true，否则返回 false。
# 以整数形式，返回满足 i < j 且 isPrefixAndSuffix(words[i], words[j]) 为 true 的下标对 (i, j) 的 数量 。

# 很短的根号暴力
# !暴力，只对存在相同长度字符串的前后缀进行检查。
# !令 L 表示字符串的总长度，复杂度 O(LsqrtL).

# 具体思路:
# !1.遍历每个单词，枚举答案的长度 i，如果当前单词的前缀和后缀相同，那么答案就加上之前出现的相同前缀的次数。
# 2. 记单词长度为L，单词总长度为∑，比较的复杂度最多为 `min(L*L,∑)`.(要么就是比较新字符串的所有前缀，也就是1+2+…+t=O(t^2), 要么就是和之前的字符串全部比一遍，也就是L。)
#    对于长度<=sqrt(∑)的字符串，总比较次数之和<=∑L*L<=∑L*sqrt(∑)=O(∑sqrt(∑)).
#    对于长度>sqrt(∑)的字符串，最多只有sqrt(∑)个，总比较次数之和<=∑L*sqrt(∑)=O(∑sqrt(∑)).


from typing import List
from collections import defaultdict


class Solution:
    def countPrefixSuffixPairs(self, words: List[str]) -> int:
        res = 0
        counter = defaultdict(int)
        visitedLen = set()  # !最多 sqrtL 种长度
        for w in words:
            for len_ in range(1, len(w) + 1):
                if len_ in visitedLen and w[:len_] == w[-len_:]:
                    res += counter[w[:len_]]
            counter[w] += 1
            visitedLen.add(len(w))
        return res
