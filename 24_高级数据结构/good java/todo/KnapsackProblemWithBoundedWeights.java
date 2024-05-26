package template.problem;

import template.utils.SequenceUtils;

import java.util.Arrays;

/**
 * <p>
 * Based on "Linear Time Algorithms for Knapsack Problems with Bounded Weights" by David Pisinger
 * </p>
 * <p>
 * O(nM) while n is the number of elements, and M is the maximum weight in those elements
 * </p>
 */
public class KnapsackProblemWithBoundedWeights {
    /**
     * <pre>
     * maximize: \sum_{i\in I} w_i
     * subject to: \sum_{i\in I} w_i <= C
     * </pre>
     *
     * @param w weights of elements
     * @param C the capacity of knapsack
     * @return
     */
    public static int solve(int[] w, int C) {
        assert C >= 0;
        int n = w.length;
        int maxW = Arrays.stream(w).max().getAsInt();
        int[] prev = new int[maxW * 2];
        int[] next = new int[maxW * 2];
        int[] prevprev = new int[maxW * 2];

        int zero = maxW - 1;
        Arrays.fill(prev, -1);
        Arrays.fill(prevprev, -1);
        int ps = 0;
        int b = 0;
        while (b < n && ps + w[b] <= C) {
            ps += w[b];
            b++;
        }
        if (b == n) {
            return ps;
        }
        prev[zero + ps - C] = b;
        for (int i = b; i <= n; i++) {
            Arrays.fill(next, -1);
            for (int j = maxW * 2 - 1; j >= 0; j--) {
                if (prev[j] == -1) {
                    continue;
                }
                //add
                if (i + 1 <= n) {
                    next[j] = Math.max(next[j], prev[j]);
                    if (j <= zero) {
                        next[j + w[i]] = Math.max(next[j + w[i]], prev[j]);
                    }
                }
                //or sub
                if (j > zero) {
                    int threshold = Math.max(0, prevprev[j]);
                    for (int k = prev[j] - 1; k >= threshold; k--) {
                        prev[j - w[k]] = Math.max(prev[j - w[k]], k);
                    }
                }
            }
            int[] tmp = prevprev;
            prevprev = prev;
            prev = next;
            next = tmp;
        }
        int best = zero;
        while (best >= 0 && prevprev[best] == -1) {
            best--;
        }

        return C + best - zero;
    }

    public static boolean[] solution(int[] w, int C) {
        assert C >= 0;
        int n = w.length;
        int maxW = Arrays.stream(w).max().getAsInt();
        int[][] dp = new int[maxW * 2][n + 1];
        int[][] from = new int[maxW * 2][n + 1];

        class Optimizer {
            public void update(int i, int j, int fi, int fj, int val) {
                if (dp[i][j] < val) {
                    dp[i][j] = val;
                    from[i][j] = fi * (n + 1) + fj;
                }
            }
        }
        Optimizer optimizer = new Optimizer();

        int zero = maxW - 1;
        SequenceUtils.deepFill(dp, -1);
        SequenceUtils.deepFill(from, -1);
        int ps = 0;
        int b = 0;
        while (b < n && ps + w[b] <= C) {
            ps += w[b];
            b++;
        }
        if (b == n) {
            boolean[] ans = new boolean[n];
            Arrays.fill(ans, true);
            return ans;
        }
        dp[zero + ps - C][b] = b;
        for (int i = b; i <= n; i++) {
            for (int j = maxW * 2 - 1; j >= 0; j--) {
                if (dp[j][i] == -1) {
                    continue;
                }
                //add
                if (i + 1 <= n) {
                    optimizer.update(j, i + 1, j, i, dp[j][i]);
                    if (j <= zero) {
                        optimizer.update(j + w[i], i + 1, j, i, dp[j][i]);
                    }
                }
                //or sub
                if (j > zero) {
                    int threshold = i == 0 ? 0 : Math.max(0, dp[j][i - 1]);
                    for (int k = dp[j][i] - 1; k >= threshold; k--) {
                        optimizer.update(j - w[k], i, j, i, k);
                    }
                }
            }
        }
        int best = zero;
        while (best >= 0 && dp[best][n] == -1) {
            best--;
        }
        boolean[] ans = new boolean[n];
        Arrays.fill(ans, 0, b, true);
        int curX = best;
        int curY = n;
        while (from[curX][curY] != -1) {
            int lx = from[curX][curY] / (n + 1);
            int ly = from[curX][curY] % (n + 1);
            if (ly == curY) {
                //remove prev
                assert dp[curX][curY] < b;
                ans[dp[curX][curY]] = false;
            } else if (lx == curX) {
                //omit
            } else {
                ans[ly] = true;
            }
            curX = lx;
            curY = ly;
        }
//        int sum = 0;
//        for (int i = 0; i < n; i++) {
//            if (ans[i]) {
//                sum += w[i];
//            }
//        }
//        assert sum == best - zero + C;
        return ans;
    }
}
