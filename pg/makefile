.PHONY: setup
setup:
	@docker-compose -f setup/docker-compose.yml up -d

clear:
	@docker-compose -f setup/docker-compose.yml down
	@rm -rf setup/pg_data
	
dev: 
	@go run *.go