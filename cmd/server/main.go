package main

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"

	"raider/internal/db"
	"raider/internal/web"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

// getEnv fetches environment variable or returns fallback
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustLoadCA(path string) *x509.CertPool {
	caCert, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read CA cert: %v", err)
	}

	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(caCert) {
		log.Fatal("failed to append CA cert")
	}
	return pool
}

func mustLoadServerCert(certFile, keyFile string) tls.Certificate {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf("failed to load server cert/key: %v", err)
	}
	return cert
}

func main() {
	_ = godotenv.Load()

	host := getEnv("WEB_HOST", "0.0.0.0")
	port := getEnv("WEB_PORT", "8080")
	dbPath := getEnv("DB_PATH", "/tmp/raider.db")
	caPath := getEnv("CA_PATH", "/tmp/ca.pem")
	certPath := getEnv("CERT_PATH", "/tmp/server.pem")
	keyPath := getEnv("KEY_PATH", "/tmp/server.key")
	staticPath := getEnv("STATIC_PATH", "/tmp/static/")

	db.InitDB(dbPath)

	engine := html.New("./internal/web/templates", ".html")
	engine.Reload(true)

	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Static("/static", staticPath)

	web.SetupRoutes(app)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{
			mustLoadServerCert(certPath, keyPath),
		},
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs:  mustLoadCA(caPath),
		MinVersion: tls.VersionTLS13,
	}

	ln, err := tls.Listen("tcp", host+":"+port, tlsConfig)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Secure Raider running at https://%s:%s\n", host, port)
	log.Fatal(app.Listener(ln))
}
