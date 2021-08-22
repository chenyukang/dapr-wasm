rustwasmc  build --enable-ext
cp ./pkg/classify_lib_bg.wasm ../image-api-go/lib/classify_bg.wasm

cargo build --target wasm32-wasi

 ./wasmedgec-tensorflow --generic-binary ./target/wasm32-wasi/debug/grayscale_bin.wasm  grayscale.so
cp ./grayscale.so ../image-api-rs/lib
cp ./grayscale.so ../image-api-go/lib

echo "finished build image-wasm-rs ..."