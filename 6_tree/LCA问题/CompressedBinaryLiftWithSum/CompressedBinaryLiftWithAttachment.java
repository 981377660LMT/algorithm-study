package template.graph;

import java.util.Arrays;
import java.util.function.IntFunction;
import java.util.function.IntPredicate;
import java.util.function.Supplier;

public interface Sum<S> {
  void add(S other);

  /**
   * copy s.data
   *
   * @param other
   */
  void copy(S other);

  S clone();
}

// 用于倍增结构优化建图、查询路径聚合值.
// https://taodaling.github.io/blog/2020/03/18/binary-lifting/
public class CompressedBinaryLiftWithAttachment<S extends Sum<S>> implements LcaOnTree, KthAncestor {
    ParentOnTree pot;
    DepthOnTree dot;
    int[] jump;
    Object[] attachments;
    Object[] singles;

    private void consider(int root) {
        if (root == -1 || jump[root] != -1) {
            return;
        }
        int p = pot.parent(root);
        consider(p);
        addLeaf(root, p);
    }

    public CompressedBinaryLiftWithAttachment(int n, DepthOnTree dot, ParentOnTree pot, Supplier<S> supplier, IntFunction<S> single) {
        this.dot = dot;
        this.pot = pot;
        jump = new int[n];
        attachments = new Object[n];
        singles = new Object[n];
        for (int i = 0; i < n; i++) {
            attachments[i] = supplier.get();
            singles[i] = single.apply(i);
        }
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
            ((S) attachments[leaf]).copy((S) singles[leaf]);
            ((S) attachments[leaf]).add((S) attachments[pId]);
            ((S) attachments[leaf]).add((S) attachments[jump[pId]]);
        } else {
            jump[leaf] = pId;
            ((S) attachments[leaf]).copy((S) singles[leaf]);
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
    public int firstTrue(int node, IntPredicate predicate, S sum) {
        while (!predicate.test(node)) {
            if (predicate.test(jump[node])) {
                sum.add((S) singles[node]);
                node = pot.parent(node);
            } else {
                sum.add((S) attachments[node]);
                if (node == jump[node]) {
                    sum.add((S) singles[node]);
                    return -1;
                }
                node = jump[node];
            }
        }
        sum.add((S) singles[node]);
        return node;
    }



    public int lastTrue(int node, IntPredicate predicate, S sum) {
        if (!predicate.test(node)) {
            return -1;
        }
        while (true) {
            if (predicate.test(jump[node])) {
                if (node == jump[node]) {
                    sum.add((S) singles[node]);
                    return node;
                }
                sum.add((S) attachments[node]);
                node = jump[node];
            } else if (predicate.test(pot.parent(node))) {
                sum.add((S) singles[node]);
                node = pot.parent(node);
            } else {
                sum.add((S) singles[node]);
                return node;
            }
        }
    }

    public int kthAncestor(int node, int k, S s) {
        int targetDepth = dot.depth(node) - k;
        return firstTrue(node, i -> dot.depth(i) <= targetDepth, s);
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
