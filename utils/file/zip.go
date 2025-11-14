package file

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Zip 文件压缩
func Zip(srcFile string, destZip string) error {
	zipFile, err := os.Create(destZip)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	err = filepath.Walk(srcFile, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = strings.TrimPrefix(path, filepath.Dir(srcFile)+"/")
		// header.Name = path
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}
		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}
		return err
	})
	if err != nil {
		return err
	}

	return err
}

// Unzip 文件解压
func Unzip(zipFile, destDir string) error {
	zipReader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	destDirAbs, err := filepath.Abs(destDir)
	if err != nil {
		return err
	}
	// Ensure destDirAbs ends with a path separator
	if !strings.HasSuffix(destDirAbs, string(os.PathSeparator)) {
		destDirAbs += string(os.PathSeparator)
	}

	for _, f := range zipReader.File {
		fpath := filepath.Join(destDir, f.Name)
		// Clean and get absolute path
		fpathAbs, err := filepath.Abs(filepath.Clean(fpath))
		if err != nil {
			return err
		}
		// Check for Zip Slip (directory traversal)
		if !strings.HasPrefix(fpathAbs, destDirAbs) {
			return 	// or: return fmt.Errorf("illegal file path: %s", fpathAbs)
		}

		if f.FileInfo().IsDir() {
			err := os.MkdirAll(fpathAbs, os.ModePerm)
			if err != nil {
				return err
			}
		} else {
			if err = os.MkdirAll(filepath.Dir(fpathAbs), os.ModePerm); err != nil {
				return err
			}

			inFile, err := f.Open()
			if err != nil {
				return err
			}
			defer inFile.Close()

			outFile, err := os.OpenFile(fpathAbs, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
