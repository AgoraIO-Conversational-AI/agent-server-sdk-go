package main

import (
    "fmt"
    "go/parser"
    "go/token"
    "os"
    "path/filepath"
    "regexp"
    "strings"
)

var (
    // `concepts` and `reference` snippets must declare whether they are runnable examples or API fragments.
    codeBlockRE = regexp.MustCompile("(?:(<!--\\s*snippet:\\s*(executable|fragment)\\s*-->)[ \\t]*\\n)?```go\\n([\\s\\S]*?)```")
    bannedPatterns = []*regexp.Regexp{
        regexp.MustCompile(`\{\{\s*owner\s*\}\}`),
        regexp.MustCompile(`\{\{\s*repo\s*\}\}`),
        regexp.MustCompile(`from agora-agent-server-sdk`),
    }
)

func main() {
    root, err := os.Getwd()
    if err != nil {
        panic(err)
    }

    files := []string{filepath.Join(root, "README.md")}
    err = filepath.Walk(filepath.Join(root, "docs"), func(path string, info os.FileInfo, walkErr error) error {
        if walkErr != nil {
            return walkErr
        }
        if info.IsDir() {
            return nil
        }
        if strings.HasSuffix(path, ".md") {
            files = append(files, path)
        }
        return nil
    })
    if err != nil {
        panic(err)
    }

    failures := []string{}
    snippetCount := 0
    fragmentCount := 0
    for _, file := range files {
        content, readErr := os.ReadFile(file)
        if readErr != nil {
            failures = append(failures, fmt.Sprintf("%s: %v", rel(root, file), readErr))
            continue
        }
        text := string(content)

        for _, pattern := range bannedPatterns {
            if pattern.MatchString(text) {
                failures = append(failures, fmt.Sprintf("%s contains banned pattern: %s", rel(root, file), pattern.String()))
            }
        }

        matches := codeBlockRE.FindAllStringSubmatch(text, -1)
        for _, match := range matches {
            annotation := match[2]
            code := strings.TrimSpace(match[3])
            if isAnnotatedSection(root, file) && annotation == "" {
                failures = append(failures, fmt.Sprintf("%s contains an unannotated go snippet", rel(root, file)))
                continue
            }

            mode := snippetModeForGoSnippet(code, annotation)
            if mode == "fragment" {
                fragmentCount++
                continue
            }
            snippetCount++
            if !validateGoSnippet(code) {
                failures = append(failures, fmt.Sprintf("%s contains an invalid Go snippet", rel(root, file)))
            }
        }
    }

    if snippetCount == 0 {
        failures = append(failures, "No Go code blocks found in README/docs markdown.")
    }

    if len(failures) > 0 {
        fmt.Fprintln(os.Stderr, "Documentation validation failed:")
        for _, failure := range failures {
            fmt.Fprintf(os.Stderr, "- %s\n", failure)
        }
        os.Exit(1)
    }

    fmt.Printf("Validated %d executable and %d fragment Go snippets across %d markdown files.\n", snippetCount, fragmentCount, len(files))
}

func validateGoSnippet(code string) bool {
    candidates := []string{code}
    candidates = append(candidates, "package main\n\n"+code)
    candidates = append(candidates, "package main\n\nfunc _snippet() {\n"+indent(code)+"\n}\n")
    if candidate, ok := buildImportAwareCandidate(code); ok {
        candidates = append(candidates, candidate)
    }

    for _, candidate := range candidates {
        if _, err := parser.ParseFile(token.NewFileSet(), "snippet.go", candidate, parser.AllErrors); err == nil {
            return true
        }
    }
    return false
}

func shouldValidateGoSnippet(code string) bool {
    return !strings.Contains(code, "...") && (strings.Contains(code, "package main") || strings.Contains(code, "func main()"))
}

func snippetModeForGoSnippet(code string, annotation string) string {
    if annotation == "fragment" {
        return "fragment"
    }
    if annotation == "executable" {
        return "executable"
    }
    if shouldValidateGoSnippet(code) {
        return "executable"
    }
    return "fragment"
}

func buildImportAwareCandidate(code string) (string, bool) {
    lines := strings.Split(code, "\n")
    importLines := []string{}
    bodyLines := []string{}

    for _, rawLine := range lines {
        trimmed := strings.TrimSpace(rawLine)
        if trimmed == "" {
            continue
        }
        if strings.HasPrefix(trimmed, "//") {
            continue
        }
        if isImportSpec(trimmed) {
            importLines = append(importLines, normalizeImportSpec(trimmed))
            continue
        }
        bodyLines = append(bodyLines, rawLine)
    }

    if len(importLines) == 0 {
        return "", false
    }

    candidate := "package main\n\nimport (\n"
    for _, line := range importLines {
        candidate += "\t" + line + "\n"
    }
    candidate += ")\n"
    if len(bodyLines) == 0 {
        return candidate, true
    }
    candidate += "\nfunc _snippet() {\n" + indent(strings.Join(bodyLines, "\n")) + "\n}\n"
    return candidate, true
}

func isImportSpec(line string) bool {
    return strings.HasPrefix(line, "\"") ||
        strings.HasPrefix(line, "import ") ||
        regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*\s+"`).MatchString(line)
}

func normalizeImportSpec(line string) string {
    trimmed := strings.TrimSpace(line)
    return strings.TrimPrefix(trimmed, "import ")
}

func indent(code string) string {
    lines := strings.Split(code, "\n")
    for i, line := range lines {
        lines[i] = "\t" + line
    }
    return strings.Join(lines, "\n")
}

func rel(root string, path string) string {
    relative, err := filepath.Rel(root, path)
    if err != nil {
        return path
    }
    return relative
}

func isAnnotatedSection(root string, path string) bool {
    normalized := filepath.ToSlash(rel(root, path))
    return strings.Contains(normalized, "docs/concepts/") || strings.Contains(normalized, "docs/reference/")
}
