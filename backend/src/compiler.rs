use json;
use pest;
use pest_consume::{match_nodes, Error, Parser};
use reqwest::ClientBuilder;
use std::time::Duration;

type Node<'i> = pest_consume::Node<'i, Rule, ()>;

mod pharos_types;

use pharos_types::{constructors::*, Code, Exp, Op1, Op2, Stmt};
// mod pharos_types //{
//     include!("pharos_types.rs");
// }

#[derive(Parser)]
#[grammar = "pharos_grammar.pest"]
pub struct pharos_parser;

pub fn parse(blocks: &json::JsonValue) -> Result<Vec<&json::JsonValue>, &'static str> {
    println!("Printing in parser");
    println!();
    let mut parsed_code = Vec::new();
    match blocks {
        json::JsonValue::Array(arr) => {
            for i in 0..arr.len() {
                let line = arr.get(i).unwrap();
                parsed_code.push(&line["text"]);
                //println!("{}", line["text"]);
            }
            println!();
            Ok(parsed_code)
        }
        _ => Err("Recieved wrong type of input"),
    }
}

#[pest_consume::parser]
impl pharos_parser {
    fn EOI(_input: Node) -> Result<(), Error<Rule>> {
        Ok(())
    }
    fn code(_input: Node) -> Result<Code, Error<Rule>> {
        match_nodes! {_input.into_children();
            [function_call(n)] => {Ok(Code::Exp{value: n})},
            [statement(n)] => {Ok(Code::Stmt{value: n})},
        }
    }
    fn expression(_input: Node) -> Result<Exp, Error<Rule>> {
        match_nodes!(_input.into_children();
        [string(r)] => { return Ok(r); },
        [number(r)] => { return Ok(r); },
        [boolean(r)] => { return Ok(r); },
        [identifier(r)] => { return Ok(r); },
        [function_call(r)] => { return Ok(r); },
        [binop_logic(r)] => {return Ok(r);},
        [uniop_logical(r)] => {return Ok(r);},
        [binop_math(r)] => {return Ok(r);}
        );
        panic!("unknown expression");
    }
    fn statement(_input: Node) -> Result<Stmt, Error<Rule>> {
        match_nodes!(_input.into_children();
        [function_block(r)] => { return Ok(r); },
        );
        panic!("unkown statement");
    }
    fn string(_input: Node) -> Result<Exp, Error<Rule>> {
        let s = _input.as_str();
        Ok(exp_string(s.to_string()))
    }
    fn number(_input: Node) -> Result<Exp, Error<Rule>> {
        let num = _input
            .as_str()
            .parse::<f64>()
            .map_err(|e| _input.error(e))?;
        //println!("inside number {:?}", num);
        Ok(exp_number(num))
    }
    fn boolean(_input: Node) -> Result<Exp, Error<Rule>> {
        let b = _input
            .as_str()
            .parse::<bool>()
            .map_err(|e| _input.error(e))?;
        Ok(exp_boolean(b))
    }
    fn identifier(_input: Node) -> Result<Exp, Error<Rule>> {
        let s = _input.as_str().to_string();
        Ok(exp_identifier(s))
    }
    fn function_call(_input: Node) -> Result<Exp, Error<Rule>> {
        let mut params = Vec::<Exp>::new();
        let x = match_nodes!(_input.into_children();
            [identifier(n)..] =>{
                let x = n.as_slice().to_owned();
                let y = x.get(0).unwrap();
                let mut identifier_name;
                match y{
                    Exp::Identifier{name} => {
                        identifier_name = y.to_owned();
                    },
                    _ => panic!("recieved something other than identifier"),
                };
                if (x.len() > 1){
                    for i in 0..x.len(){
                        params.push(x.get(i).unwrap().to_owned());
                    }
                }
                return Ok(exp_functioncall(Box::new(identifier_name), params));
            },
        );
    }
    fn bin_logical_operators(_input: Node) -> Result<Op2, Error<Rule>> {
        //println!("{}", _input.as_str());
        let o = _input.as_str();
        match o {
            ">=" | "is greater than or equal to" => {
                return Ok(Op2::GTE {});
            }
            "<=" | "is less than or equal to" => {
                return Ok(Op2::LTE {});
            }
            "!=" | "is not equal to" => {
                return Ok(Op2::NEqual {});
            }
            "==" | "is equal to" => {
                return Ok(Op2::Equal {});
            }
            ">" | "is greater than" => {
                return Ok(Op2::GT {});
            }
            "<" | "is less than" => {
                return Ok(Op2::LT {});
            }
            "||" | "or" => {
                return Ok(Op2::Or {});
            }
            "&&" | "and" => {
                return Ok(Op2::And {});
            }
            _ => panic!("recieved unimplemented bin operator"),
        };
    }
    fn binop_logic(_input: Node) -> Result<Exp, Error<Rule>> {
        let temp = _input.as_str();
        match_nodes!(_input.into_children();
        [expression(e1), bin_logical_operators(o), expression(e2)] => {
            //println!("inside binop e1 {:?}, e2 {:?}", e1, e2);
            return Ok(exp_binop(o, Box::new(e1), Box::new(e2)));
        },
        );
    }
    fn bin_math_operators(_input: Node) -> Result<Op2, Error<Rule>> {
        //println!("{}", _input.as_str());
        let o = _input.as_str();
        match o {
            "+" | "plus" => {
                return Ok(Op2::Add {});
            }
            "-" | "minus" => {
                return Ok(Op2::Sub {});
            }
            "*" | "times" => {
                return Ok(Op2::Mul {});
            }
            "/" | "divided by" => {
                return Ok(Op2::Div {});
            }
            "%" | "mod" => {
                return Ok(Op2::Mod {});
            }
            _ => panic!("recieved unimplemented bin operator"),
        };
    }
    fn binop_math(_input: Node) -> Result<Exp, Error<Rule>> {
        let temp = _input.as_str();
        match_nodes!(_input.into_children();
        [expression(e1), bin_math_operators(o), expression(e2)] => {
            //println!("inside binop e1 {:?}, e2 {:?}", e1, e2);
            return Ok(exp_binop(o, Box::new(e1), Box::new(e2)));
        },
        );
    }
    fn uniop_logical(_input: Node) -> Result<Exp, Error<Rule>> {
        //println!("{:?}", _input);
        match_nodes!(_input.into_children();
        [expression(n)] => {return Ok(exp_uniop(Op1::Not, Box::new(n)));}
        );
        Ok(exp_number(-1.0))
    }
    fn function_block(_input: Node) -> Result<Stmt, Error<Rule>> {
        use Rule::*;
        println!("{:?}", _input);
        println!();
        let contents: Vec<Node> = _input.children().collect();
        println!("{:?}", contents);
        let temp = contents.get(0).unwrap();
        let name = match temp.as_rule() {
            identifier => pharos_parser::identifier(temp.to_owned()).unwrap(),
            _ => panic!("first node is not identifier of the function"),
        };
        println!();
        println!("{:?}", name);
        let mut args = Vec::new();
        let mut body_start = 1;
        if (_input.as_str().find("inputs").is_some()) {
            for i in 1..contents.len() {
                let curr = contents.get(i).unwrap();
                let arg_name = match curr.as_rule() {
                    identifier => pharos_parser::identifier(curr.to_owned()).unwrap(),
                    _ => {
                        body_start = i;
                        break;
                    }
                };
                args.push(arg_name);
            }
        }
        let mut body = Vec::new();
        for i in body_start..contents.len() {
            let curr = contents.get(i).unwrap();
            let body_content = match curr.as_rule() {
                code => pharos_parser::code(curr.to_owned()).unwrap(),
                _ => panic!("recieved a type other than code in body"),
            };
            println!("curr {:?}", curr);
            body.push(body_content);
        }
        println!();
        println!("{:?}", args);
        println!();
        println!("{:?}", body);
        Ok(Stmt::Function {
            name: Box::new(name),
            args: args,
            body: body,
        })
    }
}

// test will always fail, used for testing parser without initializingt the server
#[test]
fn parser_grammar() {
    let pairs = pharos_parser::parse(
        Rule::pharos,
        r#"
         // basic data types //
         6
         -9
         9.23
         "foo 0 0 - - #"
         false 
         true
         printer

         // logical expressions //
         ( 2 < 5)
         ( 6 == 7)
         (2 is less than three)
         (2 is less than or equal to y)
         (false || true)
         ((false || true) || (true and false))
         (((5 > 3) && (3 < 4)) or !(5 < 6))
         not 4
         (((6 > 2) or (5 > 10) ) and ((3 > 6) and (4 < 9)))
         
         // Math expressions//
         (5 + 3)
         (6 + (8 - 9))
         (10 + (1 * (5 - 6)))
         (6 plus false)
         ((9 mod 10) - false)

         // print statements // 
         print ("hello")
         out ("output")
         output ("foo")
         out (hello)

         // if statments //

         if (3 < 4){
             print ("if")
         }
         else if (5 > 10){
             out ((5 > 10))
         }
         else{
             if (4 || 5){
                 let z = 50
                 if (9 || false ){
                     output (foo(a,b))
                 }
                 else{
                     let exit = 1
                 }
             }
         }

         //loop statments //

         loop 10 times {
             let x = 10
         }

         loop until (10 > 10){
             let z = 10
             loop 2 times{
                 print (!4)
                 loop until (3 > 2){
                     let nest = true
                 }
             }
         }

         // functions //
         function noinput {
             let foo = "true"
             return
         }
         
         function hoot inputs owl, bird {
             let fly = "away"
             return fly
         }

         // function calls // 
         hello ()
         hello (one)
         hello (one, two)
        
        "#,
    )
    .unwrap();
    for pair in pairs {
        println!("{:?} : {:?}", pair.as_str(), pair.as_rule());
    }

    println!();
    assert_eq!(true, false);
}

#[test]
fn parser_types() {
    let pairs2 = pharos_parser::parse(
        Rule::function_block,
        "function foo inputs hello, call { hoo() }",
    )
    .unwrap();
    let input = pairs2.single().unwrap();
    let result = pharos_parser::function_block(input).unwrap();
    println!("result: {:?}", result);
    println!();
    assert_eq!(true, false);
}

// this test code for Reqwest library.
// fn test_client() -> () {
//     println!("inside test");
//     let request_url = format!("https://webhook.site/5cf1c266-7597-464e-9ae9-8219fa840e17");
//     let timeout = Duration::new(5, 0);
//     let client = ClientBuilder::new().timeout(timeout).build().unwrap();
//     let response = client.head(&request_url).send();
//     println!("responded {:?}", response);
//     // let client = Client::new(rocket::ignite()).expect("valid rocket");
//     // let adress = "http://localhost:3000/".parse().unwrap();
//     // println!("got adress");
//     // let response = client
//     //     .post("/")
//     //     .remote(adress)
//     //     .body("Hello, world!")
//     //     .dispatch();
//     // println!("{:?}", response);
// }
