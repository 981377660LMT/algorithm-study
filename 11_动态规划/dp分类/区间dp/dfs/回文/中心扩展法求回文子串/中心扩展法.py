from typing import List, Tuple


def expand1(s: str) -> List[Tuple[int, int]]:
    """中心扩展法求所有回文子串 O(n^2)"""

    def expand(left: int, right: int):
        while left >= 0 and right < n and s[left] == s[right]:
            intervals.append((left, right))  # 长度为 right-left+1
            left -= 1
            right += 1

    n = len(s)
    intervals = []
    for i in range(n):
        expand(i, i)
        expand(i, i + 1)
    return intervals


def expand2(s: str) -> List[List[bool]]:
    """中心扩展法标记所有回文子串 O(n^2)

    Return:
        isPalindrome[i][j] 表示 s[i:j+1] 是否是回文串

    # !判断一个子串是否为回文串可以用马拉车算法优化到 O(n)
    """

    def expand(left: int, right: int):
        while left >= 0 and right < n and s[left] == s[right]:
            isPalindrome[left][right] = True
            left -= 1
            right += 1

    n = len(s)
    isPalindrome = [[False] * n for _ in range(n)]  # isPalindrome[i][j] 表示 s[i:j+1] 是否是回文串
    for i in range(n):
        expand(i, i)
        expand(i, i + 1)
    return isPalindrome


if __name__ == "__main__":
    print(expand1("abaccdbbd"))
    print(expand2("abaccdbbd"))
