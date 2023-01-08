# hash(s, p, m) = (val(s[0]) * p0 + val(s[1]) * p1 + ... + val(s[k-1]) * pk-1) mod m.
# 其中 val(s[i]) 表示 s[i] 在字母表中的下标，从 val('a') = 1 到 val('z') = 26 。
# !请你返回 s 中 第一个 长度为 k 的 子串 sub ，满足 hash(sub, power, modulo) == hashValue 。
# !注意题目中的hash函数，是从左小右大，而一般的hash函数是从左大(p**(len-1))右小(p**0)
# 因为 pow(power, -1, modulo) 不一定存在(power与mod不一定互质)，所以不能正向滑窗


class Solution:
    def subStrHash(self, s: str, power: int, modulo: int, k: int, hashValue: int) -> str:
        s = s[::-1]
        n = len(s)
        res, curHash = -1, 0
        for right in range(n):
            curHash = (curHash * power + ord(s[right]) - 96) % modulo
            if right >= k:
                curHash -= (ord(s[right - k]) - 96) * pow(power, k, modulo)
                curHash %= modulo
            if right >= k - 1:
                if curHash == hashValue:
                    res = right
        return s[res - k + 1 : res + 1][::-1] if res != -1 else ""


print(Solution().subStrHash(s="leetcode", power=7, modulo=20, k=2, hashValue=0))
print(Solution().subStrHash(s="fbxzaad", power=31, modulo=100, k=3, hashValue=32))
