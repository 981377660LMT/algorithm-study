package template.problem;

import template.math.DigitUtils;
import template.primitve.generated.datastructure.IntegerArrayList;
import template.utils.SortUtils;

import java.util.Arrays;
import java.util.stream.IntStream;

// 区间绝对值差的最小值.
public class IntervalNearestValuePair {
    static long inf = (long) 2e18;

    /**
     * <p>
     * Find the difference between nearest pair in interval [l, r],
     * for example 3 is the ans in [3, 11, -100, 8]
     * </p>
     * <pre>
     * time complexity: O(nlog_2nlog_2M+qlog_2n)
     * space complexity: O(n\log_2M)
     * note: M is the max absolute value in a
     * </pre>
     *
     * @param a
     * @param qs
     */
    public static void solve(long[] a, Query[] qs) {
        for (Query q : qs) {
            q.ans = inf;
        }
        Arrays.sort(qs, (x, y) -> Integer.compare(x.r, y.r));
        solve0(a, qs);
        long max = (long) 1e18;
        for (int i = 0; i < a.length; i++) {
            a[i] = max - a[i];
        }
        solve0(a, qs);
    }

    private static void solve0(long[] a, Query[] qs) {
        int n = a.length;
        IntegerArrayList[] ceil = new IntegerArrayList[n];
        IntegerArrayList buf = new IntegerArrayList(60);
        for (int i = 0; i < n; i++) {
            ceil[i] = new IntegerArrayList();
        }
        int[] indices = IntStream.range(0, n).toArray();
        SortUtils.quickSort(indices, (i, j) -> a[i] == a[j] ? Integer.compare(i, j) : -Long.compare(a[i], a[j]), 0, n);

        SegmentBS segBS = new SegmentBS(0, n - 1);
        segBS.reset(0, n - 1, inf);
        for (int i : indices) {
            long threshold = inf - 1;
            int j = segBS.query(0, i, 0, n - 1, threshold);
            buf.clear();
            while (j != -1) {
                buf.add(j);
                //x - a[i] < a[j] - x
                //2x < a[j] + a[i]
                threshold = DigitUtils.maximumIntegerLessThanDiv(a[j] + a[i], 2);
                j = segBS.query(0, j - 1, 0, n - 1, threshold);
            }
            ceil[i].addAll(buf);
            segBS.update(i, i, 0, n - 1, a[i]);
        }

        Segment segtree = new Segment(0, n - 1);
        int l = 0;
        for (int i = 0; i < n; i++) {
            for (int t = 0; t < ceil[i].size(); t++) {
                int j = ceil[i].get(t);
                segtree.update(0, j, 0, n - 1, a[j] - a[i]);
            }
            while (l < qs.length && qs[l].r == i) {
                qs[l].ans = Math.min(qs[l].ans, segtree.query(qs[l].l, qs[l].l, 0, n - 1));
                l++;
            }
        }
    }

    public static class Query {
        public int l;
        public int r;
        public long ans;

        public Query(int l, int r) {
            this.l = l;
            this.r = r;
        }
    }

    static class Segment implements Cloneable {
        private Segment left;
        private Segment right;
        private long min;

        private void modify(long x) {
            min = Math.min(min, x);
        }

        public void pushUp() {
        }

        public void pushDown() {
            left.modify(min);
            right.modify(min);
            min = inf;
        }

        public Segment(int l, int r) {
            min = inf;
            if (l < r) {
                int m = DigitUtils.floorAverage(l, r);
                left = new Segment(l, m);
                right = new Segment(m + 1, r);
                pushUp();
            } else {
            }
        }

        private boolean covered(int ll, int rr, int l, int r) {
            return ll <= l && rr >= r;
        }

        private boolean noIntersection(int ll, int rr, int l, int r) {
            return ll > r || rr < l;
        }

        public void update(int ll, int rr, int l, int r, long x) {
            if (noIntersection(ll, rr, l, r)) {
                return;
            }
            if (covered(ll, rr, l, r)) {
                modify(x);
                return;
            }
            pushDown();
            int m = DigitUtils.floorAverage(l, r);
            left.update(ll, rr, l, m, x);
            right.update(ll, rr, m + 1, r, x);
            pushUp();
        }

        public long query(int ll, int rr, int l, int r) {
            if (noIntersection(ll, rr, l, r)) {
                return inf;
            }
            if (covered(ll, rr, l, r)) {
                return min;
            }
            pushDown();
            int m = DigitUtils.floorAverage(l, r);
            return Math.min(left.query(ll, rr, l, m),
                    right.query(ll, rr, m + 1, r));
        }

        private Segment deepClone() {
            Segment seg = clone();
            if (seg.left != null) {
                seg.left = seg.left.deepClone();
            }
            if (seg.right != null) {
                seg.right = seg.right.deepClone();
            }
            return seg;
        }

        protected Segment clone() {
            try {
                return (Segment) super.clone();
            } catch (CloneNotSupportedException e) {
                throw new RuntimeException(e);
            }
        }

        private void toString(StringBuilder builder) {
            if (left == null && right == null) {
                builder.append(min).append(",");
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

    static class SegmentBS implements Cloneable {
        private SegmentBS left;
        private SegmentBS right;
        private long min;

        private void modify(long x) {
            min = x;
        }

        public void pushUp() {
            min = Math.min(left.min, right.min);
        }

        public void pushDown() {
        }

        public SegmentBS(int l, int r) {
            if (l < r) {
                int m = DigitUtils.floorAverage(l, r);
                left = new SegmentBS(l, m);
                right = new SegmentBS(m + 1, r);
                pushUp();
            } else {

            }
        }

        public void reset(int l, int r, long x) {
            if (l < r) {
                int m = DigitUtils.floorAverage(l, r);
                left.reset(l, m, x);
                right.reset(m + 1, r, x);
                pushUp();
            } else {
                min = x;
            }
        }

        private boolean covered(int ll, int rr, int l, int r) {
            return ll <= l && rr >= r;
        }

        private boolean noIntersection(int ll, int rr, int l, int r) {
            return ll > r || rr < l;
        }

        public void update(int ll, int rr, int l, int r, long x) {
            if (noIntersection(ll, rr, l, r)) {
                return;
            }
            if (covered(ll, rr, l, r)) {
                modify(x);
                return;
            }
            pushDown();
            int m = DigitUtils.floorAverage(l, r);
            left.update(ll, rr, l, m, x);
            right.update(ll, rr, m + 1, r, x);
            pushUp();
        }

        public int query(int ll, int rr, int l, int r, long threshold) {
            if (noIntersection(ll, rr, l, r) || min > threshold) {
                return -1;
            }
            if (covered(ll, rr, l, r) && l == r) {
                return l;
            }
            pushDown();
            int m = DigitUtils.floorAverage(l, r);
            int ans = right.query(ll, rr, m + 1, r, threshold);
            if (ans == -1) {
                ans = left.query(ll, rr, l, m, threshold);
            }
            return ans;
        }

        private SegmentBS deepClone() {
            SegmentBS seg = clone();
            if (seg.left != null) {
                seg.left = seg.left.deepClone();
            }
            if (seg.right != null) {
                seg.right = seg.right.deepClone();
            }
            return seg;
        }

        protected SegmentBS clone() {
            try {
                return (SegmentBS) super.clone();
            } catch (CloneNotSupportedException e) {
                throw new RuntimeException(e);
            }
        }

        private void toString(StringBuilder builder) {
            if (left == null && right == null) {
                builder.append(min).append(",");
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
