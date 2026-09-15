package mirror

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRealRepoParity 用真实仓库做 bash 版发布物的对账:本地克隆 canonical
// 仓库,以 Go 实现向本地 bare 镜像发布同一 tag,断言合成 tree 与 bash 版
// 已发布产物完全一致(证明双实现对"同 tag 内容确定性一致"的承诺)。
//
// 环境变量触发:
//
//	MIRROR_PARITY_SRC   canonical 仓库本地路径
//	MIRROR_PARITY_TAG   待对账 tag
//	MIRROR_PARITY_TREE  bash 版发布 commit 的 tree hash
//	MIRROR_PARITY_FROM / MIRROR_PARITY_TO  模块身份映射
//	MIRROR_PARITY_ORIG_COMMIT 可选;bash 版曾在 worktree 里 tag -f,
//	污染了本地 tag refs,克隆后需以 origin 的原始 commit 重打 tag
func TestRealRepoParity(t *testing.T) {
	src := os.Getenv("MIRROR_PARITY_SRC")
	wantTree := os.Getenv("MIRROR_PARITY_TREE")
	tag := os.Getenv("MIRROR_PARITY_TAG")
	from := os.Getenv("MIRROR_PARITY_FROM")
	to := os.Getenv("MIRROR_PARITY_TO")
	origCommit := os.Getenv("MIRROR_PARITY_ORIG_COMMIT")
	if src == "" || wantTree == "" || tag == "" || from == "" || to == "" {
		t.Skip("跳过实仓对账(设置 MIRROR_PARITY_* 环境变量触发)")
	}

	work := t.TempDir()
	ctx := context.Background()

	cloneDir := filepath.Join(work, "clone")
	_, err := git.PlainCloneContext(ctx, cloneDir, false, &git.CloneOptions{URL: src})
	require.NoError(t, err)

	clone, err := git.PlainOpen(cloneDir)
	require.NoError(t, err)
	tagRef := plumbing.NewTagReferenceName(tag)
	if origCommit != "" {
		// 恢复被 bash 版覆盖的本地 tag:指向原始发布 commit
		require.NoError(t, clone.Storer.RemoveReference(tagRef))
		require.NoError(t, clone.Storer.SetReference(plumbing.NewHashReference(
			tagRef, plumbing.NewHash(origCommit))))
	}

	mirrorDir := filepath.Join(work, "mirror.git")
	_, err = git.PlainInit(mirrorDir, true)
	require.NoError(t, err)
	_, err = clone.CreateRemote(&config.RemoteConfig{Name: "github", URLs: []string{mirrorDir}})
	require.NoError(t, err)

	rep, err := Publish(ctx, Options{
		RepoDir: cloneDir,
		Mapping: Mapping{Source: from, Target: to},
		Tags:    []string{tag},
	})
	require.NoError(t, err)
	require.Len(t, rep.Tags, 1)

	mirror, err := git.PlainOpen(mirrorDir)
	require.NoError(t, err)
	commit, err := resolveRefCommit(mirror, plumbing.NewTagReferenceName(tag))
	require.NoError(t, err)
	assert.Equal(t, wantTree, commit.TreeHash.String(),
		"Go 版合成 tree 应与 bash 版已发布产物逐字节一致")
	assert.Equal(t, tag, rep.Tags[0].Tag)
	t.Logf("sourceCommit=%s commit=%s replaced=%d",
		rep.Tags[0].SourceCommit, rep.Tags[0].Commit, len(rep.Tags[0].ReplacedFiles))
}
