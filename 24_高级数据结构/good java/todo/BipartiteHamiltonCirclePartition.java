package template.problem;

import template.graph.Graph;
import template.primitve.generated.datastructure.IntegerArrayList;
import template.primitve.generated.graph.IntegerDinic;
import template.primitve.generated.graph.IntegerFlow;
import template.primitve.generated.graph.IntegerFlowEdge;

import java.util.ArrayList;
import java.util.List;

// 二分图哈密尔顿回路分解
public class BipartiteHamiltonCirclePartition {
    List<IntegerFlowEdge>[] g;
    int L;
    int R;

    int left(int i) {
        return i;
    }

    int right(int i) {
        return i + L;
    }

    int src() {
        return L + R;
    }

    int sink() {
        return src() + 1;
    }

    public BipartiteHamiltonCirclePartition(int L, int R) {
        this.L = L;
        this.R = R;
        g = Graph.createGraph(sink() + 1);
        for (int i = 0; i < L; i++) {
            IntegerFlow.addFlowEdge(g, src(), left(i), 2);
        }
        for (int i = 0; i < R; i++) {
            IntegerFlow.addFlowEdge(g, right(i), sink(), 2);
        }
    }

    public void addEdge(int a, int b) {
        IntegerFlow.addFlowEdge(g, left(a), right(b), 1);
    }

    public boolean solve() {
        if (L != R) {
            return false;
        }
        IntegerDinic dinic = new IntegerDinic();
        int flow = dinic.apply(g, src(), sink(), 2 * L);
        return flow == 2 * L;
    }

    boolean[] visited;

    void dfsL(int root, IntegerArrayList seq) {
        visited[root] = true;
        seq.add(root);
        for (IntegerFlowEdge e : g[root]) {
            if (e.real && e.flow == 1 && !visited[e.to]) {
                dfsR(e.to, seq);
                return;
            }
        }
    }

    void dfsR(int root, IntegerArrayList seq) {
        visited[root] = true;
        seq.add(root);
        for (IntegerFlowEdge e : g[root]) {
            if (!e.real && e.rev.flow == 1 && !visited[e.to]) {
                dfsL(e.to, seq);
                return;
            }
        }
    }

    /**
     * left start with 0 and right start with L
     */
    public List<IntegerArrayList> getAllCircle() {
        visited = new boolean[sink() + 1];
        List<IntegerArrayList> ans = new ArrayList<>();
        IntegerArrayList seq = new IntegerArrayList(sink() + 1);
        for (int i = 0; i < L; i++) {
            if (visited[i]) {
                continue;
            }
            seq.clear();
            dfsL(i, seq);
            ans.add(seq.clone());
        }
        return ans;
    }
}
