package template.problem;

import template.graph.Graph;
import template.primitve.generated.graph.LongAugmentCallback;
import template.primitve.generated.graph.LongAugmentMinimumCostFlow;
import template.primitve.generated.graph.LongCostFlowEdge;
import template.primitve.generated.graph.LongFlow;

import java.util.ArrayList;
import java.util.List;

/**
 * <p>
 * Given directed graph.
 * You can increase the length of some edges with
 * each unit cost and you can't make such operation more than limit times.
 * Find the minimum possible distance in O(FE\log_2V+q\log_2F)
 * while F=\sum_e cost(e).
 * </p>
 */
public class MakeDistanceFarthest {
    private static long inf = Long.MAX_VALUE / 4;

    private static class LinearFunction {
        long l;
        long a;
        long b;

        public LinearFunction(long l, long a, long b) {
            this.l = l;
            this.a = a;
            this.b = b;
        }

        long getL() {
            return a * l + b;
        }

        double inverse(double y) {
            return (y - b) / a;
        }

        long apply(long x) {
            return a * x + b;
        }
    }

    private List<LongCostFlowEdge>[] g;
    LinearFunction[] fs;

    public MakeDistanceFarthest(int n) {
        g = Graph.createGraph(n);
    }

    public void addEdge(int u, int v, long len, long cost) {
        LongFlow.addCostFlowEdge(g, u, v, cost, len);
    }

    public void addLimitedEdge(int u, int v, long len, long cost, long limit) {
        LongFlow.addCostFlowEdge(g, u, v, cost, len);
        LongFlow.addCostFlowEdge(g, u, v, inf, len + limit);
    }

    /**
     * ensure there is a path from s to t
     *
     * @param mcf
     * @param s
     * @param t
     */
    public void solve(LongAugmentMinimumCostFlow mcf, int s, int t, long budgeLimit, long distLimit, long flowLimit) {
        List<LinearFunction> list = new ArrayList<>();
        LongAugmentCallback callback = new LongAugmentCallback() {
            long sumFlow = 0;
            long sumCost = 0;

            @Override
            public boolean callback(long flow, long pathCost) {
                sumFlow += flow;
                sumCost += flow * pathCost;

                if (!list.isEmpty() && list.get(list.size() - 1).l == pathCost) {
                    list.remove(list.size() - 1);
                }
                LinearFunction func = new LinearFunction(pathCost, sumFlow, -sumCost);
                list.add(func);
                return func.getL() <= budgeLimit && func.l <= distLimit;
            }
        };
        mcf.setCallback(callback);
        mcf.apply(g, s, t, flowLimit);
        fs = list.toArray(new LinearFunction[0]);
    }

    /**
     * get maximum distance with no more than x expense in O(\log_2F).
     *
     * @param x
     * @return
     */
    public double queryByExpense(long x) {
        int l = 0;
        int r = fs.length - 1;
        while (l < r) {
            int mid = (l + r) / 2;
            boolean valid = mid + 1 >= fs.length ||
                    fs[mid + 1].getL() > x;
            if (valid) {
                r = mid;
            } else {
                l = mid + 1;
            }
        }
        return fs[l].inverse(x);
    }

    /**
     * <pre>
     * get minimum expense that make distance greater than or equal to x.
     * time complexity: O(\log_2n)
     * </pre>
     */
    public long queryByShortestPath(long x) {
        int l = 0;
        int r = fs.length - 1;
        while (l < r) {
            int mid = (l + r + 1) / 2;
            boolean valid = fs[mid].l <= x;
            if (valid) {
                l = mid;
            } else {
                r = mid - 1;
            }
        }
        return fs[l].apply(x);
    }
}
