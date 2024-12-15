def v2(x):
    c = 0
    while (x & 1) == 0:
        x >>= 1
        c += 1
    return c


def solve_improved(A):
    N = len(A)
