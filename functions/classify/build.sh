rustup target add wasm32-wasi
rustup override set 1.50.0
rustwasmc  build --enable-ext

cp ./pkg/classify_bg.wasm ../../image-api-go/lib/classify_bg.wasm
echo -e "finished build functions/classify ..."


#This will triger error for wasi
cargo build --target wasm32-wasi --release
cp ./target/wasm32-wasi/release/classify.wasm  ../../image-api-go/lib/classify_bg.wasm