package command

import (
	"testing"

	"github.com/golonzovsky/kubecolor/kubectl"
	"github.com/golonzovsky/kubecolor/testutil"
)

func Test_ResolveSubcommand(t *testing.T) {
	tests := []struct {
		name                   string
		args                   []string
		conf                   *KubecolorConfig
		isOutputTerminal       func() bool
		expectedShouldColorize bool
		expectedInfo           *kubectl.SubcommandInfo
	}{
		{
			name:             "basic case",
			args:             []string{"get", "pods"},
			isOutputTerminal: func() bool { return true },
			conf: &KubecolorConfig{
				Plain:      false,
				ForceColor: false,
				KubectlCmd: "kubectl",
			},
			expectedShouldColorize: true,
			expectedInfo:           &kubectl.SubcommandInfo{Subcommand: kubectl.Get},
		},
		{
			name:             "when plain, it won't colorize",
			args:             []string{"get", "pods"},
			isOutputTerminal: func() bool { return true },
			conf: &KubecolorConfig{
				Plain:      true,
				ForceColor: false,
				KubectlCmd: "kubectl",
			},
			expectedShouldColorize: false,
			expectedInfo:           &kubectl.SubcommandInfo{Subcommand: kubectl.Get},
		},
		{
			name:             "when help, it will colorize",
			args:             []string{"get", "pods", "-h"},
			isOutputTerminal: func() bool { return true },
			conf: &KubecolorConfig{
				Plain:      false,
				ForceColor: false,
				KubectlCmd: "kubectl",
			},
			expectedShouldColorize: true,
			expectedInfo:           &kubectl.SubcommandInfo{Subcommand: kubectl.Get, Help: true},
		},
		{
			name:             "when both plain and force, plain is chosen",
			args:             []string{"get", "pods"},
			isOutputTerminal: func() bool { return true },
			conf: &KubecolorConfig{
				Plain:      true,
				ForceColor: true,
				KubectlCmd: "kubectl",
			},
			expectedShouldColorize: false,
			expectedInfo:           &kubectl.SubcommandInfo{Subcommand: kubectl.Get},
		},
		{
			name:             "when no subcommand is found, it becomes help",
			args:             []string{},
			isOutputTerminal: func() bool { return true },
			conf: &KubecolorConfig{
				Plain:      false,
				ForceColor: false,
				KubectlCmd: "kubectl",
			},
			expectedShouldColorize: true,
			expectedInfo:           &kubectl.SubcommandInfo{Help: true},
		},
		{
			name:             "when not tty, it won't colorize",
			args:             []string{"get", "pods"},
			isOutputTerminal: func() bool { return false },
			conf: &KubecolorConfig{
				Plain:      false,
				ForceColor: false,
				KubectlCmd: "kubectl",
			},
			expectedShouldColorize: false,
			expectedInfo:           &kubectl.SubcommandInfo{Subcommand: kubectl.Get},
		},
		{
			name:             "even if not tty, if force, it colorizes",
			args:             []string{"get", "pods"},
			isOutputTerminal: func() bool { return false },
			conf: &KubecolorConfig{
				Plain:      false,
				ForceColor: true,
				KubectlCmd: "kubectl",
			},
			expectedShouldColorize: true,
			expectedInfo:           &kubectl.SubcommandInfo{Subcommand: kubectl.Get},
		},
		{
			name:             "when the subcommand is unsupported, it won't colorize",
			args:             []string{"-h"},
			isOutputTerminal: func() bool { return true },
			conf: &KubecolorConfig{
				Plain:      false,
				ForceColor: false,
				KubectlCmd: "kubectl",
			},
			expectedShouldColorize: true,
			expectedInfo:           &kubectl.SubcommandInfo{Help: true},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			isOutputTerminal = tt.isOutputTerminal
			shouldColorize, info := ResolveSubcommand(tt.args, tt.conf)
			testutil.MustEqual(t, tt.expectedShouldColorize, shouldColorize)
			testutil.MustEqual(t, tt.expectedInfo, info)
		})
	}
}
