from typing import List, Optional, Tuple

MOD = int(1e9 + 7)
INF = int(1e20)


def getDigit(num: int, /, *, index: int, radix: int):
    """返回 `radix` 进制下 `num` 的 `index` 位的数字，`index` 最低位(最右)为 0 """
    assert radix >= 2 and index >= 0
    prefix = num // pow(radix, index)
    return prefix % radix


class Solution:
    def distributeCookies(self, cookies: List[int], k: int) -> int:
        """k进制枚举TLE了 这样是8**8的 肯定要剪枝"""
        n = len(cookies)
        res = int(1e20)
        for state in range(k ** n):
            groups = [0] * k
            for i in range(n):
                prefix = state // (k ** i)
                mod = prefix % k
                groups[mod] += cookies[i]
            max_ = max(groups)
            if max_ < res:
                res = max_
        return res


print(Solution().distributeCookies(cookies=[8, 15, 10, 20, 8], k=2))

