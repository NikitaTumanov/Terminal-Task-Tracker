package taskstorage

type TaskStatus int

const (
	StatusNotDone TaskStatus = iota
	StatusInProgress
	StatusDone
)

type TaskStorage struct{
	Index int
	Name string
	Status TaskStatus
}

func (t *TaskStorage) Add(task string){
	//TODO
}

func (t *TaskStorage) Update(task string){
	//TODO
}

func (t *TaskStorage) Delete(task string){
	//TODO
}

func (t *TaskStorage) UpdateStatus(task string){
	//TODO
}

func (t TaskStorage) AllTasks(task string){
	//TODO
}

func (t TaskStorage) DoneTasks(task string){
	//TODO
}

func (t TaskStorage) NotDoneTasks(task string){
	//TODO
}

func (t TaskStorage) InProgressTasks(task string){
	//TODO
}