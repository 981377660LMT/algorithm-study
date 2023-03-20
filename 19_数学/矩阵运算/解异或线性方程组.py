# 解异或线性方程组
# !解决F2上的线性方程组的高斯消元法的实现。
# 加法和乘法分别对应按位异或和按位与操作

from typing import List, Optional


def solve_linear_equation_F2_col(A: List[int], b: List[int], m: int) -> Optional[List[int]]:
    """使用高斯消元法求解线性方程组Ax = b

    矩阵A 的每一行由一个整数表示,整数的二进制表示即为矩阵行的元素.
    向量b也是一个整数列表,每个元素对应一个二进制位.
    """
    assert len(A) == len(b)

    row = len(A)
    for r in range(row):
        indexes = max(range(r, row), key=lambda x: A[x])
        if A[indexes] == 0:
            rank = r
            break
        if indexes != r:
            A[r], A[indexes] = A[indexes], A[r]
            b[r], b[indexes] = b[indexes], b[r]
        for j in range(row):
            if r != j and A[j] > A[j] ^ A[r]:
                A[j] ^= A[r]
                b[j] ^= b[r]
    else:
        rank = row

    # 向量b在rank位置之后的所有元素是否为0。如果不是，返回None，表示方程无解。
    if any(b[i] for i in range(rank, row)):
        return

    res = [0] * m
    for r in range(rank):
        res[A[r].bit_length() - 1] = b[r]
    return res


if __name__ == "__main__":

    # https://yukicoder.me/problems/no/1421
    # 在一个国家里有n个城镇，它们分别被编号为1到n,每个城镇中有若干栋房子，
    # 这个国家的国王想知道每个城镇分别有多少栋房子。
    # 有一天，m 个旅行者访问了这个国家。第j（1≤j≤m）个旅行者游览了国王的p个城镇，
    # 这些城镇的编号为Sj。每个旅行者j在离开国王时报告了以下两点：
    # !他们参观过的城市的编号Sj。
    # !他们参观过的p个城市中，房子数量的总XOR值。
    # !于是，国王决定根据m个旅行者的报告计算出每个城市的房子数量ai。
    # 但是，旅行者中可能有人提交了虚假的报告。
    # 如果存在一个与m个旅行者的所有报告相符的房屋数组合，请输出其中的一个。
    # 如果没有（旅行者的报告中存在矛盾），请输出-1。

    n, m = map(int, input().split())
    A = [0] * m
    b = [0] * m
    for i in range(m):
        _ = int(input())
        cities = list(map(int, input().split()))
        xor_ = int(input())
        for id in cities:
            A[i] |= 1 << (id - 1)
        b[i] = xor_
    res = solve_linear_equation_F2_col(A, b, n)
    if res is None:
        print(-1)
    else:
        print(*res, sep="\n")
