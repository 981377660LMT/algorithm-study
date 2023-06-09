# !KMP 求字符串循环节
# aabaabaab
# ******
#    ******
# 如果存在周期，那么 i-kmp[i] = 3 为最小周期 (i=9,kmp[i]=6)
# !当区间[l+d,r]的哈希值与[l,r-d]的哈希值相等时，那么该区间[l,r]是以 d 为循环节的**


from typing import List, Tuple
from kmp import KMP


# KMP求前缀周期/字符串循环节/字符串周期
# 我们希望知道一个 N 位字符串 S 的前缀是否具有循环节。
# 换言之，对于每一个从头开始的长度为 i（i>1）的前缀，是否由重复出现的子串 A 组成，即 AAA…A （A 重复出现 K 次,K>1）。
# 如果存在，请找出最短的循环节对应的 K 值（也就是这个前缀串的所有可能重复节中，最大的 K 值）。
def getMinCycle(s: str) -> List[Tuple[int, int]]:
    """求字符串 S 的前缀 s[:i+1] 的循环节

    Returns:
        List[Tuple[int, int]]: 前缀的长度 循环节出现的次数
    """
    res = []
    n = len(s)
    kmp = KMP(s)
    for i in range(1, n):  # 所有前缀
        period = kmp.period(i)
        if period:
            res.append((i + 1, (i + 1) // period))
    return res


count = 1
# input = lambda: sys.stdin.readline().strip()
while True:
    n = int(input())
    if n == 0:
        break

    print(f"Test case #{count}")
    count += 1

    s = input()
    res = getMinCycle(s)
    for preLen, times in res:
        print(f"{preLen} {times}")
    print()
