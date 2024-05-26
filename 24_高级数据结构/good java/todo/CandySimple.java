package template.problem;

import template.primitve.generated.datastructure.LongComparator;
import template.utils.SortUtils;

public class CandyAssignProblemSimple {
    /**
     * <pre>
     * ans[i] means how many candy transferred from i to i + 1, and it might be negative
     * O(n)
     * </pre>
     */
    public static long[] solve(long[] have, long[] req) {
        int n = have.length;
        long[] b = new long[n];
        for (int i = 0; i < n - 1; i++) {
            int prev = i - 1;
            if (prev < 0) {
                prev += n;
            }
            b[i] = b[prev] + have[i] - req[i];
        }
        long[] kthB = b.clone();
        for (int i = 0; i < n; i++) {
            kthB[i] = -kthB[i];
        }
        int half = (n + 1) / 2;
        long x = SortUtils.theKthSmallestElement(kthB, LongComparator.NATURE_ORDER,
                0, n, half);
        for (int i = 0; i < n; i++) {
            b[i] += x;
        }
        return b;
    }
}