from typing import List, Tuple

# !KMP 求`前缀`循环节


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


# 我们希望知道一个 N 位字符串 S 的前缀是否具有循环节。
# 换言之，对于每一个从头开始的长度为 i（i>1）的前缀，是否由重复出现的子串 A 组成，即 AAA…A （A 重复出现 K 次,K>1）。
# 如果存在，请找出最短的循环节对应的 K 值（也就是这个前缀串的所有可能重复节中，最大的 K 值）。
def getMinCycle(s: str) -> List[Tuple[int, int]]:
    """字符串 S 的前缀 s[:i+1] 是的循环节
    
    Returns:
        List[Tuple[int, int]]: 前缀的长度 循环节的长度
    """
    res = []
    n = len(s)
    nexts = getNext(s)
    for i in range(1, n):  # 所有前缀
        len_ = (i + 1) - nexts[i]
        if len_ and (i + 1) > len_ and (i + 1) % len_ == 0:
            res.append((i + 1, (i + 1) // len_))
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

