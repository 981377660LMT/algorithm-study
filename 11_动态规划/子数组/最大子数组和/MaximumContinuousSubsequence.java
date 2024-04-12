package template.algo;

import template.primitve.generated.datastructure.IntToLongFunction;

public class MaximumContinuousSubsequence {
    public static final long INF = (long) 2e18;

    /**
     * O(n)找到最大的一个子数组（不能为空）
     */
    public static long solve(IntToLongFunction func, int l, int r) {
        long prefix = 0;
        long max = -INF;
        for (int i = l; i <= r; i++) {
            prefix += func.apply(i);
            max = Math.max(max, prefix);
            prefix = Math.max(0, prefix);
            prefix = Math.min(prefix, INF);
        }
        return max;
    }
}
