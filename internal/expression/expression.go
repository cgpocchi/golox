package expression

import "golox/internal/token"

type Expr interface { }

type Grouping struct {
    expression Expr
}

func NewGrouping(expression Expr) (* Grouping) {
    return &Grouping{
        expression: expression,
    }
}

type Literal struct {
    value any
}

func NewLiteral(value any) (* Literal) {
    return &Literal{
        value: value,
    }
}

type Unary struct {
    operator token.Token
    right Expr
}

func NewUnary(operator token.Token, right Expr) (* Unary) {
    return &Unary{
        operator: operator,
        right: right,
    }
}

type Binary struct {
    left Expr
    operator token.Token
    right Expr
}

func NewBinary(left Expr, operator token.Token, right Expr) (* Binary) {
    return &Binary{
        left: left,
        operator: operator,
        right: right,
    }
}
