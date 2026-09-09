// The Go-built tools the shared recipes run, pinned as tool directives
// in a module of their own so their dependency graph never reaches the
// project's go.mod (book/tooling.md).
module github.com/farcloser/godolint/tools

go 1.25.9

tool (
	github.com/farcloser/godolint/cmd/godolint
	github.com/goccy/go-graphviz/cmd/dot
	github.com/google/go-licenses/v2
	github.com/vbatts/git-validation
	golang.org/x/tools/cmd/deadcode
	golang.org/x/vuln/cmd/govulncheck
)

require (
	github.com/containerd/typeurl/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/disintegration/imaging v1.6.2 // indirect
	github.com/farcloser/godolint v0.1.0 // indirect
	github.com/fatih/color v1.18.0 // indirect
	github.com/flopp/go-findfont v0.1.0 // indirect
	github.com/fogleman/gg v1.3.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/goccy/go-graphviz v0.2.5 // indirect
	github.com/goccy/go-graphviz/cmd/dot v0.0.0-20251129032125-76e04975df88 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/golang/groupcache v0.0.0-20241129210726-2c02b8208cf8 // indirect
	github.com/google/go-licenses/v2 v2.0.1 // indirect
	github.com/google/licenseclassifier/v2 v2.0.0 // indirect
	github.com/hashicorp/go-version v1.6.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jessevdk/go-flags v1.6.1 // indirect
	github.com/magefile/mage v1.15.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/moby/buildkit v0.31.1 // indirect
	github.com/otiai10/copy v1.10.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10 // indirect
	github.com/rs/zerolog v1.34.0 // indirect
	github.com/sergi/go-diff v1.2.0 // indirect
	github.com/sirupsen/logrus v1.9.4 // indirect
	github.com/spf13/cobra v1.7.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/tetratelabs/wazero v1.10.1 // indirect
	github.com/urfave/cli/v3 v3.9.0 // indirect
	github.com/vbatts/git-validation v1.2.2 // indirect
	go.opencensus.io v0.24.0 // indirect
	golang.org/x/image v0.21.0 // indirect
	golang.org/x/mod v0.39.0 // indirect
	golang.org/x/net v0.58.0 // indirect
	golang.org/x/sync v0.22.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/telemetry v0.0.0-20260811182544-a038080d80e5 // indirect
	golang.org/x/term v0.45.0 // indirect
	golang.org/x/text v0.41.0 // indirect
	golang.org/x/tools v0.49.0 // indirect
	golang.org/x/vuln v1.7.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	k8s.io/klog/v2 v2.140.0 // indirect
	mvdan.cc/sh/v3 v3.12.0 // indirect
)
