package template.problem;

import template.datastructure.DSU;
import template.primitve.generated.datastructure.IntegerComparator;
import template.utils.SortUtils;

import java.util.Arrays;
import java.util.stream.IntStream;

public class ClusterProblem {
    /**
     * for any i != j satisfy dists[i * n + j] > 0 and dists[i * n + i] = 0
     */
    public static boolean[] cluster(int n, double[] dists) {
        return cluster(n, (i, j) -> Double.compare(dists[i], dists[j]));
    }

    /**
     * for any i != j satisfy dists[i * n + j] > 0 and dists[i * n + i] = 0
     */
    public static boolean[] cluster(int n, long[] dists) {
        return cluster(n, (i, j) -> Long.compare(dists[i], dists[j]));
    }

    public static boolean[] cluster(int n, IntegerComparator comp) {
        int[] indices = IntStream.range(0, n * n).toArray();
        SortUtils.quickSort(indices, comp,
                0, indices.length);
        DSUExt dsu = new DSUExt(n);
        dsu.init();
        for (int i = 0; i < indices.length; i++) {
            int j = i;
            while (j + 1 < indices.length && comp.compare(indices[i], indices[j + 1]) == 0) {
                j++;
            }
            for (int k = i; k <= j; k++) {
                int a = indices[k] / n;
                int b = indices[k] - a * n;
                dsu.merge(a, b);
            }
            for (int k = i; k <= j; k++) {
                int a = indices[k] / n;
                dsu.addEdge(dsu.find(a));
            }
            i = j;
        }
        return dsu.dp[dsu.find(0)];
    }

    /**
     * <pre>
     * Construction: O(n^2)
     * Merge all connected component: O(n^2)
     * </pre>
     */
    static class DSUExt extends DSU {
        boolean[][] dp;
        int[] size;
        int[] edge;

        public DSUExt(int n) {
            super(n);
            edge = new int[n];
            size = new int[n];
            dp = new boolean[n][n + 1];
        }

        @Override
        public void init(int n) {
            super.init(n);
            for (int i = 0; i < n; i++) {
                Arrays.fill(dp[i], false);
                size[i] = 1;
            }
        }

        @Override
        protected void preMerge(int a, int b) {
            for (int i = size[a]; i >= 1; i--) {
                if (!dp[a][i]) {
                    continue;
                }
                dp[a][i] = false;
                for (int j = 1; j <= size[b]; j++) {
                    if (dp[b][j]) {
                        dp[a][i + j] = true;
                    }
                }
            }
            edge[a] += edge[b];
            size[a] += size[b];
        }

        public void addEdge(int x) {
            edge[x]++;
            if (edge[x] == size[x] * size[x]) {
                dp[x][1] = true;
            }
        }
    }
}


