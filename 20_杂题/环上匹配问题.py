# 环上匹配问题
# 环上线段相交问题
# https://taodaling.github.io/blog/2021/07/26/%E4%B8%80%E4%BA%9B%E5%9C%86%E7%8E%AF%E9%97%AE%E9%A2%98/
# 题目1：给定2n个人，顺时针编号为1,…,2n。
# 现在要求让环上的两两拉一条绳子，且绳子之间的交点数最多（重复的交点也要多次统计）。
# !很显然让i和i+n共拉一条绳子，这时候任意两条绳子都有交点，达到理论最大值C(2n,2)=n*(2n-1)。


# 题目2：给定2n个人，顺时针编号为1,…,2n。现在要求让环上的两两拉一条绳子，且绳子之间的交点数最多（重复的交点也要多次统计）。
# !其中我们已经存在k条绳子，拉第i条绳子的人为ai和bi。
# !实际上我们把所有未匹配的人按照顺时针排序后，第i个人和第i+n−k个人拉一条绳子即可。
# 可以发现如果(a,b)和(c,d)形成的两条绳子无交点，那么(a,d)和(c,b)就一定有交点，且与其它绳子的交点总数是不会减少的。
# 换言之最优解的情况下，未选择的绳子必须两两相交，而这时候只有一个唯一解，就是每个人和对面的人拉绳子。

from typing import List, Tuple


# https://www.luogu.com.cn/problem/CF1552C
def CF1552C(n: int, pairs: List[Tuple[int, int]]) -> int:
    def intersect(a: int, b: int, c: int, d: int) -> bool:
        """弦(a,b)和弦(c,d)是否相交."""
        return (a < c < b < d) or (c < a < d < b)

    m = len(pairs)
    allPairs = []
    for a, b in pairs:
        if a > b:
            a, b = b, a
        allPairs.append((a, b))

    used = [False] * (2 * n)
    for a, b in pairs:
        used[a] = used[b] = True
    remain = [v for v in range(2 * n) if not used[v]]
    for i in range(n - m):
        allPairs.append((remain[i], remain[i + n - m]))

    res = 0
    for i in range(n):
        a, b = allPairs[i]
        for j in range(i + 1, n):
            c, d = allPairs[j]
            res += intersect(a, b, c, d)
    return res


if __name__ == "__main__":
    T = int(input())
    for _ in range(T):
        n, k = map(int, input().split())
        pairs = []
        for i in range(k):
            a, b = map(int, input().split())
            a, b = a - 1, b - 1
            pairs.append((a, b))
        print(CF1552C(n, pairs))
