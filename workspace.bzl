load("@io_bazel_rules_go//go:def.bzl", "go_repository")

def archive_repositories():
    go_repository(
        name = "com_github_akamensky_argparse",
        importpath = "github.com/akamensky/argparse",
        commit = "95911c018170ab2092b96d15985030390b4535af",
    )
