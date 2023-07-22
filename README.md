[![Review Assignment Due Date](https://classroom.github.com/assets/deadline-readme-button-24ddc0f5d75046c5622901739e7c5dd533143b0c8e959d652212380cedb1ea36.svg)](https://classroom.github.com/a/YCCXVJKc)
[![Open in Visual Studio Code](https://classroom.github.com/assets/open-in-vscode-718a45dd9cf7e7f842a935f5ebbe5719a5e09af4491e668f4dbf3b35d5cca122.svg)](https://classroom.github.com/online_ide?assignment_repo_id=11469711&assignment_repo_type=AssignmentRepo)

[![Open in Visual Studio Code](https://classroom.github.com/assets/open-in-vscode-c66648af7eb3fe8bc4f294546bfd86ef473780cde1dea487d3c4ff354943c9ae.svg)](https://classroom.github.com/online_ide?assignment_repo_id=10359099&assignment_repo_type=AssignmentRepo)

<!-- PROJECT SHIELDS -->
<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://github.com/othneildrew/Best-README-Template">
    <img src="images/logo.png" alt="Logo" width="80" height="80">
  </a>

  <h3 align="center">Tanmay Vyas's submission</h3>

  <p align="center">
    Submission to DYTE's summer intern hiring by Tanmay Vyas 20BCE0755
  </p>
</div>

<!-- ABOUT THE PROJECT -->

## About The Project

The Task Management Web Application is a user-friendly system that allows individuals to create accounts for their organizations. The account owner, known as the super user, gains administrative privileges and can manage users, groups, roles, and tasks within the organization, leveraging Role-Based Access Control (RBAC) for authorization. It ensures strong security measures to protect against potential threats and supports easy data upload using CSV and Excel files. The application is built with GoLang, GoFiber, and GORM, and it is containerized using Docker for easy deployment. The database options include PostgreSQL.

## Resources

- [Deployed Server URL](http://dyte-232883358.ap-south-1.elb.amazonaws.com)

- [API Documentation](http://dyte-232883358.ap-south-1.elb.amazonaws.com/api)

- [Postman collection](https://www.postman.com/lively-crater-584197/workspace/dyte)
  (Deployed url & tokens are set in global)

- Password for all seeded users is `password`.

- Admin username: `admin`.

- Few student usernames: `20BCE0755`, `20BCE0001`

### Built With

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)
![AWS](https://img.shields.io/badge/AWS-%23FF9900.svg?style=for-the-badge&logo=amazon-aws&logoColor=white)
![Postman](https://img.shields.io/badge/Postman-FF6C37?style=for-the-badge&logo=Postman&logoColor=white)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- GETTING STARTED -->

## Important Notes

- While generating user accounts, password is optional. If it is not provided, a system generated password will be provided with the response.
- To test roles and groups assigned to a user routes are provided. Refer postman for it.
- Presently, 12 system generated roles will be able. You can use read roles to access it.
- Accessing the resources requires user to be authorized with roles.
- When using a new database please seed roles using the seed role route in postman.

## Getting Started

### Prerequisites

- Docker
- Docker Compose
- GoLang

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

6. The application will be accessible at http://localhost:3000.

### Using Docker Image

1. Make sure you have Docker installed on your system.
2. Clone the repository.
3. Navigate to the root directory of the project.
4. Take reference from the `.env.example` file to create the environment file `.env` and update the environment variables as needed.
5. Run the following command to build the Docker image:

```
./run.sh
```

OR

> Build the Docker image using the current directory as the build context

```
docker build -t tanmaybalkantask:latest .
```

> Run the Docker container interactively, mapping the required port (30000 in this case)
> and using the .env file from the host machine as a volume inside the container

```
docker run -it -p 3000:3000 --env-file .env tanmaybalkantask:latest
```

6. If you are facing permission issues, run the following command:

```
chmod +x run.sh
```

7. The application will be accessible at http://localhost:3000.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- CONTACT -->

## Checklist

It lists out all the features successfully implemented. The remaining tasks are presented in descending order of priority.

- [x] Secure user registration and authentication
- [x] Account Deactivation and Deletion: Allow users to deactivate or delete their accounts. I
- [x] Implement a mechanism to handle account deletion securely while
      considering data retention policies.
- [x] Role-based and Group-based access management on resources(Tasks) with ability to create custom roles and groups
- [x] Protection against vulnerabilities like SQL injection attacks. Used GORM
- [x] Support for bulk upload using CSV(Both users and tasks) making sure all the
      relationships are preserved accurately. (Added support for excel also)
- [x] Docker and containerize your application code to run, including the
      database.
- [ ] Multi-tenant functionality. (Almost done). Fully suppoers multi-tenancy but implementation of postgres schemas was something I was not able to do.
- [ ] Use a reverse-proxy of your choice. (Partilly done)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Contact

Tanmay Vyas

[![GitHub](https://img.shields.io/badge/github-%23121011.svg?style=for-the-badge&logo=github&logoColor=white)](https://github.com/Tanmay000009)
[![LinkedIn](https://img.shields.io/badge/linkedin-%230077B5.svg?style=for-the-badge&logo=linkedin&logoColor=white)](https://www.linkedin.com/in/tanmay-vyas-09/)
[![Gmail](https://img.shields.io/badge/Gmail-D14836?style=for-the-badge&logo=gmail&logoColor=white)](mailto:tanmayvyas09@gmail.com)
