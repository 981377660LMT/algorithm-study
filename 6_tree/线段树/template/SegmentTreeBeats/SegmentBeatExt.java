package template.datastructure;

import template.math.DigitUtils;
import template.primitve.generated.datastructure.IntToLongFunction;

/**
 * <p>Given a1, a2, ..., an, support four type operations</p>
 * <ul>
 *     <li>given interval l,r, for i in [l,r], update a[i]=min(a[i],x)</li>
 *     <li>given interval l,r, for i in [l,r], update a[i]=max(a[i],x)</li>
 *     <li>given interval l,r, for i in [l,r], update a[i]=a[i]+x</li>
 *     <li>given interval l,r, find a[l]+a[l+1]+...+a[r]</li>
 * </ul>
 * <p>each operation finished with O((log n)^2) time complexity</p>
 */
public class SegmentBeatExt implements Cloneable {
    private static long inf = (long) 2e18;
    private SegmentBeatExt left;
    private SegmentBeatExt right;
    private long firstLargest;
    private long secondLargest;
    private int firstLargestCnt;
    private long firstSmallest;
    private long secondSmallest;
    private int firstSmallestCnt;
    private long dirty;
    private int size;
    private long sum;

    private void setMin(long x) {
        if (firstLargest <= x) {
            return;
        }
        sum -= (firstLargest - x) * firstLargestCnt;
        firstLargest = x;

        if (firstSmallest >= x) {
            firstSmallest = x;
        }
        secondSmallest = Math.min(secondSmallest, x);
        if (secondSmallest == firstSmallest) {
            secondSmallest = inf;
        }
    }

    private void setMax(long x) {
        if (firstSmallest >= x) {
            return;
        }
        sum += (x - firstSmallest) * firstSmallestCnt;
        firstSmallest = x;

        if (firstLargest <= x) {
            firstLargest = x;
        }
        secondLargest = Math.max(secondLargest, x);
        if (secondLargest == firstLargest) {
            secondLargest = -inf;
        }
    }

    private void modify(long x) {
        dirty += x;
        sum += x * size;
        firstSmallest += x;
        firstLargest += x;
        secondSmallest += x;
        secondLargest += x;
    }

    public void pushUp() {
        firstLargest = Math.max(left.firstLargest, right.firstLargest);
        secondLargest = Math.max(left.firstLargest == firstLargest ? left.secondLargest : left.firstLargest, right.firstLargest == firstLargest ? right.secondLargest : right.firstLargest);
        firstLargestCnt = (left.firstLargest == firstLargest ? left.firstLargestCnt : 0) + (right.firstLargest == firstLargest ? right.firstLargestCnt : 0);

        firstSmallest = Math.min(left.firstSmallest, right.firstSmallest);
        secondSmallest = Math.min(left.firstSmallest == firstSmallest ? left.secondSmallest : left.firstSmallest,
                right.firstSmallest == firstSmallest ? right.secondSmallest : right.firstSmallest);
        firstSmallestCnt = (left.firstSmallest == firstSmallest ? left.firstSmallestCnt : 0) + (right.firstSmallest == firstSmallest ? right.firstSmallestCnt : 0);

        sum = left.sum + right.sum;
        size = left.size + right.size;
    }

    public void pushDown() {
        if (dirty != 0) {
            left.modify(dirty);
            right.modify(dirty);
            dirty = 0;
        }
        left.setMin(firstLargest);
        right.setMin(firstLargest);
        left.setMax(firstSmallest);
        right.setMax(firstSmallest);
    }

    public SegmentBeatExt(int l, int r, IntToLongFunction func) {
        if (l < r) {
            int m = DigitUtils.floorAverage(l, r);
            left = new SegmentBeatExt(l, m, func);
            right = new SegmentBeatExt(m + 1, r, func);
            pushUp();
        } else {
            sum = firstSmallest = firstLargest = func.apply(l);
            firstSmallestCnt = firstLargestCnt = 1;
            secondSmallest = inf;
            secondLargest = -inf;
            size = 1;
        }
    }

    private boolean covered(int ll, int rr, int l, int r) {
        return ll <= l && rr >= r;
    }

    private boolean noIntersection(int ll, int rr, int l, int r) {
        return ll > r || rr < l;
    }

    public void updateMin(int ll, int rr, int l, int r, long x) {
        if (noIntersection(ll, rr, l, r)) {
            return;
        }
        if (covered(ll, rr, l, r)) {
            if (firstLargest <= x) {
                return;
            }
            if (secondLargest < x) {
                setMin(x);
                return;
            }
        }
        pushDown();
        int m = DigitUtils.floorAverage(l, r);
        left.updateMin(ll, rr, l, m, x);
        right.updateMin(ll, rr, m + 1, r, x);
        pushUp();
    }

    public void updateMax(int ll, int rr, int l, int r, long x) {
        if (noIntersection(ll, rr, l, r)) {
            return;
        }
        if (covered(ll, rr, l, r)) {
            if (firstSmallest >= x) {
                return;
            }
            if (secondSmallest > x) {
                setMax(x);
                return;
            }
        }
        pushDown();
        int m = DigitUtils.floorAverage(l, r);
        left.updateMax(ll, rr, l, m, x);
        right.updateMax(ll, rr, m + 1, r, x);
        pushUp();
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

    public long querySum(int ll, int rr, int l, int r) {
        if (noIntersection(ll, rr, l, r)) {
            return 0;
        }
        if (covered(ll, rr, l, r)) {
            return sum;
        }
        pushDown();
        int m = DigitUtils.floorAverage(l, r);
        return left.querySum(ll, rr, l, m) +
                right.querySum(ll, rr, m + 1, r);
    }

    public long queryMax(int ll, int rr, int l, int r) {
        if (noIntersection(ll, rr, l, r)) {
            return -inf;
        }
        if (covered(ll, rr, l, r)) {
            return firstLargest;
        }
        pushDown();
        int m = DigitUtils.floorAverage(l, r);
        return Math.max(left.queryMax(ll, rr, l, m),
                right.queryMax(ll, rr, m + 1, r));
    }

    public long queryMin(int ll, int rr, int l, int r) {
        if (noIntersection(ll, rr, l, r)) {
            return -inf;
        }
        if (covered(ll, rr, l, r)) {
            return firstSmallest;
        }
        pushDown();
        int m = DigitUtils.floorAverage(l, r);
        return Math.max(left.queryMin(ll, rr, l, m),
                right.queryMin(ll, rr, m + 1, r));
    }

    private SegmentBeatExt deepClone() {
        SegmentBeatExt seg = clone();
        if (seg.left != null) {
            seg.left = seg.left.deepClone();
        }
        if (seg.right != null) {
            seg.right = seg.right.deepClone();
        }
        return seg;
    }

    @Override
    protected SegmentBeatExt clone() {
        try {
            return (SegmentBeatExt) super.clone();
        } catch (CloneNotSupportedException e) {
            throw new RuntimeException(e);
        }
    }

    private void toString(StringBuilder builder) {
        if (left == null && right == null) {
            builder.append(sum).append(",");
            return;
        }
        pushDown();
        left.toString(builder);
        right.toString(builder);
    }

    @Override
    public String toString() {
        StringBuilder builder = new StringBuilder();
        deepClone().toString(builder);
        if (builder.length() > 0) {
            builder.setLength(builder.length() - 1);
        }
        return builder.toString();
    }
}
