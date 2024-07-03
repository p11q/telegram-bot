package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"telegram_bot/lib/my_er"
	"telegram_bot/storage"
	"time"
)

const defaultPerm = 0774

type Storage struct {
	basePath string
}

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Page) (err error) {
	defer func() { err = my_er.WrapIfErr("can't save page", err) }()

	fPath := filepath.Join(s.basePath, page.UserName)

	if err := os.MkdirAll(fPath, defaultPerm); err != nil {
		return err
	}

	fName, err := fileName(page)
	if err != nil {
		return err
	}

	fPath = filepath.Join(fPath, fName)

	file, err := os.Create(fPath)
	if err != nil {
		return err
	}

	defer func() { _ = file.Close() }()

	if err := gob.NewEncoder(file).Encode((page)); err != nil {
		return err
	}
	return nil
}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() { err = my_er.WrapIfErr("can't pick random page", err) }()

	path := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, storage.ErrNoSavedPages
	}

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(files))

	file := files[n]

	return s.decoderPage(filepath.Join(path, file.Name()))
}

func (s Storage) Remove(page *storage.Page) error {
	fName, err := fileName(page)
	if err != nil {
		return my_er.Wrap("can't remove file", err)
	}

	path := filepath.Join(s.basePath, page.UserName, fName)

	if err := os.Remove(path); err != nil {
		return my_er.Wrap(fmt.Sprintf("can't remove file %s", path), err)
	}
	return nil
}

func (s Storage) IsExists(page *storage.Page) (bool, error) {
	fName, err := fileName(page)
	if err != nil {
		return false, my_er.Wrap("can't remove file", err)
	}

	path := filepath.Join(s.basePath, page.UserName, fName)

	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		return false, my_er.Wrap(fmt.Sprintf("cant't check if file %s exists", path), err)
	}

	return true, nil
}

func (s Storage) decoderPage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, my_er.Wrap("can't decode page", err)
	}

	defer func() { _ = f.Close() }()

	var p storage.Page

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, my_er.Wrap("can't decode page", err)
	}
	return &p, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
