build-oauth2:
	docker build --tag oauth2 -f oauth2/Dockerfile .

build-websmtp:
	docker build --tag websmtp -f websmtp/Dockerfile .

build-dashboard:
	docker build --tag dashboard -f dashboard/Dockerfile .
