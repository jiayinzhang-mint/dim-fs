package utils

import "os"

// WriteToFile takes in a file pointer and byte array and writes the byte array into the file
// returns error if pointer is nil or error in writing to file
func WriteToFile(f *os.File, data []byte) error {
	w := 0
	n := len(data)
	for {

		nw, err := f.Write(data[w:])
		if err != nil {
			return err
		}
		w += nw
		if nw >= n {
			return nil
		}
	}
}
