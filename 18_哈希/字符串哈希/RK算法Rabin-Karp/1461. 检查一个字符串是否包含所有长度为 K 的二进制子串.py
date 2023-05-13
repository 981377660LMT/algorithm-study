# 1461. 检查一个字符串是否包含所有长度为 K 的二进制子串
class Solution:
    def hasAllCodes(self, s: str, k: int) -> bool:
        return len(set(s[i : i + k] for i in range(len(s) - k + 1))) == 2**k


# 1016. 子串能表示从 1 到 N 数字的二进制串
# 由于S的长度只有1000，能表达的有效的32位整数只有3万多个，而不重复的更少，所以n的范围是唬人的，因此考虑遍及即可
# We return false as soon as we dont find a perticular no.
# hence the probability of returning false is very high before the time limit exceeds
# 即：因为 s 的长度最大为1000，所以实际上 n >= 1000 就可以直接返回 false 了，还没遍历完就直接返回false了
class Solution2:
    def queryString(self, s: str, n: int) -> bool:
        return all(bin(i)[2:] in s for i in range(n, n // 2, -1))
