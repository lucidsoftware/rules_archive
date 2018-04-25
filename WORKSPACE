workspace(name = "rules_archive")

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "be6f3aba03ccef66c8ce9896f9108f28cb9a4514a57440712f32af4ab8fa1938",
    strip_prefix = "rules_go-0.11.0",
    url = "https://github.com/bazelbuild/rules_go/archive/0.11.0.zip",
)

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")

go_rules_dependencies()

go_register_toolchains()

load(":workspace.bzl", "archive_repositories")

archive_repositories()
