package logging

import (
	"bufio"
	"os"
	"sync"
)

type FileHandler struct {
	NullHandler

	Level     Level
	Formatter Formatter
	FilePath  string
	FileMode  os.FileMode

	file   *os.File
	writer *bufio.Writer
	lock   sync.RWMutex
}

func (handler *FileHandler) GetLevel() Level {
	return handler.Level
}

func (handler *FileHandler) open() error {
	if handler.writer != nil {
		return nil
	}

	file, err := os.OpenFile(handler.FilePath, os.O_APPEND|os.O_WRONLY, handler.FileMode)

	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		file, err = os.Create(handler.FilePath)
		if err != nil {
			return err
		}
	}

	handler.file = file
	handler.writer = bufio.NewWriter(handler.file)

	return nil
}

func (handler *FileHandler) Close() error {
	handler.lock.Lock()
	defer handler.lock.Unlock()

	return handler.close()
}

func (handler *FileHandler) close() error {
	if handler.writer != nil {
		if err := handler.writer.Flush(); err != nil {
			return err
		}
		handler.writer = nil
	}
	if handler.file != nil {
		if err := handler.file.Close(); err != nil {
			return err
		}
		handler.file = nil
	}
	return nil
}

func (handler *FileHandler) Handle(record *Record) error {
	handler.lock.Lock()
	defer handler.lock.Unlock()

	msg := handler.Formatter.Format(record) + "\n"

	if handler.writer == nil {

		if err := handler.open(); err != nil {
			return err
		}
	}

	_, err := handler.writer.Write([]byte(msg))
	if err != nil {
		return err
	}
	handler.writer.Flush()

	return nil
}
