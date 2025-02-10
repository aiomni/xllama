package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"os"

	"github.com/aiomni/xllama/ollama"
	"github.com/aiomni/xllama/routes"
	"github.com/aiomni/xllama/server"
	"github.com/valyala/fasthttp"
)

const ollama_default_addr = "http://localhost:11434"

//go:embed web/dist/*
var webFiles embed.FS

var isDev bool

func main() {
	addrFlag := flag.String("addr", ollama_default_addr, "")
	flag.Parse()

	ctx := context.Background()
	client := ollama.NewOllamaClientWithAddr(*addrFlag)

	err := client.Heartbeat(ctx)
	if err != nil {
		fmt.Printf("Ollama Error: %v\n", err)
		fmt.Println("Please ensure that Ollama is installed and running. You can install it from https://ollama.com or use the addr flag to specify a custom server.")
		return
	}

	var files fs.FS

	if isDev {
		files = os.DirFS("./")
	} else {
		files = webFiles
	}

	s := server.Server{
		Port:       5187,
		WebFiles:   files,
		Middleware: server.WrapOllamaMiddleware(client),
		Routes: map[string]fasthttp.RequestHandler{
			"/api/version":          routes.Version,
			"/api/models":           routes.ListModels,
			"/api/chat/completions": routes.Chat,
		},
	}

	go func() {
		s.Start()
	}()

	fmt.Printf("Server started on  http://localhost:%d\n", s.Port)

	select {}
}
