from typing import List


# s = "abababababa"
# t = "aba"
# 一共有5个可重叠子串t


# s = "abababababa"
#     "aba"
#       "aba"
#         "aba"
#           "aba"
#             "aba"

# 可以直接字符串哈希O(n)
# 或者 kmp 线性复杂度，相当于普通 KMP 的最坏状态（需要遍历整个待匹配串）


def calSubstr(s: str, t: str) -> List[int]:
    """可用字符串哈希优化比较字符串相等"""
    res = []
    for start in range(len(s)):
        if s[start : start + len(t)] == t:
            res.append(start)
    return res


print(calSubstr("abababababa", "aba"))
