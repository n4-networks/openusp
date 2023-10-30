alias dc="docker compose -f ./deployments/docker-compose.yaml" 
alias dclocal="docker compose -f ./deployments/docker-compose_local.yaml" 
alias cli="docker run --env-file configs/openusp.env --network=openusp -it --rm n4networks/openusp-cli"
