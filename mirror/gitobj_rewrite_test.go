package mirror

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRewriteModulePath(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		source string
		target string
		want   string
	}{
		{
			name:   "go.mod 中 module 声明",
			data:   "module github.com/a/b\n\ngo 1.26\n",
			source: "github.com/a/b", target: "github.com/x/y",
			want: "module github.com/x/y\n\ngo 1.26\n",
		},
		{
			name:   "子包路径也被改写",
			data:   "import \"github.com/a/b/pkg/client\"\n",
			source: "github.com/a/b", target: "github.com/x/y",
			want: "import \"github.com/x/y/pkg/client\"\n",
		},
		{
			name:   "b-extra 不被误改",
			data:   "module github.com/a/b-extra\n",
			source: "github.com/a/b", target: "github.com/x/y",
			want: "module github.com/a/b-extra\n",
		},
		{
			name:   "b.cool 不被误改",
			data:   "require github.com/a/b.cool v1.0.0\n",
			source: "github.com/a/b", target: "github.com/x/y",
			want: "require github.com/a/b.cool v1.0.0\n",
		},
		{
			name:   "多次出现全部改写",
			data:   "module github.com/a/b\n\nimport \"github.com/a/b/internal/util\"\n",
			source: "github.com/a/b", target: "github.com/x/y",
			want: "module github.com/x/y\n\nimport \"github.com/x/y/internal/util\"\n",
		},
		{
			name:   "文本末尾无换行",
			data:   "module github.com/a/b",
			source: "github.com/a/b", target: "github.com/x/y",
			want: "module github.com/x/y",
		},
		{
			name:   "source 不存在时不改",
			data:   "module github.com/other/repo\n",
			source: "github.com/a/b", target: "github.com/x/y",
			want: "module github.com/other/repo\n",
		},
		{
			name:   "字符串引号内",
			data:   "\"github.com/a/b\"",
			source: "github.com/a/b", target: "github.com/x/y",
			want: "\"github.com/x/y\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := string(rewriteModulePath([]byte(tt.data), tt.source, tt.target))
			assert.Equal(t, tt.want, got)
		})
	}
}
