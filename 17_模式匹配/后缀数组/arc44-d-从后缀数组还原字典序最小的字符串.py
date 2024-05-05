# arc44-d-从后缀数组还原字典序最小的字符串
# 字符串由大写字母构成
# 如果不存在,输出-1


from typing import List, Optional


def solve(sa: List[int]) -> Optional[str]:
    n = len(sa)
    rank = [0] * (n + 1)
    for i, v in enumerate(sa):
        rank[v] = i + 1

    sb = [ord("A")] * (n + 1)
    pre = sa[0] + 1
    for c in sa[1:]:
        c += 1
        w = sb[pre]
        x = rank[pre]
        y = rank[c]
        if x > y:
            w += 1
        sb[c] = w
        pre = c

    if max(sb) > ord("Z"):
        return None
    return "".join(map(chr, sb[1:]))


if __name__ == "__main__":
    n = int(input())
    sa = list(map(int, input().split()))
    sa = [x - 1 for x in sa]
    res = solve(sa)
    if res is None:
        print(-1)
    else:
        print(res)
