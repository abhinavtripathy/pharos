use json;

pub fn parse(blocks: &json::JsonValue) -> Result<Vec<&json::JsonValue>, &'static str> {
    println!();
    println!("Printing in parser");
    let mut parsed_code = Vec::new();
    match blocks {
        json::JsonValue::Array(arr) => {
            for i in 0..arr.len() {
                let line = arr.get(i).unwrap();
                parsed_code.push(&line["text"]);
                println!("{}", line["text"]);
            }
            println!();
            Ok(parsed_code)
        }
        _ => Err("Recieved wrong type of input"),
    }
}
