// the definition of types need to be more specific
// this is only the first draft for testing purposes

#[derive(PartialEq, Debug, Clone)]
pub enum Exp {
    Undefined {},
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
    // FunctionCall {
    //     name: Box<Exp>,
    //     args: Vec<Exp>,
    // }, // TODO find a way to make it work without making function call an expression and a statement
    Function {
        name: Box<Exp>,
        args: Vec<Exp>,
        body: Vec<Code>,
    },
    Print {
        value: Box<Exp>,
    },
    Assignment {
        name: Box<Exp>,
        value: Box<Exp>,
    },
    Conditionals {
        cond: Vec<Stmt>,
    },
    If {
        cond: Box<Exp>,
        body: Vec<Code>,
    },
    Else_If {
        cond: Box<Exp>,
        body: Vec<Code>,
    },
    Else {
        body: Vec<Code>,
    },
    Loop {
        is_while: bool,
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
    use crate::compiler::pharos_types::{Code, Exp, Op1, Op2, Stmt};

    pub fn exp_undefined() -> Exp {
        Exp::Undefined {}
    }
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
    pub fn stmt_function(n: Box<Exp>, a: Vec<Exp>, b: Vec<Code>) -> Stmt {
        Stmt::Function {
            name: n,
            args: a,
            body: b,
        }
    }
    pub fn stmt_print(v: Box<Exp>) -> Stmt {
        Stmt::Print { value: v }
    }
    pub fn stmt_assignment(n: Box<Exp>, v: Box<Exp>) -> Stmt {
        Stmt::Assignment { name: n, value: v }
    }
    pub fn stmt_conditionals(c: Vec<Stmt>) -> Stmt {
        Stmt::Conditionals { cond: c }
    }
    pub fn stmt_if(c: Box<Exp>, b: Vec<Code>) -> Stmt {
        Stmt::If { cond: c, body: b }
    }
    pub fn stmt_else_if(c: Box<Exp>, b: Vec<Code>) -> Stmt {
        Stmt::Else_If { cond: c, body: b }
    }
    pub fn stmt_else(b: Vec<Code>) -> Stmt {
        Stmt::Else { body: b }
    }
    pub fn stmt_loop(i: bool, c: Box<Exp>, b: Vec<Code>) -> Stmt {
        Stmt::Loop {
            is_while: i,
            cond: c,
            body: b,
        }
    }
    pub fn stmt_return(v: Box<Exp>) -> Stmt {
        Stmt::Return { value: v }
    }
}
