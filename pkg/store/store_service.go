package store

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

type UserSecret struct {
	Name  string
	Value string
}

func SaveSecret(secret *UserSecret, path string) error {
	filePath := filepath.Join(path, secret.Name+".txt")
	return os.WriteFile(filePath, []byte(secret.Value), 0600)
}

func GetSecret(name string, path string) (UserSecret, error) {
	filePath := filepath.Join(path, name+".txt")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return UserSecret{}, err
	}
	return UserSecret{
		Name:  name,
		Value: string(data),
	}, nil
}

func SwipeSecret(secret UserSecret, deleteAfter int, path string) {
	go func() {
		log.Printf("SwipeSecret started for %s, waiting %d seconds", secret.Name, deleteAfter)
		time.Sleep(time.Duration(deleteAfter) * time.Second)

		path := filepath.Join(path, secret.Name)
		if err := os.Remove(path); err != nil {
			log.Printf("Failed to delete secret %s: %v", secret.Name, err)
		} else {
			log.Printf("Deleted secret %s", secret.Name)
		}
	}()
}
