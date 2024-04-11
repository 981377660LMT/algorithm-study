package template.problem;

import template.graph.DirectedEdge;
import template.graph.UndirectedEdge;
import template.math.DigitUtils;

import java.util.Arrays;
import java.util.List;
import java.util.function.IntToLongFunction;


// https://taodaling.github.io/blog/2019/04/30/%E5%8C%BA%E9%97%B4%E6%93%8D%E4%BD%9C%E9%97%AE%E9%A2%98/
// 每次选择出发收益最大的顶点出发，k
// 次后得到的就是最大收益。假设v
// 是初始时出发收益最大的顶点，那么如果最后的方案中不选择v
// ，记方案中被搜集的顶点形成的树为T
// ，我们可以找到T
// 中v
// 的深度最小的祖先，回退之前T
// 下的任意一个特殊顶点带来的影响，并加入v
// ，可以证明这时候我们的收益是不会减少的。
// 首先我们可以利用所有顶点的dfs序构建线段树，之后利用线段树实现区间修改和弹出全局收益最大顶点的能力（即大根堆）。
/**
 * 给定一颗树，树根为1，树上有n个顶点，第i个顶点有一笔价值wi的财富。
 * 我们现在可以派出k个人，从k个顶点出发，向根出发，搜集路上所有的财富（同一个顶点的财富只会被搜集一次）。
 * 问我们最多可以得到多少总财富？其中−10^9≤wi≤10^9,1≤k≤n≤10^6.
 * <br>
 * O(n\log_2n)
 */
public class TreePickSpecialPointProblem {
    private List<? extends DirectedEdge>[] g;
    private long[] w;
    private long[] sumOfWeight;
    private int[] parents;
    private int n;
    private LongPriorityQueueBasedOnSegment segment;
    private int[] intervalL;
    private int[] intervalR;
    private int[] inverse;
    private boolean[] visited;
    private static final long INF = (long) 2e18;

    private int order;

    private void dfs(int root, int p, long sum) {
        intervalL[root] = order++;
        inverse[intervalL[root]] = root;
        sumOfWeight[root] = w[root] + sum;
        parents[root] = p;
        for (DirectedEdge e : g[root]) {
            if (e.to == p) {
                continue;
            }
            dfs(e.to, root, sumOfWeight[root]);
        }
        intervalR[root] = order - 1;
    }

    public TreePickSpecialPointProblem(int n) {
        this.n = n;
        sumOfWeight = new long[n];
        parents = new int[n];
        segment = new LongPriorityQueueBasedOnSegment(0, n);
        visited = new boolean[n];
        intervalL = new int[n];
        intervalR = new int[n];
        inverse = new int[n];
    }

    private void prepare(List<UndirectedEdge>[] g, long[] w, int root, boolean[] choice) {
        this.g = g;
        this.w = w;
        order = 0;
        dfs(root, -1, 0);
        segment.reset(0, n, i -> i < g.length && choice[inverse[i]] ? -sumOfWeight[inverse[i]] : INF);
        Arrays.fill(visited, 0, g.length, false);
    }

    private void visit(int root, int childL, int childR) {
        if (root == -1 || visited[root]) {
            return;
        }
        visited[root] = true;
        segment.update(intervalL[root], intervalR[root], 0, n, w[root]);
        visit(parents[root], intervalL[root], intervalR[root]);
    }

    /**
     * 从中恰好选择k个特殊点
     */
    public long apply(List<UndirectedEdge>[] g, long[] w, int root, int k, boolean[] choice, boolean[] select) {
        prepare(g, w, root, choice);
        Arrays.fill(select, 0, g.length, false);
        long ans = 0;
        for (int i = 0; i < k; i++) {
            ans -= segment.minimum;
            int node = inverse[segment.pop(0, n)];
            visit(node, 0, -1);
            select[node] = true;
        }
        return ans;
    }

    /**
     * 选择任意个特殊点
     */
    public long apply(List<UndirectedEdge>[] g, long[] w, int root, boolean[] choice, boolean[] select) {
        prepare(g, w, root, choice);
        Arrays.fill(select, 0, g.length, false);
        long ans = 0;
        while (segment.minimum < 0) {
            ans -= segment.minimum;
            int node = inverse[segment.pop(0, n)];
            visit(node, 0, -1);
            select[node] = true;
        }
        return ans;
    }

    /**
     * 选择任意个特殊点
     */
    public long applyNotMoreThan(List<UndirectedEdge>[] g, long[] w, int root, boolean[] choice, boolean[] select, int k) {
        prepare(g, w, root, choice);
        Arrays.fill(select, 0, g.length, false);
        long ans = 0;
        while (segment.minimum < 0 && k > 0) {
            k--;
            ans -= segment.minimum;
            int node = inverse[segment.pop(0, n)];
            visit(node, 0, -1);
            select[node] = true;
        }
        return ans;
    }
    
    // 可以利用所有顶点的dfs序构建线段树
    // 线段树实现区间修改和弹出全局收益最大顶点的能力（即大根堆）。
    private static class LongPriorityQueueBasedOnSegment implements Cloneable {
        private LongPriorityQueueBasedOnSegment left;
        private LongPriorityQueueBasedOnSegment right;
        private long minimum;
        private long dirty;

        public void modify(long x) {
            minimum += x;
            dirty += x;
        }

        public void pushUp() {
            minimum = Math.min(left.minimum, right.minimum);
        }

        public void pushDown() {
            if (dirty != 0) {
                left.modify(dirty);
                right.modify(dirty);
                dirty = 0;
            }
        }

        public LongPriorityQueueBasedOnSegment(int l, int r) {
            if (l < r) {
                int m = DigitUtils.floorAverage(l, r);
                left = new LongPriorityQueueBasedOnSegment(l, m);
                right = new LongPriorityQueueBasedOnSegment(m + 1, r);
                pushUp();
            } else {

            }
        }

        public void reset(int l, int r, IntToLongFunction function) {
            dirty = 0;
            if (l < r) {
                int m = DigitUtils.floorAverage(l, r);
                left.reset(l, m, function);
                right.reset(m + 1, r, function);
                pushUp();
            } else {
                minimum = function.applyAsLong(l);
            }
        }

        private boolean covered(int ll, int rr, int l, int r) {
            return ll <= l && rr >= r;
        }

        private boolean noIntersection(int ll, int rr, int l, int r) {
            return ll > r || rr < l;
        }

        public void update(int ll, int rr, int l, int r, long val) {
            if (noIntersection(ll, rr, l, r)) {
                return;
            }
            if (covered(ll, rr, l, r)) {
                modify(val);
                return;
            }
            int m = DigitUtils.floorAverage(l, r);
            pushDown();
            left.update(ll, rr, l, m, val);
            right.update(ll, rr, m + 1, r, val);
            pushUp();
        }

        public int pop(int l, int r) {
            if (l == r) {
                minimum = INF;
                return l;
            }
            int m = DigitUtils.floorAverage(l, r);
            int ans;
            pushDown();
            if (left.minimum == minimum) {
                ans = left.pop(l, m);
            } else {
                ans = right.pop(m + 1, r);
            }
            pushUp();
            return ans;
        }

        private LongPriorityQueueBasedOnSegment deepClone() {
            LongPriorityQueueBasedOnSegment seg = clone();
            if (seg.left != null) {
                seg.left = seg.left.deepClone();
            }
            if (seg.right != null) {
                seg.right = seg.right.deepClone();
            }
            return seg;
        }

        protected LongPriorityQueueBasedOnSegment clone() {
            try {
                return (LongPriorityQueueBasedOnSegment) super.clone();
            } catch (CloneNotSupportedException e) {
                throw new RuntimeException(e);
            }
        }

        private void toString(StringBuilder builder) {
            if (left == null && right == null) {
                builder.append(minimum).append(",");
                return;
            }
            pushDown();
            left.toString(builder);
            right.toString(builder);
        }

        public String toString() {
            StringBuilder builder = new StringBuilder();
            deepClone().toString(builder);
            if (builder.length() > 0) {
                builder.setLength(builder.length() - 1);
            }
            return builder.toString();
        }
    }
}
