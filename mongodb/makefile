.PHONY: setup
setup:
	@mkdir -p setup/mg_data
	@docker-compose -f setup/docker-compose.yml up -d

clear:
	@docker-compose -f setup/docker-compose.yml down
	@rm -rf setup/mg_data
	
dev: 
	@go run *.go