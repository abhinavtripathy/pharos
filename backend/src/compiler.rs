use json;
use reqwest::ClientBuilder;
use std::time::Duration;
extern crate pest;

use pest::Parser;

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

// test will always fail, used for testing parser without initializingt the server
#[test]
fn testParser() {
    let pairs = pharos_parser::parse(
        Rule::pharos,
        r#"
         ( 2 < 5)
         ( 6 == 7)
         (2 is less than three)
         exit
         outside
         (false || true)
         ((false || true) || (true and false))
         (((5 > 3) && (3 < 4)) or !(5 < 6))
         not 4
         (((6 > 2) or (5 > 10) ) and ((3 > 6) and (4 < 9)))
         let x = 5
         (5 + 3)
         print "hello"
        "#,
    )
    // .unwrap();
    // let pairs = pharos_parser::parse(
    //     Rule::json,
    //     r#"
    //      [[2],2,3,2,[1,3]]
    //      (4 < (5 < 2))
    //     "#,
    // )
    .unwrap();
    //println!("{:?}", pairs);
    for pair in pairs {
        //let tokens: Vec<_> = pair.tokens().collect();
        println!("{:?}", pair.as_str());
        // if tokens.len() > 2 {
        //     // for i in 0..tokens.len() {
        //     // }
        // }
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
