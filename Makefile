container : stop
	docker build -t "projects:github-watchlists" .

run : container
	docker run -p 9000:9000 -d --name github-watchlists projects:github-watchlists

stop:
	-docker stop github-watchlists
	-docker rm github-watchlists

