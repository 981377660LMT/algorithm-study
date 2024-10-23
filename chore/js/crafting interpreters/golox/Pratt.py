# https://github.com/ravener/bantam.py/blob/main/main.py


from abc import ABCMeta, abstractmethod
from collections.abc import Iterator
from enum import Enum, auto, unique
from io import StringIO


# region TYPE
class ParseException(Exception):
    pass


class Precedence:
    # 运算符优先级
    ASSIGNMENT = 1
    CONDITIONAL = 2
    SUM = 3
    PRODUCT = 4
    EXPONENT = 5
    PREFIX = 6
    POSTFIX = 7
    CALL = 8


@unique
class TokenType(Enum):
    LEFT_PAREN = auto()
    RIGHT_PAREN = auto()
    COMMA = auto()
    ASSIGN = auto()
    PLUS = auto()
    MINUS = auto()
    ASTERISK = auto()
    SLASH = auto()
    CARET = auto()
    TILDE = auto()
    BANG = auto()
    QUESTION = auto()
    COLON = auto()
    NAME = auto()
    EOF = auto()

    def punctuator(self):
        mappings = {
            "LEFT_PAREN": "(",
            "RIGHT_PAREN": ")",
            "COMMA": ",",
            "ASSIGN": "=",
            "PLUS": "+",
            "MINUS": "-",
            "ASTERISK": "*",
            "SLASH": "/",
            "CARET": "^",
            "TILDE": "~",
            "BANG": "!",
            "QUESTION": "?",
            "COLON": ":",
        }

        return mappings.get(self.name)


class Token:
    def __init__(self, type, text):
        self.type = type
        self.text = text


# endregion


# region EXPRESSIONS AST


class Expression(metaclass=ABCMeta):
    @abstractmethod
    def print(self, builder):
        pass


# An assignment expression like "a = b".
class AssignExpression(Expression):
    def __init__(self, name, right):
        self.name = name
        self.right = right

    def print(self, builder):
        builder.write("(")
        builder.write(self.name)
        builder.write(" = ")
        self.right.print(builder)
        builder.write(")")


# A function call like "a(b, c, d)".
class CallExpression(Expression):
    def __init__(self, function, args):
        self.function = function
        self.args = args

    def print(self, builder):
        self.function.print(builder)
        builder.write("(")

        for i, arg in enumerate(self.args):
            arg.print(builder)

            if i < len(self.args) - 1:
                builder.write(", ")

        builder.write(")")


# A ternary conditional expression like "a ? b : c".
class ConditionalExpression(Expression):
    def __init__(self, condition, then_arm, else_arm):
        self.condition = condition
        self.then_arm = then_arm
        self.else_arm = else_arm

    def print(self, builder):
        builder.write("(")
        self.condition.print(builder)
        builder.write(" ? ")
        self.then_arm.print(builder)
        builder.write(" : ")
        self.else_arm.print(builder)
        builder.write(")")


# A simple variable name expression like "abc".
class NameExpression(Expression):
    def __init__(self, name):
        self.name = name

    def print(self, builder):
        builder.write(self.name)


# A binary arithmetic expression like "a + b" or "c ^ d".
class OperatorExpression(Expression):
    def __init__(self, left, operator, right):
        self.left = left
        self.operator = operator
        self.right = right

    def print(self, builder):
        builder.write("(")
        self.left.print(builder)
        builder.write(" ")
        builder.write(self.operator.punctuator())
        builder.write(" ")
        self.right.print(builder)
        builder.write(")")


# A postfix unary arithmetic expression like "a!".
class PostfixExpression(Expression):
    def __init__(self, left, operator):
        self.left = left
        self.operator = operator

    def print(self, builder):
        builder.write("(")
        self.left.print(builder)
        builder.write(self.operator.punctuator())
        builder.write(")")


# A prefix unary arithmetic expression like "!a" or "-b".
class PrefixExpression(Expression):
    def __init__(self, operator, right):
        self.operator = operator
        self.right = right

    def print(self, builder):
        builder.write("(")
        builder.write(self.operator.punctuator())
        self.right.print(builder)
        builder.write(")")


# endregion


# region PARSELETS
class PrefixParselet(metaclass=ABCMeta):
    @abstractmethod
    def parse(self, parser, token):
        pass


class InfixParselet(metaclass=ABCMeta):
    @abstractmethod
    def parse(self, parser, left, token):
        pass

    @abstractmethod
    def get_precedence(self):
        pass


# Parses assignment expressions like "a = b". The left side of an assignment
# expression must be a simple name like "a", and expressions are
# right-associative. (In other words, "a = b = c" is parsed as "a = (b = c)").
class AssignParselet(InfixParselet):
    def parse(self, parser, left, token):
        right = parser.parse_expression(Precedence.ASSIGNMENT - 1)

        if not isinstance(left, NameExpression):
            raise ParseException("The left-hand side of an assignment must be a name.")

        return AssignExpression(left.name, right)

    def get_precedence(self):
        return Precedence.ASSIGNMENT


# Generic infix parselet for a binary arithmetic operator. The only
# difference when parsing, "+", "-", "*", "/", and "^" is precedence and
# associativity, so we can use a single parselet class for all of those.
class BinaryOperatorParselet(InfixParselet):
    def __init__(self, precedence, is_right):
        self.precedence = precedence
        self.is_right = is_right  # !True if right-associative.

    def parse(self, parser, left, token):
        # To handle right-associative operators like "^", we allow a slightly
        # lower precedence when parsing the right-hand side. This will let a
        # parselet with the same precedence appear on the right, which will then
        # take *this* parselet's result as its left-hand argument.

        right = parser.parse_expression(self.precedence - (1 if self.is_right else 0))

        return OperatorExpression(left, token.type, right)

    def get_precedence(self):
        return self.precedence


# Parselet to parse a function call like "a(b, c, d)".
class CallParselet(InfixParselet):
    def parse(self, parser, left, token):
        # Parse the comma-separated arguments until we hit, ")".
        args = []

        #  There may be no arguments at all.
        if not parser.match(TokenType.RIGHT_PAREN):
            # Unfortunately no, do-while loop in Python :(
            args.append(parser.parse_expression())
            while parser.match(TokenType.COMMA):
                args.append(parser.parse_expression())
            parser.consume(TokenType.RIGHT_PAREN)

        return CallExpression(left, args)

    def get_precedence(self):
        return Precedence.CALL


# Parselet for the condition or "ternary" operator, like "a ? b : c".
class ConditionalParselet(InfixParselet):
    def parse(self, parser, left, token):
        then_arm = parser.parse_expression()
        parser.consume(TokenType.COLON)
        else_arm = parser.parse_expression(Precedence.CONDITIONAL - 1)

        return ConditionalExpression(left, then_arm, else_arm)

    def get_precedence(self):
        return Precedence.CONDITIONAL


# Parses parentheses used to group an expression, like "a * (b + c)".
class GroupParselet(PrefixParselet):
    def parse(self, parser, token):
        expression = parser.parse_expression()
        parser.consume(TokenType.RIGHT_PAREN)
        return expression


# Simple parselet for a named variable like "abc".
class NameParselet(PrefixParselet):
    def parse(self, parser, token):
        return NameExpression(token.text)


# Generic infix parselet for an unary arithmetic operator. Parses postfix
# unary "?" expressions.
class PostfixOperatorParselet(InfixParselet):
    def __init__(self, precedence):
        self.precedence = precedence

    def parse(self, parser, left, token):
        return PostfixExpression(left, token.type)

    def get_precedence(self):
        return self.precedence


# Generic prefix parselet for an unary arithmetic operator. Parses prefix
# unary "-", "+", "~", and "!" expressions.
class PrefixOperatorParselet(PrefixParselet):
    def __init__(self, precedence):
        self.precedence = precedence

    def parse(self, parser, token):
        # To handle right-associative operators like "^", we allow a slightly
        # lower precedence when parsing the right-hand side. This will let a
        # parselet with the same precedence appear on the right, which will then
        # take *this* parselet's result as its left-hand argument.
        right = parser.parse_expression(self.precedence)

        return PrefixExpression(token.type, right)

    def get_precedence(self):
        return self.precedence


# endregion


# region LEXER
# A very primitive lexer. Takes a string and splits it into a series of
# Tokens. Operators and punctuation are mapped to unique keywords. Names,
# which can be any series of letters, are turned into NAME tokens. All other
# characters are ignored (except to separate names). Numbers and strings are
# not supported. This is really just the bare minimum to give the parser
# something to work with.
class Lexer(Iterator):
    def __init__(self, text: str) -> None:
        self.index = 0
        self.text = text
        self.punctuators = {}

        # Register all of the TokenTypes that are explicit punctuators.
        for type in TokenType:
            punctuator = type.punctuator()
            if punctuator:
                self.punctuators[punctuator] = type

    def __next__(self):
        while self.index < len(self.text):
            c = self.text[self.index]
            self.index += 1

            if c in self.punctuators:
                return Token(self.punctuators[c], c)
            elif c.isalpha():
                # Handle names.
                start = self.index - 1
                while self.index < len(self.text):
                    if not self.text[self.index].isalpha():
                        break
                    self.index += 1

                name = self.text[start : self.index]
                return Token(TokenType.NAME, name)
            else:
                # Ignore all other characters (whitespace, etc.)
                pass

        # Once we've reached the end of the string, just return EOF tokens. We'll
        # just keeping returning them as many times as we're asked so that the
        # parser's lookahead doesn't have to worry about running out of tokens.
        return Token(TokenType.EOF, "")


# endregion


# region PARSER


class Parser:
    def __init__(self, tokens):
        self.tokens = tokens
        self.read = []
        self.prefix_parselets = {}
        self.infix_parselets = {}

    def register(self, token, parselet):
        if isinstance(parselet, PrefixParselet):
            self.prefix_parselets[token] = parselet
        elif isinstance(parselet, InfixParselet):
            self.infix_parselets[token] = parselet

    def parse_expression(self, precedence=0):
        token = self.consume()
        prefix = self.prefix_parselets.get(token.type)

        if not prefix:
            raise ParseException('Could not parse "{}".'.format(token.text))

        left = prefix.parse(self, token)

        while precedence < self.get_precedence():
            token = self.consume()

            infix = self.infix_parselets.get(token.type)
            left = infix.parse(self, left, token)

        return left

    def match(self, expected):
        token = self.look_ahead(0)

        if token.type != expected:
            return False

        self.consume()
        return True

    def consume(self, expected=None):
        token = self.look_ahead(0)

        if expected and token.type != expected:
            raise Exception("Expected token {} and found {}".format(expected, token.type))

        return self.read.pop(0)

    def look_ahead(self, distance):
        while distance >= len(self.read):
            self.read.append(next(self.tokens))

        return self.read[distance]

    def get_precedence(self):
        parser = self.infix_parselets.get(self.look_ahead(0).type)

        if parser:
            return parser.get_precedence()

        return 0


# Extends the generic Parser class with support for parsing the actual Bantam
# grammar.
class BantamParser(Parser):
    def __init__(self, lexer):
        super().__init__(lexer)

        # Register all of the parselets for the grammar.

        # Register the ones that need special parselets.
        self.register(TokenType.NAME, NameParselet())
        self.register(TokenType.ASSIGN, AssignParselet())
        self.register(TokenType.QUESTION, ConditionalParselet())
        self.register(TokenType.LEFT_PAREN, GroupParselet())
        self.register(TokenType.LEFT_PAREN, CallParselet())

        # Register the simple operator parselets.
        self.prefix(TokenType.PLUS, Precedence.PREFIX)
        self.prefix(TokenType.MINUS, Precedence.PREFIX)
        self.prefix(TokenType.TILDE, Precedence.PREFIX)
        self.prefix(TokenType.BANG, Precedence.PREFIX)

        # For kicks, we'll make "!" both prefix and postfix, kind of like ++.
        self.postfix(TokenType.BANG, Precedence.POSTFIX)

        self.infix_left(TokenType.PLUS, Precedence.SUM)
        self.infix_left(TokenType.MINUS, Precedence.SUM)
        self.infix_left(TokenType.ASTERISK, Precedence.PRODUCT)
        self.infix_left(TokenType.SLASH, Precedence.PRODUCT)
        self.infix_right(TokenType.CARET, Precedence.EXPONENT)

    # Registers a postfix unary operator parselet for the given token and
    # precedence.
    def postfix(self, token, precedence):
        self.register(token, PostfixOperatorParselet(precedence))

    # Registers a prefix unary operator parselet for the given token and
    # precedence.
    def prefix(self, token, precedence):
        self.register(token, PrefixOperatorParselet(precedence))

    # Registers a left-associative binary operator parselet for the given token
    # and precedence.
    def infix_left(self, token, precedence):
        self.register(token, BinaryOperatorParselet(precedence, False))

    # Registers a right-associative binary operator parselet for the given token
    # and precedence.
    def infix_right(self, token, precedence):
        self.register(token, BinaryOperatorParselet(precedence, True))


# endregion


passed = 0
failed = 0


# Parses the given chunk of code and verifies that it matches the expected
# pretty-printed result.
def test(source, expected):
    global passed
    global failed

    lexer = Lexer(source)
    parser = BantamParser(lexer)

    try:
        result = parser.parse_expression()
        builder = StringIO()
        result.print(builder)
        actual = builder.getvalue()

        if expected == actual:
            passed += 1
        else:
            failed += 1
            print("[FAIL] Expected: " + expected)
            print("         Actual: " + actual)
    except ParseException as ex:
        failed += 1
        print("[FAIL] Expected: " + expected)
        print("          Error: " + str(ex))


def main() -> None:
    global passed
    global failed

    # Function call.
    test("a()", "a()")
    test("a(b)", "a(b)")
    test("a(b, c)", "a(b, c)")
    test("a(b)(c)", "a(b)(c)")
    test("a(b) + c(d)", "(a(b) + c(d))")
    test("a(b ? c : d, e + f)", "a((b ? c : d), (e + f))")

    # Unary precedence.
    test("~!-+a", "(~(!(-(+a))))")
    test("a!!!", "(((a!)!)!)")

    # Unary and binary predecence.
    test("-a * b", "((-a) * b)")
    test("!a + b", "((!a) + b)")
    test("~a ^ b", "((~a) ^ b)")
    test("-a!", "(-(a!))")
    test("!a!", "(!(a!))")

    # Binary precedence.
    test("a = b + c * d ^ e - f / g", "(a = ((b + (c * (d ^ e))) - (f / g)))")

    # Binary associativity.
    test("a = b = c", "(a = (b = c))")
    test("a + b - c", "((a + b) - c)")
    test("a * b / c", "((a * b) / c)")
    test("a ^ b ^ c", "(a ^ (b ^ c))")

    # Conditional operator.
    test("a ? b : c ? d : e", "(a ? b : (c ? d : e))")
    test("a ? b ? c : d : e", "(a ? (b ? c : d) : e)")
    test("a + b ? c * d : e / f", "((a + b) ? (c * d) : (e / f))")

    # Grouping.
    test("a + (b + c) + d", "((a + (b + c)) + d)")
    test("a ^ (b + c)", "(a ^ (b + c))")
    test("(!a)!", "((!a)!)")

    # Show the results.
    if failed == 0:
        print("Passed all " + str(passed) + " tests.")
    else:
        print("----")
        print("Failed " + str(failed) + " out of " + str(failed + passed) + " tests.")


if __name__ == "__main__":
    main()
