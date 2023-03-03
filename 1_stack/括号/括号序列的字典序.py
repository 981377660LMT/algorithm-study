# https://oi-wiki.org/topic/bracket/
# https://www.cnblogs.com/suxxsfe/p/15986908.html
# 字典序后继
# 字典序计算


from typing import Tuple


# 给出合法的括号序列 s，
# 我们要求出按字典序升序排序的长度为 |s| 的所有合法括号序列中，
# 序列 s 的下一个合法括号序列
# !找到一个最大的i满足s[i]是左括号且[0,i-1]内左括号数量严格大于右括号数量
def nextBalancedSequence(s: str) -> Tuple[str, bool]:
    """https://oi-wiki.org/topic/bracket/"""
    n = len(s)
    depth = 0
    for i in range(n - 1, -1, -1):
        if s[i] == "(":
            depth -= 1
        else:
            depth += 1
        if s[i] == "(" and depth > 0:  # 找到最后一个左括号将其变为右括号
            depth -= 1
            open = (n - i - 1 - depth) // 2
            close = n - i - 1 - open
            next = s[:i] + ")" + "(" * open + ")" * close
            return next, True
    return "", False


# MOD = int(1e9 + 7)
# N = 105
# dp = [[0] * N for _ in range(N)]
# for i in range(N):
#     dp[i][i] = 1
# for i in range(1, N):
#     for j in range(i):
#         dp[i][j] = dp[i - 1][j + 1] + dp[i - 1][j - 1]
#         dp[i][j] %= MOD
MOD = int(1e9 + 7)
fac = [1]
ifac = [1]
for i in range(1, int(1e3) + 10):
    fac.append((fac[-1] * i) % MOD)
    ifac.append((ifac[-1] * pow(i, MOD - 2, MOD)) % MOD)


def C(n: int, k: int) -> int:
    if n < 0 or k < 0 or n < k:
        return 0
    return ((fac[n] * ifac[k]) % MOD * ifac[n - k]) % MOD


def cal(i: int, j: int) -> int:
    if (i - j) & 1:
        return 0
    return (C(i + 1, (i - j) // 2) * (j + 1) // (i + 1)) % MOD


# 给出合法的括号序列 s，求出它的字典序排名(1-based)。
# !设 f(i,j) 表示长度为 i 且存在 j 个未匹配的右括号且不存在未匹配的左括号的括号序列的个数。
# f(0, 0) = 1, f(i,j) = f(i-1,j+1) + f(i-1,j-1)  https://oeis.org/A053121
# 考虑求出字典序比 s 小的括号序列 p 的个数。
# 遍历s,如果遇到右括号,就计算把这个位置变为左括号时的合法括号的个数,
# 假设[0,i]中左括号比右括号多k个,那么就要统计右边的 f(len(s)-(i+1),k)
def getRank(s: str) -> int:
    n = len(s)
    res = 0
    diff = 0
    for i, c in enumerate(s):
        if c == ")":  # change to "("
            res += cal(n - 1 - i, diff + 1)
            res %= MOD
            diff -= 1
        else:
            diff += 1
    return res + 1


# 求长度为n的字典序第k小的合法括号序列(1-based)
# 从前往后考虑，看这个位置能不能填右括号
def getKth(n: int, k: int) -> str:
    if n & 1 or k > catalan(n):
        return ""
    diff = 0
    res = [""] * n
    mid = n // 2
    for i in range(n):
        if diff + 1 <= mid and cal(n - 1 - i, diff + 1) >= k:
            res[i] = "("
            diff += 1
        else:
            res[i] = ")"
            if diff + 1 <= mid:
                k -= cal(n - 1 - i, diff + 1)
            diff -= 1
    return "".join(res)


def catalan(n: int) -> int:
    """长为n的合法括号序列的个数"""
    if n & 1:
        return 0
    n //= 2
    return C(2 * n, n) // (n + 1)


print(nextBalancedSequence("(())"))
print(getRank("()"))
print(getRank("()()"))
print(getRank("((()))"))
print(cal(8, 6))
print(getKth(4, 2))
