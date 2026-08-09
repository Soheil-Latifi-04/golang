package main

func main() {
	todos := Todos{}
	storage := NewJSONStorage[Todos]("todos.json")
	// store, err := NewPostgresStorage[Todos](cfg.DatabaseDSN, "todos")
	// if err != nil {
	// 	panic(err)
	// }
	storage.Load(&todos)
	cmdFlags := NewCmdFlags()
	cmdFlags.Execute(&todos)
	storage.Save(todos)
}
