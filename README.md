<!-- ABOUT THE PROJECT -->

## About The Project

The GO-IAM is a user-friendly system that allows individuals to create accounts for their organizations. The account owner, known as the super user, gains administrative privileges and can manage users, groups, roles, and tasks within the organization, leveraging Role-Based Access Control (RBAC) for authorization. It ensures strong security measures to protect against potential threats and supports easy data upload using CSV and Excel files. The application is built with GoLang, GoFiber, and GORM, and it is containerized using Docker for easy deployment. The database options include PostgreSQL.

## Resources

- [Postman collection](https://warped-eclipse-990288.postman.co/workspace/b96e819e-4478-42e8-9999-dd0a4efe4053)

### Built With

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![GoFiber](https://img.shields.io/badge/gofiber-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)
![Postman](https://img.shields.io/badge/Postman-FF6C37?style=for-the-badge&logo=Postman&logoColor=white)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- GETTING STARTED -->

## Important Notes

- While generating user accounts, password is optional. If it is not provided, a system generated password will be provided with the response.
- To test roles and groups assigned to a user routes are provided. Refer postman for it.
- When using a new database please seed roles using the seed role route in postman.
- Presently, 12 system generated roles will be able. You can use read roles to access them.
- Accessing the resources requires user to be authorized with roles.
- Application is made of sub modules, so it'll be easier to migrate to microservices in future.

## Getting Started

### Prerequisites

- GoLang
- Docker (optional)
- Docker Compose (optional)

### Running Go Locally

1. Make sure you have Go installed on your system.
2. Clone the repository.
3. Navigate to the root directory of the project.
4. Take reference from the `.env.example` file to create the environment file `.env` and update the environment variables as needed.
5. Run the following command to build and run the Go application:

```
go run main.go
```

6. The application will be accessible at http://localhost:3000.

### Using Docker Compose

1. Make sure you have Docker and Docker Compose installed on your system.
2. Clone the repository.
3. Navigate to the root directory of the project.
4. Take reference from the `.env.example` file to create the environment file `.env` and update the environment variables as needed.
5. Run the following command to build and start the application in Docker containers:

```
docker-compose up -d
```

OR

```
./compose.sh
```

> If you are facing permission issues, run the following command:

```
chmod +x run.sh
```

6. The application will be accessible at http://localhost:3000.

### Using Docker Image

1. Make sure you have Docker installed on your system.
2. Clone the repository.
3. Navigate to the root directory of the project.
4. Take reference from the `.env.example` file to create the environment file `.env` and update the environment variables as needed.
5. Run the following command to build the Docker image:

> Build the Docker image using the current directory as the build context

```
docker build -t tanmaybalkantask:latest .
```

> Run the Docker container interactively, mapping the required port (30000 in this case)
> and using the .env file from the host machine as a volume inside the container

```
docker run -it -p 3000:3000 --env-file .env tanmaybalkantask:latest
```

OR

```
./run.sh
```

> If you are facing permission issues, run the following command:

```
chmod +x run.sh
```

7. The application will be accessible at http://localhost:3000.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- CONTACT -->

## Contact

Tanmay Vyas

[![GitHub](https://img.shields.io/badge/github-%23121011.svg?style=for-the-badge&logo=github&logoColor=white)](https://github.com/Tanmay000009)
[![LinkedIn](https://img.shields.io/badge/linkedin-%230077B5.svg?style=for-the-badge&logo=linkedin&logoColor=white)](https://www.linkedin.com/in/tanmay-vyas-09/)
[![Gmail](https://img.shields.io/badge/Gmail-D14836?style=for-the-badge&logo=gmail&logoColor=white)](mailto:tanmayvyas09@gmail.com)
[![Resume](https://img.shields.io/badge/Resume-000000?style=for-the-badge&logo=read-the-docs&logoColor=white)](https://drive.google.com/file/d/1lkfmeqseeSwK1GlJHEblz2ZuYzdNBRhm/view?usp=drive_link)
