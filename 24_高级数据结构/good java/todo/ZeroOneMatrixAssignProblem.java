package template.problem;

import template.math.DigitUtils;
import template.math.GCDs;

/**
 * 对于大小为n*m的01矩阵，要求每一行的和为r，每一列的和为c（nr=mc且r<=m且c<=n)。
 * 这个问题会设计一套赋值方案
 */
public class ZeroOneMatrixAssignProblem {
    int n;
    int m;
    int r;
    int c;
    int g;
    int mg;
    int rg;

    /**
     * O(log n)
     */
    public ZeroOneMatrixAssignProblem(int n, int m, int r, int c) {
        if (!(n * r == m * c && r <= m && c <= n)) {
            throw new IllegalArgumentException();
        }
        this.n = n;
        this.m = m;
        this.r = r;
        this.c = c;
        g = GCDs.gcd(r, m);
        mg = m / g;
        rg = r / g;
    }

    /**
     * O(1)
     */
    public boolean isOne(int i, int j) {
        return (i + j) % mg < rg;
    }

    /**
     * O(1)
     */
    public int valueOf(int i, int j) {
        return isOne(i, j) ? 1 : 0;
    }

    @Override
    public String toString() {
        StringBuilder builder = new StringBuilder();
        for (int i = 0; i < n; i++) {
            for (int j = 0; j < m; j++) {
                builder.append(valueOf(i, j));
            }
            builder.append('\n');
        }
        return builder.toString();
    }
}
