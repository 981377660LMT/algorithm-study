# 请你返回 好字符串 的数目。
# 好字符串 的定义为：它的长度为 n ，字典序大于等于 s1 ，字典序小于等于 s2 ，且不包含 evil 为子字符串。
from functools import lru_cache
from typing import List


MOD = int(1e9 + 7)

# 1 <= n <= 500
# 1 <= evil.length <= 50
# 所有字符串都只包含小写英文字母。

# kmp求next数组
@lru_cache(None)
def getNext(needle: str) -> List[int]:
    """kmp O(n)求 `needle`串的 `next`数组

    `next[i]`表示`[:i+1]`这一段字符串中最长公共前后缀(不含这一段字符串本身,即真前后缀)的长度
    https://www.ruanyifeng.com/blog/2013/05/Knuth%E2%80%93Morris%E2%80%93Pratt_algorithm.html
    """
    next = [0] * len(needle)
    j = 0

    for i in range(1, len(needle)):
        while j and needle[i] != needle[j]:  # 1. fallback后前进：匹配不成功j往右走
            j = next[j - 1]

        if needle[i] == needle[j]:  # 2. 匹配：匹配成功j往右走一步
            j += 1

        next[i] = j

    return next


# 不熟悉如何利用kmp的next数组优化字符串匹配
@lru_cache(None)
def cal(upper: str, evil: str) -> int:
    """字典序小于等于upper且不含evil的字符串个数"""

    @lru_cache(None)
    def dfs(index: int, isLimit: bool, hit: int) -> int:
        """当前在第pos位,isLimit表示是否贴合上界,hit表示匹配了多少个evil字符"""
        if hit == m:
            return 0
        if index == n:
            return 1
        res = 0
        up = upper[index] if isLimit else 'z'
        for cur in range(97, ord(up) + 1):
            select = chr(cur)
            nextHit = hit

            while nextHit > 0 and select != evil[nextHit]:
                nextHit = evilNext[nextHit - 1]
            if select == evil[nextHit]:
                nextHit += 1

            res += dfs(index + 1, (isLimit and select == up), nextHit)
            res %= MOD
        return res

    n, m = len(upper), len(evil)
    evilNext = getNext(evil)
    return dfs(0, True, 0)


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
# 链接：https://leetcode-cn.com/problems/find-all-good-strings
# 著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
