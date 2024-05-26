package template.problem;

import template.binary.Bits;
import template.binary.Log2;
import template.rand.Randomized;
import template.utils.Buffer;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;

/**
 * 给定序列$a_1, a_2, ... , a_n$, 对于$n^2$个二元组$(i,j)$，记$f(i,j)=a_i ^ a_j$，将所有二元组在$f$下的值从小到大排序，问第$k$大的是哪个。
 * <br>
 * 时间复杂度为$O(n(\log_2 M+\log_2 n)$，其中$M$是最大的数，空间复杂度为O(n)
 */
public class KthXorTwoElement {
    /**
     * This method will modify the argument, backing up it if necessary
     */
    public static long solve(long[] data, long k) {
        int n = data.length;

        Buffer<Interval> buffer = new Buffer<>(Interval::new, x -> {}, n * 2);

        Randomized.shuffle(data);
        Arrays.sort(data);

        List<Interval> lastLevel = new ArrayList<>(n);
        List<Interval> curLevel = new ArrayList<>(n);
        lastLevel.add(newInterval(buffer, 0, n - 1));
        int level = Log2.floorLog(data[n - 1]);
        long mask = 0;
        for (; level >= 0; level--) {
            curLevel.clear();
            for (Interval interval : lastLevel) {
                int l = interval.l;
                int r = interval.r;
                int m = r;
                while (m >= l && Bits.get(data[m], level) == 1) {
                    m--;
                }
                interval.m = m;
            }
            long total = 0;
            for (Interval interval : lastLevel) {
                total += (long) (interval.m - interval.l + 1) * (interval.relative.m - interval.relative.l + 1);
                total += (long) (interval.r - interval.m) * (interval.relative.r - interval.relative.m);
            }
            if (total < k) {
                k -= total;
                mask = Bits.set(mask, level, true);
                for (Interval interval : lastLevel) {
                    if (interval.relative == interval) {
                        if (interval.l <= interval.m && interval.m < interval.r) {
                            Interval a = newInterval(buffer, interval.l, interval.m);
                            Interval b = newInterval(buffer, interval.m + 1, interval.r);
                            a.relative = b;
                            b.relative = a;
                            curLevel.add(a);
                            curLevel.add(b);
                        }
                    } else if (interval.r >= interval.relative.r) {
                        if (interval.l <= interval.m && interval.relative.r > interval.relative.m) {
                            Interval a = newInterval(buffer, interval.l, interval.m);
                            Interval b = newInterval(buffer, interval.relative.m + 1, interval.relative.r);
                            a.relative = b;
                            b.relative = a;
                            curLevel.add(a);
                            curLevel.add(b);
                        }
                        if (interval.m < interval.r && interval.relative.m >= interval.relative.l) {
                            Interval a = newInterval(buffer, interval.m + 1, interval.r);
                            Interval b = newInterval(buffer, interval.relative.l, interval.relative.m);
                            a.relative = b;
                            b.relative = a;
                            curLevel.add(a);
                            curLevel.add(b);
                        }
                    }
                }
            } else {
                for (Interval interval : lastLevel) {
                    if (interval.relative == interval) {
                        if (interval.l <= interval.m) {
                            Interval a = newInterval(buffer, interval.l, interval.m);
                            a.relative = a;
                            curLevel.add(a);
                        }
                        if (interval.m < interval.r) {
                            Interval a = newInterval(buffer, interval.m + 1, interval.r);
                            a.relative = a;
                            curLevel.add(a);
                        }
                    } else if (interval.r >= interval.relative.r) {
                        if (interval.l <= interval.m && interval.relative.l <= interval.relative.m) {
                            Interval a = newInterval(buffer, interval.l, interval.m);
                            Interval b = newInterval(buffer, interval.relative.l, interval.relative.m);
                            a.relative = b;
                            b.relative = a;
                            curLevel.add(a);
                            curLevel.add(b);
                        }
                        if (interval.m < interval.r && interval.relative.m < interval.relative.r) {
                            Interval a = newInterval(buffer, interval.m + 1, interval.r);
                            Interval b = newInterval(buffer, interval.relative.m + 1, interval.relative.r);
                            a.relative = b;
                            b.relative = a;
                            curLevel.add(a);
                            curLevel.add(b);
                        }
                    }
                }
            }

            for (Interval interval : lastLevel) {
                buffer.release(interval);
            }

            List<Interval> tmp = curLevel;
            curLevel = lastLevel;
            lastLevel = tmp;
        }

        return mask;
    }

    private static Interval newInterval(Buffer<Interval> buffer, int l, int r) {
        Interval ans = buffer.alloc();
        ans.l = l;
        ans.r = r;
        return ans;
    }

    static class Interval {
        int l;
        int r;

        int m;
        Interval relative = this;
    }
}
