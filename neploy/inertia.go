package neploy

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	inertia "github.com/romsar/gonertia"
)

func initInertia() *inertia.Inertia {
	viteHotFile := "./public/hot"
	rootViewFile := "./resources/views/root.html"
	manifestPath := "./public/build/manifest.json"
	viteManifestPath := "./public/build/.vite/manifest.json"

	// Verificamos si existe el archivo hot (modo dev)
	_, err := os.Stat(viteHotFile)
	isDevMode := err == nil

	// Siempre intentamos mover el archivo de manifest si no existe
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		log.Println("manifest.json no existe. Intentando mover desde .vite")

		if _, err := os.Stat(viteManifestPath); os.IsNotExist(err) {
			log.Println("No se encontró el archivo .vite/manifest.json:", err)
			return nil
		}

		if err := os.Rename(viteManifestPath, manifestPath); err != nil {
			log.Println("Error al mover el archivo .vite/manifest.json:", err)
			return nil
		}

	}

	if isDevMode {
		i, err := inertia.NewFromFile(
			rootViewFile,
			inertia.WithSSR(),
		)
		if err != nil {
			log.Fatal(err)
		}
		i.ShareTemplateFunc("vite", func(entry string) (string, error) {
			content, err := os.ReadFile(viteHotFile)
			if err != nil {
				return "", err
			}
			url := strings.TrimSpace(string(content))
			if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
				url = url[strings.Index(url, ":")+1:]
			} else {
				url = "//localhost:8080"
			}
			if entry != "" && !strings.HasPrefix(entry, "/") {
				entry = "/" + entry
			}
			return url + entry, nil
		})
		return i
	}

	// Si no está en modo dev, entonces cargamos desde el archivo manifest en producción
	i, err := inertia.NewFromFile(
		rootViewFile,
		inertia.WithVersionFromFile(manifestPath),
		inertia.WithSSR(),
	)
	if err != nil {
		log.Fatal(err)
	}

	i.ShareTemplateFunc("vite", vite(manifestPath, "/build/"))
	return i
}

func vite(manifestPath, buildDir string) func(path string) (string, error) {
	f, err := os.Open(manifestPath)
	if err != nil {
		log.Fatalf("No se puede abrir el archivo manifest: %s", err)
	}
	defer f.Close()

	viteAssets := make(map[string]*struct {
		File   string `json:"file"`
		Source string `json:"src"`
	})

	err = json.NewDecoder(f).Decode(&viteAssets)
	if err != nil {
		log.Fatalf("Error al decodificar el archivo manifest: %s", err)
	}

	return func(p string) (string, error) {
		if val, ok := viteAssets[p]; ok {
			return path.Join("/", buildDir, val.File), nil
		}
		return "", fmt.Errorf("asset %q no encontrado", p)
	}
}
