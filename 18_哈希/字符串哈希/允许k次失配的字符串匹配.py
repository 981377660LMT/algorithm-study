# https://oi-wiki.org/string/hash/#%E5%85%81%E8%AE%B8-k-%E6%AC%A1%E5%A4%B1%E9%85%8D%E7%9A%84%E5%AD%97%E7%AC%A6%E4%B8%B2%E5%8C%B9%E9%85%8D
# 子串匹配

from StringHasher import useStringHasher


def match(long: str, short: str, k: int) -> int:
    """
    允许失配k次的字符串匹配.k<=5.
    求long中有多少个子串short,长度相同,允许失配k次.
    """
    n1, n2 = len(long), len(short)
    if n1 < n2:
        return 0

    h1 = useStringHasher([ord(c) for c in long])
    h2 = useStringHasher([ord(c) for c in short])

    def indexOfDiff(start: int, offset: int) -> int:
        """
        从long的下标start开始,找到long[offset:offset+n2)与short第一个不同的位置.
        如果不存在,返回-1.
        """

        if start >= offset + n2:
            return -1

        left, right = start, offset + n2
        while left <= right:
            mid = (left + right) >> 1
            if h1(offset, mid) == h2(0, mid - offset):
                left = mid + 1
            else:
                right = mid - 1
        return left if left < offset + n2 else -1

    res = 0
    for start in range(n1 - n2 + 1):
        cur = start
        ok = False
        for _ in range(k + 1):
            nextDiff = indexOfDiff(cur, start)
            if nextDiff == -1:
                ok = True
                break
            cur = nextDiff + 1
        res += ok
    return res


if __name__ == "__main__":
    import sys

    input = lambda: sys.stdin.readline().rstrip("\r\n")

    n = int(input())
    s = input()
    m = int(input())
    t = input()
    k = int(input())
    print(match(s, t, k))
