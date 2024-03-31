package template.problem;

import template.algo.DancingLink;

// 解数独
public class SudokuProblem {
    /**
     * 0 means unknown
     */
    public static int[][] solve(int[][] mat) {
        assert mat.length == 9 && mat[0].length == 9;
        boolean[][] cover = new boolean[9 * 9 * 9][9 * 9 * 4];
        for (int i = 0; i < 9; i++) {
            for (int j = 0; j < 9; j++) {
                for (int k = 0; k < 9; k++) {
                    int row = (i * 9 + j) * 9 + k;
                    int cell = (i / 3) * 3 + (j / 3);
                    cover[row][i * 9 + k] = true;
                    cover[row][9 * 9 + j * 9 + k] = true;
                    cover[row][9 * 9 * 2 + cell * 9 + k] = true;
                    if (mat[i][j] == 0 || mat[i][j] - 1 == k) {
                        cover[row][9 * 9 * 3 + i * 9 + j] = true;
                    }
                }
            }
        }
        DancingLink dl = DancingLink.newDenseInstance(cover, cover.length, cover[0].length);
        int[] selected = dl.getSolution();
        if (selected == null) {
            return null;
        }
        int[][] ans = new int[9][9];
        for (int choice : selected) {
            int k = choice % 9;
            choice /= 9;
            int j = choice % 9;
            choice /= 9;
            int i = choice;
            ans[i][j] = k + 1;
        }
        return ans;
    }
}
