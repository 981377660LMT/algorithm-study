MOD = 998244353


def char_to_parity(c):
    # Returns the parity bit pattern for A,B,C
    # A affects bit 0, B bit 1, C bit 2
    if c == "A":
        return 1  # 001 in binary
    elif c == "B":
        return 2  # 010 in binary
    elif c == "C":
        return 4  # 100 in binary
    else:
        raise ValueError("Invalid character")


def solve():
    import sys

    input = sys.stdin.readline
    N, K = map(int, input().split())
    S = input().strip()

    # dp[i][state][g] = number of ways
    # i ranges from 0 to N
    # state in 0..7 (for parity of A,B,C)
    # g in 0..K (we only need to track up to K, but K can be large)
    # We'll only track up to K good substrings; if more than K, we can cap it to K because we only need "at least K".

    # Initialize dp array
    # To save memory, we can use two-layer DP (current and next) since dp[i] depends only on dp[i-1].
    dp = [[[0] * (K + 1) for _ in range(8)] for __ in range(2)]
    dp_curr = dp[0]
    dp_next = dp[1]

    # Base case
    dp_curr[0][0] = 1  # empty prefix, parity=000, no good substrings

    # We'll process each character
    for i in range(N):
        c = S[i]
        # Compute ways_count for dp_curr
        ways_count = [0] * 8
        for st in range(8):
            for g in range(K + 1):
                if dp_curr[st][g]:
                    ways_count[st] = (ways_count[st] + dp_curr[st][g]) % MOD

        # Prepare dp_next
        for st in range(8):
            for g in range(K + 1):
                dp_next[st][g] = 0

        # Determine possible characters
        chars = []
        if c == "?":
            chars = ["A", "B", "C"]
        else:
            chars = [c]

        for old_st in range(8):
            for g in range(K + 1):
                val = dp_curr[old_st][g]
                if val == 0:
                    continue
                for ch in chars:
                    p = char_to_parity(ch)
                    new_st = old_st ^ p
                    # new good substrings formed = ways_count[new_st] + ways_count[new_st ^ 7]
                    formed = ways_count[new_st] + ways_count[new_st ^ 7]
                    formed %= MOD

                    new_g = g + formed
                    if new_g > K:
                        new_g = K  # Cap at K because we only need counts for g≥K
                    dp_next[new_st][new_g] = (dp_next[new_st][new_g] + val) % MOD

        # Swap
        dp_curr, dp_next = dp_next, dp_curr

    # Sum over all states and g≥K
    ans = 0
    for st in range(8):
        for g in range(K, K + 1):
            ans = (ans + dp_curr[st][g]) % MOD

    print(ans)


solve()
