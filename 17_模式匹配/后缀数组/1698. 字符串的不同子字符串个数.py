# 本质不同子串个数 = 所有子串个数 - 相同子串个数
from SA import useSA


def countDistinct(s: str) -> int:
    """给定一个字符串 s,返回 s 的不同子字符串的个数。

    用所有子串的个数，减去相同子串的个数，就可以得到不同子串的个数。
    计算后缀数组和高度数组。根据高度数组的定义，所有高度之和就是相同子串的个数。

    https://leetcode-cn.com/problems/number-of-distinct-substrings-in-a-string/solution/on-hou-zhui-shu-zu-by-endlesscheng-jo3p/
    https://leetcode-cn.com/problems/number-of-distinct-substrings-in-a-string/solution/python-ji-bai-100-onfu-za-du-saissuan-fa-mwz7/
    """
    _sa, _rank, height = useSA(list(map(ord, s)))
    n = len(s)
    # print(SA, RK, H)
    # [0, 1, 2, 3] [0, 1, 2, 3] [0, 0, 0, 0]
    return n * (n + 1) // 2 - sum(height)


assert countDistinct("banana") == 15
