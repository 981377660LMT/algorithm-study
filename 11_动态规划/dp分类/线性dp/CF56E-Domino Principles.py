# 给你一些多米诺骨牌，在x轴上从左到右排列在一起，
# 给定他们的坐标 x和高度 h
# 问，每一个倒向右边的时候会压倒多少个骨牌？
# !第i块骨牌能压倒的范围在 [xi+1,xi+h-1]之间
# 给定的骨牌并不是按照x从小到大的顺序排列的，被压倒骨牌可能也会压倒其他骨牌。


from typing import List, Tuple


def cf56E(dominos: List[Tuple[int, int]]) -> List[int]:
    n = len(dominos)
    dominoWithId = [(x, h, i) for i, (x, h) in enumerate(dominos)]
    dominoWithId.sort(key=lambda x: x[0])
    res = [1] * n
    maxRight = [-1] * n  # 排序后的第 i 块骨牌压不倒的第一块骨牌的位置
    for i in range(n - 2, -1, -1):
        right = i + 1
        while right != -1 and dominoWithId[i][0] + dominoWithId[i][1] > dominoWithId[right][0]:
            res[dominoWithId[i][2]] += res[dominoWithId[right][2]]
            right = maxRight[right]
        maxRight[i] = right
    return res


if __name__ == "__main__":
    import sys

    input = lambda: sys.stdin.readline().rstrip("\r\n")
    n = int(input())
    dominos = [tuple(map(int, input().split())) for _ in range(n)]
    print(*cf56E(dominos))
