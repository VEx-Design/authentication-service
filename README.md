# Authentication Service

## Introduction
The Authentication Service is designed to make a sign in with Google for users.

## Features
- Secure user authentication mechanisms
- Google log-in

## Installation Dependencies
Ensure you have Go installed. After cloning the project, install the necessary packages and libraries with the following command:

```bash
go mod tidy
```

## Environment Variables
Create a `.env` file in the root directory and configure it based on `.env.example`. This file should contain all necessary environment-specific configurations.

## Running the Project
To run the project, use the following command:

```bash
make restart
```

This command will build and start the service.

## Docker
### Building the Docker Image and Running the Container
To build the Docker image and run the container, execute:

```bash
make up
```
