package template.problem;


import template.binary.Bits;
import template.binary.Log2;

public class MaximumIndependentSet {
    /**
     * meet in middle, O(2^{n/2})
     * @param edges
     * @param weights
     * @param selections
     * @return
     */
    public static long solve(boolean[][] edges, long[] weights, boolean[] selections) {
        int n = weights.length;
        if (n > 60) {
            throw new IllegalArgumentException("Too large set");
        }
        long[] adj = new long[n];
        for (int i = 0; i < n; i++) {
            for (int j = 0; j < n; j++) {
                if (i == j) {
                    continue;
                }
                if (edges[i][j] || edges[j][i]) {
                    adj[i] = Bits.set(adj[i], j);
                }
            }
        }

        int leftHalf = (n + 1) / 2;
        int rightHalf = n - leftHalf;


        int[] subsets = new int[1 << leftHalf];
        long[] lSum = new long[1 << leftHalf];
        for (int i = 1; i < 1 << leftHalf; i++) {
            int log = Log2.floorLog(i);
            int highest = 1 << log;
            lSum[i] = lSum[i - highest] + weights[log];
            subsets[i] = subsets[i - highest];
            int possible = subsets[(int) ((i - highest) & (~adj[log]))];
            if (lSum[subsets[i]] < lSum[possible] + weights[log]) {
                subsets[i] = possible | highest;
            }
        }

        int mask = (1 << leftHalf) - 1;
        long[] rSum = new long[1 << rightHalf];
        long[] nearby = new long[1 << rightHalf];

        long solution = subsets[mask];
        long ans = lSum[subsets[mask]];
        long inf = (long) 2e18;
        for (int i = 1; i < 1 << rightHalf; i++) {
            int log = Log2.floorLog(i);
            int highest = 1 << log;
            rSum[i] = Bits.get(nearby[i - highest], leftHalf + log) == 1 ? -inf : (rSum[i - highest] + weights[leftHalf + log]);
            nearby[i] = nearby[i - highest] | adj[leftHalf + log];

            int leftSubset = subsets[(int) (mask & (~nearby[i]))];
            long cand = rSum[i] + lSum[leftSubset];
            if (cand > ans) {
                ans = cand;
                solution = ((long) i << leftHalf) | leftSubset;
            }
        }

        for (int i = 0; i < n; i++) {
            selections[i] = Bits.get(solution, i) == 1;
        }
        return ans;
    }
}
