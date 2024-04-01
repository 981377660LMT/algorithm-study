package template.problem;

import template.datastructure.RangeTree;
import template.math.DigitUtils;
import template.primitve.generated.datastructure.IntToLongFunction;

import java.util.*;


public class MinimumSegmentSubtract {
//    public static interface IntervalConsumer {
//        void consume(int l, int r, long time);
//    }

    /**
     * <p>
     * 给定数组a，允许每次选择一段连续的元素，将每个元素减少1，且要求数组的长度len满足l<=len<=r。要求经过若干次操作后，每个元素都变成0。
     * 判断是否有解，有解则输出最小操作次数，否则输出-1。
     * </p>
     * <p>
     * 要求必须满足r>=2l-1且l>=1
     * </p>
     * <p>
     * O(5n) 时空复杂度
     * </p>
     * <p>
     * 最多n次操作
     * </p>
     */
    public static long solve(int n, int l, int r, IntToLongFunction a) {
        assert r == n || r >= 2 * l - 1;
        assert l > 0;
        assert n >= 0;
        if (n == 0) {
            return 0;
        }
        r = Math.min(r, n);
        long buildingSize = 0;
        RangeMultiDeque dq = new RangeMultiDeque(r);
        List<Interval> all = new ArrayList<>(n);
        Deque<Interval> building = new ArrayDeque<>(n);
        for (int i = 0; i < n; i++) {
            //enough
            while (!building.isEmpty() && i - building.peekFirst().begin == l) {
                Interval head = building.removeFirst();
                dq.add(head);
                buildingSize -= head.time;
            }
            long ai = a.apply(i);
            long sub;
            while ((sub = dq.size + buildingSize - ai) > 0) {
                //remove some
                Interval head = dq.pop(i % r);
                if (head == null) {
                    return -1;
                }
                if (head.time <= sub) {
                    //all in
                    head.end = i;
                    all.add(head);
                } else {
                    Interval split = head.split(sub);
                    split.end = i;
                    all.add(split);
                    dq.add(head);
                }
            }
            long req = -sub;
            if (req > 0) {
                Interval newInterval = new Interval();
                newInterval.begin = i;
                newInterval.time = req;
                building.addLast(newInterval);
                buildingSize += newInterval.time;
            }
        }

        while (!building.isEmpty()) {
            Interval head = building.removeFirst();
            if (n - head.begin < l) {
                return -1;
            }
            head.end = n;
            all.add(head);
        }
        for (int i = 0; i < r; i++) {
            Interval iter = dq.ready[i];
            while (iter != null) {
                iter.end = n;
                all.add(iter);
                iter = iter.next;
            }
        }

        long sum = 0;
        for (Interval interval : all) {
//            consumer.consume(interval.begin, interval.end - 1, interval.time);
            sum += interval.time * DigitUtils.ceilDiv(interval.end - interval.begin, r);
        }
        return sum;
    }

    static class RangeMultiDeque {
        RangeTree rt;
        Interval[] ready;
        int n;
        long size;

        public RangeMultiDeque(int n) {
            this.n = n;
            rt = new RangeTree(n);
            ready = new Interval[n];
        }

        public void add(Interval interval) {
            int index = interval.begin % n;
            interval.next = ready[index];
            ready[index] = interval;
            rt.add(index);
            size += interval.time;
        }

        public Interval pop(int k) {
            int c = rt.ceil(k);
            if (c == -1 && k > 0) {
                c = rt.ceil(0);
            }
            if (c == -1) {
                return null;
            }
            Interval ans = ready[c];
            ready[c] = ans.next;
            ans.next = null;
            size -= ans.time;
            if (ready[c] == null) {
                rt.remove(c);
            }
            return ans;
        }
    }

    static class Interval {
        int begin;
        int end;
        long time;
        Interval next;

        Interval split(long alloc) {
            Interval ans = new Interval();
            ans.begin = begin;
            ans.end = end;
            ans.time = alloc;
            time -= alloc;
            return ans;
        }

        @Override
        public String toString() {
            return Arrays.toString(new long[]{begin, time});
        }
    }
}
