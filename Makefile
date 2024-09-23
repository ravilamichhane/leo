install:
	cd cli/leo && go install

test:
	leo generate controller -n task -r
	leo generate fetchers -c task
	leo generate form -c task -n CreateTask

run:
	cd cmd/web && go run .