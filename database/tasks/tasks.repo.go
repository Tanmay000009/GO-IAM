package tasksRepo

import (
	"balkantask/database"
	"balkantask/model"

	"github.com/google/uuid"
)

func GetAllTasks() ([]model.Task, error) {
	db := database.DB
	var tasks []model.Task
	err := db.Preload("Roles").Find(&tasks).Error
	return tasks, err
}

func GetTaskById(id uuid.UUID) (model.Task, error) {
	db := database.DB
	var task model.Task
	err := db.Preload("Roles").First(&task, "id = ?", id).Error
	return task, err
}

func GetTaskByName(name string) (model.Task, error) {
	db := database.DB
	var task model.Task
	err := db.Preload("Roles").First(&task, "name = ?", name).Error
	return task, err
}

func CreateTask(task *model.Task) (*model.Task, error) {
	db := database.DB
	err := db.Create(&task).Error
	return task, err
}

func UpdateTask(task model.Task) (model.Task, error) {
	db := database.DB
	err := db.Save(&task).Error
	return task, err
}

func DeleteTask(task *model.Task) error {
	db := database.DB
	err := db.Model(&task).Association("Roles").Clear()
	err = db.Model(&task).Association("Groups").Clear()
	err = db.Delete(&task).Error
	return err
}

func AddRoleToTask(task model.Task, role model.Role) (model.Task, error) {
	db := database.DB
	err := db.Model(&task).Association("Roles").Append(&role)
	return task, err

}

func DeleteRoleFromTask(task model.Task, roles model.Role) (model.Task, error) {
	db := database.DB
	err := db.Model(&task).Association("Roles").Delete(roles)
	return task, err
}
