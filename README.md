# rules_archive

Tools and Bazel rules for creating portable, reproducible archives.

## Setup

Copy to the follow to your workspace:

```
http_archive(
    name = "io_bazel_rules_go",
    sha256 = "be6f3aba03ccef66c8ce9896f9108f28cb9a4514a57440712f32af4ab8fa1938",
    strip_prefix = "rules_go-0.11.0",
    url = "https://github.com/bazelbuild/rules_go/archive/0.11.0.zip",
)

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")

go_rules_dependencies()

go_register_toolchains()

http_archive(
    name = "rules_archive",
    sha256 = "<hash>",
    strip_prefix = "rules_go-<commit>",
    url = "https://github.com/bazelbuild/rules_go/archive/<commit>.zip",
)

load("@rules_archive:workspace.bzl", "archive_repositories")

archive_repositories()
```

## Zip

### Tools

@rules_archive/zip

```
usage: zip [-h|--help] [-a|--archive "<value>" [-a|--archive "<value>" ...]]
           [-f|--file "<value>" [-f|--file "<value>" ...]] [-o|--output
           "<value>"] [-x|--compress]

           Create a zip archive from files, directories, or other archives.
           Strips timestamps. Preserves Unix permissions.

Arguments:

  -h  --help      Print help information
  -a  --archive   Zip archive to merge. Add to the root by default; use
                  name=path to add files to name instead.
  -f  --file      File or directory of files to add. Adds using the specified
                  path; use name=path to add files to name instead.
  -o  --output    Archive output. Default: -
  -x  --compress  Deflate contents
```

### Rules

TODO
