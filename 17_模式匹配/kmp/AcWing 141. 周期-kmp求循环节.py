# 我们希望知道一个 N 位字符串 S 的前缀是否具有循环节。
# 换言之，对于每一个从头开始的长度为 i（i>1）的前缀，是否由重复出现的子串 A 组成，即 AAA…A （A 重复出现 K 次,K>1）。
# 如果存在，请找出最短的循环节对应的 K 值（也就是这个前缀串的所有可能重复节中，最大的 K 值）。


from typing import List

# input = lambda: sys.stdin.readline().strip()


def getNext(shorter: str) -> List[int]:
    """kmp O(n)求 `shorter`串的 `next`数组

    `next[i]`表示`[:i+1]`这一段字符串中最长公共前后缀(不含这一段字符串本身,即真前后缀)的长度
    https://www.ruanyifeng.com/blog/2013/05/Knuth%E2%80%93Morris%E2%80%93Pratt_algorithm.html
    """
    next = [0] * len(shorter)
    j = 0

    for i in range(1, len(shorter)):
        while j and shorter[i] != shorter[j]:  # 1. fallback后前进：匹配不成功j往右走
            j = next[j - 1]

        if shorter[i] == shorter[j]:  # 2. 匹配：匹配成功j往右走一步
            j += 1

        next[i] = j

    return next


def getMinCycle(s: str) -> int:
    """字符串 S 的前缀是否具有循环节"""
    n = len(s)
    nexts = getNext(s)
    for i in range(1, n):
        len_ = (i + 1) - nexts[i]
        if len_ and (i + 1) > len_ and (i + 1) % len_ == 0:
            return i + 1
    return -1


count = 1
while True:
    n = int(input())
    if n == 0:
        break

    print(f"Test case #{count}")
    count += 1

    # next[i]表示[:i+1]这一段字符串中最长公共前后缀(不是原串)的长度
    next = getNext(input())

    for i in range(1, n):
        len_ = (i + 1) - next[i]
        if len_ and (i + 1) > len_ and (i + 1) % len_ == 0:
            print((i + 1), (i + 1) // len_)  # 长度与循环次数
    print()


# todo

