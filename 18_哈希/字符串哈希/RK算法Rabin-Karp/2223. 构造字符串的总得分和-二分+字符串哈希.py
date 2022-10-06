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


class Solution:
    def sumScores(self, s: str) -> int:
        def countPre(curLen: int, start: int) -> int:
            left, right = 1, curLen
            while left <= right:
                mid = (left + right) // 2
                if hasher(start, start + mid) == hasher(0, mid):
                    left = mid + 1
                else:
                    right = mid - 1
            return right

        n = len(s)
        hasher = useStringHasher(s)
        res = 0
        for i in range(1, n + 1):
            if s[-i] != s[0]:
                continue
            count = countPre(i, n - i)
            res += count
        return res


print(Solution().sumScores("babab"))
print(Solution().sumScores("bab"))
print(Solution().sumScores(s="azbazbzaz"))
