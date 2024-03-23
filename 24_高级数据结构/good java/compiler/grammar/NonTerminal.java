package template.compiler.grammar;

import java.util.ArrayList;
import java.util.List;

public class NonTerminal extends Symbol {
    List<Production> productions = new ArrayList<>();

    public NonTerminal(String name) {
        super(name);
    }
}
