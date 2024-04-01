package template.problem;

import template.graph.Graph;
import template.primitve.generated.graph.IntegerDinic;
import template.primitve.generated.graph.IntegerFlow;
import template.primitve.generated.graph.IntegerFlowEdge;
import template.utils.GridUtils;
import template.utils.SequenceUtils;

import java.util.List;

public class GridHamiltonCirclePartition {
    int n;
    int m;

    int idOfCell(int i, int j) {
        return i * m + j;
    }

    int idOfSrc() {
        return n * m;
    }

    int idOfDst() {
        return idOfSrc() + 1;
    }

    boolean[][] grid;
    int[][] to;

    boolean valid(int i, int j) {
        return i >= 0 && i < n && j >= 0 && j < m && grid[i][j];
    }

    /**
     * <pre>
     * O((nm)^{1.5})
     * find hamilton circle in grid.
     * (res[i][j] / m, res[i][j] % m) is the next node
     * </pre>
     *
     * @param grid
     * @return
     */
    public int[][] partition(boolean[][] grid) {
        this.grid = grid;
        n = grid.length;
        m = grid[0].length;

        to = new int[n][m];
        SequenceUtils.deepFill(to, -1);
        int[] cnts = new int[2];
        for (int i = 0; i < n; i++) {
            for (int j = 0; j < m; j++) {
                if (grid[i][j]) {
                    cnts[(i + j) % 2]++;
                }
            }
        }
        if (cnts[0] != cnts[1]) {
            return null;
        }
        g = Graph.createGraph(n);
        for (int i = 0; i < n; i++) {
            for (int j = 0; j < m; j++) {
                if (!grid[i][j]) {
                    continue;
                }
                if ((i + j) % 2 == 0) {
                    IntegerFlow.addFlowEdge(g, idOfSrc(), idOfCell(i, j), 2);

                    for (int[] d : GridUtils.DIR4) {
                        int ni = i + d[0];
                        int nj = j + d[1];
                        if (!valid(ni, nj)) {
                            continue;
                        }
                        IntegerFlow.addFlowEdge(g, idOfCell(i, j), idOfCell(ni, nj), 1);
                    }
                } else {
                    IntegerFlow.addFlowEdge(g, idOfCell(i, j), idOfDst(), 2);
                }
            }
        }

        IntegerDinic dinic = new IntegerDinic();
        int flow = dinic.apply(g, idOfSrc(), idOfDst(), 2 * cnts[0]);
        if (flow != 2 * cnts[0]) {
            return null;
        }

        for (int i = 0; i < n; i++) {
            for (int j = 0; j < m; j++) {
                if ((i + j) % 2 != 0 || to[i][j] != -1) {
                    continue;
                }
                dfsL(i, j, -1, -1);
            }
        }
        return to;
    }

    List<IntegerFlowEdge>[] g;

    void dfsL(int i, int j, int si, int sj) {
        for (IntegerFlowEdge e : g[idOfCell(i, j)]) {
            if (!e.real || e.flow == 0) {
                continue;
            }
            if (e.to / m == si && e.to % m == sj) {
                continue;
            }
            to[i][j] = e.to;
            dfsR(e.to / m, e.to % m, i, j);
            break;
        }
    }

    void dfsR(int i, int j, int si, int sj) {
        for (IntegerFlowEdge e : g[idOfCell(i, j)]) {
            if (e.real || e.rev.flow == 0) {
                continue;
            }
            if (e.to / m == si && e.to % m == sj) {
                continue;
            }
            to[i][j] = e.to;
            dfsR(e.to / m, e.to % m, i, j);
            break;
        }
    }

}
