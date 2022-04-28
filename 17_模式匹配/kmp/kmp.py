from typing import List


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


if __name__ == '__main__':
    next = getNext('aabaabaabaab')
    assert next == [0, 1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9]

