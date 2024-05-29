package template.primitve.generated.graph;


import template.primitve.generated.datastructure.IntegerDeque;
import template.primitve.generated.datastructure.IntegerDequeImpl;

import java.util.List;

public class LongSpfaMinimumCostFlow implements LongMinimumCostFlow, LongAugmentMinimumCostFlow {
    IntegerDeque deque;
    long[] dists;
    boolean[] inque;
    LongCostFlowEdge[] prev;
    List<LongCostFlowEdge>[] net;
    LongAugmentCallback callback = LongAugmentCallback.NIL;

    @Override
    public void setCallback(LongAugmentCallback callback) {
        this.callback = callback == null ? LongAugmentCallback.NIL : callback;
    }

    public LongSpfaMinimumCostFlow() {
    }

    private void prepare(int vertexNum) {
        if (dists == null || dists.length < vertexNum) {
            deque = new IntegerDequeImpl(vertexNum);
            dists = new long[vertexNum];
            inque = new boolean[vertexNum];
            prev = new LongCostFlowEdge[vertexNum];
        }
    }

    private void spfa(int s, long inf) {
        deque.clear();
        for (int i = 0; i < net.length; i++) {
            dists[i] = inf;
            inque[i] = false;
        }
        dists[s] = 0;
        prev[s] = null;
        deque.addLast(s);
        while (!deque.isEmpty()) {
            int head = deque.removeFirst();
            inque[head] = false;
            for (LongCostFlowEdge e : net[head]) {
                if (e.flow > 0 && dists[e.to] > dists[head] - e.cost) {
                    dists[e.to] = dists[head] - e.cost;
                    prev[e.to] = e;
                    if (!inque[e.to]) {
                        inque[e.to] = true;
                        deque.addLast(e.to);
                    }
                }
            }
        }
    }


    private static final long INF = Long.MAX_VALUE / 4;

    @Override
    public long[] apply(List<LongCostFlowEdge>[] net, int s, int t, long send) {
        prepare(net.length);
        long cost = 0;
        long flow = 0;
        this.net = net;
        while (flow < send) {
            spfa(t, INF);
            if (dists[s] == INF) {
                break;
            }
            int iter = s;
            long sent = send - flow;
            while (prev[iter] != null) {
                sent = Math.min(sent, prev[iter].flow);
                iter = prev[iter].rev.to;
            }
            if (!callback.callback(sent, dists[s])) {
                break;
            }
            iter = s;
            while (prev[iter] != null) {
                LongFlow.send(prev[iter], -sent);
                iter = prev[iter].rev.to;
            }
            cost += sent * dists[s];
            flow += sent;
        }
        return new long[]{flow, cost};
    }
}
