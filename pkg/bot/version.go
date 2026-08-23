package bot

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
)

var (
	// GitCommit and other build info can optionally be set via -ldflags
	GitCommit = ""
	BuildTime = ""
	RepoURL   = "https://github.com/justprox/ye-userinfo-bot"
)

// BuildInfo stores versioning and transparency metadata.
type BuildInfo struct {
	Revision   string
	CommitTime string
	Modified   bool
	GoVersion  string
	BinarySHA  string
	RepoURL    string
}

// GetBuildInfo retrieves Git VCS information and binary hash.
func GetBuildInfo() BuildInfo {
	info := BuildInfo{
		GoVersion: runtime.Version(),
		RepoURL:   RepoURL,
	}

	if GitCommit != "" {
		info.Revision = GitCommit
	}
	if BuildTime != "" {
		info.CommitTime = BuildTime
	}

	// Read VCS settings embedded by Go compiler
	if bi, ok := debug.ReadBuildInfo(); ok {
		for _, s := range bi.Settings {
			switch s.Key {
			case "vcs.revision":
				if info.Revision == "" {
					info.Revision = s.Value
				}
			case "vcs.time":
				if info.CommitTime == "" {
					info.CommitTime = s.Value
				}
			case "vcs.modified":
				info.Modified = s.Value == "true"
			}
		}
	}

	if info.Revision == "" || info.Revision == "unknown" {
		if vcsSha := strings.TrimSpace(os.Getenv("VERCEL_GIT_COMMIT_SHA")); vcsSha != "" {
			info.Revision = vcsSha
		} else {
			info.Revision = "unknown"
		}
	}

	// Calculate SHA-256 of the running binary file itself
	if exePath, err := os.Executable(); err == nil {
		if file, err := os.Open(exePath); err == nil {
			defer file.Close()
			h := sha256.New()
			if _, err := io.Copy(h, file); err == nil {
				info.BinarySHA = hex.EncodeToString(h.Sum(nil))
			}
		}
	}

	return info
}

// formatVersion returns formatted transparency and build information.
func formatVersion(lang Lang) string {
	b := GetBuildInfo()
	var sb strings.Builder

	if lang == LangRU {
		sb.WriteString("ℹ️ <b>О боте и сборке (Transparency Info):</b>\n\n")
		sb.WriteString(fmt.Sprintf("🔨 <b>Git Revision:</b> <code>%s</code>", escape(b.Revision)))
		if b.Modified {
			sb.WriteString(" <i>(modified)</i>")
		}
		sb.WriteString("\n")
		if b.CommitTime != "" {
			sb.WriteString(fmt.Sprintf("🕒 <b>Commit Date:</b> <code>%s</code>\n", escape(b.CommitTime)))
		}
		sb.WriteString(fmt.Sprintf("⚙️ <b>Go Version:</b> <code>%s</code> (%s/%s)\n", escape(b.GoVersion), runtime.GOOS, runtime.GOARCH))
		if b.BinarySHA != "" {
			sb.WriteString(fmt.Sprintf("🔒 <b>Binary SHA-256:</b> <code>%s</code>\n", escape(b.BinarySHA)))
		}
		if b.RepoURL != "" {
			sb.WriteString(fmt.Sprintf("📦 <b>Source Code:</b> <a href=\"%s\">%s</a>\n", escape(b.RepoURL), escape(b.RepoURL)))
		}
		sb.WriteString("\n💡 <i>Бот работает с нулевым логированием (Zero-Logging) и полностью открытым исходным кодом.</i>")
	} else {
		sb.WriteString("ℹ️ <b>About & Build (Transparency Info):</b>\n\n")
		sb.WriteString(fmt.Sprintf("🔨 <b>Git Revision:</b> <code>%s</code>", escape(b.Revision)))
		if b.Modified {
			sb.WriteString(" <i>(modified)</i>")
		}
		sb.WriteString("\n")
		if b.CommitTime != "" {
			sb.WriteString(fmt.Sprintf("🕒 <b>Commit Date:</b> <code>%s</code>\n", escape(b.CommitTime)))
		}
		sb.WriteString(fmt.Sprintf("⚙️ <b>Go Version:</b> <code>%s</code> (%s/%s)\n", escape(b.GoVersion), runtime.GOOS, runtime.GOARCH))
		if b.BinarySHA != "" {
			sb.WriteString(fmt.Sprintf("🔒 <b>Binary SHA-256:</b> <code>%s</code>\n", escape(b.BinarySHA)))
		}
		if b.RepoURL != "" {
			sb.WriteString(fmt.Sprintf("📦 <b>Source Code:</b> <a href=\"%s\">%s</a>\n", escape(b.RepoURL), escape(b.RepoURL)))
		}
		sb.WriteString("\n💡 <i>This bot operates with zero logging and fully open-source verifiable code.</i>")
	}

	return sb.String()
}
