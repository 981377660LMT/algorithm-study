package template.algo;

import java.util.ArrayDeque;
import java.util.Deque;

/**
 * Used to optimize such dp form:
 * <pre>
 * (Y(j) - Y(k)) / (X(j) - X(k)) >= S(i)
 * </pre>
 * while k < j < i and S(i) is an decreasing function.
 */
public class GeqSlopeOptimizer {
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

    boolean maximizeYForSameX;

    public void setMaximizeYForSameX(boolean x) {
        this.maximizeYForSameX = x;
    }

    public GeqSlopeOptimizer() {
        this(0);
    }

    public GeqSlopeOptimizer(int exp) {
        deque = new ArrayDeque<>(exp);
    }

    Deque<Point> deque;

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
        if (!deque.isEmpty() && deque.peekLast().x == x) {
            if ((deque.peekLast().y >= y) == maximizeYForSameX) {
                return null;
            } else {
                deque.removeLast();
            }
        }
        Point t1 = new Point(x, y, id);
        while (deque.size() >= 2) {
            Point t2 = deque.removeLast();
            Point t3 = deque.peekLast();
            if (slope(t3, t2) > slope(t2, t1)) {
                deque.addLast(t2);
                break;
            }
        }
        deque.addLast(t1);
        return t1;
    }


    public int getBestChoice(long s) {
        while (deque.size() >= 2) {
            Point h1 = deque.removeFirst();
            Point h2 = deque.peekFirst();
            if (slope(h2, h1) < s) {
                deque.addFirst(h1);
                break;
            }
        }
        return deque.peekFirst().id;
    }

    public void clear() {
        deque.clear();
    }

    public void since(int id) {
        while (!deque.isEmpty() && deque.peekFirst().id < id) {
            deque.removeFirst();
        }
    }
}
