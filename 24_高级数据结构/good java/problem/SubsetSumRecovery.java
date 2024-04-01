package template.problem;

import template.algo.SortSubsetSum;
import template.utils.SortUtils;

import java.util.Arrays;

class SubsetSumRecovery {
    public long[] recovery(int n, long[] sums) {
        long m = Arrays.stream(sums).min().orElse(-1);
        if (m >= 0) {
            return greedy(n, sums);
        }
        long[] mod = sums.clone();
        for (int i = 0; i < mod.length; i++) {
            mod[i] -= m;
        }
        long[] sol = greedy(n, mod);
        if (sol == null) {
            return null;
        }
        for (int i = 0; i < 1 << n; i++) {
            int num = 0;
            for (int j = 0; j < n; j++) {
                num += ((i >> j) & 1) * sol[j];
            }
            if (num == -m) {
                for (int j = 0; j < n; j++) {
                    if (((i >> j) & 1) == 1) {
                        sol[j] = -sol[j];
                    }
                }
                return sol;
            }
        }
        return null;
    }


    private long[] greedy(int n, long[] sums) {
        Arrays.sort(sums);
        if (sums[0] != 0) {
            return null;
        }
        long[] ans = new long[n];
        int tail = 0;
        long[] a = new long[1 << n];
        long[] b = new long[1 << n];
        int l = 0;
        int r = -1;
        for (int i = 1; i < sums.length; i++) {
            if (l <= r && a[l] == sums[i]) {
                l++;
                continue;
            }
            if (tail >= n) {
                return null;
            }
            long[] subset = SortSubsetSum.sortedSubsetSum(Arrays.copyOf(ans, tail));
            for (int j = 0; j < subset.length; j++) {
                subset[j] += sums[i];
            }
            SortUtils.mergeAscending(a, l, r, subset, 1, subset.length - 1, b, 0);
            r = r - l + subset.length - 1;
            l = 0;
            long[] tmp = a;
            a = b;
            b = tmp;
            ans[tail++] = sums[i];
        }
        return ans;
    }
}
