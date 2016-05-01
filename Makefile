container :
	docker build -t "projects:github-watchlists" .

up : down container
	docker run -p 9000:9000 -d --name github-watchlists projects:github-watchlists

down:
	-docker stop github-watchlists

clean:
	-docker rm github-watchlists

