# 100132. 统计美丽子字符串 II
# https://leetcode.cn/problems/count-beautiful-substrings-ii/description/


# 给你一个字符串 s 和一个正整数 k 。
# 用 vowels 和 consonants 分别表示字符串中元音字母和辅音字母的数量。
# 如果某个字符串满足以下条件，则称其为 美丽字符串 ：
# vowels == consonants，即元音字母和辅音字母的数量相等。
# (vowels * consonants) % k == 0，即元音字母和辅音字母的数量的乘积能被 k 整除。
# 返回字符串 s 中 非空美丽子字符串 的数量。


# !注意到 (x/2)^2 % k == 0，即 x^2 % (4k) == 0 .
# 一个数的x平方被n整除意味着什么?
# 对n进行质因子分解, n = p1^a1 * p2^a2 * ... * pk^ak
# !则对每个 pi, x中必须包含至少 `ceil(ai/2)` 个pi，乘积记为mul.
# !充要条件为：x是mul的倍数.
# !现在问题变成，有多少个和为 0 的子数组，其长度是 mul 的倍数？
# 剩下可以 O(n)解决.

from collections import Counter, defaultdict


VOWELS = set(["a", "e", "i", "o", "u"])


class Solution:
    def beautifulSubstrings(self, s: str, k: int) -> int:
        primes = getPrimeFactors(4 * k)
        mul = 1
        for p, a in primes.items():
            mul *= pow(p, (a + 1) // 2)

        nums = [1 if s in VOWELS else -1 for s in s]

        preSum = defaultdict(int)
        preSum[(mul - 1, 0)] = 1  # (lengthMod,sum) -> count
        res, curSum = 0, 0
        for i, v in enumerate(nums):
            curSum += v
            pair = (i % mul, curSum)
            res += preSum[pair]
            preSum[pair] += 1

        return res


def getPrimeFactors(n: int) -> "Counter[int]":
    """n 的素因子分解 O(sqrt(n))"""
    res = Counter()
    upper = int(n**0.5) + 1
    for i in range(2, upper):
        while n % i == 0:
            res[i] += 1
            n //= i
    if n > 1:
        res[n] += 1
    return res


if __name__ == "__main__":
    print(Solution().beautifulSubstrings("baeyh", 2))
