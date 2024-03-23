package template.string.re;

public class Subgraph {
    State head;
    Transfer outbound;

    public Subgraph(State head, Transfer outbound) {
        this.head = head;
        this.outbound = outbound;
    }

}
