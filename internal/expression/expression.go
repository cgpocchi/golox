package expression

import "golox/internal/token"

type Visitor[T any] interface {
    VisitBinary(expr *Binary) T
    VisitGrouping(expr *Grouping) T
    VisitLiteral(expr *Literal) T
    VisitUnary(expr *Unary) T
}

type Expr interface {
    Accept(visitor Visitor[any]) any
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

func (e *Binary) Accept(visitor Visitor[any]) any {
    return visitor.VisitBinary(e)
}


type Grouping struct {
    expression Expr
}

func NewGrouping(expression Expr) (* Grouping) {
    return &Grouping{
        expression: expression,
    }
}

func (e *Grouping) Accept(visitor Visitor[any]) any {
    return visitor.VisitGrouping(e)
}


type Literal struct {
    value any
}

func NewLiteral(value any) (* Literal) {
    return &Literal{
        value: value,
    }
}

func (e *Literal) Accept(visitor Visitor[any]) any {
    return visitor.VisitLiteral(e)
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

func (e *Unary) Accept(visitor Visitor[any]) any {
    return visitor.VisitUnary(e)
}

