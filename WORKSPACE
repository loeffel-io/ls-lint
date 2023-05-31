load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

############################################################
# http archives ############################################
############################################################

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "6dc2da7ab4cf5d7bfc7c949776b1b7c733f05e56edc4bcd9022bb249d2e2a996",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.39.1/rules_go-v0.39.1.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.39.1/rules_go-v0.39.1.zip",
    ],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "727f3e4edd96ea20c29e8c2ca9e8d2af724d8c7778e7923a854b2c80952bc405",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.30.0/bazel-gazelle-v0.30.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.30.0/bazel-gazelle-v0.30.0.tar.gz",
    ],
)

http_archive(
    name = "rules_pkg",
    sha256 = "8f9ee2dc10c1ae514ee599a8b42ed99fa262b757058f65ad3c384289ff70c4b8",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_pkg/releases/download/0.9.1/rules_pkg-0.9.1.tar.gz",
        "https://github.com/bazelbuild/rules_pkg/releases/download/0.9.1/rules_pkg-0.9.1.tar.gz",
    ],
)

http_archive(
    name = "aspect_rules_js",
    sha256 = "0b69e0967f8eb61de60801d6c8654843076bf7ef7512894a692a47f86e84a5c2",
    strip_prefix = "rules_js-1.27.1",
    url = "https://github.com/aspect-build/rules_js/releases/download/v1.27.1/rules_js-v1.27.1.tar.gz",
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

############################################################
# custom repositories ######################################
############################################################

load("//:repositories.bzl", "go_repositories")

# gazelle:repository_macro repositories.bzl%go_repositories
go_repositories()

############################################################
# go #######################################################
############################################################

go_rules_dependencies()

go_register_toolchains(version = "1.20.4")

############################################################
# gazelle ##################################################
############################################################

gazelle_dependencies()

############################################################
# rules_pkg ################################################
############################################################

load("@rules_pkg//:deps.bzl", "rules_pkg_dependencies")

rules_pkg_dependencies()

############################################################
# rules_js #################################################
############################################################

load("@aspect_rules_js//js:repositories.bzl", "rules_js_dependencies")

rules_js_dependencies()

load("@rules_nodejs//nodejs:repositories.bzl", "DEFAULT_NODE_VERSION", "nodejs_register_toolchains")

nodejs_register_toolchains(
    name = "nodejs",
    node_version = DEFAULT_NODE_VERSION,
)

load("@aspect_rules_js//npm:repositories.bzl", "npm_translate_lock", "pnpm_repository")

npm_translate_lock(
    name = "npm",
    npmrc = "//deployments/npm:.npmrc",
    pnpm_lock = "//deployments/npm:pnpm-lock.yaml",
)

load("@npm//:repositories.bzl", "npm_repositories")

npm_repositories()

pnpm_repository(name = "pnpm")

load("@aspect_bazel_lib//lib:repositories.bzl", "register_jq_toolchains")

register_jq_toolchains()

############################################################
# github cli ###############################################
############################################################

http_archive(
    name = "com_github_cli_cli_darwin_arm64",
    build_file_content = """exports_files(glob(["bin/*"]))""",
    sha256 = "38ca9a355376abd1475362cf8b3cacf2a757198fe5fe70349cb1767666abacc6",
    strip_prefix = "gh_2.29.0_macOS_arm64",
    urls = [
        "https://github.com/cli/cli/releases/download/v2.29.0/gh_2.29.0_macOS_arm64.zip",
    ],
)

http_archive(
    name = "com_github_cli_cli_linux_amd64",
    build_file_content = """exports_files(glob(["bin/*"]))""",
    sha256 = "9fe05f43a11a7bf8eacf731422452d1997e6708d4160ef0efcb13c103320390e",
    strip_prefix = "gh_2.29.0_linux_amd64",
    urls = [
        "https://github.com/cli/cli/releases/download/v2.29.0/gh_2.29.0_linux_amd64.tar.gz",
    ],
)

############################################################
# coreutils (sha256) #######################################
############################################################

http_archive(
    name = "com_github_uutils_coreutils_darwin_arm64",
    build_file_content = """exports_files(["coreutils"])""",
    sha256 = "17f7b6ef619104498860a8b3835e72bf044f0e788d3032a7b723635875047f6c",
    strip_prefix = "coreutils-0.0.18-x86_64-apple-darwin",
    urls = [
        "https://github.com/uutils/coreutils/releases/download/0.0.18/coreutils-0.0.18-x86_64-apple-darwin.tar.gz",  # only amd64
    ],
)

http_archive(
    name = "com_github_uutils_coreutils_linux_amd64",
    build_file_content = """exports_files(["coreutils"])""",
    sha256 = "eecf076c2012bdc81a8e375efc2abc45b13e57127a2ed6028cef6aa404d7215f",
    strip_prefix = "coreutils-0.0.18-x86_64-unknown-linux-gnu",
    urls = [
        "https://github.com/uutils/coreutils/releases/download/0.0.18/coreutils-0.0.18-x86_64-unknown-linux-gnu.tar.gz",
    ],
)
