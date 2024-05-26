package template.problem;

import template.algo.DoubleBinarySearch;
import template.math.DigitUtils;
import template.primitve.generated.datastructure.IntegerMinQueue;
import template.primitve.generated.datastructure.LongPreSum;
import template.utils.SequenceUtils;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;

public class IntervalPickProblem {
    public static class Interval {
        public final int l;
        public final int r;
        public boolean used;

        public Interval(int l, int r) {
            this.l = l;
            this.r = r;
        }

        @Override
        public String toString() {
            return String.format("(%d, %d)", l, r);
        }
    }

    public static Interval[] sortAndRemoveUnusedInterval(Interval[] intervals) {
        if (intervals.length == 0) {
            return intervals;
        }
        Arrays.sort(intervals, (a, b) -> Long.compare(a.l, b.l));
        int len = 1;
        for (int i = 1; i < intervals.length; i++) {
            if (intervals[i].l == intervals[len - 1].l) {
                if (intervals[i].r > intervals[len - 1].r) {
                    SequenceUtils.swap(intervals, i, len - 1);
                }
                continue;
            }
            if (intervals[i].r > intervals[len - 1].r) {
                SequenceUtils.swap(intervals, i, len);
                len++;
            }
        }
        return Arrays.copyOf(intervals, len);
    }

    /**
     * 有一个序列data[0], data[1], ... , data[n - 1], 以及m个区间
     * intervals[0], intervals[1], ..., intervals[m - 1].
     * <br>
     * 要求我们选择任意个区间，要求被选中的区间覆盖的元素的和最大。
     * <br>
     * 时间复杂度为O(n\log_2n)
     */
    public static long solve(long[] data, Interval[] intervals) {
        LongPreSum lps = new LongPreSum(i -> data[i], data.length);
        Arrays.sort(intervals, (a, b) -> a.l == b.l ? a.r - b.r : a.l - b.l);
        int n = intervals.length;
        long[] dp = new long[n];
        int[] last = new int[n];

        LongSegmentQuery query = new LongSegmentQuery();
        int m = data.length;
        LongSegment lower = new LongSegment(0, m);
        LongSegment upper = new LongSegment(0, m);

        for (int i = 0; i < n; i++) {
            Interval now = intervals[i];
            dp[i] = lps.intervalSum(now.l, now.r);
            last[i] = -1;

            query.reset();
            lower.query(0, now.l - 1, 0, m, query);
            if (query.val + lps.intervalSum(now.l, now.r) > dp[i]) {
                dp[i] = query.val + lps.intervalSum(now.l, now.r);
                last[i] = query.index;
            }
            query.reset();
            upper.query(now.l, now.r, 0, m, query);
            if (query.val + lps.prefix(now.r) > dp[i]) {
                dp[i] = query.val + lps.prefix(now.r);
                last[i] = query.index;
            }

            lower.update(now.r, now.r, 0, m, dp[i], i);
            upper.update(now.r, now.r, 0, m, dp[i] - lps.prefix(now.r), i);
        }

        int maxIndex = 0;
        for (int i = 0; i < n; i++) {
            if (dp[i] > dp[maxIndex]) {
                maxIndex = i;
            }
        }

        if (dp[maxIndex] < 0) {
            return 0;
        }
        int trace = maxIndex;
        while (trace >= 0) {
            intervals[trace].used = true;
            trace = last[trace];
        }

        return dp[maxIndex];
    }

    private static class WQSResult {
        double maxValue;
        int time;

        public WQSResult(double maxValue, int time) {
            this.maxValue = maxValue;
            this.time = time;
        }
    }

    private static class NormalWQSSolver {
        LongPreSum lps;
        int n;
        double[] dp;
        int[] time;
        DoubleSegmentQuery query;
        int m;
        DoubleSegment lower;
        DoubleSegment upper;
        long[] data;
        Interval[] intervals;

        public NormalWQSSolver(long[] data, Interval[] intervals) {
            this.data = data;
            this.intervals = intervals;
            lps = new LongPreSum(i -> data[i], data.length);
            n = intervals.length;
            dp = new double[n];
            time = new int[n];
            query = new DoubleSegmentQuery();
            m = data.length;
            lower = new DoubleSegment(0, m);
            upper = new DoubleSegment(0, m);
        }

        public WQSResult solve(double cost) {
            Arrays.fill(dp, 0);
            Arrays.fill(time, 0);
            lower.reset(0, m);
            upper.reset(0, m);

            for (int i = 0; i < n; i++) {
                Interval now = intervals[i];
                dp[i] = lps.intervalSum(now.l, now.r);
                time[i] = 1;

                query.reset();
                lower.query(0, now.l - 1, 0, m, query);
                if (query.val + lps.intervalSum(now.l, now.r) > dp[i]) {
                    dp[i] = query.val + lps.intervalSum(now.l, now.r);
                    time[i] = time[query.index] + 1;
                }
                query.reset();
                upper.query(now.l, now.r, 0, m, query);
                if (query.val + lps.prefix(now.r) > dp[i]) {
                    dp[i] = query.val + lps.prefix(now.r);
                    time[i] = time[query.index] + 1;
                }

                dp[i] -= cost;
                lower.update(now.r, now.r, 0, m, dp[i], i);
                upper.update(now.r, now.r, 0, m, dp[i] - lps.prefix(now.r), i);
            }

            int maxIndex = 0;
            for (int i = 0; i < n; i++) {
                if (dp[i] > dp[maxIndex]) {
                    maxIndex = i;
                }
            }

            if (dp[maxIndex] < 0) {
                return new WQSResult(0, 0);
            }
            return new WQSResult(dp[maxIndex], time[maxIndex]);
        }
    }


    /**
     * 有一个序列data[0], data[1], ... , data[n - 1], 以及m个区间
     * intervals[0], intervals[1], ..., intervals[m - 1].
     * <br>
     * 要求我们正好选择k个区间，要求被选中的区间覆盖的元素的和最大。
     * <br>
     */
    public static long solve(long[] data, Interval[] intervals, int k) {
        Arrays.sort(intervals, (a, b) -> a.l == b.l ? a.r - b.r : a.l - b.l);
        NormalWQSSolver solver = new NormalWQSSolver(data, intervals);
        DoubleBinarySearch dbs = new DoubleBinarySearch(1e-12, 1e-12) {
            @Override
            public boolean check(double mid) {
                return solver.solve(mid).time <= k;
            }
        };

        long sum = 0;
        for (long x : data) {
            sum += Math.abs(x);
        }
        double cost = dbs.binarySearch(-sum, sum);
        long ans = DigitUtils.round(solver.solve(cost).maxValue + k * cost);
        return ans;
    }

    private static class LongSegmentQuery {
        int index;
        long val;

        public void reset() {
            index = -1;
            val = -LongSegment.inf;
        }

        public void update(int index, long val) {
            if (this.val < val) {
                this.index = index;
                this.val = val;
            }
        }
    }

    private static class DoubleSegmentQuery {
        int index;
        double val;

        public void reset() {
            index = -1;
            val = -DoubleSegment.inf;
        }

        public void update(int index, double val) {
            if (this.val < val) {
                this.index = index;
                this.val = val;
            }
        }
    }

    private static class LongSegment implements Cloneable {
        private static long inf = (long) 2e18;

        private LongSegment left;
        private LongSegment right;
        private long val = -inf;
        private int index = -1;

        public void pushUp() {
            val = Math.max(left.val, right.val);
            if (val == left.val) {
                index = left.index;
            } else {
                index = right.index;
            }
        }

        public void pushDown() {
        }

        public LongSegment(int l, int r) {
            if (l < r) {
                int m = DigitUtils.floorAverage(l, r);
                left = new LongSegment(l, m);
                right = new LongSegment(m + 1, r);
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

        public void update(int ll, int rr, int l, int r, long x, int index) {
            if (noIntersection(ll, rr, l, r)) {
                return;
            }
            if (covered(ll, rr, l, r)) {
                if (x > val) {
                    this.val = x;
                    this.index = index;
                }
                return;
            }
            pushDown();
            int m = DigitUtils.floorAverage(l, r);
            left.update(ll, rr, l, m, x, index);
            right.update(ll, rr, m + 1, r, x, index);
            pushUp();
        }

        public void query(int ll, int rr, int l, int r, LongSegmentQuery query) {
            if (noIntersection(ll, rr, l, r)) {
                return;
            }
            if (covered(ll, rr, l, r)) {
                query.update(index, val);
                return;
            }
            pushDown();
            int m = DigitUtils.floorAverage(l, r);
            left.query(ll, rr, l, m, query);
            right.query(ll, rr, m + 1, r, query);
        }
    }

    private static class DoubleSegment implements Cloneable {
        private static double inf = 1e50;

        private DoubleSegment left;
        private DoubleSegment right;
        private double val = -inf;
        private int index = -1;

        public void pushUp() {
            val = Math.max(left.val, right.val);
            if (val == left.val) {
                index = left.index;
            } else {
                index = right.index;
            }
        }

        public void pushDown() {
        }

        public DoubleSegment(int l, int r) {
            if (l < r) {
                int m = DigitUtils.floorAverage(l, r);
                left = new DoubleSegment(l, m);
                right = new DoubleSegment(m + 1, r);
                pushUp();
            } else {
            }
        }

        public void reset(int l, int r) {
            if (l < r) {
                int m = DigitUtils.floorAverage(l, r);
                left.reset(l, m);
                right.reset(m + 1, r);
                pushUp();
            } else {
                val = -inf;
                index = -1;
            }
        }

        private boolean covered(int ll, int rr, int l, int r) {
            return ll <= l && rr >= r;
        }

        private boolean noIntersection(int ll, int rr, int l, int r) {
            return ll > r || rr < l;
        }

        public void update(int ll, int rr, int l, int r, double x, int index) {
            if (noIntersection(ll, rr, l, r)) {
                return;
            }
            if (covered(ll, rr, l, r)) {
                if (x > val) {
                    this.val = x;
                    this.index = index;
                }
                return;
            }
            pushDown();
            int m = DigitUtils.floorAverage(l, r);
            left.update(ll, rr, l, m, x, index);
            right.update(ll, rr, m + 1, r, x, index);
            pushUp();
        }

        public void query(int ll, int rr, int l, int r, DoubleSegmentQuery query) {
            if (noIntersection(ll, rr, l, r)) {
                return;
            }
            if (covered(ll, rr, l, r)) {
                query.update(index, val);
                return;
            }
            pushDown();
            int m = DigitUtils.floorAverage(l, r);
            left.query(ll, rr, l, m, query);
            right.query(ll, rr, m + 1, r, query);
        }
    }

    private static class NonNegativeWQSSolver {
        LongPreSum lps;
        double[] dp;
        int[] time;
        IntegerMinQueue left;
        IntegerMinQueue middle;
        int n;
        long[] data;
        Interval[] intervals;

        public NonNegativeWQSSolver(long[] data, Interval[] intervals) {
            this.n = intervals.length;
            this.data = data;
            this.intervals = intervals;
            lps = new LongPreSum(i -> data[i], data.length);
            dp = new double[n + 1];
            time = new int[n + 1];
            left = new IntegerMinQueue(n, (a, b) -> -Double.compare(dp[a], dp[b]));
            middle = new IntegerMinQueue(n, (a, b) -> -Double.compare(dp[a] - lps.prefix(intervals[a - 1].r), dp[b] - lps.prefix(intervals[b - 1].r)));
        }

        /**
         * 有一个序列data[0], data[1], ... , data[n - 1], 以及m个区间
         * intervals[0], intervals[1], ..., intervals[m - 1].且序列元素非负。
         * <br>
         * 要求我们选择任意个区间，要求被选中的区间覆盖的元素的和最大。
         * <br>
         * 时间复杂度为O(n)
         */
        public WQSResult solveNonNegative(double cost) {
            Arrays.fill(dp, 0);
            Arrays.fill(time, 0);
            left.clear();
            middle.clear();

            left.addLast(0);
            for (int i = 1; i <= n; i++) {
                IntervalPickProblem.Interval now = intervals[i - 1];
                while (!middle.isEmpty() && intervals[middle.peek() - 1].r < now.l) {
                    left.addLast(middle.removeFirst());
                }
                dp[i] = -1e50;
                if (dp[i] < dp[left.min()] + lps.intervalSum(now.l, now.r)) {
                    dp[i] = dp[left.min()] + lps.intervalSum(now.l, now.r);
                    time[i] = time[left.min()] + 1;
                }
                if (!middle.isEmpty() && dp[i] < dp[middle.min()] + lps.prefix(now.r) - lps.prefix(intervals[middle.min() - 1].r)) {
                    dp[i] = dp[middle.min()] + lps.prefix(now.r) - lps.prefix(intervals[middle.min() - 1].r);
                    time[i] = time[middle.min()] + 1;
                }
                dp[i] -= cost;
                middle.addLast(i);
            }

            int maxIndex = 0;
            for (int i = 1; i <= n; i++) {
                if (dp[i] > dp[maxIndex]) {
                    maxIndex = i;
                }
            }

            return new WQSResult(dp[maxIndex], time[maxIndex]);
        }
    }


    /**
     * 有一个序列data[0], data[1], ... , data[n - 1], 以及m个区间
     * intervals[0], intervals[1], ..., intervals[m - 1].
     * <br>
     * 要求所有元素非负，且intervals[i + 1].r > intervals[i - 1].r, intervals[i + 1].l > intervals[i - 1].l;
     * <br>
     * 要求我们正好选择k个区间，要求被选中的区间覆盖的元素的和最大。
     * <br>
     * 区间需要满足intervals[i].l < intervals[i + 1].l且intervals[i].r < intervals[i + 1].r
     * <br>
     */
    public static long solveNonNegative(long[] data, IntervalPickProblem.Interval[] intervals, int k, int round) {
        if (round <= 0) {
            throw new IllegalArgumentException();
        }
        long sum = Arrays.stream(data).sum();
        NonNegativeWQSSolver solver = new NonNegativeWQSSolver(data, intervals);
        WQSResult result = null;

        double l = 0;
        double r = sum;
        double m = 0;
        while (round-- > 0) {
            m = (l + r) / 2;
            result = solver.solveNonNegative(m);
            if (result.time > k) {
                l = m;
            } else if (result.time < k) {
                r = m;
            } else {
                break;
            }
        }
        long ans = DigitUtils.round(result.maxValue + k * m);
        return ans;
    }

}
