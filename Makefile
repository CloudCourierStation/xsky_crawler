help:
	@grep -hE '^[ a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-17s\033[0m %s\n", $$1, $$2}'

# ==============================================================================
# Targets

.PHONY: cr-run
cr-run: ;$(info $(M)...Begin to run campus recruitment crawler.)  @ ## Run campus recruitment crawler
	go run cmd/cmd.go cr --url https://xskydata.jobs.feishu.cn/school --current 1 --limit 100


.PHONY: unittest
unittest: ;$(info $(M)...Begin to benchtest xsky crawler.) @ ## run xsky crawler unit test
	go test   -v ./...  --cover


.PHONY: benchtest
benchtest: ;$(info $(M)...Begin to benchtest xsky crawler.) @ ## run xsky crawler bench test
	go test -bench=. -v ./... -run=^$ -test.benchmem

.PHONY: clear
clear: ;$(info $(M)...clear running output file.) @ ## run xsky crawler clean task
	rm -rf *.json

