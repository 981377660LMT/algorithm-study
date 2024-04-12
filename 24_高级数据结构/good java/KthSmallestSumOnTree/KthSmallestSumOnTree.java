package template.graph;

import java.util.*;

/**
 * 给定一颗树，每条边上有一个非负权重，要求从根出发，找到路径权重最小的k条路径
 */
public class KthSmallestSumOnTree {
    /**
     * O(klog k)
     *
     * @param root
     * @param k
     */
    public static List<State> kthSmallestSumOnTree(Vertex root, int k) {
        List<State> ans = new ArrayList<>(k);
        PriorityQueue<State> pq = new PriorityQueue<>(2 * k, Comparator.comparingLong(x -> x.sum));
        pq.add(new State(root, null, 0));
        while (!pq.isEmpty() && ans.size() < k) {
            State state = pq.remove();
            ans.add(state);
            //child or bro
            if (state.iterator.hasNext()) {
                Edge e = state.iterator.next();
                pq.add(new State(e.to, state, state.sum + e.weight));
            }
            if (state.parent != null && state.parent.iterator.hasNext()) {
                Edge e = state.parent.iterator.next();
                pq.add(new State(e.to, state.parent, state.parent.sum + e.weight));
            }
        }
        return ans;
    }

    /**
     * O(klog k)
     *
     * @param root
     * @param k
     */
    public List<EdgeState> kthSmallestSumOnTreeWithEdge(Vertex root, int k) {
        List<EdgeState> ans = new ArrayList<>(k);
        PriorityQueue<EdgeState> pq = new PriorityQueue<>(2 * k, Comparator.comparingLong(x -> x.sum));
        pq.add(new EdgeState(root, null, 0, null));
        while (!pq.isEmpty() && ans.size() < k) {
            EdgeState state = pq.remove();
            ans.add(state);
            //child or bro
            if (state.iterator.hasNext()) {
                Edge e = state.iterator.next();
                pq.add(new EdgeState(e.to, state, state.sum + e.weight, e));
            }
            if (state.parent != null && state.parent.iterator.hasNext()) {
                Edge e = state.parent.iterator.next();
                pq.add(new EdgeState(e.to, state.parent, state.parent.sum + e.weight, e));
            }
        }
        return ans;
    }

    public static class State {
        public Vertex v;
        Iterator<Edge> iterator;
        public State parent;
        public long sum;

        public State(Vertex v, State parent, long sum) {
            this.v = v;
            this.parent = parent;
            this.sum = sum;
            iterator = v.children();
        }

        public State(Vertex v, Iterator<Edge> iterator, State parent, long sum) {
            this.v = v;
            this.iterator = iterator;
            this.parent = parent;
            this.sum = sum;
        }
    }

    public static class EdgeState extends State {
        public Edge edge;


        public EdgeState(Vertex v, State parent, long sum, Edge e) {
            super(v, parent, sum);
            this.edge = e;
        }
    }

    public interface Vertex {
        Iterator<Edge> children();
    }

    public static class Edge {
        public Vertex to;
        public long weight;
    }
}

