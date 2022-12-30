# 1147. 段式回文

# 给你一个字符串 text，在确保它满足段式回文的前提下，
# 请你返回 段 的 最大数量 k
# 输入：text = "ghiabcdefhelloadamhelloabcdefghi"
# 输出：7
# 解释：我们可以把字符串拆分成 "(ghi)(abcdef)(hello)(adam)(hello)(abcdef)(ghi)"。

# 字符串哈希+贪心(头尾一起，能分就分) O(n)
# https://leetcode.cn/problems/longest-chunked-palindrome-decomposition/solution/zi-fu-chuan-ha-xi-tan-xin-tou-wei-yi-qi-79czc/


class Solution:
    def longestDecomposition(self, text: str) -> int:
        def cal(left: int, right: int) -> int:
            if left >= right:
                return 0
            if left + 1 == right:
                return 1
            i, j = left, right
            while i < j:
                i += 1
                j -= 1
                if H(left, i) == H(j, right):
                    return 2 + cal(i, j)
            return 1

        H = useStringHasher(text)
        return cal(0, len(text))


from typing import Sequence


def useStringHasher(s: Sequence[str], mod=10**11 + 7, base=1313131, offset=0):
    n = len(s)
    prePow = [1] * (n + 1)
    preHash = [0] * (n + 1)
    for i in range(1, n + 1):
        prePow[i] = (prePow[i - 1] * base) % mod
        preHash[i] = (preHash[i - 1] * base + ord(s[i - 1]) - offset) % mod

    def sliceHash(left: int, right: int):
        """切片 `s[left:right]` 的哈希值"""
        if left >= right:
            return 0
        left += 1
        return (preHash[right] - preHash[left - 1] * prePow[right - left + 1]) % mod

    return sliceHash


print(Solution().longestDecomposition("ghiabcdefhelloadamhelloabcdefghi"))
print(Solution().longestDecomposition("aaa"))
