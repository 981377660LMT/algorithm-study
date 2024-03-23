package template.algo;

import java.util.Arrays;

public class UndoKnapsack {
    public long[] dp;
    public static final long INF = (long) 2e18;

    public UndoKnapsack(int size) {
        dp = new long[size + 1];
        Arrays.fill(dp, -INF);
        dp[0] = 0;
    }

    /**
     * apply take O(dp.length) and undo take O(1)
     */
    public UndoOperation add(int weight, long value) {
        return new UndoOperation() {
            long[] next = new long[dp.length];

            @Override
            public void apply() {
                for (int i = 0; i < dp.length; i++) {
                    next[i] = dp[i];
                    if (i >= weight) {
                        next[i] = Math.max(next[i], dp[i - weight] + value);
                    }
                }

                swap();
            }

            void swap() {
                long[] tmp = dp;
                dp = next;
                next = tmp;
            }

            @Override
            public void undo() {
                swap();
            }
        };
    }
}
