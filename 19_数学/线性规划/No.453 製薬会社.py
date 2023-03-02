# https://yukicoder.me/problems/771
# https://yukicoder.me/submissions/728212
# https://yukicoder.me/submissions/781385


from typing import List


def simplex(A: List[List[float]], b: List[int], c: List[int]) -> float:
    """
    单纯形法求解线性规划问题
    A*x <= b,x >= 0 的条件下，最大化 c*x 的值
    """

    def make_tableau(A, b, c):
        m = len(A)

        tableau = []
        for i in range(m):
            row = A[i] + [int(j == i) for j in range(m)] + [b[i]]
            tableau.append(row)
        row = c + [0] * (m + 1)
        tableau.append(row)

        return tableau

    def pivot_index(tableau):
        cN = tableau[-1][:-1]
        piv_col = -1
        for col, x in enumerate(cN):
            if x > 0:
                piv_col = col
                break
        if piv_col == -1:
            return False, -1, -1

        a = [tableau[i][piv_col] for i in range(len(tableau) - 1)]
        b = [tableau[i][-1] for i in range(len(tableau) - 1)]
        thetas = [bi / ai if ai > 0 else float("inf") for ai, bi in zip(a, b)]
        piv_row = thetas.index(min(thetas))

        return True, piv_row, piv_col

    def step(tableau, piv_row, piv_col):
        h = len(tableau)
        w = len(tableau[0])
        piv = tableau[piv_row][piv_col]

        for j in range(w):
            tableau[piv_row][j] /= piv

        for i in range(h):
            if i == piv_row:
                continue
            d = tableau[i][piv_col]
            for j in range(w):
                tableau[i][j] -= d * tableau[piv_row][j]

    tableau = make_tableau(A, c, b)

    while True:
        improved, piv_row, piv_col = pivot_index(tableau)
        if not improved:
            break
        step(tableau, piv_row, piv_col)

    return -tableau[-1][-1]


if __name__ == "__main__":
    # 制药公司制作两种药品A和B，
    # 每一个A需要3/4kg薬品C和1/4kg薬品D，
    # 每一个B需要2/7kg薬品C和5/7kg薬品D，
    # 每个A售价1000，每个B售价2000，
    # 给定薬品C和薬品D的数量，最大化的销售额
    C, D = map(int, input().split())
    A = [[3 / 4, 2 / 7], [1 / 4, 5 / 7]]
    b = [1000, 2000]
    c = [C, D]

    print(simplex(A, b, c))
