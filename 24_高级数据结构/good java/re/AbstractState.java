package template.string.re;

import java.util.ArrayList;
import java.util.List;
import java.util.stream.Collectors;

public abstract class AbstractState implements State {
    List<Transfer> adj = new ArrayList<>();
    int id = -1;

    void setId(int id) {
        this.id = id;
    }

    @Override
    public List<Transfer> adj() {
        return adj;
    }

    @Override
    public void register(Transfer s) {
        adj.add(s);
    }

    @Override
    public int id() {
        assert id >= 2;
        return id;
    }

    @Override
    public String toString() {
        if (adj.isEmpty()) {
            return "";
        }
        return adj.stream().map(x -> id() + "-->" + x.toString()).collect(Collectors.joining("\n")) + "\n";
    }
}
