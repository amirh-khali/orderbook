package scheduler

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-co-op/gocron/v2"
)

var scheduler gocron.Scheduler

func Init() {
	s, err := gocron.NewScheduler()
	if err != nil {
		log.Fatalf("failed to create scheduler: %s\n", err)
	}

	scheduler = s
	log.Printf("scheduler created %v\n", scheduler)
}

func AddJob(d time.Duration, function any) {
	j, err := scheduler.NewJob(
		gocron.DurationJob(d),
		gocron.NewTask(function),
	)
	if err != nil {
		log.Fatalf("failed to create job: %s\n", err)
	}
	log.Println("a new job added to the scheduler with id: ", j.ID())
}

func waitSignal() os.Signal {
	ch := make(chan os.Signal, 2)
	signal.Notify(
		ch,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)
	for {
		sig := <-ch
		switch sig {
		case syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM:
			return sig
		}
	}
}

func Start() {
	scheduler.Start()
	log.Println("scheduler is started!")

	sig := waitSignal()
	log.Println("scheduler is stopped! cause: ", sig.String())

	_ = scheduler.Shutdown()
}
