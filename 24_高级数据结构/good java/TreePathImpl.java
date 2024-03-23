package template.graph;

public class TreePathImpl implements TreePath {
    DepthOnTree dot;
    KthAncestor ancestor;
    LcaOnTree lca;
    int a;
    int b;
    int c;

    public TreePathImpl(DepthOnTree dot, KthAncestor kthAncestor, LcaOnTree lca) {
        this.dot = dot;
        this.ancestor = kthAncestor;
        this.lca = lca;
    }

    public void init(int a, int b) {
        this.a = a;
        this.b = b;
        c = lca.lca(a, b);
    }

    public int length() {
        return dot.depth(a) + dot.depth(b) - 2 * dot.depth(c);
    }

    /**
     * a is 0-th, k <= length()
     * <p>
     * O(log_2n)
     */
    public int kthNodeOnPath(int k) {
        if (k <= dot.depth(a) - dot.depth(c)) {
            return ancestor.kthAncestor(a, k);
        }
        return ancestor.kthAncestor(b, length() - k);
    }

    @Override
    public boolean onPath(int u) {
        return lca.lca(u, c) == c && (lca.lca(u, a) == u || lca.lca(u, b) == u);
    }
}
