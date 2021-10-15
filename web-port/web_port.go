package main

import (
	"bytes"
	"context"
	b64 "encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	dapr "github.com/dapr/go-sdk/client"
)

func GetEnvValue(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func daprClientSend(image []byte, w http.ResponseWriter) {
	// create the client
	client, err := dapr.NewClientWithPort(GetEnvValue("DAPR_GRPC_PORT", "50001"))
	//client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}
	defer client.Close()

	ctx := context.Background()

	content := &dapr.DataContent{
		ContentType: "text/plain",
		Data:        image,
	}

	resp, err := client.InvokeMethodWithContent(ctx, "image-api-go", "/api/image", "post", content)
	if err != nil {
		panic(err)
	}
	log.Printf("dapr-wasmedge-go method api/image has invoked, response: %s", string(resp))
	fmt.Printf("Image classify result: %q\n", resp)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", string(resp))
	//client.Close()
}

func httpClientSend(image []byte, w http.ResponseWriter) {
	client := &http.Client{}

	// http://localhost:<daprPort>/v1.0/invoke/<appId>/method/<method-name>
	req, err := http.NewRequest("POST", "http://localhost:3502/v1.0/invoke/image-api-rs/method/api/image", bytes.NewBuffer(image))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "text/plain")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "image/png")
	res := b64.StdEncoding.EncodeToString([]byte(body))
	//fmt.Print(string(body))
	fmt.Fprintf(w, "%s", res)
}

type mock struct{}

func (m *mock) WriteHeader(statusCode int) {
}

func (m *mock) Header() http.Header {
	return nil
}
func (m *mock) Write(b []byte) (int, error) {
	return len(b), nil
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	println("imageHandler ....")
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		println("error: ", err.Error())
		panic(err)
	}
	api := r.Header.Get("api")
	if api == "go" {
		daprClientSend(body, w)
		for i := 0; i < 10; i++ {
			go func() {
				for i := 0; i < 10; i++ {
					daprClientSend(body, &mock{})
				}
			}()
		}
	} else {
		httpClientSend(body, w)
	}
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	title, _ := filepath.EvalSymlinks("." + r.URL.Path)
	types := map[string]string{
		".html": "text/html",
		".css":  "text/css",
		".js":   "application/javascript",
		".ico":  "image/vnd.microsoft.icon",
	}
	content, _ := loadFile(title)
	w.Header().Set("Content-Type", "text/html")
	for key, typ := range types {
		if strings.HasSuffix(title, key) {
			w.Header().Set("Content-Type", typ)
			break
		}
	}
	if content == nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		fmt.Fprintf(w, "%s", content)
	}
}

func loadFile(path string) ([]byte, error) {
	println("loading page: {}", path)
	body, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func main() {
	http.HandleFunc("/static/", staticHandler)
	http.HandleFunc("/api/hello", imageHandler)
	println("listen to 8080 ...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
