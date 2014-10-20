package scs


//
// Manages concurrently running tasks.
//
// The manager is able to limit the number of concurrently
// running tasks to a maximum number.
//
// It can also limit the amount of tasks that can be
// waiting in the queue to be started.
//
type ExecutionManager struct {
    sem     chan byte
    done    chan byte
    waiting chan Task
    running uint32
}

//
// Construct a new ExecutionManager
func NewExecutionManager (max_queued, max_running uint16) {
    return &ExecutionManager {
        sem:        make(chan byte, max_running),
        done:       make(chan byte, max_running),
        waiting:    make(chan byte, max_queued),
        running:    0,
    }
}

//
//
func (em *ExecutionManager) drainDone () {
    done := uint32(0)
    for cont := true; cont; {
        select {
        case <-em.done:
            done++
        default:
            cont = false
        }
    }
    em.running -= done
}

//
//
func (em *ExecutionManager) jobSemaphoreDown () {
    defer func () {
        em.running++
    }

    select {
    case em.sem <- 1:
        return
    default:
        em.drainDone()
    }
    em.sem <- 1
}

//
//
func (em *ExecutionManager) startTask (t Task) {
    em.jobSemaphoreDown()
    go em.wrapTaskReturn(t)
}

//
// Add a task to the waiting list to be executed.
//
// When there are less than max_running routines executing,
// t will start executing.
//
func (em *ExecutionManager) Execute (t Task) {
   em.startTask(t)
}

///////////////////////
// Subroutine part
///////////////////////

//
//
func (em *ExecutionManager) jobSemaphoreUp () {
    em.done <- 1
    <-em.sem
}

//
//
func (em* ExecutionManager) wrapTaskReturn (t Task) {
    t()
    em.jobSemaphoreUp()
}

