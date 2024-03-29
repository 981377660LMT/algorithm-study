package template.graph;

import java.util.ArrayDeque;
import java.util.ArrayList;
import java.util.Deque;
import java.util.List;
import java.util.TreeSet;

public class DifferSystem {
  private static class Node {
    List<Edge> adj = new ArrayList<>();
    long dist;
    boolean inque;
    int times;
    int id;

    @Override
    public String toString() {
      return "a" + id;
    }
  }

  private static class Edge {
    final Node to;
    final long len;

    private Edge(Node next, long len) {
      this.to = next;
      this.len = len;
    }
  }

  Node[] nodes;
  Deque<Node> deque;
  int n;
  boolean allPos = true;

  public DifferSystem(int n) {
    this.n = n;
    deque = new ArrayDeque<>(n);
    nodes = new Node[n];
    for (int i = 0; i < n; i++) {
      nodes[i] = new Node();
      nodes[i].id = i;
    }
  }

  public void clear(int n) {
    this.n = n;
    allPos = true;
    for (int i = 0; i < n; i++) {
      nodes[i].adj.clear();
    }
  }

  public void lessThanOrEqualTo(int i, int j, long d) {
    nodes[j].adj.add(new Edge(nodes[i], d));
    allPos = allPos && d >= 0;
  }

  public void greaterThanOrEqualTo(int i, int j, long d) {
    lessThanOrEqualTo(j, i, -d);
  }

  public void equalTo(int i, int j, long d) {
    greaterThanOrEqualTo(i, j, d);
    lessThanOrEqualTo(i, j, d);
  }

  public void lessThan(int i, int j, long d) {
    lessThanOrEqualTo(i, j, d - 1);
  }

  public void greaterThan(int i, int j, long d) {
    greaterThanOrEqualTo(i, j, d + 1);
  }

  boolean hasSolution;

  private boolean spfa() {
    if (allPos) {
      dijkstra();
      return true;
    }
    while (!deque.isEmpty()) {
      Node head = deque.removeFirst();
      head.inque = false;
      if (head.times >= n) {
        return false;
      }
      for (Edge edge : head.adj) {
        Node node = edge.to;
        if (node.dist <= edge.len + head.dist) {
          continue;
        }
        node.dist = edge.len + head.dist;
        if (node.inque) {
          continue;
        }
        node.times++;
        node.inque = true;
        deque.addLast(node);
      }
    }
    return true;
  }

  public long possibleSolutionOf(int i) {
    return nodes[i].dist;
  }

  private void prepare(long initDist) {
    deque.clear();
    for (int i = 0; i < n; i++) {
      nodes[i].dist = initDist;
      nodes[i].times = 0;
      nodes[i].inque = false;
    }
  }

  public boolean hasSolution() {
    prepare(0);
    for (int i = 0; i < n; i++) {
      nodes[i].inque = true;
      deque.addLast(nodes[i]);
    }
    hasSolution = spfa();
    return hasSolution;
  }

  public static final long INF = (long) 2e18;

  /**
   * Find max(ai - aj), if INF is returned, it means no constraint between ai and aj
   */
  public long findMaxDifferenceBetween(int i, int j) {
    runSpfaSince(j);
    return nodes[i].dist;
  }

  /**
   * Find min(ai - aj), if INF is returned, it means no constraint between ai and aj
   */
  public long findMinDifferenceBetween(int i, int j) {
    long r = findMaxDifferenceBetween(j, i);
    if (r == INF) {
      return INF;
    }
    return -r;
  }

  /**
   * After invoking this method, the value of i is max(ai - aj)
   */
  public boolean runSpfaSince(int j) {
    prepare(INF);
    deque.clear();
    deque.addLast(nodes[j]);
    nodes[j].dist = 0;
    nodes[j].inque = true;
    hasSolution = spfa();
    return hasSolution;
  }


  private void dijkstra() {
    TreeSet<Node> set = new TreeSet<>(
        (a, b) -> a.dist == b.dist ? Integer.compare(a.id, b.id) : Long.compare(a.dist, b.dist));
    set.addAll(deque);
    while (!set.isEmpty()) {
      Node head = set.pollFirst();
      for (Edge e : head.adj) {
        if (e.to.dist <= head.dist + e.len) {
          continue;
        }
        set.remove(e.to);
        e.to.dist = head.dist + e.len;
        set.add(e.to);
      }
    }
  }

  @Override
  public String toString() {
    StringBuilder builder = new StringBuilder();
    for (int i = 0; i < n; i++) {
      for (Edge edge : nodes[i].adj) {
        builder.append(edge).append('\n');
      }
    }
    builder.append("-------------\n");
    if (!hasSolution) {
      builder.append("impossible");
    } else {
      for (int i = 0; i < n; i++) {
        builder.append("a").append(i).append("=").append(nodes[i].dist).append('\n');
      }
    }
    return builder.toString();
  }
}
