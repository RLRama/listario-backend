# listario-backend

This project demonstrates efficient backend RESTful API design using idiomatic Go, built to power a task management web app.

# Tech stack

- Go
- GORM
- Zerolog
- Iris

# Features

- JWT-based middleware authentication
- Clean architecture (handlers, services, repositories)
- Swagger (OpenAPI) UI for API routes manual testing
- Interface usage, dependency injection, error handling and validations, environment-based configuration

# Getting started

## Native setup

1. Clone, move into directory and create your .env

```bash
git clone https://github.com/RLRama/listario-backend.git
cd listario-backend
cp .env.example .env
nano .env # then fill with valid Postgres DB and API properties
```

2. You can either:
   * Build from source and run executable:
     
     * ```bash
       go build -o listario-backend main.go
       ```
       
     * ```bash
       ./main.go # run generated executable generated previously
       ```
   * Run binary:
      
      * ```bash
        go run main.go
        ```

3. Access the Swagger at http://localhost:8080/swagger/index.html#/ you can change the port on your .env

# Environment variables

- See [.env.example](./.env.example)

# Key takeaways

- Used Go interfaces to abstract database layer cleanly
- Learned clean error handling
- Learned basic architecture patterns and Go idioms

# Roadmap

- [x] Basic task CRUD endpoints
- [x] User system (authentication, sessions, etc.)
- [x] Live deployment (on [Koyeb](https://modest-sibley-rlrama-ba015418.koyeb.app/swagger/index.html#/))
- [ ] Docker containerization
- [ ] Advanced task features (searching, exporting, priorities, etc.)
- [ ] Email notifications

# License

Project licensed under [MIT](./LICENSE) license.

# Contact

> Ramiro Ignacio Rios Lopez
> 
> [LinkedIn](https://www.linkedin.com/in/rlrama/) - [Email](mailto:rl.ramiro11@gmail.com) - [GitHub](https://github.com/RLRama) - [Personal website](https://rlrama.onrender.com/)
