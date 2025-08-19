# Fiber DI Server Template

This is a template for a Fiber server with dependency injection using Go modules and the fx framework.


## Motives

This project using the [Fx](https://uber-go.github.io/fx/) framework to handle the dependency injection and the [Fiber](https://github.com/gofiber/fiber) framework to handle the HTTP requests.

It's simply a template for a server that uses the fx framework to handle the dependency injection and the fiber framework to handle the HTTP requests.

For internal modules you can use the fx framework to handle the dependency injection.
Check an example in `src/internal/fx.go` which is implements all modules that are related to the internal part of the server.

## Tech Stack

- [Fiber](https://github.com/gofiber/fiber)
- [Fx](https://uber-go.github.io/fx/)
- [Goose](https://github.com/pressly/goose)

## Adding a new module

If you want to add a new module to your project you can follow the instructions in `src/bootstrap/fx.go`

Or in can if you want to connect module on your own like `gRPC`

Simple create a folder in `src/bootstrap` and add the following file `fx.go` to it

And done! Now you can use the module in your project, or create a merge request to add the module to the template

## Adding a new Flow

To add for an example User module you can follow the next steps, because this project using 

`Controller -> Service -> Repository`

You should follow this steps

1. Create a folder for your new module in `src/internal`
2. Create a `controller.go` to accept incoming request inside of this folder
3. Add `service` for business logic
4. Add `repository` to handle data base logic
5. After that you can add the new module to `src/internal/<module_name>fx.go` and write this code inside 
```go 
var Module = fx.Module("<module_name>", fx.Options(
	fx.Provide(
		fx.Annotate(
			NewYourController, // Function which is returning your new controller
			fx.As(new(IYourNewController)), // And cast it to real one you've added 
		), 
		// You can also add more providers if you will use them in other modules
		fx.Annotate(
            NewYourService, 
			fx.As(new(IYourNewService)), 
		),
		fx.Annotate(
			NewYourRepository,
			fx.As(new(IYourNewRepository)), 
		),
	),
))
```
7. Then you should add the new module to `src/internal/fx.go` as follows

```go
var Module = fx.Module("internal", fx.Options(
	yourNewModule.Module,

	// Please leave this line at the end
	routes.Module,
))
```

## Controller

Controllers are responsible for handling incoming requests and sending responses back to the client.

### To create a new controller in your new module
Firstly create a folder in `src/internal/<module_name>` and add the following file `controller.go` to it
And create a method `SetupRoutes` in `controller.go` interface

#### controller.go
```go
type IYourController interface {
	SetupRoutes(router fiber.Router)
}

type YourController struct {
	Logger *zap.Logger
	Service IYourService // Inject your service
}

func NewYourController(logger *zap.Logger, service IYourService) IYourController {
	logger.Info("🚀 Your controller initialized")
	return &YourController{
		Logger: logger,
		Service: service,
	}
}

func (c *YourController) SetupRoutes(router fiber.Router) {
    // Group your routes
    group := router.Group("/your-resource")
    
    // Define your endpoints
    group.Get("/", c.ListResources)
    group.Post("/", middleware.RequestValidator(model.CreateResourceRequest{}), c.CreateResource)
    group.Get("/:id", middleware.ParamsValidator(model.ResourceParams{}), c.GetResource)
}

// Example handler using validated request data from ctx.Locals
func (c *YourController) CreateResource(ctx *fiber.Ctx) error {
    // Get validated request payload from ctx.Locals
    payload := ctx.Locals("payload").(model.CreateResourceRequest)
    
    // Call your service
    result, err := c.Service.CreateResource(ctx.Context(), payload)
    if err != nil {
        return utils.InternalServerError(ctx, "create_failed", "Failed to create resource", nil)
    }
    
    return ctx.Status(fiber.StatusCreated).JSON(result)
}
```

#### Adding a new controller routes to the main router

Navigate to `src/internal/routes/routes.go`
And in `registerRoutes` add the pointer to your new controller as
`yourNewController *yourModule.IYourNewController`

And below write
```go
package routes 

func registerRoutes(r *Route, 
    yourNewController yourModule.IYourNewController, // Replace with your actual controller
	log *zap.Logger) {
    api := r.Fiber.Group("/api")

    yourNewController.SetupRoutes(api)
	
	/* Other code */
}
```

After that you can simply go to `src/internal/fx.go` and add this module 

## Service

Services are responsible for the business logic of the application.

### To create a new service in your module

Create a file named `service.go` in your module directory:

```go
// service.go
package yourmodule

import (
	"context"
	"go.uber.org/zap"
)

// Define service interface
type IYourService interface {
	GetResource(ctx context.Context, id string) (*model.Resource, error)
	ListResources(ctx context.Context, filters model.ResourceFilters) ([]*model.Resource, error)
	CreateResource(ctx context.Context, data model.CreateResourceRequest) (*model.Resource, error)
	UpdateResource(ctx context.Context, id string, data model.UpdateResourceRequest) (*model.Resource, error)
	DeleteResource(ctx context.Context, id string) error
}

// Implement service
type YourService struct {
	Logger     *zap.Logger
	Repository IYourRepository
}

// Constructor for the service
func NewYourService(logger *zap.Logger, repo IYourRepository) IYourService {
	logger.Info("🚀 Your service initialized")
	return &YourService{
		Logger:     logger,
		Repository: repo,
	}
}

// Service methods implementation
func (s *YourService) GetResource(ctx context.Context, id string) (*model.Resource, error) {
	// Add business logic here
	return s.Repository.FindByID(ctx, id)
}

func (s *YourService) CreateResource(ctx context.Context, data model.CreateResourceRequest) (*model.Resource, error) {
	// Implement business logic
	// Example: Transform DTO to entity
	resource := &model.Resource{
		// Map fields from request to entity
	}
	
	// Call repository
	return s.Repository.Create(ctx, resource)
}

// Implement other methods...
```

## Repository

Repositories are responsible for handling the data access logic of the application.

### To create a new repository in your module

Create a file named `repository.go` in your module directory:

```go
// repository.go
package yourmodule

import (
	"context"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Repository interface
type IYourRepository interface {
	FindByID(ctx context.Context, id string) (*model.Resource, error)
	FindAll(ctx context.Context, filters model.ResourceFilters) ([]*model.Resource, error)
	Create(ctx context.Context, resource *model.Resource) (*model.Resource, error)
	Update(ctx context.Context, resource *model.Resource) (*model.Resource, error)
	Delete(ctx context.Context, id string) error
}

// Repository implementation
type YourRepository struct {
	DB     *sqlx.DB
	Logger *zap.Logger
}

// Constructor for the repository
func NewYourRepository(db *sqlx.DB, logger *zap.Logger) IYourRepository {
	logger.Info("🚀 Your repository initialized")
	return &YourRepository{
		DB:     db,
		Logger: logger,
	}
}

// Repository methods implementation
func (r *YourRepository) FindByID(ctx context.Context, id string) (*model.Resource, error) {
	// SQL query implementation
	query := `SELECT * FROM resources WHERE id = $1`
	
	var resource model.Resource
	err := r.DB.GetContext(ctx, &resource, query, id)
	if err != nil {
		r.Logger.Error("Failed to find resource", zap.Error(err), zap.String("id", id))
		return nil, err
	}
	
	return &resource, nil
}

func (r *YourRepository) Create(ctx context.Context, resource *model.Resource) (*model.Resource, error) {
	// SQL query implementation
	query := `INSERT INTO resources (...) VALUES (...) RETURNING *`
	
	// Implementation details
	return resource, nil
}

// Implement other methods...
```

## Standard API Response Format

This project uses a standardized response format for all API responses to ensure consistency across the application. The response structure is defined in `src/internal/utils/response.go`.

### Response Structure

Every API response follows this structure:

```json
{
  "success": true,
  "data": { ... },
  "error": null,
  "meta": { ... }
}
```

- `success`: Boolean indicating whether the request was successful
- `data`: Contains the response payload (only present in successful responses)
- `error`: Contains error details (only present in error responses)
- `meta`: Contains metadata such as pagination information (optional)

### Error Response Structure

Error responses include detailed information about what went wrong:

```json
{
  "success": false,
  "data": null,
  "error": {
    "code": "validation_error",
    "message": "Invalid input parameters",
    "details": { ... }
  },
  "meta": null
}
```

### Using the Response Utilities

The `utils` package provides helper functions to generate standardized responses:

```go
// Success response with data
return utils.Success(c, user)

// Success response with pagination metadata
return utils.SuccessWithMeta(c, users, page, pageSize, totalPages, totalCount)

// Error responses
return utils.BadRequest(c, "validation_error", "Invalid input", validationErrors)
return utils.NotFound(c, "user_not_found", "User not found", nil)
return utils.InternalServerError(c, err)
return utils.Unauthorized(c, "Invalid token")
return utils.Forbidden(c, "Insufficient permissions")

// Special status responses
return utils.Created(c, newUser)
return utils.NoContent(c)
return utils.SuccessOnly(c)
```

### Builder Pattern

For more complex responses, you can use the builder pattern:

```go
response := utils.NewResponse().
    WithData(data).
    WithMeta(page, pageSize, totalPages, totalCount)
return response.Send(c, fiber.StatusOK)
```

### Benefits of Standardized Responses

Using standardized responses across your API provides several benefits:

1. Consistency for API consumers
2. Simplified error handling on the client side
3. Clear separation between successful and error responses
4. Ability to include additional metadata without breaking changes
5. Easier documentation and testing

Always use these utility functions for returning responses from your controllers to maintain a consistent API experience.

## Working with Validation and ctx.Locals

This project uses middleware validation to automatically validate incoming requests and store the validated data in the request context using `ctx.Locals()`.

### Request Validation Flow

1. Validation middleware is applied to routes that need request validation
2. Middleware validates request body, query parameters, or URL parameters
3. Validated data is stored in the context using `ctx.Locals()`
4. Handler functions can access the validated data using the appropriate key

### Available Validators

The project provides three validators in `src/internal/middleware/validator.go`:

1. **RequestValidator**: Validates and parses request bodies
   ```go
   // Usage
   router.Post("/users", middleware.RequestValidator(model.CreateUserRequest{}), controller.CreateUser)
   
   // Access in handler
   payload := ctx.Locals("payload").(model.CreateUserRequest)
   ```

2. **QueryValidator**: Validates and parses query parameters
   ```go
   // Usage
   router.Get("/users", middleware.QueryValidator(model.ListUsersQuery{}), controller.ListUsers)
   
   // Access in handler
   query := ctx.Locals("query").(model.ListUsersQuery)
   ```

3. **ParamsValidator**: Validates and parses URL parameters
   ```go
   // Usage
   router.Get("/users/:id", middleware.ParamsValidator(model.UserParams{}), controller.GetUser)
   
   // Access in handler
   params := ctx.Locals("params").(model.UserParams)
   ```

### Example Implementation

```go
// Define validation struct in model package
type CreateUserRequest struct {
    Name     string `json:"name" validate:"required,min=2,max=50"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

// In controller, use validation middleware
func (c *UserController) SetupRoutes(router fiber.Router) {
    users := router.Group("/users")
    users.Post("/", middleware.RequestValidator(model.CreateUserRequest{}), c.CreateUser)
}

// In handler, access validated data
func (c *UserController) CreateUser(ctx *fiber.Ctx) error {
    // Get validated data from ctx.Locals
    req := ctx.Locals("payload").(model.CreateUserRequest)
    
    // Use validated data
    user, err := c.Service.CreateUser(ctx.Context(), req)
    if err != nil {
        return utils.InternalServerError(ctx, "create_failed", "Failed to create user", nil)
    }
    
    return ctx.Status(fiber.StatusCreated).JSON(user)
}
```

## Probe Server for Kubernetes Health Checks

This project includes a dedicated probe server specifically designed for Kubernetes health checks and monitoring. The probe server runs as a separate HTTP server on its own port, providing endpoints for liveness, readiness, and startup checks.

### Probe Server Endpoints

The probe server exposes the following endpoints:

- `/healthz` (Liveness Probe): Returns 200 OK if the application is running
- `/readyz` (Readiness Probe): Returns 200 OK if the application is ready to serve traffic
- `/startupz` (Startup Probe): Returns 200 OK if the application has completed its initial startup

### Configuring the Probe Server

You can configure the probe server using environment variables:

```dotenv
# Enable/disable the probe server
ENABLE_PROBE_SERVER=true

# Set the port for the probe server
PROBE_PORT=8081
```

### Using with Kubernetes

Configure your Kubernetes deployment to use the probe endpoints:

```yaml
livenessProbe:
  httpGet:
    path: /healthz
    port: 8081
  initialDelaySeconds: 5
  periodSeconds: 10

readinessProbe:
  httpGet:
    path: /readyz
    port: 8081
  initialDelaySeconds: 5
  periodSeconds: 10

startupProbe:
  httpGet:
    path: /startupz
    port: 8081
  failureThreshold: 30
  periodSeconds: 10
```

### Probe States and Application Lifecycle

The probe server automatically manages the health check states based on your application lifecycle:

1. When your application starts, the startup probe returns 200 OK when the Fiber server is initialized.
2. The readiness probe returns 200 OK after initialization is complete and the application is ready to serve traffic.
3. During shutdown, the readiness probe will return non-200 status to allow Kubernetes to stop routing traffic to the pod.
4. The liveness probe will continue to return 200 OK until the application fully shuts down.

### Implementing Custom Health Checks

You can customize the probe server behavior by accessing the ProbeServer instance in your own components:

```go
type YourService struct {
    Logger *zap.Logger
    ProbeServer *probe.ProbeServer
}

func NewYourService(logger *zap.Logger, probeServer *probe.ProbeServer) IYourService {
    return &YourService{
        Logger: logger,
        ProbeServer: probeServer,
    }
}

// Mark service as not ready when a critical dependency fails
func (s *YourService) HandleDependencyFailure() {
    s.ProbeServer.MarkNotReady()
    s.Logger.Error("Service marked as not ready due to dependency failure")
}

// Mark service as ready when recovered
func (s *YourService) HandleDependencyRecovered() {
    s.ProbeServer.MarkReady()
    s.Logger.Info("Service marked as ready after dependency recovery")
}
```

## Swagger Documentation

This project uses [Swagger](https://swagger.io/) for API documentation via the [gofiber/swagger](https://github.com/gofiber/swagger) middleware.

### Setting Up Swagger

1. Install the required dependencies:

```bash
# Install swag CLI
go install github.com/swaggo/swag/cmd/swag@latest

# Install swagger middleware for Fiber
go get -u github.com/gofiber/swagger
```

2. Configure your API main.go file with general API information:

```go
// @title Fiber API
// @version 1.0
// @description This is a sample server for a Fiber API.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email your-email@example.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api
func main() {
    // ...
}
```

### Writing Swagger Documentation

Add Swagger annotations to your controller handlers. Here's an example:

```go
// CreateUser godoc
// @Summary Create a new user
// @Description Creates a new user with the given input data
// @Tags users
// @Accept json
// @Produce json
// @Param request body model.CreateUserRequest true "User creation request"
// @Success 201 {object} model.User
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /users [post]
func (c *UserController) CreateUser(ctx *fiber.Ctx) error {
    // Implementation...
}
```

### Generating Swagger Documentation

Generate or update the Swagger specification files by running:

```bash
swag init -g main.go
```

This will generate the required files in the `docs` directory.

### Accessing Swagger UI

Once your server is running, the Swagger UI is available at `/swagger` endpoint. It provides an interactive documentation where you can:

- Browse all available endpoints
- Test API endpoints directly from the UI
- View request/response schemas
- Try different parameters

Swagger documentation is enabled by default in development mode. You can control this via the `SWAGGER_ENABLED` environment variable.

## Running with Air (Hot Reload)

[Air](https://github.com/cosmtrek/air) provides live reloading for Go applications, which is very useful during development.

### Installing Air

```bash
# Using go install
go install github.com/cosmtrek/air@latest

# Or using curl (Linux/macOS)
curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
```

### Configuring Air

Create a file named `.air.toml` in your project root with the following configuration:

```toml
# .air.toml
root = "."
tmp_dir = "tmp"

[build]
  cmd = "go build -o ./tmp/main ."
  bin = "./tmp/main"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "postgres_data", "clickhouse_data"]
  exclude_file = []
  exclude_regex = ["_test.go", "_templ.go"]
  exclude_unchanged = true
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  kill_delay = "0s"
  log = "build-errors.log"
  send_interrupt = false
  stop_on_error = true

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  time = false

[misc]
  clean_on_exit = false
```

### Running the Application with Air

Start your application with automatic reloading using:

```bash
# If air is in your PATH
air

# Or if you installed it locally
./bin/air
```

With Air running, any changes to your Go files will trigger an automatic rebuild and restart of your application.

## Working with diffrent Modules

If you want to use some DataBase migration tools in your project you need to install goose
```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

You can follow instructions in `migrations/README.md`

After that you can simply add the following module to your project and use it in your code

### Postgres

1. Go to `src/bootstrap/fx.go` and add the following line
```go
package bootstrap

var Module = fx.Options(
	// Core modules like Logger and Config
    postgres.Module, 
	// Other modules which will be using Postgres...
)
```
2. For localhost development please add the following lines to your `.env` file. The config represents the database connection configuration
```dotenv
# Database configuration
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_NAME=reelsmarket
```

### Click House

Go to `src/bootstrap/fx.go` and add the following line
```go
package bootstrap

var Module = fx.Options(
    // Core modules like Logger and Config
    clickhouse.Module, 
    // Other modules which will be using ClickHouse...
)
```

### Redis

Go to `src/bootstrap/fx.go` and add the following line
```go
package bootstrap

var Module = fx.Options(
    // Core modules like Logger and Config
    redis.Module, 
    // Other modules which will be using Redis...
)
```

### Kafka

Go to `src/bootstrap/fx.go` and add the following line
```go
package bootstrap

var Module = fx.Options(
    // Core modules like Logger and Config
    kafka.Module, 
    // Other modules which will be using Kafka...
)
```

For localhost development please add the following lines to your `.env` file:
```dotenv
# Kafka configuration
KAFKA_BROKERS=localhost:9092
KAFKA_USERNAME=
KAFKA_PASSWORD=
KAFKA_TIMEOUT_SECONDS=10
KAFKA_ASYNC=true
KAFKA_BATCH_SIZE=100
KAFKA_BATCH_TIMEOUT_MS=1000
KAFKA_ALLOW_AUTO_TOPIC_CREATION=true
```

#### Working with Kafka

The Kafka module provides methods for producing and consuming messages:

```go
// Producing messages
func (s *YourService) SendMessage(ctx context.Context, key string, value any) error {
    // Serialize your value to JSON
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }
    
    return s.Kafka.WriteMessage(ctx, "your-topic", []byte(key), data)
}

// Consuming messages
func (s *YourService) StartConsumer(ctx context.Context) {
    s.Kafka.ConsumeMessages(ctx, "your-topic", "your-consumer-group", func(ctx context.Context, msg kafka.Message) error {
        // Process the message
        var data YourDataStruct
        if err := json.Unmarshal(msg.Value, &data); err != nil {
            s.Logger.Error("Failed to unmarshal message", zap.Error(err))
            return err
        }
        
        // Handle the message
        return s.ProcessData(ctx, data)
    })
}
```

You can find more information for configs in `src/bootstrap/config/config.go`

## Usage

### Standard Run

```bash
go run main.go
```

### Run with Hot Reload (Development)

```bash
air
```

### Docker

```bash
# Build the Docker image
docker build -t fiber-app .

# Run the container
docker run -p 8080:8080 fiber-app
```

---
written by @Xusk947 (e.g. Aziz)