package template.problem;

import java.util.ArrayList;
import java.util.List;
import java.util.TreeSet;

// 动态最短路
public class DirectedModifiableShortestPath {
    static long inf = (long) 2e18;
    int curTag = 0;
    Segment seg;
    List<Edge> edges;
    Node[] nodes;
    Node src;
    Node dst;

    public void addEdge(int u, int v, long cost) {
        assert !prepared;
        Edge e = new Edge();
        e.a = nodes[u];
        e.b = nodes[v];
        e.w = cost;
        e.a.adj.add(e);
        e.b.adj.add(e);
        edges.add(e);
    }

    boolean prepared;

    private void prepare() {
        if (prepared) {
            return;
        }
        prepared = true;

        TreeSet<Node> pq = new TreeSet<>((a, b) -> a.distToSrc == b.distToSrc ? a.id - b.id : Long.compare(a.distToSrc, b.distToSrc));

        src.distToSrc = 0;
        pq.add(src);
        while (!pq.isEmpty()) {
            Node head = pq.pollFirst();
            for (Edge e : head.adj) {
                if(e.a != head){
                    continue;
                }
                Node node = e.other(head);
                if (node.distToSrc > head.distToSrc + e.w) {
                    pq.remove(node);
                    node.distToSrc = head.distToSrc + e.w;
                    pq.add(node);
                }
            }
        }

        dst.distToDst = 0;
        pq = new TreeSet<>((a, b) -> a.distToDst == b.distToDst ? a.id - b.id : Long.compare(a.distToDst, b.distToDst));
        pq.add(dst);
        while (!pq.isEmpty()) {
            Node head = pq.pollFirst();
            for (Edge e : head.adj) {
                if (e.b != head) {
                    continue;
                }
                Node node = e.other(head);
                if (node.distToDst > head.distToDst + e.w) {
                    pq.remove(node);
                    node.distToDst = head.distToDst + e.w;
                    pq.add(node);
                }
            }
        }

        for (Node trace = src; trace != dst; ) {
            Edge next = null;
            for (Edge e : trace.adj) {
                if(e.a != trace){
                    continue;
                }
                Node node = e.other(trace);
                if (node.distToDst + e.w == trace.distToDst) {
                    next = e;
                    break;
                }
            }
            next.tag = ++curTag;
            trace = next.other(trace);
        }

        seg = new Segment(1, curTag);
        for (Edge e : edges) {
            if (e.tag != -1) {
                continue;
            }
            update(e.a, e.b, e.w);
        }

    }

    public long minDist() {
        prepare();
        return dst.distToSrc;
    }

    /**
     * O(1)
     */
    public long queryOnAddEdge(int a, int b, long dist) {
        prepare();
        long ans = dst.distToSrc;
        ans = Math.min(ans, nodes[a].distToSrc + dist + nodes[b].distToDst);
        return ans;
    }

    /**
     * O(log n)
     */
    public long queryOnModifyEdge(int i, long cost) {
        prepare();
        Edge e = edges.get(i);
        long w = cost;
        long ans = e.a.distToSrc + e.b.distToDst + w;
        if (e.tag == -1) {
            ans = Math.min(ans, dst.distToSrc);
        } else {
            long val = seg.query(e.tag, e.tag, 1, curTag);
            ans = Math.min(ans, val);
        }
        return ans;
    }

    /**
     * O(log n)
     */
    public long queryOnDeleteEdge(int i) {
        return queryOnModifyEdge(i, inf);
    }

    public DirectedModifiableShortestPath(int n, int m, int s, int t) {
        nodes = new Node[n];
        for (int i = 0; i < n; i++) {
            nodes[i] = new Node();
            nodes[i].id = i;
            nodes[i].distToSrc = nodes[i].distToDst = inf;
        }
        edges = new ArrayList<>(m);
        src = nodes[s];
        dst = nodes[t];
    }

    private void update(Node a, Node b, long w) {
        long dist = a.distToSrc + b.distToDst + w;
        int l = prev(a);
        int r = post(b);
        seg.update(l + 1, r - 1, 1, curTag, dist);
    }

    private int post(Node root) {
        if (root.r == Integer.MIN_VALUE) {
            root.r = curTag + 1;
            for (Edge e : root.adj) {
                if(e.a != root){
                    continue;
                }
                Node node = e.other(root);
                if (node.distToDst + e.w == root.distToDst) {
                    if (e.tag != -1) {
                        root.r = e.tag;
                    } else {
                        root.r = post(node);
                    }
                    break;
                }
            }
        }
        return root.r;
    }

    private int prev(Node root) {
        if (root.l == Integer.MIN_VALUE) {
            root.l = -1;
            for (Edge e : root.adj) {
                if(e.b != root){
                    continue;
                }
                Node node = e.other(root);
                if (node.distToSrc + e.w == root.distToSrc) {
                    if (e.tag != -1) {
                        root.l = e.tag;
                    } else {
                        root.l = prev(node);
                    }
                    break;
                }
            }
        }
        return root.l;
    }


    static class Edge {
        Node a;
        Node b;
        long w;
        int tag = -1;

        public Node other(Node x) {
            return x == a ? b : a;
        }

    }

    static class Node {
        List<Edge> adj = new ArrayList<>();
        long distToSrc;
        long distToDst;
        int l = Integer.MIN_VALUE;
        int r = Integer.MIN_VALUE;
        int id;

        public String toString() {
            return "" + id;
        }

    }

    static class Segment implements Cloneable {
        private Segment left;
        private Segment right;
        static long inf = (long) 2e18;
        long min = inf;

        private void modify(long x) {
            min = Math.min(min, x);
        }

        public void pushUp() {

        }

        public void pushDown() {
            left.modify(min);
            right.modify(min);
            min = inf;
        }

        public Segment(int l, int r) {
            if (l < r) {
                int m = (l + r) >> 1;
                left = new Segment(l, m);
                right = new Segment(m + 1, r);
                pushUp();
            } else {

            }
        }

        private boolean covered(int ll, int rr, int l, int r) {
            return ll <= l && rr >= r;
        }

        private boolean noIntersection(int ll, int rr, int l, int r) {
            return ll > r || rr < l;
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
            int m = (l + r) >> 1;
            left.update(ll, rr, l, m, x);
            right.update(ll, rr, m + 1, r, x);
            pushUp();
        }

        public long query(int ll, int rr, int l, int r) {
            if (noIntersection(ll, rr, l, r)) {
                return inf;
            }
            if (covered(ll, rr, l, r)) {
                return min;
            }
            pushDown();
            int m = (l + r) >> 1;
            return Math.min(left.query(ll, rr, l, m),
                    right.query(ll, rr, m + 1, r));
        }

        private Segment deepClone() {
            Segment seg = clone();
            if (seg.left != null) {
                seg.left = seg.left.deepClone();
            }
            if (seg.right != null) {
                seg.right = seg.right.deepClone();
            }
            return seg;
        }

        protected Segment clone() {
            try {
                return (Segment) super.clone();
            } catch (CloneNotSupportedException e) {
                throw new RuntimeException(e);
            }
        }

        private void toString(StringBuilder builder) {
            if (left == null && right == null) {
                builder.append(min).append(",");
                return;
            }
            pushDown();
            left.toString(builder);
            right.toString(builder);
        }

        public String toString() {
            StringBuilder builder = new StringBuilder();
            deepClone().toString(builder);
            if (builder.length() > 0) {
                builder.setLength(builder.length() - 1);
            }
            return builder.toString();
        }
    }
}
