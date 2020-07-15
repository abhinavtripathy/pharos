#![feature(proc_macro_hygiene, decl_macro)]
#![allow(warnings)]
use json;
use reqwest::ClientBuilder;
use rocket::http::Method;
use rocket::response::content;
use rocket::*;
use rocket_contrib::json::Json;
use rocket_cors;
use rocket_cors::{AllowedHeaders, AllowedOrigins, Error};
use serde::Deserialize;
use std::time::Duration;

mod compiler;

#[macro_use]
extern crate pest_derive;

// this is to establish a format in which we can pass json data between backend and front end
#[derive(Debug, PartialEq, Eq, Deserialize, Responder)]
struct Data {
    raw: String, // this attribute is to store the json data in strigified form
}

// front end will send teh code with this post request and the backend will respond with the result in a json object
#[post("/post", format = "json", data = "<code>")]
fn retrieve(code: Json<Data>) -> content::Json<Data> {
    // parsing the recieved code and printing it for debugging purposes. This part of the code is not essential
    let retrived_code = json::parse(&code.raw).unwrap();
    let retrived_code_blocks = &retrived_code["raw"]["blocks"]; // getting the array containing the code from the json input
    for i in 0..retrived_code_blocks.len() {
        println!("{}", &retrived_code_blocks[i]["text"]);
    }

    println!("Request Recieved"); // acknowledging that a request has been recieved

    let parsed_code = compiler::parse(retrived_code_blocks).unwrap(); // sending the recieved code to the compiler.

    // sending the result as a json response back to the client
    let mut parsed_code_vec = Vec::new();
    for i in 0..parsed_code.len() {
        parsed_code_vec.push(parsed_code.get(i).unwrap().dump());
    }
    let parsed_code_vec_string = json::stringify(parsed_code_vec);
    println!();
    println!("Recieved parsed code {:?}", parsed_code_vec_string);
    let temp = Data {
        raw: parsed_code_vec_string,
    };
    content::Json(temp)
    //let temp = content::Json(&parsed_code_vec_string);
}

// this is a simple code to show the server has started when the link to the current server is opened
#[get("/")]
fn start() -> &'static str {
    "Server started"
}

fn main() -> Result<(), Error> {
    // for the current purposes we will only allow the host of the fornt end server, this part of the code needs to be changed for production version
    let allowed_origins = AllowedOrigins::some_exact(&["http://localhost:3000/"]);

    let headers = rocket_cors::CorsOptions {
        allowed_origins,
        allowed_methods: vec![Method::Post, Method::Get]
            .into_iter()
            .map(From::from)
            .collect(),
        allowed_headers: AllowedHeaders::All,
        allow_credentials: true,
        ..Default::default()
    }
    .to_cors()?;

    // starting the server and adding the different request formats and restrictions
    rocket::ignite()
        .mount("/", routes![start])
        .mount("/", routes![retrieve])
        .attach(headers)
        .launch();

    Ok(())
}
