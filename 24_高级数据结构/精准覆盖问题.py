#
#  Exact Cover
#
#  Description:
#  We are given a family of sets F on [0,n).
#  The exact cover problem is to find a subfamily of F
#  such that each k in [0,n) is covered exactly once.
#  For example, if F consists from
#  {1,2}, {1,2,3}, {3,4}, {5,6}, {3,5,6},
#  then the exact cover is
#  {1,2}, {3,4}, {5,6}.
#
#  Algorithm:
#  Knuth's algorithm X is the following recursive algorithm:
#
#  select some k in [0,n)
#  for each subset S that covers k
#  select S and remove all conflicting sets
#  recursion
#
# To implement this algorithm efficiently, we can use
# a data structure, which is called dancing links.
#
#  Verified:
#  SPOJ 1428: EASUDOKU
#  SPOJ 1110: SUDOKU
#
# 精确覆盖/精准覆盖问题


from math import sqrt
from typing import List, Optional, Tuple


def exactCover(sets: List[List[int]], n: Optional[int] = None) -> List[int]:
    """精确覆盖问题.
    sets所有元素在[0,n)内.
    求出sets的一个子集,使得[0,n)内每个元素恰好出现一次.
    返回这个子集的索引.
    不存在返回空列表.
    """

    def remove(x: int) -> None:
        L[R[x]] = L[x]
        R[L[x]] = R[x]
        i = D[x]
        while i != x:
            j = R[i]
            while j != i:
                U[D[j]] = U[j]
                D[U[j]] = D[j]
                S[C[j]] -= 1
                j = R[j]
            i = D[i]

    def resume(x: int) -> None:
        i = U[x]
        while i != x:
            j = L[i]
            while j != i:
                U[D[j]] = j
                D[U[j]] = j
                S[C[j]] += 1
                j = L[j]
            i = U[i]
        L[R[x]] = x
        R[L[x]] = x

    def rec() -> bool:
        if R[n] == n:  # type: ignore
            return True
        col = R[n]  # type: ignore
        i = R[n]  # type: ignore
        while i != n:
            if S[i] < S[col]:
                col = i
            i = R[i]
        if S[col] == 0:
            return False
        remove(col)
        i = D[col]
        while i != col:
            res.append(A[i])
            j = R[i]
            while j != i:
                remove(C[j])
                j = R[j]
            if rec():
                return True
            j = L[i]
            while j != i:
                resume(C[j])
                j = L[j]
            res.pop()
            i = D[i]
        resume(col)
        return False

    if n is None:
        n = 0
        for s in sets:
            for x in s:
                if x + 1 > n:
                    n = x + 1
    M = (n + 1) * (1 + len(sets)) + 10  # row * col
    L, R, U, D = [0] * M, [0] * M, [0] * M, [0] * M
    S, C, A = [0] * M, [0] * M, [0] * M
    for i in range(n + 1):
        L[i] = i - 1
        R[i] = i + 1
        D[i] = U[i] = C[i] = i
    L[0] = n
    R[n] = 0
    p = n + 1
    for row in range(len(sets)):
        for i in range(len(sets[row])):
            col = sets[row][i]
            C[p] = col
            A[p] = row
            S[col] += 1
            D[p] = D[col]
            U[p] = col
            D[col] = U[D[p]] = p
            if i == 0:
                L[p] = R[p] = p
            else:
                L[p] = p - 1
                R[p] = R[p - 1]
                R[p - 1] = L[R[p]] = p
            p += 1

    res = []
    rec()
    return res


def sudoku(board: List[List[str]], inplace=True, whiteSpace=".") -> Tuple[List[List[str]], bool]:
    """给定一个n*n的数独棋盘,求解数独.
    whiteSpace 表示需要填写的位置,'1'~'n'表示已经填写的数字.
    !Not Verified.
    """

    def getId(a: int, b: int, c: int) -> int:
        return w * w * w * w * a + w * w * b + c

    def addSet(i: int, j: int, k: int) -> None:
        sets.append([])
        sets[-1].append(getId(0, i, j))
        sets[-1].append(getId(1, i, k))
        sets[-1].append(getId(2, j, k))
        sets[-1].append(getId(3, w * (i // w) + (j // w), k))
        ns.append((i, j, k))

    n = len(board)
    w = int(sqrt(n))  # w*w内n个数需要为1~n
    sets = []
    ns = []
    for i in range(w * w):
        for j in range(w * w):
            if board[i][j] == whiteSpace:
                for k in range(w * w):
                    addSet(i, j, k)
            else:
                addSet(i, j, int(board[i][j]) - 1)

    x = exactCover(sets)
    if not x:
        return [], False
    if not inplace:
        board = [list(row) for row in board]
    for a in x:
        i, j, k = ns[a]
        board[i][j] = str(k + 1)
    return board, True


if __name__ == "__main__":
    n = 3
    sets = [[0, 1, 2], [1], [2], [3]]
    print(exactCover(sets))

    class Solution2:
        def solveSudoku(self, board: List[List[str]]) -> None:
            """
            Do not return anything, modify board in-place instead.
            """
            sudoku(board, True, whiteSpace=".")

    # P4929 【模板】舞蹈链（DLX）
    # https://www.luogu.com.cn/problem/P4929
    ROW, COL = map(int, input().split())
    sets = [[] for _ in range(ROW)]
    for i in range(ROW):
        row = list(map(int, input().split()))
        for j in range(COL):
            if row[j] == 1:
                sets[i].append(j)
    rows = exactCover(sets, COL)
    if not rows:
        print("No Solution!")
    else:
        print(*[v + 1 for v in rows], end="\n")
