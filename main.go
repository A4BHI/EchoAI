package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Delta struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Echo struct {
	Model   string  `json:"model"`
	Message []Delta `json:"messages"`
	Stream  bool    `json:"stream"`
}

type response struct {
	Choices []struct {
		Delta Delta `json:"delta"`
	} `json:"choices"`
}

func send(w http.ResponseWriter, r *http.Request, apikey string, url string, history *[]Delta) {
	if r.Method == http.MethodPost {
		// bytess, _ := io.ReadAll(r.Body)
		var Assistant_Res string
		// qs := string(bytess)
		qs, _ := io.ReadAll(r.Body)
		fmt.Println(string(qs))
		*history = append(*history, Delta{Role: "user", Content: string(qs)})
		req := Echo{
			Model:   "nvidia/llama-3.1-nemotron-ultra-253b-v1",
			Message: *history,
			Stream:  true,
		}

		jsonbytes, _ := json.Marshal(req)

		client, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonbytes))

		client.Header.Set("Authorization", "Bearer "+apikey)
		client.Header.Set("Content-Type", "Application/JSON")

		res, _ := http.DefaultClient.Do(client)

		scanner := bufio.NewScanner(res.Body)
		ress := response{}
		for scanner.Scan() {
			lines := scanner.Text()

			if lines == "" || lines == "data: [DONE]" {
				continue
			}

			if strings.HasPrefix(lines, "data:") {
				lines, _ = strings.CutPrefix(lines, "data:")
			}

			json.Unmarshal([]byte(lines), &ress)

			for i := 0; i < len(ress.Choices); i++ {
				fmt.Print(strings.Trim(ress.Choices[i].Delta.Content, "#*"))
				Assistant_Res += strings.Trim(ress.Choices[i].Delta.Content, "#*")
			}

		}
		w.Header().Set("content-type", "text/plain")
		fmt.Fprint(w, Assistant_Res)

		*history = append(*history, Delta{Role: "assistant", Content: Assistant_Res})

	}
}
func main() {
	history := []Delta{}
	history = append(history, Delta{Role: "system", Content: "You are a helpful assistant who gives answer straight to the point  without making a huge answer."})

	apikey := "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"
	url := "https://integrate.api.nvidia.com/v1/chat/completions"

	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		send(w, r, apikey, url, &history)
	})
	http.Handle("/", http.FileServer(http.Dir("echo")))
	http.ListenAndServe(":8080", nil)
}
