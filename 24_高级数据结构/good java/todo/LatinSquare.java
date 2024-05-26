package template.problem;

import template.math.DigitUtils;
import template.utils.SequenceUtils;

/**
 * <pre>
 * 拉丁方是一个n^2的矩阵，矩阵的每个元素都是0到n-1之间的整数，且每一行每一列都是一个置换。
 * 支持O(1)的左右上下滚动，行逆，列逆，转置操作
 * 支持O(1)单个单元格查询操作
 * </pre>
 */
public class LatinSquare {
    int[][] ij;
    int[][] iv;
    int[][] jv;
    int n;
    ElementModifyWrapper[] cast = new ElementModifyWrapper[]{
            new ElementModifyWrapper(new Element() {
                @Override
                public int apply(int i, int j, int v) {
                    return i;
                }

                @Override
                public int index() {
                    return 0;
                }
            }),
            new ElementModifyWrapper(new Element() {
                @Override
                public int apply(int i, int j, int v) {
                    return j;
                }

                @Override
                public int index() {
                    return 1;
                }
            }),
            new ElementModifyWrapper(new Element() {
                @Override
                public int apply(int i, int j, int v) {
                    return v;
                }

                @Override
                public int index() {
                    return 2;
                }
            }),
    };

    public LatinSquare(int[][] mat) {
        this(mat, true);
    }

    public LatinSquare(int[][] mat, boolean copy) {
        n = mat.length;
        ij = mat;
        ij = copy ? new int[n][n] : mat;
        iv = new int[n][n];
        jv = new int[n][n];
        for (int i = 0; i < n; i++) {
            for (int j = 0; j < n; j++) {
                assert mat[i][j] >= 0 && mat[i][j] < n;
                int v = mat[i][j];
                ij[i][j] = v;
                iv[i][v] = j;
                jv[j][v] = i;
            }
        }
    }

    /**
     * O(1)
     */
    public int get(int i, int j) {
        int type = cast[2].index();
        int[][] table;
        int inv;
        if (type == 0) {
            table = jv;
            //jv
            inv = cast[0].index() == 1 ? 0 : 1;
        } else if (type == 1) {
            table = iv;
            //iv
            inv = cast[0].index() == 0 ? 0 : 1;
        } else {
            table = ij;
            inv = cast[0].index() == 0 ? 0 : 1;
        }
        if (inv == 1) {
            int tmp = i;
            i = j;
            j = tmp;
        }
        int mi = DigitUtils.mod(i - cast[0 ^ inv].mod, n);
        int mj = DigitUtils.mod(j - cast[1 ^ inv].mod, n);
        int ans = DigitUtils.mod(table[mi][mj] + cast[2].mod, n);
        return ans;
    }

    public void up(long x) {
        down(-x);
    }

    public void left(long x) {
        right(-x);
    }

    public void down(long x) {
        cast[0].modify(x, n);
    }

    public void right(long x) {
        cast[1].modify(x, n);
    }

    public void rowInv() {
        SequenceUtils.swap(cast, 1, 2);
    }

    public void colInv() {
        SequenceUtils.swap(cast, 0, 2);
    }

    public void plus(long x) {
        cast[2].modify(x, n);
    }

    public void transpose() {
        SequenceUtils.swap(cast, 0, 1);
    }

    static interface Element {
        int apply(int i, int j, int v);

        int index();
    }

    static class ElementModifyWrapper implements Element {
        public ElementModifyWrapper(Element e) {
            this.e = e;
        }

        int mod;

        @Override
        public int apply(int i, int j, int v) {
            return e.apply(i, j, v) + mod;
        }

        public void modify(long x, int n) {
            mod = DigitUtils.mod(mod + x, n);
        }

        @Override
        public int index() {
            return e.index();
        }

        Element e;
    }

    @Override
    public String toString() {
        StringBuilder ans = new StringBuilder("\n");
        for (int i = 0; i < n; i++) {
            for (int j = 0; j < n; j++) {
                ans.append(get(i, j)).append(' ');
            }
            ans.append('\n');
        }
        return ans.toString();
    }
}
