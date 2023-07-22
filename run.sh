#!/bin/bash

# Build the Docker image using the current directory as the build context
docker build -t tanmaybalkantask:latest .

# Run the Docker container interactively, mapping the required port (30000 in this case)
# and using the .env file from the host machine as a volume inside the container
docker run -it -p 3000:3000 --env-file .env tanmaybalkantask:latest
