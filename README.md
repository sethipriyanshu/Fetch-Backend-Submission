# Fetch Backend Internship Challenge
### Description:- A REST API that helps keep track of points and point transactions.
### Language Used - Go 1.22.4 windows/amd64
### Framework - Echo V4
### Database - PostgreSQL
### Author - Priyanshu Sethi (psethi1818@gmail.com)


How to Run this app using docker?

Step 1 - Have Docker Installed
Step 2- Build the Docker Image: Open a terminal and navigate to the root directory of your project. Run the following command to build the Docker image:    docker build -t my-go-app .
Step 3 - Run the Docker Container: Start a container from the Docker image you just built. This command will run the container in detached mode (-d) and map port 8000 of the container to port 8000 on your host machine (-p 8000:8000) :  docker run -d -p 8000:8000 --name my-go-app-container my-go-app
Step 4 - Verify by using docker ps
Step 5 - program is up and running on http://localhost:8000
