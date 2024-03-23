package template.algo;

import java.util.ArrayDeque;
import java.util.Deque;

/**
 * Used to optimize such dp form:
 * <pre>
 * (Y(j) - Y(k)) / (X(j) - X(k)) <= S(i)
 * </pre>
 * while k < j < i and S(i) is an increasing function.
 */
public class DoubleLeqSlopeOptimizer {
    private static class Point {
        final double x;
        final double y;
        final int id;

        private Point(double x, double y, int id) {
            this.x = x;
            this.y = y;
            this.id = id;
        }
    }

    Deque<Point> deque;

    public DoubleLeqSlopeOptimizer() {
        deque = new ArrayDeque<>(0);
    }

    public DoubleLeqSlopeOptimizer(int exp) {
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
        return (b.y - a.y) / (b.x - a.x);
    }

    public Point add(double y, double x, int id) {
        Point t1 = new Point(x, y, id);
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

    public int getBestChoice(double s) {
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

    public void clear() {
        deque.clear();
    }
}

