package template.problem;

import template.primitve.generated.datastructure.LongComparator;
import template.primitve.generated.datastructure.LongPriorityQueue;
import template.rand.Randomized;
import template.utils.SequenceUtils;

import java.util.Arrays;

public class RecoveryFromPairSum {
    /**
     * <pre>
     * given n(n-1)/2 pair sum of n (n >= 2) elements, find any satisfied result or null for impossible case.
     * O(n^3\log_2n)
     * result is ordered
     * </pre>
     */
    public static long[] recover(int n, long[] sums) {
        assert sums.length == n * (n - 1) / 2;
        Randomized.shuffle(sums);
        Arrays.sort(sums);
        if (n == 0) {
            return new long[0];
        }
        if (n == 1) {
            return new long[1];
        }
        if (n == 2) {
            long[] ans = new long[]{0, sums[0]};
            if (ans[0] > ans[1]) {
                SequenceUtils.swap(ans, 0, 1);
            }
            return ans;
        }
        long[] res = new long[n];
        LongPriorityQueue pq = new LongPriorityQueue(sums.length, LongComparator.NATURE_ORDER);
        long s12 = sums[0];
        long s13 = sums[1];
        for (int i = 0; i < n - 2; i++) {
            long s23 = sums[i + 2];
            long first = s12 + s13 - s23;
            if (first % 2 != 0) {
                continue;
            }
            first /= 2;
            if (tryCase(sums, first, res, pq)) {
                return res;
            }
        }
        return null;
    }

    private static boolean tryCase(long[] sums, long first, long[] res, LongPriorityQueue pq) {
        int wpos = 0;
        res[wpos++] = first;
        pq.clear();
        for (long x : sums) {
            if (!pq.isEmpty()) {
                if (pq.peek() < x) {
                    return false;
                }
                if (pq.peek() == x) {
                    pq.pop();
                    continue;
                }
            }
            if (wpos >= res.length) {
                return false;
            }
            long newElement = x - first;
            for (int j = 1; j < wpos; j++) {
                pq.add(res[j] + newElement);
            }
            res[wpos++] = newElement;
        }
        assert pq.isEmpty();
        return true;
    }
}
