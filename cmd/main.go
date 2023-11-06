package main

import (
	"fmt"
	"io/fs"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// todo do not read the complete file
// config other programming languages
// - file extension(s), version ID incl. const, when to stop scanning (after 1st func in go)

const AutoVersion = "0.1.38"

var logger *slog.Logger

type LanguageDefinition struct {
	Language          string   // Name of the programming language
	FileExtensions    []string // Searching for files with this extension
	VersionDefinition string   // for a text starting with
	ScanStop          string   // stop searching after finding this, empty string results in NO scan stop
}

type SourceFile struct {
	FilePath string // the path and filename of the source file
	LangDef  int    // index of language of file in LanguageDefinition
}

func buildLanguageDefinitions() []LanguageDefinition {
	return []LanguageDefinition{
		{
			Language:          "Go",
			FileExtensions:    []string{"go"},
			VersionDefinition: "const AutoVersion",
			ScanStop:          "func",
		},
		{
			Language:          "Javascript",
			FileExtensions:    []string{"js", "mjs", "cjs"},
			VersionDefinition: "const AutoVersion",
			ScanStop:          "function",
		},
		{
			Language:          "Python",
			FileExtensions:    []string{"py", "pyw"},
			VersionDefinition: "AUTO_VERSION",
			ScanStop:          "def",
		},
		{
			Language:          "Typescript",
			FileExtensions:    []string{"ts", "tsx", "mts", "cts"},
			VersionDefinition: "const AutoVersion",
			ScanStop:          "func",
		},
		{
			Language:          "C-sharp",
			FileExtensions:    []string{"cs", "csx"},
			VersionDefinition: "const AutoVersion",
			ScanStop:          "func",
		},
		{
			Language:          "Java",
			FileExtensions:    []string{"java"},
			VersionDefinition: "public static final String AUTO_VERSION",
			ScanStop:          "public class", // todo protected class, private class or just search for class
		},
		{
			Language:          "Rust",
			FileExtensions:    []string{"rs", "rlib"},
			VersionDefinition: "const AUTO_VERSION: &str",
			ScanStop:          "fn",
		},
		{
			Language:          "Kotlin",
			FileExtensions:    []string{"kt"},
			VersionDefinition: "const val AUTO_VERSION",
			ScanStop:          "fun",
		},
		{
			Language:          "C",
			FileExtensions:    []string{"c"},
			VersionDefinition: "#define AUTOVERSION",
			ScanStop:          "",
		},
		{
			Language:          "C++",
			FileExtensions:    []string{"cc", "cpp", "cxx", "c++"},
			VersionDefinition: "const char AUTOVERSION[]",
			ScanStop:          "",
		},
		//{
		//	FileExtensions:    []string{"cc", "cpp", "cxx", "c++"},
		//	VersionDefinition: "const AutoVersion = ",
		//	ScanStop:          "func",
		//},
		//{
		//	FileExtensions:    []string{"php"},
		//	VersionDefinition: "const AutoVersion = ",
		//	ScanStop:          "func",
		//},
	}
}

func main() {
	langDefs := buildLanguageDefinitions()
	filesToScan := handleParameter(langDefs)

	file, err := os.OpenFile(prepareLogDir(true), os.O_RDWR|os.O_CREATE, 0664)
	if err != nil {
		os.Exit(1)
	}
	defer file.Close()
	logger = slog.New(slog.NewJSONHandler(file, nil))
	//logger = slog.New(slog.NewTextHandler(file, nil))

	sourceFiles := createSourceFileList(filesToScan, langDefs)
	for _, sourceFile := range sourceFiles {
		err = increaseVersionNumberInFile(sourceFile, langDefs[sourceFile.LangDef])
	}
}

// isFileSupported checks if given filePath is a source code file we support.
// Returns true/false and the index in the LanguageDefinition of the language
func isFileSupported(filePath string, langDefs []LanguageDefinition) (bool, int) {
	for i, l := range langDefs {
		for _, e := range l.FileExtensions {
			if strings.HasSuffix(filePath, e) {
				return true, i
			}
		}
	}
	return false, 0
}

// increaseVersionNumberInFile scanning a given file for the language dependent version string
// and increasing it by one if found. The scanning stops at a language dependent stop word.
func increaseVersionNumberInFile(sourceFile SourceFile, langDef LanguageDefinition) error {
	file, err := os.ReadFile(sourceFile.FilePath)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(file), "\n")
	// todo make number of space after const variable
	for i, line := range lines {
		line := strings.TrimSpace(line)
		switch {
		case langDef.ScanStop != "" && strings.HasPrefix(line, langDef.ScanStop):
			return nil
		case strings.HasPrefix(line, langDef.VersionDefinition):
			pre := lines[i]
			lines[i] = increaseVersionNumber(line)
			logger.Info("changed", slog.String("from", pre), slog.String("to", lines[i]), slog.String("source file", sourceFile.FilePath))
			//fmt.Println(lines[i])
			//fmt.Println(sourceFile.FilePath)
			err = os.WriteFile(sourceFile.FilePath, []byte(strings.Join(lines, "\n")), 0644) // todo use same permissions as file read
			if err != nil {
				logger.Error("writing file %s\n", err)
			}
			return nil
		}
	}
	return nil
}

// increaseVersionNumber increases the version number in the parameter string by one.
func increaseVersionNumber(line string) string {
	r := regexp.MustCompile(`(\d+)\D*";?$`)
	match := r.FindSubmatch([]byte(line))
	if len(match) == 0 {
		return line
	}
	version, err := strconv.Atoi(string(match[1]))
	if err != nil {
		logger.Error("cannot convert to int", err)
		return ""
	}
	version++
	idxs := r.FindSubmatchIndex([]byte(line))
	ns := line[:idxs[2]] + strconv.Itoa(version) + line[idxs[2]+len(match[1]):]
	return ns
}

// handleParameter handles all program options like --VersionAutoCounter, --debug
// returns title and description no options were used
func handleParameter(langDefs []LanguageDefinition) []string {
	if len(os.Args) == 1 {
		return []string{}
	}
	params := os.Args[1]
	switch params {
	case "--help":
		usage(langDefs)
		os.Exit(0)
	case "--version":
		fmt.Println(AutoVersion)
		os.Exit(0)
	case "--lang":
		printSearchText(langDefs)
		//printSupportedLanguages(langDefs)
		os.Exit(0)
	//case "--search-text":
	//	printSearchText(langDefs)
	//	os.Exit(0)
	default:
		return os.Args[1:]
	}
	return nil
}

func usage(langDefs []LanguageDefinition) {
	fmt.Println("autoversion, a program to increase a version number in a source file.")
	fmt.Println("Version:", AutoVersion)
	fmt.Println("")
	fmt.Println("Without parameter autoversion recursively scans source files of supported languages")
	fmt.Println("and increases a version number by one if found. See the README for examples.")
	fmt.Println("autoversion is typically called from a git pre-commit hook or a makefile.")
	fmt.Println("")
	fmt.Println("autversion uses no environment variables.")
	fmt.Println("autversion logs to", prepareLogDir(false))
	fmt.Println("Logs are overwritten after each run!")
	fmt.Println("")
	fmt.Println("Supported languages:")
	printSupportedLanguages(langDefs)
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("autoversion --help     Prints this text")
	fmt.Println("autoversion --version  Prints version information")
	fmt.Println("autoversion --lang     Prints detailed information re. supported languages")
	fmt.Println("autoversion            Recursively scan all files if language is supported")
	fmt.Println("autoversion [FILE ...] Scans all given files if language is supported")
	fmt.Println("")
}

func printSupportedLanguages(langDefs []LanguageDefinition) {
	var supportedLanguages []string
	for _, langDef := range langDefs {
		supportedLanguages = append(supportedLanguages, langDef.Language)
	}
	slices.Sort(supportedLanguages)
	for i, lang := range supportedLanguages {
		if i == len(langDefs)-1 {
			fmt.Printf(" and %s\n", lang)
		} else {
			if i == 0 {
				fmt.Printf("%s", lang)
			} else {
				fmt.Printf(", %s", lang)
			}
		}
	}
}

func printSearchText(langDefs []LanguageDefinition) {
	var longestLanguageName int
	var supportedLanguages []string
	for _, langDef := range langDefs {
		if len(langDef.Language) > longestLanguageName {
			longestLanguageName = len(langDef.Language)
		}
	}
	version := " = \"1.2.3\""
	for _, langDef := range langDefs {
		if len(langDef.Language) > longestLanguageName {
			longestLanguageName = len(langDef.Language)
		}
		supportedLanguages = append(supportedLanguages, fmt.Sprintf("Language: %s %s %s%s", langDef.Language, strings.Repeat(" ", longestLanguageName-len(langDef.Language)), langDef.VersionDefinition, version))
	}
	slices.Sort(supportedLanguages)
	for _, lang := range supportedLanguages {
		fmt.Printf("%s\n", lang)
	}
}

// createSourceFileList returns a list of all files to be scanned.
// Either by checking all parameter given to the program for supported file extension
// or by recursively scanning the filesystem beginning at the cwd for supported files.
// The latter is done only if the parameter list is empty
func createSourceFileList(filesPathToScan []string, langDefs []LanguageDefinition) []SourceFile {
	var sourceFiles []SourceFile
	if len(filesPathToScan) != 0 {
		for _, filePath := range filesPathToScan {
			relevant, idx := isFileSupported(filePath, langDefs)
			if relevant {
				sourceFiles = append(sourceFiles, SourceFile{
					FilePath: filePath,
					LangDef:  idx,
				})
			}
		}
		return sourceFiles
	}
	cwd, err := os.Getwd()
	if err != nil {
		logger.Error("cannot get current working dir", err)
		os.Exit(1)
	}
	err = filepath.WalkDir(cwd, func(path string, dirEntry fs.DirEntry, err error) error {
		if !dirEntry.Type().IsDir() {
			relevant, idx := isFileSupported(dirEntry.Name(), langDefs)
			if relevant {
				sourceFiles = append(sourceFiles, SourceFile{
					FilePath: path,
					LangDef:  idx,
				})
			}
		}
		return nil
	})
	if err != nil {
		logger.Error("cannot traverse dirs", err)
		os.Exit(1)
	}
	return sourceFiles
}

// prepareLogDir returns the logFilePath
// If the parameter is true it also creates ~/.config/autoversion/logs if necessary
func prepareLogDir(createDirs bool) string {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		os.Exit(1)
	}
	logFileDir := filepath.Join(userHomeDir, ".config", "autoversion", "logs")
	if createDirs {
		err = os.MkdirAll(logFileDir, 0755)
		if err != nil {
			os.Exit(1)
		}
	}
	logFileName := filepath.Join(logFileDir, "autoversion.log")
	return logFileName
}
