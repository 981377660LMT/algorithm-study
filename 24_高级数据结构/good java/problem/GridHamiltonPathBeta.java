package template.problem;

import template.utils.CloneSupportObject;
import template.utils.Debug;
import template.utils.SequenceUtils;

public class GridHamiltonPathBeta {
    Debug debug = new Debug(false);

    public char[][] solve(int n, int m, int sx, int sy, int dx, int dy) {
        PathBuilder ans = new PathBuilder(n, m, sx, sy, dx, dy);
        if (!acceptable(ans)) {
            return null;
        }
        solve(ans);
        assert ans.validatePath();
        debug.debug("ans", ans);
        return ans.build();
    }

    static HamiltonPathBF[][][][] hp = new HamiltonPathBF[6][5][5][4];

    static {
        for (int i = 1; i <= 3; i++) {
            for (int j = 1; j <= 3; j++) {
                for (int x = 0; x < i; x++) {
                    for (int y = 0; y < j; y++) {
                        hp[i][j][x][y] = new HamiltonPathBF(i, j, x, y);
                    }
                }
            }
        }
        for (int i = 0; i < 5; i++) {
            for (int j = 0; j < 4; j++) {
                hp[5][4][i][j] = new HamiltonPathBF(5, 4, i, j);
            }
        }
    }

    private int color(int x, int y) {
        return (x + y) & 1;
    }

    private int color(int[] x) {
        return color(x[0], x[1]);
    }

    private boolean acceptable(PathBuilder mat) {
        mat = mat.clone();
        if (mat.length() < mat.width()) {
            mat.rotate();
        }
        int n = mat.length();
        int m = mat.width();
        int[] src = mat.source();
        int[] dst = mat.dst();
        if (n * m % 2 == 0) {
            if ((src[0] + src[1]) % 2 == (dst[0] + dst[1]) % 2) {
                return false;
            }
        } else {
            if ((src[0] + src[1]) % 2 != 0 || (dst[0] + dst[1]) % 2 != 0) {
                return false;
            }
        }
        if (m == 1 && !((src[0] == 0 && dst[0] == n - 1) || (src[0] == n - 1 && dst[0] == 0))) {
            return false;
        }
        if (m == 2) {
            if (src[1] == dst[1]) {
            } else if (src[0] == 0 && dst[0] == 0 || src[0] == n - 1 && dst[0] == n - 1) {
            } else if (Math.abs(src[0] - dst[0]) <= 1) {
                return false;
            }
        }
        if (m == 3 && n % 2 == 0) {
            if (src[0] < dst[0]) {
                if (color(src) != color(0, 0)) {
                    if (src[0] < dst[0] - 1 || src[1] == 1) {
                        return false;
                    }
                }
            } else if(src[0] > dst[0]){
                if (color(src) != color(n - 1, 0)) {
                    if (src[0] > dst[0] + 1 || src[1] == 1) {
                        return false;
                    }
                }
            }
        }
        return true;
    }

    /**
     * DLLL
     * RRRU
     *
     * @param mat
     * @param low
     */
    private void addLoop(PathBuilder mat, int low, boolean inv) {
        int m = mat.width();
        if (!inv) {
            for (int i = 0; i < m; i++) {
                mat.setChar(low, i, 'R');
                mat.setChar(low + 1, i, 'L');
            }
            mat.setChar(low + 1, 0, 'D');
            mat.setChar(low, m - 1, 'U');
        } else {
            for (int i = 0; i < m; i++) {
                mat.setChar(low, i, 'L');
                mat.setChar(low + 1, i, 'R');
            }
            mat.setChar(low + 1, m - 1, 'D');
            mat.setChar(low, 0, 'U');
        }
    }

    private boolean strip(PathBuilder mat) {
        int n = mat.length();
        int m = mat.width();
        if (n <= 2) {
            return false;
        }
        int[] src = mat.source();
        int[] dst = mat.dst();
        //try strip
        if (src[0] >= 2 && dst[0] >= 2) {
            PathBuilder sub = mat.subregion(0, m - 1, 2, n - 1, src[0], src[1],
                    dst[0], dst[1]);
            if (acceptable(sub)) {
                solve(sub);

                for (int i = 0; i < m; i++) {
                    char c = mat.getChar(2, i);
                    if (c == 'L') {
                        //find
                        addLoop(mat, 0, true);
                        mat.setChar(2, i, 'D');
                        mat.setChar(1, i - 1, 'U');
                        break;
                    } else if (c == 'R') {
                        addLoop(mat, 0, false);
                        mat.setChar(2, i, 'D');
                        mat.setChar(1, i + 1, 'U');
                        break;
                    }
                }
                assert mat.validatePath();
                return true;
            }
        }
        if (src[0] < n - 2 && dst[0] < n - 2) {
            PathBuilder sub = mat.subregion(0, m - 1, 0, n - 3, src[0], src[1],
                    dst[0], dst[1]);
            if (acceptable(sub)) {
                solve(sub);
                for (int i = 0; i < m; i++) {
                    char c = mat.getChar(n - 3, i);
                    if (c == 'L') {
                        //find
                        addLoop(mat, n - 2, false);
                        mat.setChar(n - 3, i, 'U');
                        mat.setChar(n - 2, i - 1, 'D');
                        break;
                    } else if (c == 'R') {
                        addLoop(mat, n - 2, true);
                        mat.setChar(n - 3, i, 'U');
                        mat.setChar(n - 2, i + 1, 'D');
                        break;
                    }
                }
                assert mat.validatePath();
                return true;
            }
        }

        return false;
    }

    private boolean split(PathBuilder mat) {
        int n = mat.length();
        int m = mat.width();
        if (n <= 1) {
            return false;
        }
        int[] src = mat.source();
        int[] dst = mat.dst();
        int B = Math.min(src[0], dst[0]);
        int T = Math.max(src[0], dst[0]);

        for (int i = B; i < T; i++) {
            PathBuilder bot = mat.subregion(0, m - 1, 0, i, src[0], src[1], dst[0], dst[1]);
            PathBuilder top = mat.subregion(0, m - 1, i + 1, n - 1, src[0], src[1], dst[0], dst[1]);
            for (int j = 0; j < m; j++) {
                if (src[0] < dst[0]) {
                    bot.setDst(i, j);
                    top.setSource(0, j);
                    if (acceptable(bot) && acceptable(top)) {
                        solve(bot);
                        solve(top);
                        bot.setChar(i, j, 'U');
                        return true;
                    }
                } else {
                    top.setDst(0, j);
                    bot.setSource(i, j);
                    if (acceptable(bot) && acceptable(top)) {
                        solve(bot);
                        solve(top);
                        top.setChar(0, j, 'D');
                        return true;
                    }
                }
            }
        }
        return false;
    }

    private void solve(PathBuilder mat) {
        mat = mat.clone();
        assert acceptable(mat);
        if (mat.length() < mat.width()) {
            mat.rotate();
        }
        int n = mat.length();
        int m = mat.width();
        int[] src = mat.source();
        int[] dst = mat.dst();
        assert src[0] >= 0 && src[0] < n && src[1] >= 0 && src[1] < m;
        assert dst[0] >= 0 && dst[0] < n && dst[1] >= 0 && dst[1] < m;
        if (n <= 3 && m <= 3 || n == 5 && m == 4) {
            char[][] res = hp[n][m][src[0]][src[1]].res[dst[0]][dst[1]];
            if (res == null) {
                throw new IllegalStateException();
            }
            for (int i = 0; i < n; i++) {
                for (int j = 0; j < m; j++) {
                    mat.setChar(i, j, res[i][j]);
                }
            }
            assert mat.validatePath();
            return;
        }
        //strip
        if (strip(mat) || strip(mat.rotate())) {
            assert mat.validatePath();
            return;
        }
        if (split(mat) || split(mat.rotate())) {
            assert mat.validatePath();
            return;
        }
        throw new IllegalStateException();
    }
}

class HamiltonPathBF {
    char[][][][] res;
    char[][] to;
    int n;
    int m;

    private char[][] deepClone() {
        char[][] ans = to.clone();
        for (int i = 0; i < ans.length; i++) {
            ans[i] = ans[i].clone();
        }
        return ans;
    }


    private void dfs(int i, int j, int len) {
        if (to[i][j] != 0) {
            return;
        }
        if (len == n * m) {
            //finish
            if (res[i][j] == null) {
                res[i][j] = deepClone();
            }
            return;
        }
        if (i > 0) {
            to[i][j] = 'D';
            dfs(i - 1, j, len + 1);
        }
        if (j > 0) {
            to[i][j] = 'L';
            dfs(i, j - 1, len + 1);
        }
        if (i + 1 < n) {
            to[i][j] = 'U';
            dfs(i + 1, j, len + 1);
        }
        if (j + 1 < m) {
            to[i][j] = 'R';
            dfs(i, j + 1, len + 1);
        }
        to[i][j] = 0;
    }

    public HamiltonPathBF(int n, int m, int sx, int sy) {
        res = new char[n][m][][];
        to = new char[n][m];
        this.n = n;
        this.m = m;
        dfs(sx, sy, 1);
    }
}

class PathBuilder extends CloneSupportObject<PathBuilder> {
    char[][] mat;
    int l;
    int r;
    int b;
    int t;
    int sx;
    int sy;
    int dx;
    int dy;
    boolean rotate;

    public void setSource(int x, int y) {
        if (rotate) {
            x ^= y;
            y ^= x;
            x ^= y;
        }
        sx = x + b;
        sy = y + l;
    }

    public void setDst(int x, int y) {
        if (rotate) {
            x ^= y;
            y ^= x;
            x ^= y;
        }
        dx = x + b;
        dy = y + l;
    }

    public int[] source() {
        int[] ans = new int[]{sx - b, sy - l};
        if (rotate) {
            SequenceUtils.swap(ans, 0, 1);
        }
        return ans;
    }

    public int[] dst() {
        int[] ans = new int[]{dx - b, dy - l};
        if (rotate) {
            SequenceUtils.swap(ans, 0, 1);
        }
        return ans;
    }

    private int rawLength() {
        return t - b + 1;
    }

    private int rawWidth() {
        return r - l + 1;
    }

    public int length() {
        return rotate ? rawWidth() : rawLength();
    }

    public int width() {
        return rotate ? rawLength() : rawWidth();
    }

    public char[][] build() {
        int L = length();
        int W = width();
        char[][] ans = new char[L][W];
        for (int i = 0; i < L; i++) {
            for (int j = 0; j < W; j++) {
                ans[i][j] = getChar(i, j);
            }
        }
        return ans;
    }

    @Override
    public String toString() {
        StringBuilder ans = new StringBuilder("\n");
        int[] source = source();
        int[] dst = dst();
        ans.append(source[0]).append(' ').append(source[1]).append('\n');
        ans.append(dst[0]).append(' ').append(dst[1]).append('\n');
        char[][] data = build();
        SequenceUtils.reverse(data);
        for (char[] s : data) {
            for (char c : s) {
                ans.append(c == 0 ? '.' : c);
            }
            ans.append('\n');
        }
        return ans.toString();
    }

    public PathBuilder(int n, int m, int sx, int sy, int dx, int dy) {
        mat = new char[n][m];
        l = 0;
        r = m - 1;
        b = 0;
        t = n - 1;
        this.sx = sx;
        this.sy = sy;
        this.dx = dx;
        this.dy = dy;
    }

    public PathBuilder(char[][] mat, int l, int r, int b, int t, int sx, int sy, int dx, int dy, boolean rotate) {
        this.mat = mat;
        this.l = l;
        this.r = r;
        this.b = b;
        this.t = t;
        this.sx = sx;
        this.sy = sy;
        this.dx = dx;
        this.dy = dy;
        this.rotate = rotate;
    }

    public PathBuilder subregion(int L, int R, int B, int T, int sx, int sy, int dx, int dy) {
        if (rotate) {
            L ^= B;
            B ^= L;
            L ^= B;

            R ^= T;
            T ^= R;
            R ^= T;

            sx ^= sy;
            sy ^= sx;
            sx ^= sy;

            dx ^= dy;
            dy ^= dx;
            dx ^= dy;
        }
        return new PathBuilder(mat, l + L, l + R, b + B, b + T, sx + b, sy + l, dx + b, dy + l, rotate);
    }

    public PathBuilder rotate() {
        rotate = !rotate;
        return this;
    }

    private char rotateChar(char c) {
        switch (c) {
            case 'L':
                return 'D';
            case 'R':
                return 'U';
            case 'D':
                return 'L';
            case 'U':
                return 'R';
            default:
                return c;
        }
    }

    public char getChar(int i, int j) {
        if (rotate) {
            int tmp = i;
            i = j;
            j = tmp;
        }
        assert valid(i + b, j + l);
        char ans = mat[i + b][j + l];
        if (rotate) {
            ans = rotateChar(ans);
        }
        return ans;
    }

    public void setChar(int i, int j, char c) {
        if (rotate) {
            int tmp = i;
            i = j;
            j = tmp;
            c = rotateChar(c);
        }
        assert valid(i + b, j + l);
        mat[i + b][j + l] = c;
    }

    public void setNext(int i, int j, int ni, int nj) {
        assert Math.abs(i - ni) == 1 && j == nj || i == ni && Math.abs(j - nj) == 1;
        char c;
        if (i < ni) {
            c = 'U';
        } else if (i > ni) {
            c = 'D';
        } else if (j < nj) {
            c = 'R';
        } else {
            c = 'L';
        }
        setChar(i, j, c);
    }

    private char inverse(char c) {
        switch (c) {
            case 'L':
                return 'R';
            case 'R':
                return 'L';
            case 'U':
                return 'D';
            case 'D':
                return 'U';
            default:
                throw new IllegalArgumentException();
        }
    }

    public void inverse(int L, int R, int B, int T) {
        for (int i = B; i <= T; i++) {
            for (int j = L; j <= R; j++) {
                char c = getChar(i, j);
                c = inverse(c);
                setChar(i, j, c);
            }
        }
    }

    private boolean valid(int i, int j) {
        return b <= i && i <= t && l <= j && j <= r;
    }

    public boolean validatePath() {
        boolean[][] visited = new boolean[length()][width()];
        int[] source = source();
        int[] dst = dst();
        while (!(source[0] == dst[0] && source[1] == dst[1])) {
            if (source[0] < 0 || source[0] >= visited.length || source[1] < 0
                    || source[1] >= visited[0].length) {
                return false;
            }
            if (visited[source[0]][source[1]]) {
                return false;
            }
            visited[source[0]][source[1]] = true;
            switch (getChar(source[0], source[1])) {
                case 'L':
                    source[1]--;
                    break;
                case 'R':
                    source[1]++;
                    break;
                case 'D':
                    source[0]--;
                    break;
                case 'U':
                    source[0]++;
                    break;
            }
        }
        visited[source[0]][source[1]] = true;
        for (boolean[] arr : visited) {
            for (boolean x : arr) {
                if (!x) {
                    return false;
                }
            }
        }
        return true;
    }
}