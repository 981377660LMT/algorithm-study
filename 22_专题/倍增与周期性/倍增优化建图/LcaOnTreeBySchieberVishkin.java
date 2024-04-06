package template.graph;

import template.binary.Log2;

import java.util.List;
import java.util.function.IntPredicate;

/**
 * https://github.com/indy256/codelibrary/blob/master/java/graphs/lca/LcaSchieberVishkin.java
 */
// Answering LCA queries in O(1) with O(n) preprocessing
public class LcaOnTreeBySchieberVishkin implements LcaOnTree{
    int[] parent;
    int[] preOrder;
    int[] i;
    int[] head;
    int[] a;
    int time;

    void dfs1(List<? extends DirectedEdge>[] tree, int u, int p) {
        parent[u] = p;
        i[u] = preOrder[u] = time++;
        for (DirectedEdge e : tree[u]) {
            int v = e.to;
            if (v == p) continue;
            dfs1(tree, v, u);
            if (Integer.lowestOneBit(i[u]) < Integer.lowestOneBit(i[v])) {
                i[u] = i[v];
            }
        }
        head[i[u]] = u;
    }

    void dfs2(List<? extends DirectedEdge>[] tree, int u, int p, int up) {
        a[u] = up | Integer.lowestOneBit(i[u]);
        for (DirectedEdge e : tree[u]) {
            int v = e.to;
            if (v == p) continue;
            dfs2(tree, v, u, a[u]);
        }
    }

    public void init(List<? extends DirectedEdge>[] tree, int root) {
        init(tree, i -> i == root);
    }

    public void init(List<? extends DirectedEdge>[] tree, IntPredicate isRoot) {
        time = 0;
        for (int i = 0; i < tree.length; i++) {
            if (isRoot.test(i)) {
                dfs1(tree, i, -1);
                dfs2(tree, i, -1, 0);
            }
        }
    }

    public LcaOnTreeBySchieberVishkin(int n) {
        preOrder = new int[n];
        i = new int[n];
        head = new int[n];
        a = new int[n];
        parent = new int[n];
    }

    public LcaOnTreeBySchieberVishkin(List<? extends DirectedEdge>[] tree, IntPredicate isRoot) {
        this(tree.length);
        init(tree, isRoot);
    }

    public LcaOnTreeBySchieberVishkin(List<? extends DirectedEdge>[] tree, int root) {
        this(tree, i -> i == root);
    }

    private int enterIntoStrip(int x, int hz) {
        if (Integer.lowestOneBit(i[x]) == hz)
            return x;
        int hw = 1 << Log2.floorLog(a[x] & (hz - 1));
        return parent[head[i[x] & -hw | hw]];
    }

    public int lca(int x, int y) {
        int hb = i[x] == i[y] ? Integer.lowestOneBit(i[x]) : (1 << Log2.floorLog(i[x] ^ i[y]));
        int hz = Integer.lowestOneBit(a[x] & a[y] & -hb);
        int ex = enterIntoStrip(x, hz);
        int ey = enterIntoStrip(y, hz);
        return preOrder[ex] < preOrder[ey] ? ex : ey;
    }
}
