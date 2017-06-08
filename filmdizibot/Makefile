all: release

release:
	@goxc
	@rmdir debian/

deploy: release
	@scp release/0.1/filmdizibot*.deb ilber:
	@ssh ilber 'sudo dpkg -i filmdizibot*.deb'
	@ssh ilber 'sudo service filmdizibot restart'

.PHONY: all build vet test release deploy issues
