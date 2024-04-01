package template.problem;

import template.binary.Bits;
import template.utils.SequenceUtils;

public class FillGridWithL {
    int[][] mat;
    int[][] color;
    int[][] dirs = new int[][]{
            {1, 0},
            {-1, 0},
            {0, 1},
            {0, -1}
    };

    int indicator = 0;
    int[][] ans;

    private void paint1(int i, int j) {
        assert ans[i][j] == 0;
        assert ans[i - 1][j] == 0;
        assert ans[i][j + 1] == 0;
        ans[i][j] = ans[i - 1][j] = ans[i][j + 1] = ++indicator;
    }

    private void paint2(int i, int j) {
        assert ans[i][j] == 0;
        assert ans[i - 1][j] == 0;
        assert ans[i][j - 1] == 0;
        ans[i][j] = ans[i - 1][j] = ans[i][j - 1] = ++indicator;
    }

    private void paint3(int i, int j) {
        assert ans[i][j] == 0;
        assert ans[i + 1][j] == 0;
        assert ans[i][j - 1] == 0;
        ans[i][j] = ans[i + 1][j] = ans[i][j - 1] = ++indicator;
    }

    private void paint4(int i, int j) {
        assert ans[i][j] == 0;
        assert ans[i + 1][j] == 0;
        assert ans[i][j + 1] == 0;
        ans[i][j] = ans[i + 1][j] = ans[i][j + 1] = ++indicator;
    }

    public int[][] solve(int n, int m) {
        indicator = 0;
        boolean flip = false;
        if (n % 3 == 0 && m % 3 != 0) {
            int tmp = m;
            m = n;
            n = tmp;
            flip = !flip;
        }
        if (n % 6 == 0) {
            int tmp = m;
            m = n;
            n = tmp;
            flip = !flip;
        }
        if (n <= 1 || m <= 1 || m % 3 != 0) {
            return null;
        }
        ans = new int[n][m];
        if (n % 2 == 0) {
            for (int i = 0; i < n; i += 2) {
                for (int j = 0; j < m; j += 3) {
                    paint4(i, j);
                    paint2(i + 1, j + 2);
                }
            }
        } else if (m % 2 == 0) {
            for (int i = 0; i < n - 3; i += 2) {
                for (int j = 0; j < m; j += 3) {
                    paint4(i, j);
                    paint2(i + 1, j + 2);
                }
            }
            for (int i = 0; i < m; i += 2) {
                paint4(n - 3, i);
                paint2(n - 1, i + 1);
            }
        } else if (m >= 9 && n >= 5) {
            paint4(0, 0);
            paint4(2, 0);
            paint2(4, 1);
            paint3(1, 2);
            paint3(0, 3);
            paint1(4, 2);
            for (int i = 3; i + 2 <= m - 2; i += 2) {
                paint4(2, i);
                paint1(4, i + 1);
            }
            paint3(3, m - 1);
            paint2(2, m - 1);
            paint4(0, m - 2);
            for (int i = 4; i < m - 2; i += 3) {
                paint1(1, i);
                paint3(0, i + 2);
            }
            for (int i = 5; i < n; i += 2) {
                for (int j = 0; j < m; j += 3) {
                    paint4(i, j);
                    paint2(i + 1, j + 2);
                }
            }
        } else {
            return null;
        }
        if (flip) {
            ans = transpose(ans);
        }
        return ans;
    }

    private int[][] transpose(int[][] mat) {
        int n = mat.length;
        int m = mat[0].length;
        int[][] ans = new int[m][n];
        for (int i = 0; i < n; i++) {
            for (int j = 0; j < m; j++) {
                ans[j][i] = mat[i][j];
            }
        }
        return ans;
    }

    private void setColor(int i, int j, int v, int fi, int fj, int c) {
        if (v != mat[i][j]) {
            return;
        }
        color[i][j] = c;
        for (int[] d : dirs) {
            int x = d[0] + i;
            int y = d[1] + j;
            if (!valid(x, y)) {
                continue;
            }
            if (x == fi && y == fj) {
                continue;
            }
            setColor(x, y, v, i, j, c);
        }
    }

    private boolean valid(int i, int j) {
        return i >= 0 && j >= 0 && i < mat.length && j < mat[0].length;
    }

    private int findColor(int i, int j, int v, int fi, int fj) {
        if (v != mat[i][j]) {
            if (color[i][j] == -1) {
                return 0;
            }
            return 1 << color[i][j];
        }
        int ans = 0;
        for (int[] d : dirs) {
            int x = d[0] + i;
            int y = d[1] + j;
            if (!valid(x, y)) {
                continue;
            }
            if (x == fi && y == fj) {
                continue;
            }
            ans |= findColor(x, y, v, i, j);
        }
        return ans;
    }

    private int[][] color(int[][] mat) {
        int n = mat.length;
        int m = mat[0].length;
        this.mat = mat;
        color = new int[n][m];
        SequenceUtils.deepFill(color, -1);
        for (int i = 0; i < n; i++) {
            for (int j = 0; j < m; j++) {
                if (color[i][j] == -1) {
                    int bits = findColor(i, j, mat[i][j], -1, -1);
                    int c = 0;
                    while (Bits.get(bits, c) == 1) {
                        c++;
                    }
                    setColor(i, j, mat[i][j], -1, -1, c);
                }
            }
        }
        return color;
    }
}
