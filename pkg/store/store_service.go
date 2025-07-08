package store

import (
	"fmt"
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

func SwipeSecret(secret UserSecret, deleteAfter int) {
	filePath := filepath.Join(os.TempDir(), secret.Name+".txt")
	time.AfterFunc(time.Duration(deleteAfter)*time.Second, func() {
		err := os.Remove(filePath)
		if err != nil {
			fmt.Printf("Failed to delete %s: %v\n", filePath, err)
		} else {
			fmt.Printf("Secret %s deleted.\n", secret.Name)
		}
	})
}
