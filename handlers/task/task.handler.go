package taskHandler

import (
	rolesRepo "balkantask/database/roles"
	taskRepo "balkantask/database/tasks"
	"balkantask/model"
	orgSchema "balkantask/schemas/org"
	taskSchema "balkantask/schemas/task"
	userSchema "balkantask/schemas/user"
	"balkantask/utils/roles"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

func GetAllTasks(c *fiber.Ctx) error {
	_, orgOK := c.Locals("org").(orgSchema.OrgResponse)
	user, userOK := c.Locals("user").(userSchema.UserResponse)

	if !orgOK && !userOK && !roles.UserIsAuthorized(user.Roles, user.Groups, []roles.Role{roles.TasksWriteAccess, roles.OrgFullAccess, roles.OrgWriteAccess, roles.TasksFullAccess, roles.OrgReadAccess, roles.TasksReadAccess}) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Forbidden",
		})
	}

	tasks, err := taskRepo.GetAllTasks()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "data": tasks})
}

func GetTaskById(c *fiber.Ctx) error {
	id_ := c.Params("id")
	id, err := uuid.Parse(id_)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid ID",
			"status":  "error",
		})
	}

	_, orgOK := c.Locals("org").(orgSchema.OrgResponse)
	user, userOK := c.Locals("user").(userSchema.UserResponse)

	if !orgOK && !userOK && !roles.UserIsAuthorized(user.Roles, user.Groups, []roles.Role{roles.TasksWriteAccess, roles.OrgFullAccess, roles.OrgWriteAccess, roles.TasksFullAccess, roles.OrgReadAccess, roles.TasksReadAccess}) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Forbidden",
		})
	}

	task, err := taskRepo.GetTaskById(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "data": task})
}

func CreateTask(c *fiber.Ctx) error {
	var task taskSchema.CreateTask

	if err := c.BodyParser(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	_, orgOK := c.Locals("org").(orgSchema.OrgResponse)
	user, userOK := c.Locals("user").(userSchema.UserResponse)

	if !orgOK && !userOK && !roles.UserIsAuthorized(user.Roles, user.Groups, []roles.Role{roles.RoleWriteAccess, roles.OrgFullAccess, roles.OrgWriteAccess, roles.TasksFullAccess}) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Forbidden",
		})
	}

	errors := model.ValidateStruct(task)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation Error",
			"status":  "error",
			"errors":  errors,
		})
	}

	// Check if the roles (id) exist in the database
	rolesExist, err := rolesRepo.GetRolesByIds(task.RoleIds)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Role IDs",
			"status":  "error",
		})
	}

	// Check if the roles (name) exist in the database
	rolesExist2, err := rolesRepo.GetRolesByNames(task.RoleNames)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Role IDs",
			"status":  "error",
		})
	}

	rolesExist = append(rolesExist, rolesExist2...)
	rolesExist = roles.RemoveDuplicates(rolesExist)

	taskExists, err := taskRepo.GetTaskByName(task.Name)
	if err != nil && err.Error() != "record not found" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"status":  "error",
		})
	}

	if taskExists.ID != uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Task already exists",
			"status":  "error",
		})
	}

	// Verify if all roles were found
	if len(rolesExist) != len(task.RoleIds) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Role IDs",
			"status":  "error",
		})
	}

	newTask := model.Task{
		Name:  task.Name,
		Roles: rolesExist,
	}

	createdTask, err := taskRepo.CreateTask(&newTask)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Internal Server Error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   createdTask,
	})
}

func DeleteTaskById(c *fiber.Ctx) error {

	id_ := c.Params("id")
	id, err := uuid.Parse(id_)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid ID",
			"status":  "error",
		})
	}

	_, orgOK := c.Locals("org").(orgSchema.OrgResponse)
	user, userOK := c.Locals("user").(userSchema.UserResponse)

	if !orgOK && !userOK && !roles.UserIsAuthorized(user.Roles, user.Groups, []roles.Role{roles.TasksWriteAccess, roles.OrgFullAccess, roles.OrgWriteAccess, roles.TasksFullAccess}) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Forbidden",
		})
	}

	taskExists, err := taskRepo.GetTaskById(id)

	if err != nil || taskExists.ID == uuid.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Task Not Found",
			"status":  "false",
		})
	}

	err = taskRepo.DeleteTask(&taskExists)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Internal Server Error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   true,
	})
}

func AddRoleToTask(c *fiber.Ctx) error {
	var input taskSchema.AddOrDeleteRole
	err := c.BodyParser(&input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request",
			"status":  "error",
		})
	}

	_, orgOK := c.Locals("org").(orgSchema.OrgResponse)
	user, userOK := c.Locals("user").(userSchema.UserResponse)

	if !orgOK && !userOK && !roles.UserIsAuthorized(user.Roles, user.Groups, []roles.Role{roles.OrgFullAccess, roles.TasksFullAccess, roles.OrgWriteAccess, roles.TasksWriteAccess}) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
			"status":  "error",
		})
	}

	errors := model.ValidateStruct(input)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation Error",
			"status":  "error",
			"errors":  errors,
		})
	}

	// Check if the input contains Task ID or Task Name
	var task model.Task
	if input.TaskId != uuid.Nil {
		task, err = taskRepo.GetTaskById(input.TaskId)
	} else if input.TaskName != "" {
		task, err = taskRepo.GetTaskByName(input.TaskName)
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Task ID or Task Name is required",
			"status":  "error",
		})
	}

	if err != nil || task.ID == uuid.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Task Not Found",
			"status":  "false",
		})
	}

	// Check if the input contains Role ID or Role Name
	var role model.Role
	if input.RoleId != uuid.Nil {
		role, err = rolesRepo.GetRoleById(input.RoleId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Role doesn't exist",
				"status":  "error",
			})
		}
	} else if input.RoleName != "" {
		role, err = rolesRepo.GetRoleByName(input.RoleName)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Role doesn't exist",
				"status":  "error",
			})
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Role ID or Role Name is required",
			"status":  "error",
		})
	}

	if roles.TaskHasRole(task.Roles, []model.Role{role}) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Task already has the role",
			"status":  "error",
		})
	}

	task, err = taskRepo.AddRoleToTask(task, role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"status":  "error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Role added to Task",
		"status":  "success",
		"data":    task,
	})
}

func DeleteRoleFromTask(c *fiber.Ctx) error {
	var input taskSchema.AddOrDeleteRole
	err := c.BodyParser(&input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request",
			"status":  "error",
		})
	}

	_, orgOK := c.Locals("org").(orgSchema.OrgResponse)
	user, userOK := c.Locals("user").(userSchema.UserResponse)

	if !orgOK && !userOK && !roles.UserIsAuthorized(user.Roles, user.Groups, []roles.Role{roles.OrgFullAccess, roles.TasksFullAccess, roles.OrgWriteAccess, roles.TasksWriteAccess}) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
			"status":  "error",
		})
	}

	errors := model.ValidateStruct(input)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation Error",
			"status":  "error",
			"errors":  errors,
		})
	}

	// Check if the input contains Task ID or Task Name
	var task model.Task
	if input.TaskId != uuid.Nil {
		task, err = taskRepo.GetTaskById(input.TaskId)
	} else if input.TaskName != "" {
		task, err = taskRepo.GetTaskByName(input.TaskName)
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Task ID or Task Name is required",
			"status":  "error",
		})
	}

	if err != nil || task.ID == uuid.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Task Not Found",
			"status":  "false",
		})
	}

	// Check if the input contains Role ID or Role Name
	var role model.Role
	if input.RoleId != uuid.Nil {
		role, err = rolesRepo.GetRoleById(input.RoleId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Role doesn't exist",
				"status":  "error",
			})
		}
	} else if input.RoleName != "" {
		role, err = rolesRepo.GetRoleByName(input.RoleName)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Role doesn't exist",
				"status":  "error",
			})
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Role ID or Role Name is required",
			"status":  "error",
		})
	}

	if !roles.TaskHasRole(task.Roles, []model.Role{role}) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Task does not have the role",
			"status":  "error",
		})
	}

	// Delete the role from the user
	task, err = taskRepo.DeleteRoleFromTask(task, role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"status":  "error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Task removed from user",
		"status":  "success",
		"data":    task,
	})
}

func SeedTasksFromExcel(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid file",
			"status":  "error",
		})
	}

	uploadedFile, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to read uploaded file",
			"status":  "error",
		})
	}
	// Close the file after the function returns
	defer uploadedFile.Close()

	_, orgOK := c.Locals("org").(orgSchema.OrgResponse)
	user, userOK := c.Locals("user").(userSchema.UserResponse)

	if !orgOK && !userOK && !roles.UserIsAuthorized(user.Roles, user.Groups, []roles.Role{roles.TasksWriteAccess, roles.OrgFullAccess, roles.OrgWriteAccess, roles.TasksFullAccess}) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
			"status":  "error",
		})
	}

	// Create a temporary file to save the uploaded content
	tempFile, err := os.CreateTemp("", "upload-*.xlsx")
	// CreateTemp function, it generates a unique temporary file name by replacing the asterisk (*) with a random string.
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create temporary file",
			"status":  "error",
		})
	}
	defer os.Remove(tempFile.Name())

	// Save the uploaded content into the temporary file
	_, err = io.Copy(tempFile, uploadedFile)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to save uploaded file",
			"status":  "error",
		})
	}

	// Open the temporary file using excelize
	xlsx, err := excelize.OpenFile(tempFile.Name())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to read Excel file",
			"status":  "error",
		})
	}

	// Define the columns to read from the Excel file (adjust the column numbers accordingly)
	taskNameCol := 1
	roleNamesCol := 2

	rows, err := xlsx.GetRows("Sheet1")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to read Excel file",
			"status":  "error",
		})
	}

	var createdTasks []model.Task

	for rowIndex, row := range rows {
		if rowIndex == 0 {
			continue
		}

		// Check if the row has enough columns
		if len(row) < roleNamesCol {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": fmt.Sprintf("Insufficient columns in row %d", rowIndex+1),
				"status":  "error",
			})
		}

		taskName := row[taskNameCol-1]
		roleNames := strings.Split(row[roleNamesCol-1], ",")

		// Trim spaces from role names
		for i := range roleNames {
			roleNames[i] = strings.TrimSpace(roleNames[i])
		}

		// Retrieve the roles from the database based on role names
		rolesExist, err := rolesRepo.GetRolesByNames(roleNames)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": fmt.Sprintf("Invalid role names in row %d", rowIndex+1),
				"status":  "error",
			})
		}

		newTask := model.Task{
			Name:  taskName,
			Roles: rolesExist,
		}

		createdTask, err := taskRepo.CreateTask(&newTask)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": fmt.Sprintf("Failed to create task in row %d", rowIndex+1),
			})
		}

		createdTasks = append(createdTasks, *createdTask)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Tasks created successfully",
		"status":  "success",
		"data":    createdTasks,
	})
}

func SeedTasksFromCSV(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid file",
			"status":  "error",
		})
	}

	uploadedFile, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to read uploaded file",
			"status":  "error",
		})
	}
	// Close the file after the function returns
	defer uploadedFile.Close()

	_, orgOK := c.Locals("org").(orgSchema.OrgResponse)
	user, userOK := c.Locals("user").(userSchema.UserResponse)

	if !orgOK && !userOK && !roles.UserIsAuthorized(user.Roles, user.Groups, []roles.Role{roles.TasksWriteAccess, roles.OrgFullAccess, roles.OrgWriteAccess, roles.TasksFullAccess}) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
			"status":  "error",
		})
	}

	// Create a temporary file to save the uploaded content
	tempFile, err := os.CreateTemp("", "upload-*.csv")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create temporary file",
			"status":  "error",
		})
	}
	defer os.Remove(tempFile.Name())

	// Save the uploaded content into the temporary file
	_, err = io.Copy(tempFile, uploadedFile)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to save uploaded file",
			"status":  "error",
		})
	}

	// Open the temporary file using os
	csvFile, err := os.Open(tempFile.Name())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to read CSV file",
			"status":  "error",
		})
	}
	defer csvFile.Close()

	// Create a new CSV reader
	reader := csv.NewReader(csvFile)

	// Skip the header row
	_, err = reader.Read()
	if err != nil && err != io.EOF {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to read CSV file",
			"status":  "error",
		})
	}

	var createdTasks []model.Task

	for rowIndex := 1; ; rowIndex++ {
		row, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Failed to read CSV file",
				"status":  "error",
			})
		}

		// Check if the row has enough columns
		if len(row) < 2 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": fmt.Sprintf("Insufficient columns in row %d", rowIndex),
				"status":  "error",
			})
		}

		taskNames := row[0]
		roleNames := strings.Split(row[1], " ")
		// Trim spaces from role names
		for i := range roleNames {
			roleNames[i] = strings.TrimSpace(roleNames[i])
		}

		// Retrieve the roles from the database based on role names
		rolesExist, err := rolesRepo.GetRolesByNames(roleNames)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": fmt.Sprintf("Invalid role names in row %d", rowIndex),
				"status":  "error",
			})
		}

		newTask := model.Task{
			Name:  taskNames,
			Roles: rolesExist,
		}

		createdTask, err := taskRepo.CreateTask(&newTask)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": fmt.Sprintf("Failed to create task in row %d", rowIndex),
			})
		}

		createdTasks = append(createdTasks, *createdTask)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Tasks created successfully",
		"status":  "success",
		"data":    createdTasks,
	})
}

func TestUserTask(c *fiber.Ctx) error {
	var task taskSchema.TestTask

	if err := c.BodyParser(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	user, userOK := c.Locals("user").(userSchema.UserResponse)

	var taskExists model.Task
	var err error

	if task.TaskId != uuid.Nil {
		taskExists, err = taskRepo.GetTaskById(task.TaskId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid Task ID",
				"status":  "error",
			})
		}
	} else if task.TaskName != "" {
		taskExists, err = taskRepo.GetTaskByName(task.TaskName)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid Task Name",
				"status":  "error",
			})
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Task ID or Task Name is required",
			"status":  "error",
		})
	}

	if !userOK && !roles.UserHasTaskAuthorization(user.Roles, user.Groups, taskExists) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Forbidden",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "User has role",
		"data":   true,
	})
}
