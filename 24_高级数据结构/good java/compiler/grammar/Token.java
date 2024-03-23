package template.compiler.grammar;

import template.compiler.ast.AstNode;
import template.compiler.context.Context;

public class Token implements AstNode<String> {
    String type;
    String content;

    public Token(String type, String content) {
        this.type = type;
        this.content = content;
    }

    public String getType() {
        return type;
    }

    public String getContent() {
        return content;
    }

    @Override
    public String eval(Context ctx) {
        return content;
    }
}
