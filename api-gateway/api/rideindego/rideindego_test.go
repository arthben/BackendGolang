package rideindego

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/arthben/BackendGolang/api-gateway/internal/config"
	"github.com/arthben/BackendGolang/api-gateway/internal/database"
)

var (
	db database.DBService
)

func TestSearch(t *testing.T) {
	service := NewService(db)
	at := time.Date(2024, 11, 21, 0, 0, 0, 0, time.UTC)
	resp, _, err := service.Search(at, "3009")
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	fmt.Printf("resp: %v\n", resp)
}

func TestFetch(t *testing.T) {
	t.Run("test function FetchData", func(t *testing.T) {
		service := NewService(db)
		httpCode, err := service.FetchAndStore()
		if err != nil {
			t.Errorf("Expected no error, but error occur %s\n", err)
			return
		}

		if httpCode != http.StatusOK {
			t.Errorf("Expected status %d but got response %d", http.StatusOK, httpCode)
			return
		}
	})
}

func init() {
	os.Chdir("../..")

	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(0)
	}

	repo, err := database.NewPool(cfg)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(0)
	}

	db = repo
}
