# 性价比差的那个物品不会用超过lcm/性价比好的物品的价格这么多件，
# 否则这个lcm的总价格可以替成性价比高的物品

from typing import Tuple

INF = int(4e18)


def max2(a: int, b: int) -> int:
    return a if a > b else b


def optimizeTwoSelection(v1: int, c1: int, v2: int, c2: int, w: int) -> Tuple[int, int, int]:
    """
    给定两种物品,每种物品分数为vi,价格为ci.
    选择任意数量的物品,使得总分数>=w,且总价格最小.
    返回最小价格 cost, 以及选择的物品数量 x1, x2.

    即：
    x1, x2 >= 0
    v1 * x1 + v2 * x2 >= w
    最小化 c1 * x1 + c2 * x2.

    时间复杂度 `O(v1+v2).`
    """
    assert v1 >= 1 and v2 >= 1 and c1 >= 1 and c2 >= 1, "v1, v2, c1, c2 must be positive."
    x1, x2, cost = 0, 0, INF
    # !要么第一种物品选择个数<=v2, 要么第二种物品选择个数<=v1.
    for i1 in range(v2 + 1):
        cur = i1 * c1
        remain = max2(0, w - i1 * v1)
        i2 = (remain + v2 - 1) // v2
        cur += i2 * c2
        if cur < cost:
            x1, x2, cost = i1, i2, cur
    for i2 in range(v1 + 1):
        cur = i2 * c2
        remain = max2(0, w - i2 * v2)
        i1 = (remain + v1 - 1) // v1
        cur += i1 * c1
        if cur < cost:
            x1, x2, cost = i1, i2, cur
    return x1, x2, cost


if __name__ == "__main__":
    # E - Sensor Optimization Dilemma 2
    # https://atcoder.jp/contests/abc374/tasks/abc374_e
    # 线性规划特殊情况?
    N, X = map(int, input().split())
    A, P, B, Q = [], [], [], []
    for _ in range(N):
        a, p, b, q = map(int, input().split())
        A.append(a)
        P.append(p)
        B.append(b)
        Q.append(q)

    def check(mid: int) -> bool:
        """每个项目生成mid个产品需要的费用<=X."""
        res = 0
        for a, p, b, q in zip(A, P, B, Q):
            _, _, c = optimizeTwoSelection(a, p, b, q, mid)
            res += c
        return res <= X

    left, right = 1, int(2e9)
    ok = False
    while left <= right:
        mid = (left + right) // 2
        if check(mid):
            left = mid + 1
            ok = True
        else:
            right = mid - 1
    print(right)
