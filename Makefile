
#
#
#

INSTALL=install
INSTBINFLAGS=-c -m 755
INSTPRVBINFLAGS=-c -m 750
INSTLIBFLAGS=-c -m 644
INSTPRVLIBFLAGS=-c -m 600
INSTDIRFLAGS=-d -m 755
INSTMANFLAGS=-c -m 644
INSTPRVDIRFLAGS=-d -m 750

RM=rm -f
PREFIX=/usr
BINDIR=$(DESTDIR)$(PREFIX)/bin
ETCDIR=$(DESTDIR)/etc
OUTDIR=./bin
BUILDTMP=./buildtmp
PKGNAME=hostsysinfo
NAME=hostsysinfo
SRCFILES=cmd/$(NAME)/$(NAME).go cmd/$(NAME)/app.go cmd/$(NAME)/utilfunc.go cmd/$(NAME)/main.go

BUILDBASE=/var/tmp/build
BUILDDIR=$(BUILDBASE)/$(PKGNAME)

all: build-dirs build build-man

build:
	go build -o $(OUTDIR)/$(NAME) $(SRCFILES)

build-dirs:
	-mkdir -p $(OUTDIR)

run:
	go run $(SRCFILES)

install: install-dirs install-bin install-config

install-dirs:
	$(INSTALL) $(INSTDIRFLAGS) $(BINDIR)
	$(INSTALL) $(INSTDIRFLAGS) $(ETCDIR)
	$(INSTALL) $(INSTDIRFLAGS) $(ETCDIR)/$(PKGNAME)

install-bin:
	$(INSTALL) $(INSTBINFLAGS) -s bin/$(NAME) $(BINDIR)

install-config:
	$(INSTALL) $(INSTLIBFLAGS) config/site-config.yaml $(ETCDIR)/$(PKGNAME)

uninstall:
	$(RM) $(BINDIR)/$(PKGNAME)

compile:
	echo "Compiling for every OS and Platform"
	GOOS=freebsd GOARCH=386 go build -o $(BINDIR)/$(NAME)-freebsd-386 $(NAME).go
	GOOS=linux GOARCH=386 go build -o $(BINDIR)/$(NAME)-linux-386 $(NAME).go
	GOOS=windows GOARCH=386 go build -o $(BINDIR)/$(NAME)-windows-386 $(NAME).go



clean:
	go clean
	$(RM) *~
	$(RM) $(OUTDIR)/*

build-man:
	$(INSTALL) $(INSTDIRFLAGS) $(BUILDTMP)
	gzip --best -n -c man/hostsysinfo.1 > $(BUILDTMP)/hostsysinfo.1.gz

clean-deb:
	rm -rf $(BUILDDIR)

build-deb:
	$(INSTALL) $(INSTDIRFLAGS) $(BUILDDIR)
	$(INSTALL) $(INSTDIRFLAGS) $(BUILDDIR)/bin
	$(INSTALL) $(INSTDIRFLAGS) $(BUILDDIR)/etc
	$(INSTALL) $(INSTDIRFLAGS) $(BUILDDIR)/etc/$(PKGNAME)
	$(INSTALL) $(INSTDIRFLAGS) $(BUILDDIR)/DEBIAN
	$(INSTALL) $(INSTDIRFLAGS) $(BUILDDIR)/usr/share/doc
	$(INSTALL) $(INSTDIRFLAGS) $(BUILDDIR)/usr/share/doc/$(PKGNAME)
	$(INSTALL) $(INSTDIRFLAGS) $(BUILDDIR)/usr/share/man/man1
	$(INSTALL) $(INSTBINFLAGS) -s bin/$(NAME) $(BUILDDIR)/bin/$(NAME)
	$(INSTALL) $(INSTLIBFLAGS) config/site-config.yaml $(BUILDDIR)/etc/hostsysinfo
	$(INSTALL) $(INSTLIBFLAGS) debian/control $(BUILDDIR)/DEBIAN/control
	$(INSTALL) $(INSTLIBFLAGS) debian/conffiles $(BUILDDIR)/DEBIAN
	$(INSTALL) $(INSTLIBFLAGS) debian/copyright $(BUILDDIR)/usr/share/doc/$(PKGNAME)
	$(INSTALL) $(INSTLIBFLAGS) debian/changelog $(BUILDDIR)/usr/share/doc/$(PKGNAME)
	gzip --best -n $(BUILDDIR)/usr/share/doc/$(PKGNAME)/changelog.Debian
	$(INSTALL) $(INSTLIBFLAGS) $(BUILDTMP)/$(PKGNAME).1.gz $(BUILDDIR)/usr/share/man/man1
	sudo chown -R root:root $(BUILDDIR)
	cd $(BUILDBASE) ; dpkg-deb --build $(PKGNAME)



test-deb:
	dpkg -c $(BUILDBASE)/$(PKGNAME).deb
	lintian $(BUILDBASE)/$(PKGNAME).deb


debuild:
	debuild --prepend-path=/usr/local/go/bin -us -uc
debuild-clean:
	debuild --prepend-path=/usr/local/go/bin clean

distclean: clean
