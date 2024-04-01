package template.algo;

import java.util.ArrayDeque;
import java.util.Deque;

/**
 * Used to optimize such dp form:
 * 
 * <pre>
 * (Y(j) - Y(k)) / (X(j) - X(k)) <= S(i)
 * </pre>
 * 
 * while k < j < i and S(i) is an increasing function.
 */
public class LeqSlopeOptimizer {
  private static class Point {
    final long x;
    final long y;
    final int id;

    private Point(long x, long y, int id) {
      this.x = x;
      this.y = y;
      this.id = id;
    }
  }

  Deque<Point> deque;
  boolean maximizeYForSameX;

  public void setMaximizeYWhenConflict(boolean x) {
    this.maximizeYForSameX = x;
  }

  public LeqSlopeOptimizer() {
    deque = new ArrayDeque<>(0);
  }

  public LeqSlopeOptimizer(int exp) {
    deque = new ArrayDeque<>(exp);
  }

  private double slope(Point a, Point b) {
    if (b.x == a.x) {
      if (b.y == a.y) {
        return 0;
      } else if (b.y > a.y) {
        return 1e50;
      } else {
        return 1e-50;
      }
    }
    return (double) (b.y - a.y) / (b.x - a.x);
  }

  public Point add(long y, long x, int id) {
    Point t1 = new Point(x, y, id);
    if (!deque.isEmpty() && deque.peekLast().x == x) {
      if ((deque.peekLast().y >= y) == maximizeYForSameX) {
        return null;
      } else {
        deque.removeLast();
      }
    }
    while (deque.size() >= 2) {
      Point t2 = deque.removeLast();
      Point t3 = deque.peekLast();
      if (slope(t3, t2) < slope(t2, t1)) {
        deque.addLast(t2);
        break;
      }
    }
    deque.addLast(t1);
    return t1;
  }

  public void since(int id) {
    while (!deque.isEmpty() && deque.peekFirst().id < id) {
      deque.removeFirst();
    }
  }

  public int getBestChoice(long s) {
    while (deque.size() >= 2) {
      Point h1 = deque.removeFirst();
      Point h2 = deque.peekFirst();
      if (slope(h2, h1) > s) {
        deque.addFirst(h1);
        break;
      }
    }
    return deque.peekFirst().id;
  }

  public boolean isEmpty() {
    return deque.isEmpty();
  }

  public void clear() {
    deque.clear();
  }
}
