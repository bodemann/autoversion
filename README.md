
autoversion
-----------
autoversion increases a version number in a source file. It scans for a version number definition, e.g. `const AutoVersion = "1.9.3"` and increases the last part of the number by one, in this case to `const AutoVersion = "1.9.4"`. autoversion supports a wide variety of programming languages. 

autoversion is usually triggered from a pre commit hook or from a makefile. 

The program is developed in **pure** [Go](https://go.dev/). Github shows several other languages which happen to be 
**test files** for different languages.

Licensed under MIT

## Description

Developer often want to embed a version number into their executable. It has many advantages to
manage this number automatically. 

One way to achieve this is by embedding the output of `git rev-parse` into the executable during
the linking phase. In [Go](https://go.dev/) this can be done by the following script:

```
now=$(date +'%Y-%m-%d_%T')
go build -ldflags "-X main.sha1ver=`git rev-parse HEAD` -X main.buildTime=$now"
```
While this method makes it easy to link an executable to its source code revision it is not 
easy to understand for humans. Sometimes there is a need for an easy to read and understand number 
that automatically increases.

autoversion does this by scanning specifically given or all source files for a language specific definition and increments
it. In effect autoversion changes the source code. autoversion can easily be added to the pre-commit 
hook of git, see [here](https://verdantfox.com/blog/how-to-use-git-pre-commit-hooks-the-hard-way-and-the-easy-way) 
for a description. An alternative would be calling it from a makefile. 

In case of a Go codebase autoversion will scan all .go files searching for `const AutoVersion = ` and will
increase the last number of the following version string. For example 
- `const AutoVersion = "16"` will turn into `const AutoVersion = "17"`
- `const AutoVersion = "1.2"` will turn into `const AutoVersion = "1.3"`
- `const AutoVersion = "0.1.0`" will turn into `const AutoVersion = "0.1.1"` 
- `const AutoVersion = "v1.22.99`" will turn into `const AutoVersion = "v1.22.100"`
- `const AutoVersion = "version 0.9.32 NOT FOR RELEASE`" will turn into `const AutoVersion = "version 0.9.33 NOT FOR RELEASE"`
- `const AutoVersion = "I don't understand the concept of a number"` will **not** change at all

autoversion assumes e.g. `const AutoVersion` to be defined near the beginning of the file. For 
efficiency reasons it will stop scanning the source file as early as possible. In case of Go this
happens with the first `func` found. 

For each programming language supported autoversion stores a
list of file extensions, the appropriate definition, e.g. `public static final String AUTO_VERSION` 
for [Java](https://www.java.com/) and a text which will terminate the scanning, e.g. `function` for 
[Javascript](https://www.ecma-international.org/publications-and-standards/standards/ecma-262/).

A complete list of supported languages and the search string for the AUTOVERSION text can be printed
with **autoversion --lang**

## Installation

The easiest way is to download a suitable binary from the release page on github and copy that to your 
path. If you have [Task](https://taskfile.dev/) installed you can build it with `task build`. See the
Taskfile.yml. It is easy to understand.

## Alternatives

The same result could be achieved by using a [sed](https://www.gnu.org/software/sed/) (stream editor) script. It is just a bit harder to make this consistent over different development platforms.   


## Command Line Parameter

| Parameter                                                                                 | Result                                                                                         |
|-------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------|
| autoversion --help                                                                        | Prints a help text                                                                             |
| autoversion --lang                                                                   | Prints detailed information re. supported languages                                            |
| autoversion --version                                                                     | Prints version number                                                                          |
| autoversion [FILE ...]<br>e.g.: autoversion ~/GoProjects/myproject/cmd/**server**/main.go | Scan only files given as parameter                                                             |
| autoversion| Scan recursively all source files it supports and upgrade the AutoVersion definition if found. |

## Logging

Logging goes to ~/.config/autoversion/logs/autoversion.log<br> 
**Logs are overwritten after each run!**

## Performance notes

For efficiency reasons the current implementation scans until a language specific defined stop word, `func` in 
case of go. This works since it is expected to find a version definition before the first occurrence of func. 
This performance improvement over scanning the complete file is negligible though since autoversion reads the complete file 
with os.ReadFile instead of only parts of it.

If this imposes a problem I am open to implement a more efficient solution. Please open an issue on github.