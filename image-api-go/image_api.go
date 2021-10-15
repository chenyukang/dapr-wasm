package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/second-state/WasmEdge-go/wasmedge"

	"github.com/dapr/go-sdk/service/common"
)

func testMemory(count int) {
	file, err := os.Open("./test.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	image := make([]byte, size)

	// read file into bytes
	buffer := bufio.NewReader(file)
	_, err = buffer.Read(image)
	if err != nil {
		fmt.Println("read error")
		return
	}

	println("image size: ", len(image))

	for i := 0; i < count; i++ {
		/// Set not to print debug info
		wasmedge.SetLogErrorLevel()

		/// Create configure
		var conf = wasmedge.NewConfigure(wasmedge.REFERENCE_TYPES)
		defer conf.Delete()
		conf.AddConfig(wasmedge.WASI)

		/// Create VM with configure
		var vm = wasmedge.NewVMWithConfig(conf)
		defer vm.Delete()

		/// Init WASI (test)
		var wasi = vm.GetImportObject(wasmedge.WASI)
		wasi.InitWasi(
			os.Args[1:],     /// The args
			os.Environ(),    /// The envs
			[]string{".:."}, /// The mapping directories
			[]string{},      /// The preopens will be empty
		)

		/// Register WasmEdge-tensorflow and WasmEdge-image
		var tfobj = wasmedge.NewTensorflowImportObject()
		var tfliteobj = wasmedge.NewTensorflowLiteImportObject()
		vm.RegisterImport(tfobj)
		vm.RegisterImport(tfliteobj)
		var imgobj = wasmedge.NewImageImportObject()
		vm.RegisterImport(imgobj)

		/// Instantiate wasm

		vm.LoadWasmFile("./lib/classify_bg.wasm")
		vm.Validate()
		vm.Instantiate()

		res, err := vm.ExecuteBindgen("infer", wasmedge.Bindgen_return_array, image)
		ans := string(res.([]byte))
		if err != nil {
			println("error: ", err.Error())
		}

		fmt.Printf("Image classify result: %q\n", ans)
	}
}

func imageHandlerWASI(_ context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
	image := in.Data

	/// Set not to print debug info
	wasmedge.SetLogErrorLevel()

	/// Create configure
	var conf = wasmedge.NewConfigure(wasmedge.REFERENCE_TYPES)
	defer conf.Delete()
	conf.AddConfig(wasmedge.WASI)

	/// Create VM with configure
	var vm = wasmedge.NewVMWithConfig(conf)
	defer vm.Delete()

	/// Init WASI (test)
	var wasi = vm.GetImportObject(wasmedge.WASI)
	wasi.InitWasi(
		os.Args[1:],     /// The args
		os.Environ(),    /// The envs
		[]string{".:."}, /// The mapping directories
		[]string{},      /// The preopens will be empty
	)

	/// Register WasmEdge-tensorflow and WasmEdge-image
	var tfobj = wasmedge.NewTensorflowImportObject()
	var tfliteobj = wasmedge.NewTensorflowLiteImportObject()
	vm.RegisterImport(tfobj)
	vm.RegisterImport(tfliteobj)
	var imgobj = wasmedge.NewImageImportObject()
	vm.RegisterImport(imgobj)

	/// Instantiate wasm

	vm.LoadWasmFile("./lib/classify_bg.wasm")
	vm.Validate()
	vm.Instantiate()

	res, err := vm.ExecuteBindgen("infer", wasmedge.Bindgen_return_array, image)
	ans := string(res.([]byte))
	if err != nil {
		println("error: ", err.Error())
	}

	fmt.Printf("Image classify result: %q\n", ans)
	out = &common.Content{
		Data:        []byte(ans),
		ContentType: in.ContentType,
		DataTypeURL: in.DataTypeURL,
	}
	return out, nil
}

// currently don't use it, only for demo
func imageHandler(_ context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
	image := string(in.Data)
	cmd := exec.Command("./lib/wasmedge-tensorflow-lite", "./lib/classify.so")
	cmd.Stdin = strings.NewReader(image)

	var std_out bytes.Buffer
	cmd.Stdout = &std_out
	cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	res := std_out.String()
	fmt.Printf("Image classify result: %q\n", res)
	out = &common.Content{
		Data:        []byte(res),
		ContentType: in.ContentType,
		DataTypeURL: in.DataTypeURL,
	}
	return out, nil
}

func main() {
	time.Sleep(time.Second * 20)
	testMemory(20)
	time.Sleep(time.Second * 20)
}
