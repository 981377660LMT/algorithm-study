from collections import defaultdict
import heapq


def main():
    import sys

    N = int(input())
    A = list(map(int, input().split()))

    # 数値1〜20の出現位置を収集（1-based）
    pos_dict = [[] for _ in range(21)]
    for idx, num in enumerate(A):
        pos_dict[num].append(idx + 1)

    # 各数値のペアをpでソート
    pairs = [[] for _ in range(21)]
    for num in range(1, 21):
        pos = pos_dict[num]
        for i in range(len(pos) - 1):
            for j in range(i + 1, len(pos)):
                p, q = pos[i], pos[j]
                pairs[num].append((p, q))
        # ソートはpの昇順
        pairs[num].sort()

    INF = N + 2
    dp = [INF] * (1 << 20)
    dp[0] = 0

    for mask in range(1 << 20):
        if dp[mask] == INF:
            continue
        last_q = dp[mask]
        for num in range(1, 21):
            bit = num - 1
            if not (mask & (1 << bit)):
                # Binary search to find the first pair with p > last_q
                pl = 0
                pr = len(pairs[num])
                # We need to find the smallest index where p > last_q
                while pl < pr:
                    mid = (pl + pr) // 2
                    if pairs[num][mid][0] > last_q:
                        pr = mid
                    else:
                        pl = mid + 1
                if pl < len(pairs[num]):
                    p, q = pairs[num][pl]
                    new_mask = mask | (1 << bit)
                    if dp[new_mask] > q:
                        dp[new_mask] = q

    # Find the mask with the maximum number of set bits and dp[mask] < INF
    max_k = 0
    for mask in range(1 << 20):
        if dp[mask] < INF:
            count = mask.bit_count()
            if count > max_k:
                max_k = count
    print(max_k * 2)


if __name__ == "__main__":
    main()
