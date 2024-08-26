package template.datastructure;

import template.math.DigitUtils;
import template.primitve.generated.datastructure.IntToIntFunction;
import template.primitve.generated.datastructure.IntToLongFunction;

/**
 * <p>
 * Given a1, a2, ..., an, support two type operations
 * </p>
 * <ul>
 * <li>given interval l,r, for i in [l,r], update a[i]=min(a[i],x)</li>
 * <li>given interval l,r, find a[l]+a[l+1]+...+a[r]</li>
 * </ul>
 * <p>
 * each operation finished with O(log n) time complexity
 * </p>
 */
public class SegmentBeat implements Cloneable {
  private SegmentBeat left;
  private SegmentBeat right;
  private long first;
  private long second;
  private int firstCnt;
  private static long inf = (long) 2e18;
  private long sum;

  private void setMin(long x) {
    if (first <= x) {
      return;
    }
    sum -= (first - x) * firstCnt;
    first = x;
  }

  public void pushUp() {
    first = Math.max(left.first, right.first);
    second = Math.max(left.first == first ? left.second : left.first,
        right.first == first ? right.second : right.first);
    firstCnt = (left.first == first ? left.firstCnt : 0) + (right.first == first ? right.firstCnt : 0);
    sum = left.sum + right.sum;
  }

  public void pushDown() {
    left.setMin(first);
    right.setMin(first);
  }

  public SegmentBeat(int l, int r, IntToLongFunction func) {
    if (l < r) {
      int m = DigitUtils.floorAverage(l, r);
      left = new SegmentBeat(l, m, func);
      right = new SegmentBeat(m + 1, r, func);
      pushUp();
    } else {
      sum = first = func.apply(l);
      second = -inf;
      firstCnt = 1;
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
      if (first <= x) {
        return;
      }
      if (second < x) {
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
      return first;
    }
    pushDown();
    int m = DigitUtils.floorAverage(l, r);
    return Math.max(left.queryMax(ll, rr, l, m),
        right.queryMax(ll, rr, m + 1, r));
  }

  private SegmentBeat deepClone() {
    SegmentBeat seg = clone();
    if (seg.left != null) {
      seg.left = seg.left.deepClone();
    }
    if (seg.right != null) {
      seg.right = seg.right.deepClone();
    }
    return seg;
  }

  @Override
  protected SegmentBeat clone() {
    try {
      return (SegmentBeat) super.clone();
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
