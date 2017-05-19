all: release

release:
	@goxc
	@rmdir debian/

deploy: release
	@scp release/0.1/sevil*.deb ilber:
	@ssh ilber 'sudo dpkg -i sevil*.deb'
	@ssh ilber 'sudo service sevil restart'

.PHONY: all build vet test release deploy issues
