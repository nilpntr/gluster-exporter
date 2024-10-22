build-dev:
	mkdir -p bin
	go build -tags dev -ldflags "-s -w -X 'github.com/nilpntr/gluster-exporter/cmd.appVersion=dev'" -o bin/gluster-exporter github.com/nilpntr/gluster-exporter

run: build-dev
	./bin/gluster-exporter --log.level debug --gluster.volumes _all

version: build-dev
	@./bin/gluster-exporter version
