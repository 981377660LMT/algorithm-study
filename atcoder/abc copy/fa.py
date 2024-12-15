def f(x):

    while (x & 1) == 0:
        x >>= 1
    return x


def brute_force(A):
    N = len(A)
    ans = 0
    for i in range(N):
        for j in range(i, N):
            s = A[i] + A[j]
            ans += f(s)
    return ans


def v2(x):
    c = 0
    while (x & 1) == 0:
        x >>= 1
        c += 1
    return c


def solve_optimized(A):
    N = len(A)
    MAX_T = 25
    MAX_R = 24

    groups = [[] for _ in range(MAX_T)]
    for x in A:
        t = v2(x)
        O = x >> t
        groups[t].append(O)

    sumO = [0] * (MAX_T)
    cntO = [0] * (MAX_T)
    diag = 0
    for t in range(MAX_T):
        cntO[t] = len(groups[t])
        for o in groups[t]:
            sumO[t] += o
        diag += sumO[t]

    full_matrix = 0

    for t in range(MAX_T):
        g = groups[t]
        m = len(g)
        if m == 0:
            continue
        if m == 1:
            full_matrix += sumO[t]
            continue

        g.sort()
        C = [0] * (MAX_R + 2)
        S = [0] * (MAX_R + 2)

        for r in range(1, MAX_R + 1):
            modBase = 1 << r
            freq = [0] * (modBase)
            sumVal = [0] * (modBase)
            for o in g:
                rem = o & (modBase - 1)
                need = (-rem) & (modBase - 1)
                pairCount = freq[need]
                pairSum = sumVal[need]
                if pairCount > 0:
                    C[r] += pairCount
                    S[r] += pairSum + o * pairCount
                freq[rem] += 1
                sumVal[rem] += o

        for r in range(MAX_R, 0, -1):
            C[r] -= C[r + 1]
            S[r] -= S[r + 1]

        group_i_j_sum = 0
        for r in range(1, MAX_R + 1):
            group_i_j_sum += S[r] >> r

        full_matrix += sumO[t] + 2 * group_i_j_sum

    for ti in range(MAX_T):
        if cntO[ti] == 0:
            continue
        for tj in range(ti + 1, MAX_T):
            if cntO[tj] == 0:
                continue
            val = cntO[tj] * sumO[ti] + ((1 << (tj - ti))) * cntO[ti] * sumO[tj]
            full_matrix += 2 * val

    ans = (full_matrix + diag) // 2
    return ans


if __name__ == "__main__":
    for _ in range(1000):
        import random

        N = 10
        A = [random.randint(1, 1000) for _ in range(N)]
        res1, res2 = brute_force(A), solve_optimized(A)
        assert res1 == res2, (res1, res2, A)
    print("All tests passed")
