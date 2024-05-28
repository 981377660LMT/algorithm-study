package template.graph;

import java.util.Comparator;
import java.util.List;
import java.util.TreeSet;

/**
 * maintain set S, O(log_2n) for following operation:
 * <ul>
 *    <li>add v to S</li>
 *    <li>remove v from S</li>
 *    <li>S = emptyset</li>
 *    <li>calculate the size of minimum spanning graph contains all vertexes in S</li>
 * </ul>
 */
public class MinSpanningGraphOnTree {
    private List<? extends DirectedEdge>[] g;
    private int[] depth;
    private LcaOnTreeBySchieberVishkin lca;
    private int[] index2dfn;
    private int[] dfn2index;
    private int order;
    private TreeSet<Integer> S = new TreeSet<>(Comparator.comparingInt(x -> index2dfn[x]));
    private int size;

    void dfs(int root, int p) {
        index2dfn[root] = order++;
        dfn2index[index2dfn[root]] = root;
        depth[root] = p == -1 ? 0 : depth[p] + 1;
        for (DirectedEdge e : g[root]) {
            if (e.to == p) {
                continue;
            }
            dfs(e.to, root);
        }
    }

    public int dist(int a, int b) {
        int c = lca.lca(a, b);
        return depth[a] + depth[b] - depth[c] * 2;
    }

    public MinSpanningGraphOnTree(int n) {
        depth = new int[n];
        index2dfn = new int[n];
        dfn2index = new int[n];
        lca = new LcaOnTreeBySchieberVishkin(n);
    }

    public void init(List<? extends DirectedEdge>[] g) {
        this.g = g;
        lca.init(g, 0);
        order = 0;
        dfs(0, -1);
        clear();
    }

    public void clear() {
        S.clear();
        size = 0;
    }

    public boolean contains(int v) {
        return S.contains(v);
    }

    public void add(int v) {
        assert !S.contains(v);
        Integer floor = S.floor(v);
        Integer ceil = S.ceiling(v);
        if (floor != null) {
            size += dist(floor, v);
        }
        if (ceil != null) {
            size += dist(ceil, v);
        }
        if (floor != null && ceil != null) {
            size -= dist(floor, ceil);
        }
        S.add(v);
    }

    public void remove(int v) {
        assert S.contains(v);
        S.remove(v);
        Integer floor = S.floor(v);
        Integer ceil = S.ceiling(v);
        if (floor != null) {
            size -= dist(floor, v);
        }
        if (ceil != null) {
            size -= dist(ceil, v);
        }
        if (floor != null && ceil != null) {
            size += dist(floor, ceil);
        }
    }

    public int size() {
        return S.size();
    }

    public int minSpanningGraph() {
        if (size() == 0) {
            return 0;
        }
        int ans = size + dist(S.first(), S.last());
        assert ans % 2 == 0;
        return ans / 2 + 1;
    }

    @Override
    public String toString() {
        return S.toString() + "=" + minSpanningGraph();
    }
}
