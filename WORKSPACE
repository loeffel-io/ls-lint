load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

############################################################
# http archives ############################################
############################################################

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "91585017debb61982f7054c9688857a2ad1fd823fc3f9cb05048b0025c47d023",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.42.0/rules_go-v0.42.0.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.42.0/rules_go-v0.42.0.zip",
    ],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "d3fa66a39028e97d76f9e2db8f1b0c11c099e8e01bf363a923074784e451f809",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.33.0/bazel-gazelle-v0.33.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.33.0/bazel-gazelle-v0.33.0.tar.gz",
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
    sha256 = "a949d56fed8fa0a8dd82a0a660acc949253a05b2b0c52a07e4034e27f11218f6",
    strip_prefix = "rules_js-1.33.1",
    url = "https://github.com/aspect-build/rules_js/releases/download/v1.33.1/rules_js-v1.33.1.tar.gz",
)

http_archive(
    name = "rules_oci",
    sha256 = "21a7d14f6ddfcb8ca7c5fc9ffa667c937ce4622c7d2b3e17aea1ffbc90c96bed",
    strip_prefix = "rules_oci-1.4.0",
    url = "https://github.com/bazel-contrib/rules_oci/releases/download/v1.4.0/rules_oci-v1.4.0.tar.gz",
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

go_register_toolchains(version = "1.21.3")

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
# rules oci ################################################
############################################################

load("@rules_oci//oci:dependencies.bzl", "rules_oci_dependencies")

rules_oci_dependencies()

load("@rules_oci//oci:repositories.bzl", "LATEST_CRANE_VERSION", "oci_register_toolchains")

oci_register_toolchains(
    name = "oci",
    crane_version = LATEST_CRANE_VERSION,
)

load("@rules_oci//oci:pull.bzl", "oci_pull")

oci_pull(
    name = "distroless_base",
    digest = "sha256:ccaef5ee2f1850270d453fdf700a5392534f8d1a8ca2acda391fbb6a06b81c86",
    image = "gcr.io/distroless/base",
    platforms = [
        "linux/amd64",
        "linux/arm64",
    ],
)

############################################################
# github cli ###############################################
############################################################

http_archive(
    name = "com_github_cli_cli_darwin_arm64",
    build_file_content = """exports_files(glob(["bin/*"]))""",
    sha256 = "b198e750c164a85df371819c1f89c958c0eb4839af2cc07c8b8d03077d5a8255",
    strip_prefix = "gh_2.37.0_macOS_arm64",
    urls = [
        "https://github.com/cli/cli/releases/download/v2.37.0/gh_2.37.0_macOS_arm64.zip",
    ],
)

http_archive(
    name = "com_github_cli_cli_linux_amd64",
    build_file_content = """exports_files(glob(["bin/*"]))""",
    sha256 = "a1d79712c20af6fe36e388e762bda60c3fc8fb4af2f0cff666143f2c7d3baa88",
    strip_prefix = "gh_2.37.0_linux_amd64",
    urls = [
        "https://github.com/cli/cli/releases/download/v2.37.0/gh_2.37.0_linux_amd64.tar.gz",
    ],
)

############################################################
# coreutils (sha256) #######################################
############################################################

http_archive(
    name = "com_github_uutils_coreutils_darwin_arm64",
    build_file_content = """exports_files(["coreutils"])""",
    sha256 = "1f8c3640697d90be3cbb3be6a2afac144a4991dd98d71f0a2ad758595786b15d",
    strip_prefix = "coreutils-0.0.22-x86_64-apple-darwin",
    urls = [
        "https://github.com/uutils/coreutils/releases/download/0.0.22/coreutils-0.0.22-x86_64-apple-darwin.tar.gz",  # only amd64
    ],
)

http_archive(
    name = "com_github_uutils_coreutils_linux_amd64",
    build_file_content = """exports_files(["coreutils"])""",
    sha256 = "5f990d7cd0bf062a3cd75188f84ee4a79bc26572e85aece36024cf9bd2091be6",
    strip_prefix = "coreutils-0.0.22-x86_64-unknown-linux-gnu",
    urls = [
        "https://github.com/uutils/coreutils/releases/download/0.0.22/coreutils-0.0.22-x86_64-unknown-linux-gnu.tar.gz",
    ],
)
