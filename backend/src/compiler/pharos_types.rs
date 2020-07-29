// the definition of types need to be more specific
// this is only the first draft for testing purposes

#[derive(PartialEq, Debug, Clone)]
pub enum Exp {
    Number { value: f64 },   // Done
    Boolean { value: bool }, // Done
    Str { value: String },
    Identifier { name: String },                     // Done
    FunctionCall { name: Box<Exp>, args: Vec<Exp> }, // have to build Exp before we evaluate this
    binop { op: Op2, e1: Box<Exp>, e2: Box<Exp> },
    uniop { op: Op1, e1: Box<Exp> },
}

//  re evaluate the structure of Block
// Code should only be vec of statements, this type is only temporary
#[derive(PartialEq, Debug, Clone)]
pub enum Code {
    Exp { value: Exp },
    Stmt { value: Stmt },
}

#[derive(PartialEq, Debug, Clone)]
pub enum Stmt {
    FunctionCall {
        name: Box<Exp>,
        args: Vec<Exp>,
    }, // TODO find a way to make it work without making function call an expression and a statement
    Function {
        name: Box<Exp>,
        args: Vec<Exp>,
        body: Vec<Code>,
    },
    Print {
        value: String,
    },
    Assignment {
        name: String,
        value: Box<Exp>,
    },
    If {
        cond: Box<Exp>,
        body: Vec<Code>,
    },
    Loop {
        cond: Box<Exp>,
        body: Vec<Code>,
    },
    Return {
        value: Box<Exp>,
    },
}

#[derive(PartialEq, Debug, Clone)]
pub enum Op1 {
    Not,
}

#[derive(PartialEq, Debug, Clone)]
pub enum Op2 {
    Add,
    Sub,
    Mul,
    Div,
    Mod,
    Equal,
    NEqual,
    GT,
    LT,
    GTE,
    LTE,
    And,
    Or,
}

// TODO: create constructors for all statements
pub mod constructors {
    use crate::compiler::pharos_types::{Exp, Op1, Op2};

    pub fn exp_number(n: f64) -> Exp {
        Exp::Number { value: n }
    }
    pub fn exp_boolean(b: bool) -> Exp {
        Exp::Boolean { value: b }
    }
    pub fn exp_string(s: String) -> Exp {
        Exp::Str { value: s }
    }
    pub fn exp_identifier(n: String) -> Exp {
        Exp::Identifier { name: n }
    }
    pub fn exp_functioncall(n: Box<Exp>, a: Vec<Exp>) -> Exp {
        Exp::FunctionCall { name: n, args: a }
    }
    pub fn exp_binop(o: Op2, ex1: Box<Exp>, ex2: Box<Exp>) -> Exp {
        Exp::binop {
            op: o,
            e1: ex1,
            e2: ex2,
        }
    }
    pub fn exp_uniop(o: Op1, ex1: Box<Exp>) -> Exp {
        Exp::uniop { op: o, e1: ex1 }
    }
}
