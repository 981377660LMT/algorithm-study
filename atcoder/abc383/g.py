import sys

input = sys.stdin.readline
from collections import deque


def solve():
    N, K = map(int, input().split())
    A = list(map(int, input().split()))

    if K == 1:
        A.sort(reverse=True)
        prefix = [0]
        for x in A:
            prefix.append(prefix[-1] + x)
        M = N
        print(" ".join(map(str, prefix[1 : M + 1])))
        return

    length = N - K + 1
    current_sum = sum(A[:K])
    B = [current_sum]
    for i in range(K, N):
        current_sum += A[i] - A[i - K]
        B.append(current_sum)

    M = N // K

    dp_prev = [0] * (length + 1)
    results = [0] * (M + 1)

    dp_curr = [-(10**15)] * (length + 1)
    for i in range(1, length + 1):
        if i == 1:
            dp_curr[i] = B[i - 1]
        else:
            dp_curr[i] = max(dp_curr[i - 1], B[i - 1])
    dp_prev = dp_curr
    results[1] = dp_prev[length]

    for j in range(2, M + 1):
        dp_curr = [-(10**15)] * (length + 1)

        dq = deque()

        prefix_max_curr = -(10**15)

        for i in range(1, length + 1):
            if i > 1:
                prefix_max_curr = max(prefix_max_curr, dp_curr[i - 1])
            else:
                prefix_max_curr = -(10**15)

            if i - K >= 0:
                val_to_add = dp_prev[i - K]
                while dq and dp_prev[dq[-1]] <= val_to_add:
                    dq.pop()
                dq.append(i - K)

            while dq and dq[0] < i - K:
                dq.popleft()

            candidate = prefix_max_curr
            if dq:
                candidate = max(candidate, dp_prev[dq[0]] + B[i - 1])

            dp_curr[i] = candidate

        dp_prev = dp_curr
        results[j] = dp_prev[length]

    print(" ".join(map(str, results[1:])))


if __name__ == "__main__":
    solve()
