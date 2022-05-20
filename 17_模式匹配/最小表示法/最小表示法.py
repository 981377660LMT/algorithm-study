def findIsomorphic(s: str, isMin=True) -> str:
    """返回字符串最小/最大表示法"""
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

    while i1 < n and i2 < n and same < n:
        diff = compare(s[(i1 + same) % n], s[(i2 + same) % n])

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

    res = i2 if i1 > i2 else i1
    return s[res:] + s[:res]


if __name__ == '__main__':
    s = 'bcaijab'
    print(findIsomorphic(s))
    print(findIsomorphic(s, False))

