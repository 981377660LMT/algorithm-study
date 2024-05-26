package template.problem;

import template.primitve.generated.datastructure.LongArrayList;
import template.utils.Pair;

import java.util.Arrays;
import java.util.Comparator;

public class MaximumArea {
    /**
     * <pre>
     * 在一个有[l,r,b,t]确定的大厅中，有n根柱子。要求找到大厅内部面积最大的一个矩形，内部不包含柱子。
     * return [l, r, b, t] for maximum area rectangle which not include given points
     * </pre>
     */
    public static Pair<long[], Long> rectCover(long l, long r, long b, long t,
                                               long[] xs, long[] ys) {
        int n = xs.length;
        for (int i = 0; i < n; i++) {
            xs[i] = Math.max(xs[i], l);
            xs[i] = Math.min(xs[i], r);
            ys[i] = Math.max(ys[i], b);
            ys[i] = Math.min(ys[i], t);
        }
        long[][] pts = new long[n + 2][2];
        LongArrayList allX = new LongArrayList(n + 2);
        allX.add(l);
        allX.add(r);
        allX.addAll(xs);
        allX.unique();
        long[] allXArray = allX.getData();
        for (int i = 0; i < n; i++) {
            pts[i][0] = allX.binarySearch(xs[i]);
            pts[i][1] = ys[i];
        }
        pts[n][0] = 0;
        pts[n][1] = b;
        pts[n + 1][0] = 0;
        pts[n + 1][1] = t;
        n += 2;

        Arrays.sort(pts, Comparator.comparingLong(x -> x[1]));
        int m = allX.size();
        int[] L = new int[m];
        int[] R = new int[m];
        int[] cnt = new int[m];
        long bestArea = 0;
        long[] ans = new long[4];
        for (int i = 0; i < n; i++) {
            int to = i;
            while (to + 1 < n && pts[to + 1][1] == pts[i][1]) {
                to++;
            }
            i = to;
            for (int j = i + 1; j < n; j++) {
                cnt[(int) pts[j][0]]++;
            }
            long ll = l;
            long rr = l;
            {
                int left = 0;
                for (int j = 1; j < m; j++) {
                    if (cnt[j] > 0 || j == m - 1) {
                        long cand = allXArray[j] - allXArray[left];
                        if (rr - ll < cand) {
                            rr = allXArray[j];
                            ll = allXArray[left];
                        }
                        L[j] = left;
                        R[left] = j;
                        left = j;
                    }
                }
            }
            for (int j = n - 1; j > i; j--) {
                int x = (int) pts[j][0];
                cnt[x]--;
                if (cnt[x] == 0) {
                    //merge
                    L[R[x]] = L[x];
                    R[L[x]] = R[x];
                    long cand = allXArray[R[x]] - allXArray[L[x]];
                    if (rr - ll < cand) {
                        ll = allXArray[L[x]];
                        rr = allXArray[R[x]];
                    }
                }
                long cand = (rr - ll) * (pts[j][1] - pts[i][1]);
                if (cand > bestArea) {
                    bestArea = cand;
                    ans[0] = ll;
                    ans[1] = rr;
                    ans[2] = pts[i][1];
                    ans[3] = pts[j][1];
                }
            }
        }

        return new Pair<>(ans, bestArea);
    }
}
