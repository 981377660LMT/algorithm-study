package template.primitve.generated.datastructure;

/**
 * maintain a0, ..., an and b0, ... , bn, supporting follow operations:
 * <ul>
 * <li>add new element a[n + 1] and b[n + 1] (a[n + 1] should be the largest element in a)</li>
 * <li>given x, for all i, set ai = min(ai, x)</li>
 * <li>find max(a[i] + b[i])</li>
 * </ul>
 * <p>
 * Each operation done in O(1) in average
 */
public class MinPairContainer {
  long[] as;
  long[] bs;
  int size = 0;

  public static final long inf = (long) 4e18;

  public MinPairContainer(int n) {
    as = new long[n];
    bs = new long[n];
  }

  public void clear() {
    size = 0;
  }

  private void insert(long a, long b) {
    as[size] = a;
    bs[size] = b;
    size++;
  }

  private void pop() {
    size--;
    // if (size > 0) {
    // cs[size - 1] = Math.min(cs[size - 1], cs[size]);
    // as[size - 1] = Math.min(as[size - 1], cs[size - 1]);
    // }
  }

  private long event(long a0, long b0, long b1) {
    // min(a0) + b0 >= min(a1) + b1
    // a1 >= a0
    if (b0 <= b1) {
      return -inf;
    }
    // a0 + b0 >= delta + b1
    return a0 + b0 - b1;
  }

  private void pushDown() {
    as[size - 2] = Math.min(as[size - 2], as[size - 1]);
  }

  public void add(long a, long b) {
    assert size == 0 || a >= as[size - 1];
    while (size > 1) {
      pushDown();
      long e1 = event(as[size - 2], bs[size - 2], bs[size - 1]);
      long e2 = event(as[size - 1], bs[size - 1], b);
      if (e1 >= e2) {
        pop();
      } else {
        break;
      }
    }
    if (size == 0 || as[size - 1] + bs[size - 1] < a + b) {
      insert(a, b);
    }
  }

  public void update(long a) {
    if (size == 0) {
      return;
    }
    as[size - 1] = Math.min(as[size - 1], a);
    while (size > 1) {
      pushDown();
      long e = event(as[size - 2], bs[size - 2], bs[size - 1]);
      if (e >= as[size - 1]) {
        pop();
      } else {
        break;
      }
    }
  }

  public long query() {
    return as[size - 1] + bs[size - 1];
  }

  public boolean empty() {
    return size == 0;
  }
}
