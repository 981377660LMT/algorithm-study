# 100132. 统计美丽子字符串 II
# https://leetcode.cn/problems/count-beautiful-substrings-ii/description/

# 给你一个字符串 s 和一个正整数 k 。
# 用 vowels 和 consonants 分别表示字符串中元音字母和辅音字母的数量。
# 如果某个字符串满足以下条件，则称其为 美丽字符串 ：
# vowels == consonants，即元音字母和辅音字母的数量相等。
# (vowels * consonants) % k == 0，即元音字母和辅音字母的数量的乘积能被 k 整除。
# 返回字符串 s 中 非空美丽子字符串 的数量。

# 1. 预处理合法的元音长度，存入bitset中；
# 2. 遍历数组时记录前缀和以及对应的位置，把位置存入bitset中；
# 3. 查询时，这两个bitset交集的大小即为答案。


from collections import defaultdict


VOWELS = set(["a", "e", "i", "o", "u"])


class Solution:
    def beautifulSubstrings(self, s: str, k: int) -> int:
        """bitset保存前缀和对应位置."""
        n = len(s)
        nums = [1 if s in VOWELS else -1 for s in s]
        ok = 0  # 允许的长度集合
        for v in range(1, n // 2 + 1):
            if (v * v) % k == 0:
                ok |= 1 << (2 * v)

        res = 0
        preSum, curSum = defaultdict(int), 0
        preSum[0] = 1 << n
        for i, v in enumerate(nums, start=1):
            curSum += v
            prePos = preSum[curSum]
            intersection = ok & (prePos >> (n - i))
            res += intersection.bit_count()
            preSum[curSum] |= 1 << (n - i)

        return res
