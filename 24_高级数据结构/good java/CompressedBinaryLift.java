package template.graph;

import java.util.Arrays;
import java.util.function.IntPredicate;


/**
 * 更加优秀的树上倍增
 * https://codeforces.com/blog/entry/74847
 * <p>
 * O(log_2n) for all operation
 */
public class CompressedBinaryLift implements LcaOnTree, KthAncestor {
    ParentOnTree pot;
    DepthOnTree dot;
    int[] jump;

    private void consider(int root) {
        if (root == -1 || jump[root] != -1) {
            return;
        }
        int p = pot.parent(root);
        consider(p);
        addLeaf(root, p);
    }

    public CompressedBinaryLift(int n, DepthOnTree dot, ParentOnTree pot) {
        this.dot = dot;
        this.pot = pot;
        jump = new int[n];
        Arrays.fill(jump, -1);
        for (int i = 0; i < n; i++) {
            consider(i);
        }
    }

    private void addLeaf(int leaf, int pId) {
        if (pId == -1) {
            jump[leaf] = leaf;
        } else if (dot.depth(pId) - dot.depth(jump[pId]) == dot.depth(jump[pId]) - dot.depth(jump[jump[pId]])) {
            jump[leaf] = jump[jump[pId]];
        } else {
            jump[leaf] = pId;
        }
    }

    public int firstTrue(int node, IntPredicate predicate) {
        while (!predicate.test(node)) {
            if (predicate.test(jump[node])) {
                node = pot.parent(node);
            } else {
                if (node == jump[node]) {
                    return -1;
                }
                node = jump[node];
            }
        }
        return node;
    }

    public int lastTrue(int node, IntPredicate predicate) {
        if (!predicate.test(node)) {
            return -1;
        }
        while (true) {
            if (predicate.test(jump[node])) {
                if (node == jump[node]) {
                    return node;
                }
                node = jump[node];
            } else if (predicate.test(pot.parent(node))) {
                node = pot.parent(node);
            } else {
                return node;
            }
        }
    }

    public int kthAncestor(int node, int k) {
        int targetDepth = dot.depth(node) - k;
        return firstTrue(node, i -> dot.depth(i) <= targetDepth);
    }

    public int lca(int a, int b) {
        if (dot.depth(a) > dot.depth(b)) {
            a = kthAncestor(a, dot.depth(a) - dot.depth(b));
        } else if (dot.depth(a) < dot.depth(b)) {
            b = kthAncestor(b, dot.depth(b) - dot.depth(a));
        }
        while (a != b) {
            if (jump[a] == jump[b]) {
                a = pot.parent(a);
                b = pot.parent(b);
            } else {
                a = jump[a];
                b = jump[b];
            }
        }
        return a;
    }

}
