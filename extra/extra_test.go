package extra

import (
	"pimp-my-shell/localio"
	"runtime"
	"testing"
)

func TestInstallExtraPackages(t *testing.T) {
	dirs, err := localio.NewDirectories()
	if err != nil {
		t.Errorf("failed to create Directories type: %v", err)
	}
	osType := runtime.GOOS
	switch osType {
	case "linux":
		if err = localio.DownloadAndInstallLatestVersionOfGolang(dirs.HomeDir); err != nil {
			t.Errorf("couldn't download and install golang: %v", err)
		}
		if err = localio.RunCommandPipeOutput("go version"); err != nil {
			t.Errorf("couldn't get go version: %v", err)
		}
	default:
		//DoNothing
	}

	type args struct {
		osType   string
		dirs     *localio.Directories
		packages *localio.InstalledPackages
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Test InstallExtraPackages darwin 1", args{
			osType: "darwin",
			dirs:   dirs,
			packages: &localio.InstalledPackages{
				AptInstalledPackages: nil,
				BrewInstalledPackages: &localio.BrewInstalled{
					Names: []string{"bat", "lsd", "gnu-sed", "gotop", "yamllint", "git-delta"}, CaskFullNames: []string{"gotop"}, Taps: []string{"homebrew/core", "cjbassi/gotop"},
				},
			}}, false},
		{"Test InstallExtraPackages darwin lots of packages already installed 2", args{
			osType: "darwin",
			dirs:   dirs,
			packages: &localio.InstalledPackages{
				AptInstalledPackages: nil,
				BrewInstalledPackages: &localio.BrewInstalled{
					Names: []string{"aom", "apr", "apr-util", "argon2", "aspell", "assimp", "autoconf", "bdw-gc", "binwalk", "boost", "brotli", "c-ares", "ca-certificates",
						"cairo", "cheat", "cmake", "cointop", "coreutils", "cscope", "curl", "dbus", "deployer", "docbook", "docbook-xsl", "double-conversion", "exiftool",
						"fontconfig", "freetds", "freetype", "fribidi", "gcc", "gd", "gdbm", "gdk-pixbuf", "gettext", "ghostscript", "git-quick-stats",
						"github-markdown-toc", "glib", "gmp", "gnu-getopt", "gnu-sed", "gnupg", "gnutls", "gobject-introspection", "graphite2", "graphviz", "gts", "guile", "gulp-cli",
						"harfbuzz", "helm", "htop", "hunspell", "icu4c", "ilmbase", "imagemagick", "imath", "iproute2mac", "ipython", "isl", "jansson", "jasper", "jbig2dec", "jemalloc",
						"jpeg", "jq", "kind", "krb5", "lastpass-cli", "libarchive", "libassuan", "libb2", "libde265", "libev", "libevent", "libffi", "libgcrypt", "libgpg-error", "libheif",
						"libidn", "libidn2", "libimagequant", "libksba", "liblqr", "libmaxminddb", "libmetalink", "libmpc", "libnghttp2", "libomp", "libpng", "libpq", "libproxy", "libpthread-stubs",
						"libraqm", "librsvg", "libsmi", "libsodium", "libssh", "libssh2", "libtasn1", "libtiff", "libtool", "libunistring", "libusb", "libuv", "libx11", "libxau", "libxcb", "libxdmcp",
						"libxext", "libxrender", "libyaml", "libzip", "little-cms2", "lnav", "lsd", "lua", "lz4", "lzo", "m4", "macvim", "md4c", "mpdecimal", "mpfr", "msodbcsql17", "mssql-tools", "mysql",
						"ncurses", "neofetch", "netpbm", "nettle", "nghttp2", "nmap", "node", "npth", "nspr", "nss", "numpy", "oniguruma", "openblas", "openexr", "openjdk", "openjpeg", "openldap", "openssl@1.1",
						"p11-kit", "p7zip", "packer", "pango", "pcre", "pcre2", "php", "php@7.3", "pillow", "pinentry", "pixman", "pkg-config", "poppler", "popt", "protobuf", "putty", "pyenv", "pyenv-virtualenv",
						"python@3.10", "python@3.8", "python@3.9", "qt", "qt@5", "readline", "reattach-to-user-namespace", "rsync", "rtmpdump", "ruby", "screenresolution", "shared-mime-info", "shellcheck", "six",
						"sqlite", "ssdeep", "swaks", "tcl-tk", "telnet", "terraform", "tidy-html5", "tmux", "tree", "unbound", "unixodbc", "unrar", "unzip", "utf8proc", "watch", "webp", "wget", "x265", "xmlto",
						"xorgproto", "xxhash", "xz", "yamllint", "zeromq", "zstd"},
					CaskFullNames: []string{"font-meslo-lg-nerd-font", "wireshark"},
					Taps:          []string{"hashicorp/tap", "homebrew/core", "microsoft/mssql-release"},
				},
			}}, false},
		{"Test InstallExtraPackages Linux 3", args{
			osType: "linux",
			dirs:   dirs,
			packages: &localio.InstalledPackages{
				AptInstalledPackages:  &localio.AptInstalled{Name: []string{"bat", "lsd", "delta"}},
				BrewInstalledPackages: nil,
			}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InstallExtraPackages(tt.args.osType, tt.args.dirs, tt.args.packages); (err != nil) != tt.wantErr {
				t.Errorf("InstallExtraPackages() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
