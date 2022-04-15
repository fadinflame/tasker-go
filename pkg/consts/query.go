package consts

const CreateTaskQuery = "INSERT INTO tasks (title, text) VALUES (?, ?);"
const GetTaskQuery = "SELECT * FROM tasks WHERE id = ?;"
const UpdateTaskQuery = "UPDATE tasks SET title = ?, text = ?, is_completed = ? WHERE id = ?;"
const DeleteTaskQuery = "DELETE FROM tasks WHERE id = ?;"

const DatabasePath = "./main.sqlite"