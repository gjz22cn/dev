package db

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var G_dbConn *sql.DB

func OpenDb() (*sql.DB, error) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/pdl1"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Println("Error opening database:", err)
		return nil, err
	}
	G_dbConn = db
	return db, err
}

func CloseDb() {
	if G_dbConn != nil {
		G_dbConn.Close()
	}
}

func CreateTables() {
	projectTable := `CREATE TABLE IF NOT EXISTS projects (
		id INT AUTO_INCREMENT,
		name VARCHAR(255) NOT NULL,
		status varchar(64) DEFAULT NULL,
		PRIMARY KEY (id)
	)`

	milestoneTable := `CREATE TABLE IF NOT EXISTS milestones (
		id INT AUTO_INCREMENT,
		project_id INT,
		name VARCHAR(255) NOT NULL,
		responsible VARCHAR(255),
		planned_completion_date DATE,
		status varchar(64) DEFAULT NULL,
		PRIMARY KEY (id),
		FOREIGN KEY (project_id) REFERENCES projects(id)
	)`

	taskTable := `CREATE TABLE IF NOT EXISTS tasks (
		id INT AUTO_INCREMENT,
		milestone_id INT,
		name VARCHAR(255) NOT NULL,
		responsible VARCHAR(255),
		start_date DATE,
		end_date DATE,
		status varchar(64) DEFAULT NULL,
		PRIMARY KEY (id),
		FOREIGN KEY (milestone_id) REFERENCES milestones(id)
	)`

	subtaskTable := `CREATE TABLE IF NOT EXISTS subtasks (
		id INT AUTO_INCREMENT,
		task_id INT,
		name VARCHAR(255) NOT NULL,
		start_date DATE,
		end_date DATE,
		status varchar(64) DEFAULT NULL,
		PRIMARY KEY (id),
		FOREIGN KEY (task_id) REFERENCES tasks(id)
	)`

	_, err := G_dbConn.Exec(projectTable)
	if err != nil {
		panic(err)
	}

	_, err = G_dbConn.Exec(milestoneTable)
	if err != nil {
		panic(err)
	}

	_, err = G_dbConn.Exec(taskTable)
	if err != nil {
		panic(err)
	}

	_, err = G_dbConn.Exec(subtaskTable)
	if err != nil {
		panic(err)
	}
}

// Project CRUD handlers
func CreateProject(c *gin.Context) {
	var project struct {
		Name string `json:"name"`
	}
	if err := c.BindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := G_dbConn.Exec("INSERT INTO projects (name) VALUES (?)", project.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := res.LastInsertId()
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func GetProject(c *gin.Context) {
	id := c.Param("id")
	var project struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	err := G_dbConn.QueryRow("SELECT id, name FROM projects WHERE id = ?", id).Scan(&project.ID, &project.Name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	c.JSON(http.StatusOK, project)
}

func UpdateProject(c *gin.Context) {
	id := c.Param("id")
	var project struct {
		Name string `json:"name"`
	}
	if err := c.BindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := G_dbConn.Exec("UPDATE projects SET name = ? WHERE id = ?", project.Name, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project updated"})
}

func DeleteProject(c *gin.Context) {
	id := c.Param("id")

	_, err := G_dbConn.Exec("DELETE FROM projects WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project deleted"})
}

// Milestone CRUD handlers
func CreateMilestone(c *gin.Context) {
	var milestone struct {
		ProjectID             int    `json:"project_id"`
		Name                  string `json:"name"`
		Responsible           string `json:"responsible"`
		PlannedCompletionDate string `json:"planned_completion_date"`
	}
	if err := c.BindJSON(&milestone); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := G_dbConn.Exec("INSERT INTO milestones (project_id, name, responsible, planned_completion_date) VALUES (?, ?, ?, ?)", milestone.ProjectID, milestone.Name, milestone.Responsible, milestone.PlannedCompletionDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := res.LastInsertId()
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func GetMilestone(c *gin.Context) {
	id := c.Param("id")
	var milestone struct {
		ID                    int    `json:"id"`
		ProjectID             int    `json:"project_id"`
		Name                  string `json:"name"`
		Responsible           string `json:"responsible"`
		PlannedCompletionDate string `json:"planned_completion_date"`
	}

	err := G_dbConn.QueryRow("SELECT id, project_id, name, responsible, planned_completion_date FROM milestones WHERE id = ?", id).Scan(&milestone.ID, &milestone.ProjectID, &milestone.Name, &milestone.Responsible, &milestone.PlannedCompletionDate)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Milestone not found"})
		return
	}

	c.JSON(http.StatusOK, milestone)
}

func UpdateMilestone(c *gin.Context) {
	id := c.Param("id")
	var milestone struct {
		Name                  string `json:"name"`
		Responsible           string `json:"responsible"`
		PlannedCompletionDate string `json:"planned_completion_date"`
	}
	if err := c.BindJSON(&milestone); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := G_dbConn.Exec("UPDATE milestones SET name = ?, responsible = ?, planned_completion_date = ? WHERE id = ?", milestone.Name, milestone.Responsible, milestone.PlannedCompletionDate, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Milestone updated"})
}

func DeleteMilestone(c *gin.Context) {
	id := c.Param("id")

	_, err := G_dbConn.Exec("DELETE FROM milestones WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Milestone deleted"})
}

// Task CRUD handlers
func CreateTask(c *gin.Context) {
	var task struct {
		MilestoneID int    `json:"milestone_id"`
		Name        string `json:"name"`
		Responsible string `json:"responsible"`
		StartDate   string `json:"start_date"`
		EndDate     string `json:"end_date"`
	}
	if err := c.BindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := G_dbConn.Exec("INSERT INTO tasks (milestone_id, name, responsible, start_date, end_date) VALUES (?, ?, ?, ?, ?)", task.MilestoneID, task.Name, task.Responsible, task.StartDate, task.EndDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := res.LastInsertId()
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func GetTask(c *gin.Context) {
	id := c.Param("id")
	var task struct {
		ID          int    `json:"id"`
		MilestoneID int    `json:"milestone_id"`
		Name        string `json:"name"`
		Responsible string `json:"responsible"`
		StartDate   string `json:"start_date"`
		EndDate     string `json:"end_date"`
	}

	err := G_dbConn.QueryRow("SELECT id, milestone_id, name, responsible, start_date, end_date FROM tasks WHERE id = ?", id).Scan(&task.ID, &task.MilestoneID, &task.Name, &task.Responsible, &task.StartDate, &task.EndDate)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var task struct {
		Name        string `json:"name"`
		Responsible string `json:"responsible"`
		StartDate   string `json:"start_date"`
		EndDate     string `json:"end_date"`
	}
	if err := c.BindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := G_dbConn.Exec("UPDATE tasks SET name = ?, responsible = ?, start_date = ?, end_date = ? WHERE id = ?", task.Name, task.Responsible, task.StartDate, task.EndDate, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task updated"})
}

func DeleteTask(c *gin.Context) {
	id := c.Param("id")

	_, err := G_dbConn.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}

// Subtask CRUD handlers
func CreateSubtask(c *gin.Context) {
	var subtask struct {
		TaskID    int    `json:"task_id"`
		Name      string `json:"name"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}
	if err := c.BindJSON(&subtask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := G_dbConn.Exec("INSERT INTO subtasks (task_id, name, start_date, end_date) VALUES (?, ?, ?, ?)", subtask.TaskID, subtask.Name, subtask.StartDate, subtask.EndDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := res.LastInsertId()
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func GetSubtask(c *gin.Context) {
	id := c.Param("id")
	var subtask struct {
		ID        int    `json:"id"`
		TaskID    int    `json:"task_id"`
		Name      string `json:"name"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	err := G_dbConn.QueryRow("SELECT id, task_id, name, start_date, end_date FROM subtasks WHERE id = ?", id).Scan(&subtask.ID, &subtask.TaskID, &subtask.Name, &subtask.StartDate, &subtask.EndDate)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subtask not found"})
		return
	}

	c.JSON(http.StatusOK, subtask)
}

func UpdateSubtask(c *gin.Context) {
	id := c.Param("id")
	var subtask struct {
		Name      string `json:"name"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}
	if err := c.BindJSON(&subtask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := G_dbConn.Exec("UPDATE subtasks SET name = ?, start_date = ?, end_date = ? WHERE id = ?", subtask.Name, subtask.StartDate, subtask.EndDate, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subtask updated"})
}

func DeleteSubtask(c *gin.Context) {
	id := c.Param("id")

	_, err := G_dbConn.Exec("DELETE FROM subtasks WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subtask deleted"})
}

// 数据结构定义
type LabelInfo struct {
	Icon     string `json:"icon"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
}

type TaskInfo struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	BgColor     string `json:"bgColor"`
	Description string `json:"description"`
}

type ProjInfo struct {
	Id    string     `json:"id"`
	Label LabelInfo  `json:"label"`
	Tasks []TaskInfo `json:"data"`
}

type ProjectInfo struct {
	Id          uint    `json:"id"`
	ProjectName string  `json:"text"`
	StartDate   string  `json:"startDate"`
	Duration    string  `json:"duration"`
	Progress    float32 `json:"progress"`
	Open        bool    `json:"open"`
	Priority    uint    `json:"pri"`
	Owner       string  `json:"owner"`
}

func GetDoingProjectsTasks(c *gin.Context) {
	query := `
	SELECT 
		p.id AS project_id, p.name AS project_name, p.pm AS pm_name, p.rd AS rd_name, p.tester AS tester_name,
		m.id AS milestone_id, m.name AS milestone_name,
		t.id AS task_id, t.name AS task_name, t.start_date, t.end_date, t.owner AS task_owner, t.status AS task_status,
		s.id AS subtask_id, s.name AS subtask_name, s.start_date AS subtask_start_date, s.end_date AS subtask_end_date, s.status AS subtask_status
	FROM 
		projects p
	INNER JOIN 
		milestones m ON p.id = m.project_id
	LEFT JOIN 
		tasks t ON m.id = t.milestone_id
	LEFT JOIN 
		subtasks s ON t.id = s.task_id
	WHERE 
		p.status = 'doing' AND m.status = 'doing'
	ORDER BY 
		p.name, m.name;
	`

	rows, err := G_dbConn.Query(query)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var projects []ProjInfo
	projectMap := make(map[int]int)

	for rows.Next() {
		var (
			projectID        int
			projectName      string
			pm               sql.NullString
			rd               sql.NullString
			tester           sql.NullString
			milestoneID      int
			milestoneName    string
			taskID           sql.NullInt32
			taskName         sql.NullString
			taskStartDate    sql.NullString
			taskEndDate      sql.NullString
			taskOwner        sql.NullString
			taskStatus       sql.NullString
			subtaskID        sql.NullInt32
			subtaskName      sql.NullString
			subtaskStartDate sql.NullString
			subtaskEndDate   sql.NullString
			subtaskStatus    sql.NullString
		)

		if err := rows.Scan(&projectID, &projectName, &pm, &rd, &tester, &milestoneID, &milestoneName,
			&taskID, &taskName, &taskStartDate, &taskEndDate, &taskOwner, &taskStatus,
			&subtaskID, &subtaskName, &subtaskStartDate, &subtaskEndDate, &subtaskStatus); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Determine the background color based on the conditions
		var bgColor string
		if taskStatus.String == "finished" {
			bgColor = "deepskyblue"
		} else if taskStatus.String == "plan" {
			bgColor = "lightseagreen"
		} else {
			endDate, _ := time.Parse("2006-01-02", taskEndDate.String)
			today := time.Now()
			if endDate.After(today) {
				bgColor = "darkorange"
			} else {
				bgColor = "red"
			}
		}

		if index, exists := projectMap[projectID]; exists {
			project := &projects[index]
			project.Tasks = append(project.Tasks, TaskInfo{
				Id:          fmt.Sprintf("%d", taskID.Int32),
				Title:       milestoneName,
				Subtitle:    taskName.String + "(" + taskOwner.String + ")",
				StartDate:   taskStartDate.String,
				EndDate:     taskEndDate.String,
				BgColor:     bgColor,
				Description: "",
			})
		} else {
			// log.Println("creating new project", projectID, projectName)
			newProject := ProjInfo{
				Id:    fmt.Sprintf("%d", projectID),
				Label: LabelInfo{Title: projectName, Subtitle: pm.String + "," + rd.String + "," + tester.String},
				Tasks: []TaskInfo{},
			}

			projects = append(projects, newProject)
			index := len(projects) - 1
			projectMap[projectID] = index

			projects[index].Tasks = append(projects[index].Tasks, TaskInfo{
				Id:          fmt.Sprintf("%d", taskID.Int32),
				Title:       milestoneName,
				Subtitle:    taskName.String + "(" + taskOwner.String + ")",
				StartDate:   taskStartDate.String,
				EndDate:     taskEndDate.String,
				BgColor:     bgColor,
				Description: "",
			})
		}
	}

	c.JSON(http.StatusOK, projects)
}
