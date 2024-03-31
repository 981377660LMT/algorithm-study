package template.algo;

import template.rand.Randomized;

import java.util.*;

// https://zhuanlan.zhihu.com/p/617477033
// api:
//   newSparseInstance
//   newDenseInstance
//   getSolution
//   TODO: 精准覆盖习题练习
public class DancingLink {
    int[] ans;
    Node[] rowDummy;
    Node[] colDummy;
    Node colHead;
    int[] rowSize;
    int[] colSize;
    Deque<Node> dq;
    Deque<Node> ansDq;

    public static DancingLink newSparseInstance(int[][] pos, int n, int m) {
        DancingLink ans = new DancingLink();
        ans.init(n, m);
        Arrays.sort(pos, Comparator.<int[]>comparingInt(x -> x[0]).thenComparingInt(x -> x[1]));
        for (int[] xy : pos) {
            ans.add(xy[0], xy[1]);
        }
        ans.dq = new ArrayDeque<>(pos.length + n + m);
        ans.dance();
        return ans;
    }

    public static DancingLink newDenseInstance(boolean[][] mat, int n, int m) {
        DancingLink ans = new DancingLink();
        ans.init(n, m);
        int cnt = 0;
        for (int i = 0; i < n; i++) {
            for (int j = 0; j < m; j++) {
                if (mat[i][j]) {
                    cnt++;
                    ans.add(i, j);
                }
            }
        }
        ans.dq = new ArrayDeque<>(cnt + n + m);
        ans.dance();
        return ans;
    }

    /**
     * null means no solution
     *
     * @return
     */
    public int[] getSolution() {
        return ans;
    }

    private DancingLink() {
    }
    
    private Node newNode(int i, int j) {
        Node ans = new Node(i, j);
        return ans;
    }

    private boolean check(Node node, Node skip) {
        return node.l.r == node && node.r.l == node &&
                node.d.u == node && node.u.d == node || node == skip;
    }


    private void add(int i, int j) {
        Node node = newNode(i, j);
        node.l = rowDummy[i].l;
        rowDummy[i].l.r = node;
        node.r = rowDummy[i];
        rowDummy[i].l = node;

        node.u = colDummy[j].u;
        colDummy[j].u.d = node;
        node.d = colDummy[j];
        colDummy[j].u = node;

        rowSize[i]++;
        colSize[j]++;
    }

    private void remove(Node node) {
        node.l.r = node.r;
        node.r.l = node.l;
        node.u.d = node.d;
        node.d.u = node.u;

        assert check(node.l, node);
        assert check(node.r, node);
        assert check(node.u, node);
        assert check(node.d, node);
        // assert check(node);

        if (node.row >= 0) {
            rowSize[node.row]--;
        }
        if (node.col >= 0) {
            colSize[node.col]--;
        }
        node.deleted = true;
    }

    private void recover(Node node) {
        node.l.r = node;
        node.r.l = node;
        node.u.d = node;
        node.d.u = node;

        assert check(node.l, null);
        assert check(node.r, null);
        assert check(node.u, null);
        assert check(node.d, null);
        assert check(node, null);

        if (node.row >= 0) {
            rowSize[node.row]++;
        }
        if (node.col >= 0) {
            colSize[node.col]++;
        }
        node.deleted = false;
    }

    private void init(int n, int m) {
        colSize = new int[m];
        rowSize = new int[n];
        rowDummy = new Node[n];
        colDummy = new Node[m];
        for (int i = 0; i < n; i++) {
            rowDummy[i] = newNode(i, -1);
        }
        for (int i = 0; i < n; i++) {
            int next = i + 1;
            if (next == n) {
                next = 0;
            }
            rowDummy[i].d = rowDummy[next];
            rowDummy[next].u = rowDummy[i];
        }
        for (int i = 0; i < m; i++) {
            colDummy[i] = newNode(-1, i);
        }
        for (int i = 0; i + 1 < m; i++) {
            int next = i + 1;
            colDummy[i].r = colDummy[next];
            colDummy[next].l = colDummy[i];
        }
        colHead = newNode(-1, -1);
        colHead.l = colDummy[m - 1];
        colDummy[m - 1].r = colHead;
        colHead.r = colDummy[0];
        colDummy[0].l = colHead;
        ansDq = new ArrayDeque<>(n);
    }

    private void dfs0(Node root) {
        if (root.deleted) {
            return;
        }
        remove(root);
        dq.addLast(root);
        if (root.col == -1 || root.row == -1) {
            return;
        }
        if (root.u != root) {
            dfs1(root.u);
        }
        if (root.d != root) {
            dfs1(root.d);
        }
        if (root.l != root) {
            dfs0(root.l);
        }
        if (root.r != root) {
            dfs0(root.r);
        }
    }

    private void dfs1(Node root) {
        if (root.deleted) {
            return;
        }
        remove(root);
        dq.addLast(root);
        if (root.col == -1 || root.row == -1) {
            return;
        }
        if (root.u != root) {
            dfs1(root.u);
        }
        if (root.d != root) {
            dfs1(root.d);
        }
        if (root.l != root) {
            dfs2(root.l);
        }
        if (root.r != root) {
            dfs2(root.r);
        }
    }

    private void dfs2(Node root) {
        if (root.deleted) {
            return;
        }
        remove(root);
        dq.addLast(root);
        if (root.col == -1 || root.row == -1) {
            return;
        }
        if (root.l != root) {
            dfs2(root.l);
        }
        if (root.r != root) {
            dfs2(root.r);
        }
    }

    private void removeRow(int i) {
        ansDq.addLast(rowDummy[i]);
        dfs0(rowDummy[i].r);
    }

    private void undo(int size) {
        ansDq.removeLast();
        while (dq.size() > size) {
            recover(dq.removeLast());
        }
    }

    private boolean dance() {
        if (colHead.r == colHead) {
            ans = ansDq.stream().mapToInt(x -> x.row).toArray();
            Randomized.shuffle(ans);
            Arrays.sort(ans);
            return true;
        }
        Node bestCol = colHead.r;
        for (Node node = colHead.r; node != colHead; node = node.r) {
            if (colSize[node.col] < colSize[bestCol.col]) {
                bestCol = node;
            }
        }
        if (colSize[bestCol.col] == 0) {
            return false;
        }
        Node bestRow = bestCol.d;
        for (Node node = bestCol.d; node != bestCol; node = node.d) {
            if (rowSize[node.row] > rowSize[bestRow.row]) {
                bestRow = node;
            }
        }
        Node iter = bestRow;
        while (true) {
            iter = iter.d;
            if (iter.row == -1) {
                iter = iter.d;
            }
            int now = dq.size();
            removeRow(iter.row);
            if (dance()) {
                return true;
            }
            undo(now);

            if (iter == bestRow) {
                break;
            }
        }

        return false;
    }

    private static class Node {
        Node u, d, l, r;
        int row, col;
        boolean deleted;

        public Node(int row, int col) {
            this.row = row;
            this.col = col;

            u = d = l = r = this;
        }

        @Override
        public String toString() {
            return "(" + row + "," + col + ")";
        }
    }
}
