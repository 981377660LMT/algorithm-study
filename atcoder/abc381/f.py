import bisect
import sys

input = lambda: sys.stdin.readline().rstrip("\r\n")

INF = int(4e18)


def longest1122Subsequence(n, a):
    val_to_idx = {}
    unique_vals = sorted(set(a))
    for idx, val in enumerate(unique_vals):
        val_to_idx[val] = idx

    num_vals = len(unique_vals)
    if num_vals > 20:
        num_vals = 20

    positions = [[] for _ in range(num_vals)]
    for idx, val in enumerate(a):
        if val in val_to_idx:
            positions[val_to_idx[val]].append(idx)

    dp = [INF] * (1 << num_vals)
    dp[0] = -1

    for mask in range(1 << num_vals):
        if dp[mask] == INF:
            continue
        for v in range(num_vals):
            if not (mask & (1 << v)) and len(positions[v]) >= 2:
                pos_list = positions[v]
                i = bisect.bisect_right(pos_list, dp[mask])
                if i + 1 < len(pos_list):
                    pos1 = pos_list[i]
                    pos2 = pos_list[i + 1]
                    new_mask = mask | (1 << v)
                    if pos2 < dp[new_mask]:
                        dp[new_mask] = pos2
    max_length = 0
    for mask in range(1 << num_vals):
        if dp[mask] != INF:
            count = mask.bit_count()
            max_length = max(max_length, count * 2)
    return max_length


n = int(input())
a = list(map(int, input().split()))
print(longest1122Subsequence(n, a))
