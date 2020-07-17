use json;
use pest_consume::{match_nodes, Error, Parser};
use reqwest::ClientBuilder;
use std::time::Duration;

extern crate pest;

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
impl pharos_parser {}

// test will always fail, used for testing parser without initializingt the server
#[test]
fn testParser() {
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
         ((9 modulus 10) - false)

         // print statements // 
         print "hello"
         out "output"
         output "foo"

         // if statments //
         if (3 < 4){
             print "if"
         }
         else if (5 > 10){
             out "not possible"
         }
         else{
             if (4 || 5){
                 let z = 50
                 if (9 || false ){
                     output "another nest"
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
                 print "nest"
                 loop until (3 > 2){
                     let nest = true
                 }
             }
         }

         // functions //
         function noinput {
             let foo = "true"
         }
         
         function hoot inputs owl, bird {
             let fly = "away"
         }

         // function calls // 
         hello ()
         hello (one)
         hello (one, two)
        
        "#,
    )
    .unwrap();
    // .unwrap();
    // let pairs = pharos_parser::parse(
    //     Rule::json,
    //     r#"
    //      [[2],2,3,2,[1,3]]
    //      (4 < (5 < 2))
    //     "#,
    // )
    //println!("{:?}", pairs);
    for pair in pairs {
        println!("{:?} : {:?}", pair.as_str(), pair.as_rule());
    }

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
