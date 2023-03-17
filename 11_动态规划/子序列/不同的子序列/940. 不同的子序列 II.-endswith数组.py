"""
给定一个字符串 s,计算 s 的 `不同的子序列`` 的个数
s 仅由小写英文字母组成

!s的本质不同子序列个数
"""


# 总结：

# !方法1
# dp[i][char] 表示前i个字符中以char结尾的子序列的个数
# !转移时直接在 f[i-1][] 对应的子序列末尾加上 s[i],也可以单独加上 s[i] 作为一个子序列

# !方法2
# dp[i] 表示前i个字符中可以组成的不同子序列的个数
# https://leetcode.cn/problems/distinct-subsequences-ii/solution/bu-tong-de-zi-xu-lie-ii-by-leetcode/
# !dp[i] = 2*dp[i-1] - dp[last[s[i]]] (当前选或不选)

from itertools import accumulate
from typing import Any, Sequence


MOD = int(1e9 + 7)


def countSubseq(s: Sequence[Any]) -> int:
    """O(n)求不同子序列个数(非空)"""
    n = len(s)
    dp = [0] * (n + 1)
    dp[0] = 1
    last = dict()
    for i, char in enumerate(s):
        dp[i + 1] = 2 * dp[i] % MOD
        if char in last:
            dp[i + 1] -= dp[last[char]]
        last[char] = i
    return (dp[n] - 1) % MOD


def distinctSubseqII(s: str) -> int:
    """O(n*26)"""
    endswith = [0] * 26
    for char in s:
        endswith[ord(char) - ord("a")] = (sum(endswith) + 1) % MOD
    return sum(endswith) % (MOD)


print(distinctSubseqII(s="aba"))
print(countSubseq(s="aba"))
# 输出：6
# 解释：6 个不同的子序列分别是 "a", "b", "ab", "ba", "aa" 以及 "aba"。


if __name__ == "__main__":
    # https://yukicoder.me/problems/no/1493
    # 每次可以将相邻的两个数换成xor
    # 问可以得到的数组的个数模1e9+7
    # !合并操作变为删除操作

    n = int(input())
    nums = list(map(int, input().split()))

    nums = list(accumulate(nums, lambda x, y: x ^ y))
    nums.pop()
    print((countSubseq(nums) + 1) % MOD)
