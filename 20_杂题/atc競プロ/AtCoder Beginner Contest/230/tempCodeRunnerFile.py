        preIndex = preSum[curSum]
        dp[i] = (2 * dp[i - 1] - dp[preIndex]) % MOD