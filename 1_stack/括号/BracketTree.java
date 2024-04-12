// 1.区间最长合法括号
// 2.包含[l,r]的最短合法括号区间

package template.datastructure;


import template.graph.*;
import template.primitve.generated.datastructure.IntToIntegerFunction;
import template.primitve.generated.datastructure.IntegerDequeImpl;
import template.primitve.generated.datastructure.IntegerIterator;
import template.utils.Sum;

import java.util.ArrayList;
import java.util.List;

/**
 * It maintain the bracket sequence and help you to find the longest/number of correct substring.
 */
public class BracketTree {
    int L;
    int R;
    Node root;

    /**
     * 1 for ( and -1 for )
     * <br/>
     * O(n)
     */
    public BracketTree(int n, IntToIntegerFunction func) {
        IntegerDequeImpl dq = new IntegerDequeImpl(2 * n + 2);
        int ps = 0;
        L = 0;
        for (int i = 0; i < n; i++) {
            int v = func.apply(i);
            ps += v;
            dq.addLast(v);
            if (ps < 0) {
                ps++;
                dq.addFirst(1);
                L++;
            }
        }
        while (ps > 0) {
            ps--;
            dq.addLast(-1);
        }
        dq.addFirst(1);
        dq.addLast(-1);
        L++;

        R = L + n - 1;
        m = dq.size();
        seq = new int[m];
        int wpos = 0;
        for (IntegerIterator iterator = dq.iterator(); iterator.hasNext(); ) {
            seq[wpos++] = iterator.next();
        }
        all = new ArrayList<>(seq.length / 2);
        rpos = 0;

        build(null);
        root = all.get(0);
        open = new Node[m];
        close = new Node[m];
        for (Node node : all) {
            open[node.l] = node;
            close[node.r] = node;
        }

        assert all.size() == m / 2;
        ParentOnTree pot = new ParentOnTreeByFunction(all.size(), i -> {
            Node node = all.get(i);
            if (node.p == null) {
                return -1;
            }
            return node.p.id;
        });
        dot = new DepthOnTreeByParent(all.size(), pot);
        bl = new CompressedBinaryLiftWithAttachment<>(all.size(), dot, pot, SumImpl::new, i -> all.get(i).sum);
    }

    /**
     * res[0] for longest, res[1] for count
     * <br>
     * O(\log_2n)
     *
     * @param l
     * @param r
     * @param res
     * @return
     */
    public long[] query(int l, int r, long[] res) {
        if (res == null) {
            res = new long[2];
        } else {
            res[0] = res[1] = 0;
        }
        l += L;
        r += L;
        Node front = open[l - 1];
        if (front == null) {
            front = close[l - 1];
        }
        Node back = close[r + 1];
        if (back == null) {
            back = open[r + 1];
        }
        int lcaId = bl.lca(front.id, back.id);
        Node lca = all.get(lcaId);

        buf0.init();
        buf1.init();

        if (front == lca && back == lca) {
            res[0] = Math.max(res[0], (lca.r - 1) - (lca.l + 1) + 1);
            if (front == open[l - 1] && back == close[r + 1]) {
                res[1] += lca.cover - 1;
            }
            return res;
        } else if (front == lca) {
            bl.kthAncestor(back.id, dot.depth(back.id) - dot.depth(lca.id) - 1, buf1);
            res[0] = buf1.maxRight;
            res[1] = buf1.sumRight;
            if (back == close[r + 1]) {
                res[1] += back.cover - 1;
                res[0] = Math.max(res[0], back.length() - 2);
            }
            return res;
        } else if (back == lca) {
            bl.kthAncestor(front.id, dot.depth(front.id) - dot.depth(lca.id) - 1, buf0);
            res[0] = buf0.maxLeft;
            res[1] = buf0.sumLeft;
            if (front == open[l - 1]) {
                res[1] += front.cover - 1;
                res[0] = Math.max(res[0], front.length() - 2);
            }
            return res;
        }
        if (front == open[l - 1]) {
            res[1] += front.cover - 1;
            res[0] = Math.max(res[0], front.length() - 2);
        }
        if (back == close[r + 1]) {
            res[1] += back.cover - 1;
            res[0] = Math.max(res[0], back.length() - 2);
        }

        int fdelta;
        if ((fdelta = dot.depth(front.id) - dot.depth(lca.id)) >= 2) {
            front = all.get(bl.kthAncestor(front.id, fdelta - 2, buf0)).p;
        }
        int bdelta;
        if ((bdelta = dot.depth(back.id) - dot.depth(lca.id)) >= 2) {
            back = all.get(bl.kthAncestor(back.id, bdelta - 2, buf1)).p;
        }
        res[0] = Math.max(res[0], buf0.maxLeft);
        res[0] = Math.max(res[0], buf1.maxRight);
        res[1] += buf0.sumLeft + buf1.sumRight;

        res[0] = Math.max(res[0], (back.l - 1) - (front.r + 1) + 1);
        res[1] += choose2(back.pIndex - front.pIndex - 1);
        res[1] += interval(lca.ps, front.pIndex + 1, back.pIndex - 1);
        return res;
    }

    /**
     * find an valid interval [L, R] that contains [l, r] and it's the smallest one, null for not exist
     * <br>
     * O(\log_2n)
     */
    public int[] minValidSequence(int l, int r, int[] res) {
        if (res == null) {
            res = new int[2];
        } else {
            res[0] = res[1] = 0;
        }
        l += L;
        r += L;
        Node front = open[l];
        if (front == null) {
            front = close[l];
        }
        Node back = close[r];
        if (back == null) {
            back = open[r];
        }
        int lcaId = bl.lca(front.id, back.id);
        Node lca = all.get(lcaId);
        if (front == lca || back == lca) {
            res[0] = lca.l;
            res[1] = lca.r;
        } else {
            front = all.get(bl.kthAncestor(front.id, dot.depth(front.id) - dot.depth(lca.id) - 1));
            back = all.get(bl.kthAncestor(back.id, dot.depth(back.id) - dot.depth(lca.id) - 1));
            res[0] = front.l;
            res[1] = back.r;
        }
        if (res[0] < L || res[1] > R) {
            return null;
        }
        res[0] -= L;
        res[1] -= L;
        return res;
    }

    private long interval(long[] ps, int l, int r) {
        if (l > r) {
            return 0;
        }
        long ans = ps[r];
        if (l > 0) {
            ans -= ps[l - 1];
        }
        return ans;
    }

    SumImpl buf0 = new SumImpl();
    SumImpl buf1 = new SumImpl();

    DepthOnTree dot;
    int m;
    List<Node> all;
    Node[] open;
    Node[] close;
    CompressedBinaryLiftWithAttachment<SumImpl> bl;

    static class Node {
        List<Node> adj = new ArrayList<>();
        int pIndex;
        int l;
        int r;
        int id;
        Node p;
        long[] ps;

        SumImpl sum = new SumImpl();
        long cover;

        public int length() {
            return r - l + 1;
        }
    }

    int[] seq;
    int rpos;

    Node newNode() {
        Node node = new Node();
        node.id = all.size();
        all.add(node);
        return node;
    }

    void build(Node p) {
        Node root = newNode();
        root.p = p;
        root.l = rpos++;
        while (seq[rpos] != -1) {
            build(root);
        }
        root.r = rpos++;
        if (p != null) {
            root.pIndex = p.adj.size();
            p.adj.add(root);
        }

        int n = root.adj.size();
        root.ps = new long[n];
        long ps = 0;
        for (int i = 0; i < n; i++) {
            Node cur = root.adj.get(i);
            cur.sum.maxRight = Math.max(0, cur.l - 1 - (root.l + 1) + 1);
            cur.sum.sumRight += choose2(i);
            cur.sum.sumRight += ps;
            ps += cur.cover;
        }
        ps = 0;
        for (int i = n - 1; i >= 0; i--) {
            Node cur = root.adj.get(i);
            cur.sum.maxLeft = Math.max(0, (root.r - 1) - (cur.r + 1) + 1);
            cur.sum.sumLeft += choose2(n - 1 - i);
            cur.sum.sumLeft += ps;
            ps += cur.cover;
        }

        root.cover = 1 + choose2(n);
        for (int i = 0; i < n; i++) {
            Node node = root.adj.get(i);
            root.cover += node.cover;
            root.ps[i] = node.cover;
            if (i > 0) {
                root.ps[i] += root.ps[i - 1];
            }
        }
    }

    private long choose2(long n) {
        return n * (n - 1) / 2;
    }

    static class SumImpl implements Sum<SumImpl> {
        long maxLeft;
        long maxRight;
        long sumLeft;
        long sumRight;

        @Override
        public void add(SumImpl right) {
            maxLeft = Math.max(maxLeft, right.maxLeft);
            maxRight = Math.max(maxRight, right.maxRight);
            sumLeft += right.sumLeft;
            sumRight += right.sumRight;
        }

        @Override
        public void copy(SumImpl right) {
            maxLeft = right.maxLeft;
            maxRight = right.maxRight;
            sumLeft = right.sumLeft;
            sumRight = right.sumRight;
        }

        @Override
        public SumImpl clone() {
            SumImpl ans = new SumImpl();
            ans.copy(this);
            return ans;
        }

        void init() {
            maxLeft = maxRight = sumLeft = sumRight = 0;
        }

        @Override
        public String toString() {
            return "SumImpl{" +
                    "maxLeft=" + maxLeft +
                    ", maxRight=" + maxRight +
                    ", sumLeft=" + sumLeft +
                    ", sumRight=" + sumRight +
                    '}';
        }
    }
}
