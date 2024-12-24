package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type SQLRequest struct {
	Query string `json:"query"`
}

type SQLResponse struct {
	Results any    `json:"results,omitempty"`
	Error   string `json:"error,omitempty"`
}

func handlerDuckDB(w http.ResponseWriter, r *http.Request) {
	var req SQLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "duckdb-")
	if err != nil {
		http.Error(w, "Failed to create temporary directory", http.StatusInternalServerError)
		return
	}
	defer func() {
		os.RemoveAll(tempDir)
	}()

	// 使用管道将 SQL 内容传递给 duckdb 命令
	cmd := exec.Command("duckdb", "-json") // -d: set the database directory
	cmd.Dir = tempDir                      // 设置工作目录
	cmd.Stdin = strings.NewReader(req.Query)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Command execution failed: %v, output: %s", err, string(output))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var results any
	if err := json.Unmarshal(output, &results); err != nil {
		http.Error(w, "Failed to parse command output", http.StatusInternalServerError)
		return
	}

	resp := SQLResponse{Results: results}
	json.NewEncoder(w).Encode(resp)
}
