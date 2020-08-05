dev:
	@docker-compose down && \
		docker-compose \
			-f docker-compose.yml \
			-f docker-compose.dev.yml \
			up -d --remove-orphans --build
build:
	docker exec webserver go build -o /go/bin/book github.com/yfedoruck/book/cmd/book && \
	docker-compose restart server

web:
	@docker stop webserver && \
		docker-compose \
			-f docker-compose.yml \
			-f docker-compose.dev.yml \
			build server && \
		docker start webserver

deb:
	@docker-compose down && \
			docker-compose \
				-f docker-compose.yml \
				-f docker-compose.debug.yml \
				up -d --remove-orphans --build
heroku:
	heroku container:login && \
	heroku container:push --app nameless-brook-78889 web && \
	heroku container:release --app nameless-brook-78889 web && \
	heroku logs --tail --app nameless-brook-78889