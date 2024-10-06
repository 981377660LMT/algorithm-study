# 线性规划特殊情况?


from typing import Tuple

INF = int(4e18)


if __name__ == "__main__":
    N, X = map(int, input().split())
    A, P, B, Q = [], [], [], []
    for _ in range(N):
        a, p, b, q = map(int, input().split())
        A.append(a)
        P.append(p)
        B.append(b)
        Q.append(q)

    def optimizeTwoSelection(v1: int, c1: int, v2: int, c2: int, w: int) -> Tuple[int, int, int]:
        """
        给定两种物品,每种物品分数为vi,价格为ci.
        可以选择任意数量的物品,使得总分数>=w,且总价格最小.
        返回最小价格,以及选择的物品数量.

        即：
        x1, x2 >= 0
        v1 * x1 + v2 * x2 >= w
        最小化 c1 * x1 + c2 * x2.

        返回 x1, x2, cost.

        时间复杂度O(v1+v2).
        """
        x1, x2, cost = INF, 0, 0
        return x1, x2, cost

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
