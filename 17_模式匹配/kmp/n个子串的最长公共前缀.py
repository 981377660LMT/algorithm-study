# 一个字符串按要求分割为k个子串，求这k个子串的最长公共前缀的长度


def solve(s: str, k: int) -> int:
    def check(mid: int) -> bool:
        pre = s[:mid]
        count = s.split(pre)
        return len(count) >= k + 1

    def check2(mid: int) -> bool:
        """可以用z函数加前缀和预处理达到O(1)check"""
        ...

    left, right = 1, len(s) + 5
    while left <= right:
        mid = (left + right) // 2
        if check(mid):
            left = mid + 1
        else:
            right = mid - 1
    return right


assert solve("ababcab", 3) == 2
