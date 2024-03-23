package template.compiler.grammar;

public abstract class Symbol {
    String name;

    public Symbol(String name) {
        this.name = name;
    }

    public String getName() {
        return name;
    }
}
