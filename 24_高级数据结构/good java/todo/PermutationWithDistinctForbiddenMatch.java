package template.problem;

import template.math.DigitUtils;

/**
 * 带禁止位置的排列数
 */
public class PermutationWithDistinctForbiddenMatch {
    private int[][] dp;
    private int[][] dp2;

    /**
     * O(n^2)
     * @param mod
     * @param n
     */
    public PermutationWithDistinctForbiddenMatch(int mod, int n) {
        this.dp = new int[n + 1][n + 1];
        dp[0][0] = 1;
        for (int i = 1; i <= n; i++) {
            for (int j = 0; j <= i; j++) {
                dp[i][j] = (int) ((long)dp[i - 1][j] * (j + j) % mod);
                if (j + 1 <= n) {
                    dp[i][j] = DigitUtils.modplus(dp[i][j], dp[i - 1][j + 1], mod);
                }
                if (j > 0) {
                    dp[i][j] = (int) ((dp[i][j] + (long)dp[i - 1][j - 1] * j % mod * j) % mod);
                }
            }
        }

        dp2 = new int[n + 1][n + 1];
        for (int i = 0; i <= n; i++) {
            for (int j = 0; j <= n; j++) {
                if (j == 0) {
                    dp2[i][j] = dp[i][j];
                    continue;
                }
                if (j > 0) {
                    dp2[i][j] = (int) ((dp2[i][j] + (long)dp2[i][j - 1] * j) % mod);
                }
                if (i > 0) {
                    dp2[i][j] = (int) ((dp2[i][j] + (long)dp2[i - 1][j] * i) % mod);
                }
            }
        }

//        System.err.println(Arrays.deepToString(dp));
//        System.err.println(Arrays.deepToString(dp2));
    }

    /**
     * 排列$i+j$个数($i,j<=n$)，其中对于$k<=i$,第k个数不允许为$k$。
     * @return 方案数
     */
    public int get(int i, int j) {
        return dp2[i][j];
    }
}