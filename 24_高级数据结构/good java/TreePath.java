package template.graph;

public interface TreePath {
    int length();
    int kthNodeOnPath(int k);
    boolean onPath(int u);
}
