# D. Send More Money
# !解方程 s1+s2=s3
# 注意暴力枚举的方法没有回溯快

from itertools import permutations
from typing import List, Tuple

# 口算难题


def sendMoreMoney(words: List[str], result: str) -> Tuple[Tuple[int, int, int], bool]:
    """求方程的解，返回解和是否有解"""
    charset = list(set("".join(words) + result))
    if len(charset) > 10:
        return (0, 0, 0), False

    mp = {c: i for i, c in enumerate(charset)}
    allWords = words + [result]
    # 全排列枚举每个字母对应的数字
    for perm in permutations(range(10), len(charset)):
        # 排除前导零
        if any(perm[mp[w[0]]] == 0 for w in allWords):
            continue
        res = [0] * len(allWords)
        for i, s in enumerate(allWords):
            for c in s:
                res[i] = res[i] * 10 + perm[mp[c]]
        if sum(res[:-1]) == res[-1]:
            return tuple(res), True
    return (0, 0, 0), False


if __name__ == "__main__":

    import sys

    sys.setrecursionlimit(int(1e9))
    input = lambda: sys.stdin.readline().rstrip("\r\n")
    s1, s2, s3 = input(), input(), input()
    res, ok = sendMoreMoney([s1, s2], s3)
    if not ok:
        print("UNSOLVABLE")
        exit(0)
    print(*res, sep="\n")
