module github.com/cidverse/reposync

//go:platform linux/amd64
//go:platform darwin/amd64
//go:platform windows/amd64

go 1.22.0

toolchain go1.23.4

require (
	github.com/cidverse/cidverseutils/core v0.0.0-20241217232401-ec977fb5b6d7
	github.com/cidverse/cidverseutils/zerologconfig v0.1.1
	github.com/cidverse/go-rules v0.0.0-20231112122021-075e5e6f8abc
	github.com/cidverse/go-vcsapp v0.0.0-20241210224119-4aceacff5672
	github.com/go-git/go-git/v5 v5.12.0
	github.com/gosimple/slug v1.14.0
	github.com/rs/zerolog v1.33.0
	github.com/spf13/cobra v1.8.1
	gopkg.in/yaml.v3 v3.0.1
)

require (
	cel.dev/expr v0.19.1 // indirect
	dario.cat/mergo v1.0.1 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/ProtonMail/go-crypto v1.1.3 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.0 // indirect
	github.com/bradleyfalzon/ghinstallation/v2 v2.12.0 // indirect
	github.com/cidverse/go-ptr v0.0.0-20240331160646-489e694bebbf // indirect
	github.com/cloudflare/circl v1.5.0 // indirect
	github.com/cyphar/filepath-securejoin v0.3.6 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/go-git/gcfg v1.5.1-0.20230307220236-3a3c6141e376 // indirect
	github.com/go-git/go-billy/v5 v5.6.0 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.1 // indirect
	github.com/golang/groupcache v0.0.0-20241129210726-2c02b8208cf8 // indirect
	github.com/google/cel-go v0.22.1 // indirect
	github.com/google/go-github/v66 v66.0.0 // indirect
	github.com/google/go-github/v67 v67.0.0 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/gosimple/unidecode v1.0.1 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.7 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/kevinburke/ssh_config v1.2.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pjbgf/sha1cd v0.3.0 // indirect
	github.com/sergi/go-diff v1.3.2-0.20230802210424-5b0b94c5c0d3 // indirect
	github.com/skeema/knownhosts v1.3.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stoewer/go-strcase v1.3.0 // indirect
	github.com/xanzy/go-gitlab v0.115.0 // indirect
	github.com/xanzy/ssh-agent v0.3.3 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/exp v0.0.0-20241217172543-b2144cdd0a67 // indirect
	golang.org/x/net v0.33.0 // indirect
	golang.org/x/oauth2 v0.24.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/time v0.8.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20241216192217-9240e9c98484 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241216192217-9240e9c98484 // indirect
	google.golang.org/protobuf v1.36.0 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
)

exclude github.com/antlr4-go/antlr/v4 v4.13.1
