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
type ExecutionManager interface {
    Join ()
    Execute (Task)
}

type Task func ()

type executionManager struct {
    sem     chan byte
    done    chan byte
    waiting chan Task
    running uint32
}

//
// Construct a new ExecutionManager
func NewExecutionManager (max_queued, max_running uint16) ExecutionManager {
    return &executionManager {
        sem:        make(chan byte, max_running),
        done:       make(chan byte, max_running),
        waiting:    make(chan Task, max_queued),
        running:    0,
    }
}

//
//
func (em *executionManager) drainDone () {
    for cont := true; cont; {
        select {
        case <-em.done:
            em.running--
        default:
            cont = false
        }
    }
}

//
//
func (em *executionManager) jobSemaphoreDown () {
    defer func () {
        em.running++
    }()

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
func (em *executionManager) startTask (t Task) {
    em.jobSemaphoreDown()
    go em.wrapTaskReturn(t)
}

//
// Add a task to the waiting list to be executed.
//
// When there are less than max_running routines executing,
// t will start executing. Otherwise this method will
// block until there are available execution slots.
//
func (em *executionManager) Execute (t Task) {
   em.startTask(t)
}

//
// Wait for all running tasks to finish.
func (em *executionManager) Join () {
    for ; em.running > 0; em.running-- {
        <-em.done
    }
}

///////////////////////
// Subroutine part
///////////////////////

//
//
func (em *executionManager) jobSemaphoreUp () {
    em.done <- 1
    <-em.sem
}

//
//
func (em *executionManager) wrapTaskReturn (t Task) {
    t()
    em.jobSemaphoreUp()
}

