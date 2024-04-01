package template.problem;

import template.polynomial.FastFourierTransform;
import template.utils.CollectionUtils;
import template.utils.PrimitiveBuffers;

import java.util.ArrayList;
import java.util.List;

/**
 * O(n(\log_2n)^2), for all x calculate how many pairs (a, b) in a tree that dist(a, b) == x
 */
public class FrequencyTableOfTreeDistance {
    Node[] nodes;

    public FrequencyTableOfTreeDistance(int n) {
        nodes = new Node[n];
        for (int i = 0; i < n; i++) {
            nodes[i] = new Node();
            nodes[i].vsize = 1;
        }
    }

    public void addEdge(int a, int b) {
        assert a != b;
        addEdge(nodes[a], nodes[b], 1);
    }

    public long[] solve() {
        int n = nodes.length;
        makeTree(nodes[0], null);
        frequencyTable = new long[n];
        dfsForSize(nodes[0], null, null);
        dac(nodes[0]);
        return frequencyTable;
    }

    static class Edge {
        Node a;
        Node b;
        int w;
        int effect;

        public void setFather(Node root) {
            if (b == root) {
                Node tmp = a;
                a = b;
                b = tmp;
            }
        }

        Node other(Node x) {
            return a == x ? b : a;
        }
    }

    static class Node {
        int vsize;
        List<Edge> adj = new ArrayList<>(3);
        int depth;
        int size;
    }

    private static void addEdge(Node a, Node b, int w) {
        Edge e = new Edge();
        addEdge(a, b, w, e);
    }

    private static void addEdge(Node a, Node b, int w, Edge e) {
        e.a = a;
        e.b = b;
        e.w = w;
        a.adj.add(e);
        b.adj.add(e);
    }

    private static Edge replace(Edge e, Node root, Node newNode) {
        if (e.a == root) {
            e.a = newNode;
        } else {
            e.b = newNode;
        }
        return e;
    }

    private static void replace(Node root, Node newNode) {
        Edge e = CollectionUtils.pop(root.adj);
        replace(e, root, newNode);
        newNode.adj.add(e);
    }

    private static void makeTree(Node root, Edge p) {
        root.adj.remove(p);
        if (root.adj.size() > 2) {
            Node last = new Node();
            replace(root, last);
            replace(root, last);
            while (root.adj.size() > 1) {
                Node newNode = new Node();
                replace(root, newNode);
                addEdge(newNode, last, 0);
                last = newNode;
            }
            addEdge(root, last, 0);
        }
        if (p != null) {
            root.adj.add(p);
        }

        for (Edge e : root.adj) {
            if (e == p) {
                continue;
            }
            Node node = e.other(root);
            makeTree(node, e);
        }
    }

    private static void dfsForSize(Node root, Edge p, double[] depthTable) {
        root.size = 1;
        root.depth = p == null ? 0 : p.other(root).depth + p.w;
        if (depthTable != null) {
            depthTable[root.depth] += root.vsize;
        }
        for (Edge e : root.adj) {
            if (e == p) {
                continue;
            }
            e.setFather(root);
            Node node = e.other(root);
            dfsForSize(node, e, depthTable);
            root.size += node.size;
        }
    }

    static Edge min(Edge a, Edge b) {
        if (a == null) {
            return b;
        }
        if (b == null) {
            return a;
        }
        return a.effect < b.effect ? a : b;
    }

    private static Edge findCentroid(Node root, Edge p, int total) {
        if (p != null) {
            p.effect = Math.max(root.size, total - root.size);
        }
        Edge best = p;
        for (Edge e : root.adj) {
            if (e == p) {
                continue;
            }
            Node node = e.other(root);
            Edge res = findCentroid(node, e, total);
            best = min(res, best);
        }
        return best;
    }


    long[] frequencyTable;

    private void dac(Node root) {
        Edge centroid = findCentroid(root, null, root.size);
        if (centroid == null) {
            return;
        }
        double[] a = PrimitiveBuffers.allocDoublePow2(root.size - centroid.b.size);
        double[] b = PrimitiveBuffers.allocDoublePow2(centroid.b.size);
        centroid.a.adj.remove(centroid);
        centroid.b.adj.remove(centroid);
        dfsForSize(centroid.a, null, a);
        dfsForSize(centroid.b, null, b);
        double[] conv = FastFourierTransform.convolution(a, b);
        for (int i = 0; i < conv.length && i + centroid.w < frequencyTable.length; i++) {
            frequencyTable[i + centroid.w] += Math.round(conv[i]);
        }
        PrimitiveBuffers.release(a, b, conv);
        dac(centroid.a);
        dac(centroid.b);
    }
}
