def findIsomorphic(seq, isMin=True):
    """返回字符串最小表示法/最大表示法"""
    if len(seq) <= 1:
        return seq

    def compare(s1, s2) -> int:
        if s1 == s2:
            return 0
        if isMin:
            return 1 if s1 > s2 else -1
        else:
            return 1 if s1 < s2 else -1

    n = len(seq)
    i1, i2, same = 0, 1, 0

    while i1 < n and i2 < n and same < n:
        diff = compare(seq[(i1 + same) % n], seq[(i2 + same) % n])

        if diff == 0:
            same += 1
            continue
        elif diff > 0:
            i1 += same + 1
        else:
            i2 += same + 1

        if i1 == i2:
            i2 += 1

        same = 0

    res = min(i1, i2)
    return seq[res:] + seq[:res]


def findSuffix(s: str, isMin=True) -> str:
    """返回字典序最小的/最大的后缀

    双指针,指针l记录字典序最大后缀的首位下标,指针r向后扫描并与指针l进行比较
    !注意这里不能循环位移
    """
    if len(s) <= 1:
        return s

    def compare(s1: str, s2: str) -> int:
        if s1 == s2:
            return 0
        if isMin:
            return 1 if s1 > s2 else -1
        else:
            return 1 if s1 < s2 else -1

    n = len(s)
    i1, i2, same = 0, 1, 0

    while i2 + same < n:  # 注意不能循环位移
        diff = compare(s[i1 + same], s[i2 + same])

        if diff == 0:
            same += 1
            continue
        elif diff > 0:
            i1 += same + 1
        else:
            i2 += same + 1

        if i1 == i2:
            i2 += 1

        same = 0

    return s[i1:]


if __name__ == "__main__":
    s = "bcaijab"
    print(findIsomorphic(s))
    print(findIsomorphic(s, False))

    print(findSuffix(s))
    print(findSuffix(s, False))
