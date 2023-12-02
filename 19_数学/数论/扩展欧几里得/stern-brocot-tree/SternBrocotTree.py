# https://maspypy.github.io/library/nt/stern_brocot_tree.hpp
# 正有理数的树状结构(有理数集拓展到实数集)
# Stern–Brocot Tree
# SB树???
# https://oi-wiki.org/math/number-theory/stern-brocot/
#               1/1
#            /       \
#         1/2         2/1
#       /    \        /    \
#     1/3    2/3     3/2     3/1
#    /  \    /  \    /  \     /  \
#  1/4  2/5 3/5 3/4 4/3 5/3  5/2 4/1
#  / \  / \ / \ / \ / \ / \  / \ / \

# https://atcoder.jp/contests/abc294/editorial/6017


from math import ceil
from typing import List, Tuple


R = Tuple[int, int]  # Rational
Path = List[int]  # [0,1,2,0,1,...] 数字构成的路径,表示在R/L方向上走几步,初始为R,然后L,R,L,R...交替


class SternBrocotTree:
    """维护正有理数的平衡二叉树."""

    @staticmethod
    def getPathAndRange(x: R) -> Tuple[Path, R, R]:
        """求x在树中的路径和范围."""
        path = []
        left, right = (0, 1), (1, 0)  # 根节点为(1, 1)
        detL = left[0] * x[1] - left[1] * x[0]
        detR = right[0] * x[1] - right[1] * x[0]
        detM = detL + detR
        while True:
            if detM == 0:
                break
            k = ceil(-detM / detR)
            path.append(k)
            left = (left[0] + k * right[0], left[1] + k * right[1])
            detL += k * detR
            detM += k * detR
            if detM == 0:
                break
            k = ceil(detM / -detL)
            path.append(k)
            right = (right[0] + k * left[0], right[1] + k * left[1])
            detR += k * detL
            detM += k * detL
        return path, left, right

    @staticmethod
    def getPath(x: R) -> Path:
        """从根节点到x的路径."""
        path, _, _ = SternBrocotTree.getPathAndRange(x)
        return path

    @staticmethod
    def getRange(x: R) -> Tuple[R, R]:
        """求x的前驱和后继(分母小于等于x的分母)."""
        _, left, right = SternBrocotTree.getPathAndRange(x)
        return left, right

    @staticmethod
    def isInSubtree(x: R, y: R) -> bool:
        """判断x是否在y的子树中."""
        _, l, r = SternBrocotTree.getPathAndRange(y)
        ok1 = x[0] * l[1] - x[1] * l[0] > 0
        ok2 = r[0] * x[1] - r[1] * x[0] > 0
        return ok1 and ok2

    @staticmethod
    def fromPath(path: Path) -> Tuple[int, int]:
        """从路径找到对应的有理数结点."""
        left, right = (0, 1), (1, 0)
        for i in range(len(path)):
            k = path[i]
            if i & 1:
                right = (right[0] + k * left[0], right[1] + k * left[1])
            else:
                left = (left[0] + k * right[0], left[1] + k * right[1])
        return left[0] + right[0], left[1] + right[1]

    @staticmethod
    def getChild(x: R) -> Tuple[R, R]:
        """求x的左右孩子."""
        _, l, r = SternBrocotTree.getPathAndRange(x)
        lc = (l[0] + x[0], l[1] + x[1])
        rc = (x[0] + r[0], x[1] + r[1])
        return lc, rc

    @staticmethod
    def lca(x: R, y: R) -> R:
        """求x和y的最近公共祖先."""
        px = SternBrocotTree.getPathAndRange(x)[0]
        py = SternBrocotTree.getPathAndRange(y)[0]
        res = []
        for i in range(min(len(px), len(py))):
            k = px[i] if px[i] < py[i] else py[i]
            res.append(k)
            if k < px[i] or k < py[i]:
                break
        return SternBrocotTree.fromPath(res)

    @staticmethod
    def kthAncestor(x: R, k: int) -> R:
        """求x的k级祖先, 如果不存在则返回(-1, -1)."""
        left = (0, 1)
        right = (1, 0)
        mid = (1, 1)
        detL = left[0] * x[1] - left[1] * x[0]
        detR = right[0] * x[1] - right[1] * x[0]
        detM = detL + detR
        while True:
            if detM == 0 or k == 0:
                break
            tmp = ceil(-detM / detR)
            min_ = k if k < tmp else tmp
            left = (left[0] + min_ * right[0], left[1] + min_ * right[1])
            mid = (left[0] + right[0], left[1] + right[1])
            detL += min_ * detR
            detM += min_ * detR
            k -= min_
            if detM == 0 or k == 0:
                break
            tmp = ceil(detM / -detL)
            min_ = k if k < tmp else tmp
            right = (right[0] + min_ * left[0], right[1] + min_ * left[1])
            mid = (left[0] + right[0], left[1] + right[1])
            detR += min_ * detL
            detM += min_ * detL
            k -= min_
        if k == 0:
            return mid
        return -1, -1

    @staticmethod
    def toString(path: Path) -> str:
        """将路径转换为字符串."""
        sb = []
        c = "L"
        for x in path:
            if x == 0:
                continue
            c = "L" if c == "R" else "R"
            sb.append(c * x)
        return "".join(sb)


# 法里级数 中 a/b 第一次出现的位置的前驱和后继(正有理数)
# a/b = 19/12 → (x1/y1, x2/y2) = (11/7, 8/5) → 返回 (11,7,8,5)
def farey(a: int, b: int) -> Tuple[int, int, int, int]:
    """
    求法雷级数中某一项的的前驱和后继(分母小于等于b)
    https://zhuanlan.zhihu.com/p/323538981
    """
    assert a > 0 and b > 0
    if a == b:
        return 0, 1, 1, 0
    q = (a - 1) // b
    x1, y1, x2, y2 = farey(b, a - q * b)
    return q * x2 + y2, x2, q * x1 + y1, x1


if __name__ == "__main__":

    def yosupoJudge() -> None:
        T = int(input())
        for _ in range(T):
            s = input().split()
            if s[0] == "DECODE_PATH":
                path = []
                for char, num in zip(s[2::2], s[3::2]):
                    if not path and char == "L":
                        path.append(0)
                    path.append(int(num))
                a, b = SternBrocotTree.fromPath(path)
                print(a, b)
            elif s[0] == "ENCODE_PATH":
                a, b = map(int, s[1:])
                path = SternBrocotTree.getPath((a, b))
                res = []
                for i in range(len(path)):
                    if path[i] == 0:
                        continue
                    x = "R" if i % 2 == 0 else "L"
                    res.append(x)
                    res.append(str(path[i]))
                print(len(res) // 2, " ".join(res))

            elif s[0] == "LCA":
                a, b, c, d = map(int, s[1:])
                e, f = SternBrocotTree.lca((a, b), (c, d))
                print(e, f)
            elif s[0] == "ANCESTOR":
                k, a, b = map(int, s[1:])
                x, y = SternBrocotTree.kthAncestor((a, b), k)
                if x == -1:
                    print(-1)
                else:
                    print(x, y)
            elif s[0] == "RANGE":
                a, b = map(int, s[1:])
                x, y = SternBrocotTree.getRange((a, b))
                print(x, y)

    yosupoJudge()

    SBT = SternBrocotTree
    # diy
    print(SBT.toString(SBT.getPath((4, 3))))

    # getPath
    assert SBT.getPath((1, 1)) == []
    assert SBT.getPath((1, 2)) == [0, 1]
    assert SBT.getPath((2, 1)) == [1]
    assert SBT.getPath((1, 3)) == [0, 2]
    assert SBT.getPath((2, 3)) == [0, 1, 1]
    assert SBT.getPath((3, 2)) == [1, 1]
    assert SBT.getPath((3, 1)) == [2]
    assert SBT.getPath((1, 4)) == [0, 3]
    assert SBT.getPath((2, 5)) == [0, 2, 1]
    assert SBT.getPath((3, 5)) == [0, 1, 1, 1]
    assert SBT.getPath((3, 4)) == [0, 1, 2]
    assert SBT.getPath((4, 3)) == [1, 2]
    assert SBT.getPath((5, 3)) == [1, 1, 1]
    assert SBT.getPath((5, 2)) == [2, 1]
    assert SBT.getPath((4, 1)) == [3]

    # getRange
    assert SBT.getRange((1, 1)) == ((0, 1), (1, 0))
    assert SBT.getRange((1, 2)) == ((0, 1), (1, 1))
    assert SBT.getRange((2, 1)) == ((1, 1), (1, 0))
    assert SBT.getRange((1, 3)) == ((0, 1), (1, 2))
    assert SBT.getRange((2, 3)) == ((1, 2), (1, 1))
    assert SBT.getRange((3, 2)) == ((1, 1), (2, 1))
    assert SBT.getRange((3, 1)) == ((2, 1), (1, 0))
    assert SBT.getRange((1, 4)) == ((0, 1), (1, 3))
    assert SBT.getRange((2, 5)) == ((1, 3), (1, 2))
    assert SBT.getRange((3, 5)) == ((1, 2), (2, 3))
    assert SBT.getRange((3, 4)) == ((2, 3), (1, 1))
    assert SBT.getRange((4, 3)) == ((1, 1), (3, 2))
    assert SBT.getRange((5, 3)) == ((3, 2), (2, 1))
    assert SBT.getRange((5, 2)) == ((2, 1), (3, 1))
    assert SBT.getRange((4, 1)) == ((3, 1), (1, 0))

    # child
    assert SBT.getChild((1, 1)) == ((1, 2), (2, 1))
    assert SBT.getChild((1, 2)) == ((1, 3), (2, 3))
    assert SBT.getChild((2, 1)) == ((3, 2), (3, 1))
    assert SBT.getChild((1, 3)) == ((1, 4), (2, 5))
    assert SBT.getChild((2, 3)) == ((3, 5), (3, 4))
    assert SBT.getChild((3, 2)) == ((4, 3), (5, 3))
    assert SBT.getChild((3, 1)) == ((5, 2), (4, 1))

    # la
    NG = (-1, -1)
    assert SBT.kthAncestor((1, 1), 0) == (1, 1)
    assert SBT.kthAncestor((1, 1), 1) == NG
    assert SBT.kthAncestor((3, 4), 0) == (1, 1)
    assert SBT.kthAncestor((3, 4), 1) == (1, 2)
    assert SBT.kthAncestor((3, 4), 2) == (2, 3)
    assert SBT.kthAncestor((3, 4), 3) == (3, 4)
    assert SBT.kthAncestor((3, 4), 4) == NG
    assert SBT.kthAncestor((3, 5), 0) == (1, 1)
    assert SBT.kthAncestor((3, 5), 1) == (1, 2)
    assert SBT.kthAncestor((3, 5), 2) == (2, 3)
    assert SBT.kthAncestor((3, 5), 3) == (3, 5)
    assert SBT.kthAncestor((3, 5), 4) == NG

    import random
    from math import gcd

    def get_random():
        while True:
            x = random.randint(1, 1 << 60)
            y = random.randint(1, 1 << 60)
            if gcd(x, y) > 1:
                continue
            return (x, y)

    for _ in range(10000):
        x = get_random()
        l, r = SBT.getRange(x)
        assert x[0] * l[1] - x[1] * l[0] == 1
        assert r[0] * x[1] - r[1] * x[0] == 1
        assert l[0] + r[0] == x[0] and l[1] + r[1] == x[1]
        l, r = SBT.getChild(x)
        assert x[0] * l[1] - x[1] * l[0] == 1
        assert r[0] * x[1] - r[1] * x[0] == 1
        P = SBT.getPath(x)
        assert SBT.fromPath(P) == x
        y = get_random()
        z = SBT.lca(x, y)
        assert SBT.isInSubtree(x, z)
        assert SBT.isInSubtree(y, z)
        l, r = SBT.getChild(z)
        assert not SBT.isInSubtree(x, l) or not SBT.isInSubtree(y, l)
        assert not SBT.isInSubtree(x, r) or not SBT.isInSubtree(y, r)
