# 请你返回 好字符串 的数目。
# 好字符串 的定义为：
# 它的长度为 n ，字典序大于等于 s1 ，字典序小于等于 s2 ，且不包含 evil 为子字符串。
# 1 <= n <= 500
# 1 <= evil.length <= 50
# 所有字符串都只包含小写英文字母。
# 链接：https://leetcode-cn.com/problems/find-all-good-strings


from functools import lru_cache
from typing import List, Optional
from typing import Generic, List, Optional, Sequence, TypeVar

MOD = int(1e9 + 7)

T = TypeVar("T", int, str)


class _KMP(Generic[T]):
    """单模式串匹配"""

    @staticmethod
    def getNext(pattern: Sequence[T]) -> List[int]:
        next = [0] * len(pattern)
        j = 0
        for i in range(1, len(pattern)):
            while j and pattern[i] != pattern[j]:
                j = next[j - 1]
            if pattern[i] == pattern[j]:
                j += 1
            next[i] = j
        return next

    __slots__ = ("next", "_pattern")

    def __init__(self, pattern: Sequence[T]):
        self._pattern = pattern
        self.next = self.getNext(pattern)

    def match(self, s: Sequence[T], start=0) -> List[int]:
        res = []
        pos = 0
        for i in range(start, len(s)):
            pos = self.move(pos, s[i])
            if self.isMatched(pos):
                res.append(i - len(self._pattern) + 1)
                pos = self.next[pos - 1]  # rollback
        return res

    def move(self, pos: int, char: T) -> int:
        assert 0 <= pos < len(self._pattern)
        while pos and char != self._pattern[pos]:
            pos = self.next[pos - 1]  # rollback
        if char == self._pattern[pos]:
            pos += 1
        return pos

    def isMatched(self, pos: int) -> bool:
        return pos == len(self._pattern)

    def period(self, i: Optional[int] = None) -> int:
        """
        求字符串 S 的前缀 s[:i+1] 的最短周期(0<=i<n)
        如果不存在周期, 返回0.
        """
        if i is None:
            i = len(self._pattern) - 1
        assert 0 <= i < len(self._pattern)
        res = (i + 1) - self.next[i]
        if res and (i + 1) > res and (i + 1) % res == 0:
            return res
        return 0


def cal(upper: str, evil: str) -> int:
    """字典序小于等于upper且不含evil的字符串个数"""

    @lru_cache(None)
    def dfs(index: int, isLimit: bool, pos: int) -> int:
        """当前在第pos位,isLimit表示是否贴合上界,pos表示匹配了多少个evil字符(KMP的状态)"""
        if index == n:
            return 1
        res = 0
        up = upper[index] if isLimit else "z"
        for cur in range(97, ord(up) + 1):
            select = chr(cur)
            nextPos = kmp.move(pos, select)
            if kmp.isMatched(nextPos):
                continue
            res += dfs(index + 1, (isLimit and select == up), nextPos)
            res %= MOD
        return res

    n = len(upper)
    kmp = _KMP(evil)
    res = dfs(0, True, 0)
    dfs.cache_clear()
    return res


class Solution:
    def findGoodStrings(self, n: int, s1: str, s2: str, evil: str) -> int:
        return (cal(s2, evil) - cal(s1, evil) + int(evil not in s1)) % MOD


print(Solution().findGoodStrings(n=2, s1="gx", s2="gz", evil="x"))
print(Solution().findGoodStrings(n=8, s1="leetcode", s2="leetgoes", evil="leet"))
print(Solution().findGoodStrings(n=2, s1="aa", s2="da", evil="b"))
print(Solution().findGoodStrings(n=8, s1="pzdanyao", s2="wgpmtywi", evil="sdka"))
# 500543753

# 输出：51
# 解释：总共有 25 个以 'a' 开头的好字符串："aa"，"ac"，"ad"，...，"az"。还有 25 个以 'c' 开头的好字符串："ca"，"cc"，"cd"，...，"cz"。最后，还有一个以 'd' 开头的好字符串："da"。

# 来源：力扣（LeetCode）

# 著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
