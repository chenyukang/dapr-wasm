use std::io::{self, Read};
mod infer;

pub fn infer(image_data: &[u8]) -> String {
    infer::infer_internal(image_data)
}

pub fn main() {
    let mut buf = Vec::new();
    io::stdin().read_to_end(&mut buf).unwrap();
    let res = infer::infer_internal(&buf);
    println!("{}", res);
}
