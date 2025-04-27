# https://leetcode.cn/problems/maximum-product-of-word-lengths/solutions/1104441/zui-da-dan-ci-chang-du-cheng-ji-by-leetc-lym9/
# 给你一个字符串数组 words ，找出并返回 length(words[i]) * length(words[j]) 的最大值，并且这两个单词不含有公共字母。
# 如果不存在这样的两个单词，返回 0 。
#
# !排序剪枝
# !set 减少重复字母的位运算

from collections import defaultdict
from typing import List


class Solution:
    def maxProduct(self, words: List[str]) -> int:
        mask2len = defaultdict(int)
        for w in words:
            m = 0
            for ch in set(w):
                m |= 1 << (ord(ch) - ord("a"))
            mask2len[m] = max(mask2len[m], len(w))

        res = 0
        items = sorted(mask2len.items(), key=lambda x: -x[1])
        for i, (mi, li) in enumerate(items):
            if li * li <= res:
                break
            for mj, lj in items[i + 1 :]:
                if li * lj <= res:
                    break
                if mi & mj == 0:
                    res = li * lj
                    break

        return res


if __name__ == "__main__":
    sol = Solution()
    print(sol.maxProduct(["abcw", "baz", "foo", "bar", "xtfn", "abcdef"]))  # 16
    print(sol.maxProduct(["a", "ab", "abc", "d", "cd", "bcd", "abcd"]))  # 4
    print(sol.maxProduct(["a", "aa", "aaa", "aaaa"]))  # 0
