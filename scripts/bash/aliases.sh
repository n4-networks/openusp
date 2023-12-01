alias dc="docker compose -f $(pwd)/deployments/docker-compose.yaml" 
alias dclocal="docker compose -f $(pwd)/deployments/docker-compose_local.yaml" 
alias cli="docker run --env-file $(pwd)/configs/openusp.env --network=openusp -it --rm n4networks/openusp-cli"
