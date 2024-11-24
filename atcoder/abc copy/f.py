import sys
import threading
import bisect


def main():
    import sys

    sys.setrecursionlimit(1 << 25)
    input = sys.stdin.readline

    N, Q = map(int, sys.stdin.readline().split())
    S = sys.stdin.readline().strip()

    # Precompute consecutive '1's ending at each position
    left_ones = [0] * N
    for i in range(N):
        if S[i] == "1":
            left_ones[i] = left_ones[i - 1] + 1 if i > 0 else 1
        else:
            left_ones[i] = 0

    # Precompute consecutive '2's starting at each position
    right_twos = [0] * N
    for i in range(N - 1, -1, -1):
        if S[i] == "2":
            right_twos[i] = right_twos[i + 1] + 1 if i < N - 1 else 1
        else:
            right_twos[i] = 0

    # Record positions of '/' and calculate maximum lengths for each
    positions_slash = []
    lengths = []

    for i in range(N):
        if S[i] == "/":
            positions_slash.append(i)
            # Get the length of consecutive '1's ending at position i - 1
            left_k = left_ones[i - 1] if i > 0 else 0
            # Get the length of consecutive '2's starting at position i + 1
            right_k = right_twos[i + 1] if i + 1 < N else 0
            k = min(left_k, right_k)
            length = 2 * k + 1
            lengths.append(length)

    # Build a Sparse Table for Range Maximum Query (RMQ)
    M = len(lengths)
    if M > 0:
        log_table = [0] * (M + 1)
        for i in range(2, M + 1):
            log_table[i] = log_table[i >> 1] + 1

        K = log_table[M] + 1
        st = [[0] * K for _ in range(M)]
        for i in range(M):
            st[i][0] = lengths[i]

        for j in range(1, K):
            for i in range(M - (1 << j) + 1):
                st[i][j] = max(st[i][j - 1], st[i + (1 << (j - 1))][j - 1])

    # Process queries
    for _ in range(Q):
        L_str, R_str = sys.stdin.readline().split()
        L = int(L_str) - 1
        R = int(R_str) - 1

        idx_left = bisect.bisect_left(positions_slash, L)
        idx_right = bisect.bisect_right(positions_slash, R)

        if idx_left >= idx_right:
            print(0)
        else:
            l = idx_left
            r = idx_right - 1
            length = r - l + 1
            if M > 0:
                k = log_table[length]
                max_length = max(st[l][k], st[r - (1 << k) + 1][k])
                print(max_length)
            else:
                print(0)


if __name__ == "__main__":
    threading.Thread(target=main).start()
