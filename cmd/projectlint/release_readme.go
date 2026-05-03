package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// Строка в корневом README: Текущий релиз: `v0.1.12` (пример формата).
var readmeReleaseRE = regexp.MustCompile("Текущий релиз:\\s*`(v[0-9]+\\.[0-9]+\\.[0-9]+)`")

// lintReadmeReleaseVersusGit сравнивает semver из README с максимальным локальным тегом v*.*.*.
// Если README отстаёт от git — это finding и явное сообщение в stderr (документация не источник истины).
// Отключение: переменная окружения PROJECTLINT_SKIP_README_RELEASE не пустая.
func lintReadmeReleaseVersusGit(root string) []finding {
	if strings.TrimSpace(os.Getenv("PROJECTLINT_SKIP_README_RELEASE")) != "" {
		return nil
	}

	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil
	}
	if _, err := os.Stat(filepath.Join(absRoot, ".git")); err != nil {
		return nil
	}

	readmePath := filepath.Join(absRoot, "README.md")
	data, err := os.ReadFile(readmePath)
	if err != nil {
		return nil
	}

	lines := strings.Split(string(data), "\n")
	var readmeLine int
	var readmeVersion string
	for i, line := range lines {
		m := readmeReleaseRE.FindStringSubmatch(line)
		if m != nil {
			readmeLine = i + 1
			readmeVersion = m[1]
			break
		}
	}
	if readmeLine == 0 {
		return nil
	}

	maxTag, ok := gitMaxSemverTag(absRoot)
	if !ok || maxTag == "" {
		return nil
	}

	if semverCmp(readmeVersion, maxTag) >= 0 {
		return nil
	}

	fmt.Fprintf(os.Stderr, "projectlint: README отстаёт от git — в README %s, последний локальный тег %s (строка %d README.md). Документация не источник истины: не выводить версию только из текста README, сверяться с git fetch --tags и тегами.\n",
		readmeVersion, maxTag, readmeLine)

	msg := "версия «Текущий релиз» в README (" + readmeVersion + ") отстаёт от последнего git-тега (" + maxTag + "): документация не источник истины — сделайте git fetch --tags и выставьте в README версию по фактическим тегам или следующий согласованный релиз"

	rel, err := filepath.Rel(absRoot, readmePath)
	if err != nil || rel == "" || rel == "." {
		rel = "README.md"
	}

	return []finding{{
		path:    rel,
		line:    readmeLine,
		message: msg,
	}}
}

func gitMaxSemverTag(gitRoot string) (string, bool) {
	cmd := exec.Command("git", "-C", gitRoot, "tag", "-l", "v*")
	out, err := cmd.Output()
	if err != nil {
		return "", false
	}

	var best string
	var bestParts [3]int
	var has bool

	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		t := strings.TrimSpace(line)
		if t == "" {
			continue
		}
		parts, ok := parseSemverV(t)
		if !ok {
			continue
		}
		if !has || semverLess(bestParts, parts) {
			has = true
			best = t
			bestParts = parts
		}
	}
	return best, has
}

func parseSemverV(tag string) ([3]int, bool) {
	if !strings.HasPrefix(tag, "v") {
		return [3]int{}, false
	}
	parts := strings.Split(tag[1:], ".")
	if len(parts) != 3 {
		return [3]int{}, false
	}
	var n [3]int
	for i := range 3 {
		v, err := strconv.Atoi(parts[i])
		if err != nil || v < 0 {
			return [3]int{}, false
		}
		n[i] = v
	}
	return n, true
}

func semverLess(a, b [3]int) bool {
	for i := range 3 {
		if a[i] != b[i] {
			return a[i] < b[i]
		}
	}
	return false
}

// semverCmp сравнивает теги vX.Y.Z: -1 если a<b, 0 если равны, 1 если a>b.
func semverCmp(aTag, bTag string) int {
	ap, okA := parseSemverV(aTag)
	bp, okB := parseSemverV(bTag)
	if !okA || !okB {
		return 0
	}
	switch {
	case semverLess(ap, bp):
		return -1
	case semverLess(bp, ap):
		return 1
	default:
		return 0
	}
}
