def maxSumTwoNoOverlap(A, L, M):
    # 数组可依据索引j划分为两部分，L在i左边，M在i右边或L在i右边，M在i左边。
    # 对于每个j，按照两种情况计算最大的L和M和，并更新答案
    n = len(A)
    # 计算前缀和，方便后边求数组和
    for i in range(1, n):
        A[i] += A[i - 1]
    print(A)
    ans = A[L + M - 1]
    Lmax = A[L - 1]
    Mmax = A[M - 1]
    # i代表当前位于右边的数组的末尾索引
    for i in range(L + M, n):
        # 当L在M前时，i代表M的最后一个索引,此时M已确定
        Lmax = max(Lmax, A[i - M] - A[i - M - L])
        ans1 = Lmax + A[i] - A[i - M]
        # 当L在M后时，i代表L的最后一个索引，此时L已确定
        Mmax = max(Mmax, A[i - L] - A[i - L - M])
        ans2 = Mmax + A[i] - A[i - L]
        ans = max(ans, ans1, ans2)
    return ans


print(maxSumTwoNoOverlap([0, 6, 5, 2, 2, 5, 1, 9, 4], 1, 2))
