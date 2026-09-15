package mirror

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/filemode"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

// materializeTree 把 storer 中以 treeHash 为根的树写成目录内容,
// 供 go build 门禁编译。gitlink(子模块)跳过;symlink 原样重建。
func materializeTree(st storer.EncodedObjectStorer, treeHash plumbing.Hash, dir string) error {
	enc, err := st.EncodedObject(plumbing.TreeObject, treeHash)
	if err != nil {
		return err
	}
	tree := &object.Tree{}
	if err := tree.Decode(enc); err != nil {
		return err
	}

	for _, entry := range tree.Entries {
		path := filepath.Join(dir, entry.Name)
		switch entry.Mode {
		case filemode.Dir:
			if err := os.MkdirAll(path, 0o755); err != nil {
				return err
			}
			if err := materializeTree(st, entry.Hash, path); err != nil {
				return err
			}
		case filemode.Regular, filemode.Executable:
			data, err := storerBlobBytes(st, entry.Hash)
			if err != nil {
				return err
			}
			mode := os.FileMode(0o644)
			if entry.Mode == filemode.Executable {
				mode = 0o755
			}
			if err := os.WriteFile(path, data, mode); err != nil {
				return err
			}
		case filemode.Symlink:
			data, err := storerBlobBytes(st, entry.Hash)
			if err != nil {
				return err
			}
			_ = os.Remove(path)
			if err := os.Symlink(string(data), path); err != nil {
				return fmt.Errorf("重建 symlink %s 失败: %w", entry.Name, err)
			}
		default:
			// gitlink(子模块):go get 不涉及子模块,跳过
		}
	}
	return nil
}

func storerBlobBytes(st storer.EncodedObjectStorer, h plumbing.Hash) ([]byte, error) {
	enc, err := st.EncodedObject(plumbing.BlobObject, h)
	if err != nil {
		return nil, err
	}
	r, err := enc.Reader()
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return io.ReadAll(r)
}
