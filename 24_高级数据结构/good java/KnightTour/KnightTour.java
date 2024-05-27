package template.problem;

import template.rand.Randomized;
import template.utils.SequenceUtils;

import java.util.stream.IntStream;

public class KnightTour {
    private static int[][] dirs = new int[][]{
            {1, 2},
            {-1, 2},
            {1, -2},
            {-1, -2},
            {2, 1},
            {-2, 1},
            {-2, -1},
            {2, -1}
    };

    private int n;
    private int m;
    public int[][] mat;
    private int step;
    int[] ids;

    public KnightTour(int n, int m, int r, int c) {
        this.n = n;
        this.m = m;
        mat = new int[n][m];
        ids = IntStream.range(0, n * m).toArray();
        while (true) {
            SequenceUtils.deepFill(mat, -1);
            Randomized.shuffle(ids);
            step = 0;
            if (dfs(r, c)) {
                break;
            }
        }
    }

    private boolean possible(int i, int j) {
        if (i < 0 || j < 0 || i >= n || j >= m || mat[i][j] != -1) {
            return false;
        }
        return true;
    }

    public int degree(int i, int j) {
        int ans = 0;
        for (int[] dir : dirs) {
            if (possible(dir[0] + i, dir[1] + j)) {
                ans++;
            }
        }
        return ans;
    }

    private boolean dfs(int i, int j) {
        mat[i][j] = step++;
        int bestX = -1;
        int bestY = -1;
        int bestDeg = Integer.MAX_VALUE;
        int bestId = -1;
        for (int[] dir : dirs) {
            int x = i + dir[0];
            int y = j + dir[1];
            if (!possible(x, y)) {
                continue;
            }
            int d = degree(x, y);
            if (d < bestDeg) {
                bestDeg = d;
                bestId = -1;
            }
            if (d == bestDeg && bestId < ids[x * m + y]) {
                bestId = ids[x * m + y];
                bestX = x;
                bestY = y;
            }
        }
        if (bestDeg == Integer.MAX_VALUE) {
            //is it end
            return step == n * m;
        }
        return dfs(bestX, bestY);
    }


}
