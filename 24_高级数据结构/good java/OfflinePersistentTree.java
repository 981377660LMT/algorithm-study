package template.algo;

import java.util.ArrayList;
import java.util.List;

public class OfflinePersistentTree {
    Node[] nodes;
    int now = 0;
    int version = 0;

    /**
     * the begin version is 0
     */
    public OfflinePersistentTree(int op) {
        nodes = new Node[op + 1];
        nodes[0] = new Node();
        nodes[0].operation = UndoOperation.NIL;
    }

    public void apply(UndoOperation op) {
        assert version <= now;
        nodes[++now] = new Node();
        nodes[now].operation = op;
        nodes[version].adj.add(nodes[now]);
        version = now;
    }

    public void switchVersion(int v) {
        version = v;
    }

    public void solve() {
        dfs(nodes[0]);
    }

    private static void dfs(Node root) {
        root.operation.apply();
        for (Node node : root.adj) {
            dfs(node);
        }
        root.operation.undo();
    }

    private static class Node {
        List<Node> adj = new ArrayList<>();
        UndoOperation operation;
    }
}
