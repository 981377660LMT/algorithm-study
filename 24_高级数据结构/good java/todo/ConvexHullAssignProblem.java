package template.problem;

import template.primitve.generated.datastructure.IntegerArrayList;
import template.primitve.generated.datastructure.IntegerDequeImpl;

/**
 * Given n number a_1, ..., a_n and b_1, ... , b_n, now you have know the value of b_1 and b_n, ans
 * you're asked to answer the values of b that satisfy for all 1<i<n that b_i=max(a_i,
 * (b_{i-1}+b_{i+1})/2)
 */
public class ConvexHullAssignProblem {
  public static double[] solve(double[] a) {
    int n = a.length;
    IntegerDequeImpl dq = new IntegerDequeImpl(n);
    for (int i = 0; i < n; i++) {
      while (dq.size() >= 2) {
        int tail = dq.removeLast();
        int tail2 = dq.peekLast();

        if ((a[i] - a[tail]) / (i - tail) < (a[tail] - a[tail2]) / (tail - tail2)) {
          dq.addLast(tail);
          break;
        }
      }

      dq.addLast(i);
    }

    IntegerArrayList convexHull = new IntegerArrayList(dq.size());
    convexHull.addAll(dq.iterator());

    double[] y = new double[n + 1];
    int left = 0;
    for (int i = 0; i < n; i++) {
      while (left + 1 < convexHull.size() && convexHull.get(left + 1) <= i) {
        left++;
      }
      int last = convexHull.get(left);
      if (last == i) {
        y[i] = a[i];
      } else {
        int next = convexHull.get(left + 1);
        y[i] = (double) (i - last) / (next - last) * (a[next] - a[last]) + a[last];
      }
    }

    return y;
  }
}
