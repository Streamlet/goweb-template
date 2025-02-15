package utility

import (
	"io"
	"io/fs"
	"path/filepath"
	"regexp"
	"time"
)

func NewSsiFS(originalFs, includeFs fs.FS, supportSsiExt []string, includePattern string) fs.FS {
	supportSsiExtMap := map[string]bool{}
	for _, ext := range supportSsiExt {
		supportSsiExtMap[ext] = true
	}
	return &ssiFs{
		originalFs:     originalFs,
		includeFs:      includeFs,
		supportSsiExt:  supportSsiExtMap,
		includePattern: regexp.MustCompile(includePattern),
	}
}

type ssiFs struct {
	originalFs     fs.FS
	includeFs      fs.FS
	supportSsiExt  map[string]bool
	includePattern *regexp.Regexp
}

func (sf *ssiFs) Open(name string) (fs.File, error) {
	file, err := sf.originalFs.Open(name)
	if err != nil {
		return nil, err
	}
	if _, ok := sf.supportSsiExt[filepath.Ext(name)]; !ok {
		return file, nil
	}
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	replacedBytes := sf.includePattern.ReplaceAllFunc(bytes, func(match []byte) []byte {
		groups := sf.includePattern.FindSubmatch(match)
		if len(groups) <= 1 {
			return match
		}
		includeFile, err := sf.includeFs.Open(string(groups[1]))
		if err != nil {
			return nil
		}
		includeBytes, err := io.ReadAll(includeFile)
		if err != nil {
			return nil
		}
		return includeBytes
	})
	return newVirtualFile(file, replacedBytes), nil
}

func newVirtualFile(file fs.File, bytes []byte) fs.File {
	return &virtualFile{
		file:    file,
		content: bytes,
		read:    0,
	}
}

type virtualFile struct {
	file    fs.File
	content []byte
	read    int
}

func (vf *virtualFile) Stat() (fs.FileInfo, error) {
	fileInfo, err := vf.file.Stat()
	if err != nil {
		return nil, err
	}
	return newVirtualFileInfo(fileInfo, int64(len(vf.content))), nil
}

func (vf *virtualFile) Read(bytes []byte) (int, error) {
	bufferLen := len(bytes)
	remainLen := len(vf.content) - vf.read
	if remainLen <= 0 {
		return 0, io.EOF
	}
	copyLen := 0
	if remainLen <= bufferLen {
		copyLen = remainLen
	} else {
		copyLen = bufferLen
	}
	copy(bytes, vf.content[vf.read:vf.read+copyLen])
	vf.read += copyLen
	return copyLen, nil
}

func (vf *virtualFile) Close() error {
	return vf.file.Close()
}

func newVirtualFileInfo(fileInfo fs.FileInfo, newSize int64) fs.FileInfo {
	return &virtualFileInfo{
		fileInfo: fileInfo,
		size:     newSize,
	}
}

type virtualFileInfo struct {
	fileInfo fs.FileInfo
	size     int64
}

func (vfi *virtualFileInfo) Name() string {
	return vfi.fileInfo.Name()
}

func (vfi *virtualFileInfo) Size() int64 {
	return vfi.size
}

func (vfi *virtualFileInfo) Mode() fs.FileMode {
	return vfi.fileInfo.Mode()
}

func (vfi *virtualFileInfo) ModTime() time.Time {
	return vfi.fileInfo.ModTime()
}

func (vfi *virtualFileInfo) IsDir() bool {
	return vfi.fileInfo.IsDir()
}

func (vfi *virtualFileInfo) Sys() any {
	return vfi.fileInfo.Sys()
}
