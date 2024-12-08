def main():
    import sys

    input = sys.stdin.readline

    N, X, K = map(int, input().split())
    items = [tuple(map(int, input().split())) for _ in range(N)]

    from collections import defaultdict

    color_groups = defaultdict(list)
    for p, u, c in items:
        color_groups[c].append((p, u))
    distinct_colors = list(color_groups.keys())
    distinct_colors.sort()

    INF = float("-inf")
    color_dp_list = []
    for c in distinct_colors:
        group = color_groups[c]
        dp_color = [INF] * (X + 1)
        dp_color[0] = 0
        for p, u in group:
            for w in range(X - p, -1, -1):
                if dp_color[w] != INF:
                    new_val = dp_color[w] + u
                    if new_val > dp_color[w + p]:
                        dp_color[w + p] = new_val
        color_dp_list.append(dp_color)

    dp = [INF] * (X + 1)
    dp[0] = 0

    for dp_color in color_dp_list:
        new_dp = [INF] * (X + 1)
        for w in range(X + 1):
            if dp[w] != INF:
                if dp[w] > new_dp[w]:
                    new_dp[w] = dp[w]

        for cst in range(1, X + 1):
            if dp_color[cst] == INF:
                continue
            util = dp_color[cst]
            gain = util + K
            for w in range(X - cst + 1):
                if dp[w] != INF:
                    val = dp[w] + gain
                    if val > new_dp[w + cst]:
                        new_dp[w + cst] = val

        dp = new_dp

    ans = max(dp)
    print(ans)


if __name__ == "__main__":
    main()
