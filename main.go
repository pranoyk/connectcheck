package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Hello! Welcome to the Go HTTP service.\n")
}

func googleCallHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get("https://www.google.com")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error calling google.com: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading response: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Successfully called google.com\n")
	fmt.Fprintf(w, "Status Code: %d\n", resp.StatusCode)
	fmt.Fprintf(w, "Content Length: %d bytes\n", len(body))
	
	preview := 200
	if len(body) < preview {
		preview = len(body)
	}
	fmt.Fprintf(w, "\nFirst %d characters of response:\n%s\n", preview, string(body[:preview]))
}

func delayAPIOnShutdown() error {
	log.Println("========================================")
	log.Println("Shutdown signal received. Making call to httpbin.org/delay/10 (this will take 10 seconds)...")
	log.Println("========================================")
	
	startTime := time.Now()
	
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Using httpbin.org/delay/10 for a guaranteed 10-second response time
	resp, err := client.Get("https://httpbin.org/delay/10")
	if err != nil {
		log.Printf("ERROR: Failed to call httpbin.org: %v", err)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ERROR: Failed to read response: %v", err)
		return err
	}

	elapsed := time.Since(startTime)
	log.Println("========================================")
	log.Printf("SUCCESS: Called httpbin.org/delay/10")
	log.Printf("  Status Code: %d", resp.StatusCode)
	log.Printf("  Response Length: %d bytes", len(body))
	log.Printf("  Time taken: %v", elapsed)
	log.Println("========================================")
	return nil
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/google", googleCallHandler)

	port := ":8080"
	server := &http.Server{
		Addr: port,
	}

	// Channel to listen for shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, 
		syscall.SIGINT,  // Ctrl+C
		syscall.SIGTERM, // Kubernetes sends SIGTERM
		syscall.SIGQUIT,
		syscall.SIGHUP,
	)

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %s", port)
		log.Printf("Endpoints available:")
		log.Printf("  - /hello")
		log.Printf("  - /google")
		
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for shutdown signal
	sig := <-sigChan
	log.Printf("Received signal: %v", sig)

	// Call YouTube before shutting down
	if err := delayAPIOnShutdown(); err != nil {
		log.Printf("Warning: delayAPIOnShutdown call failed, but continuing shutdown: %v", err)
	}

	// Gracefully shutdown the server with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Println("Shutting down server gracefully...")
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("Server stopped")
}