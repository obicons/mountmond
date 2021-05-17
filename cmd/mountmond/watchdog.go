package main

import (
	"log"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"github.com/obicons/mountmond/mtab"
)

type WatchDog struct {
	// stores the paths that need mounted (e.g. /mount/foo)
	NeededMountPaths []string

	// stores the command to execute when mount is missing
	MissingMountToCmd map[string]string

	// how often should we check /etc/mtab?
	TimeBetweenChecks time.Duration

	// the location of /etc/mtab
	MTabPath string

	// for internal use only -- used to shutdown this goroutine
	shutdownChan chan bool

	// for internal use only -- manages active commands
	activeCmds map[string]*exec.Cmd
}

/*
 * Builds and initializes a new WatchDog.
 * missingMountCmds - stores commands execute when mounts are missing.
 * Post: !running(NewWatchDog(missingMountCmds)).
 */
func NewWatchDog(missingMountCmds map[string]string) *WatchDog {
	return &WatchDog{
		MissingMountToCmd: missingMountCmds,
		MTabPath:          "/etc/mtab",
		TimeBetweenChecks: time.Second * 15,
		shutdownChan:      make(chan bool),
		activeCmds:        make(map[string]*exec.Cmd),
	}
}

/*
 * Begins the WatchDog.
 * Pre: !running(w).
 * Post: running(w).
 */
func (w *WatchDog) Start() {
	mTabFile, err := os.Open(w.MTabPath)
	if err != nil {
		log.Fatalf("WatchDog.Start(): %s", err)
	}
	go w.work(mTabFile)
}

/*
 * Shutdown the watchdog service.
 * Pre: running(w).
 * Post: !running(w).
 */
func (w *WatchDog) Shutdown() {
	w.shutdownChan <- true
	<-w.shutdownChan
}

/*
 * Periodically monitors /etc/mtab for missing mounts.
 * If a mount is missing, executes its missing command.
 */
func (w *WatchDog) work(mTabFile *os.File) {
	tCh := time.NewTicker(w.TimeBetweenChecks)
	keepGoing := true
	for keepGoing {
		mTabFile.Seek(0, os.SEEK_SET)
		select {
		case <-tCh.C:
			w.checkInOnChildren()
			w.checkForMissingMounts(mTabFile)
		case <-w.shutdownChan:
			keepGoing = false
		}
	}
	tCh.Stop()
	w.reapChildren()
	mTabFile.Close()
	w.shutdownChan <- true
}

/*
 * Reaps child processes that have died, checks the status of the others.
 */
func (w *WatchDog) checkInOnChildren() {
	for mount, proc := range w.activeCmds {
		pid, _ := syscall.Wait4(proc.Process.Pid, nil, syscall.WNOHANG, nil)
		if pid != 0 {
			log.Printf("%s's command terminated.\n", mount)
			delete(w.activeCmds, mount)
		}
	}
}

/*
 * Checks if a mount is missing. If it is, executes the designated command.
 */
func (w *WatchDog) checkForMissingMounts(mTabFile *os.File) {
	presenceSet := make(map[string]bool)
	tabChan := make(chan mtab.MTabEntry)

	go mtab.ReadMTab(mTabFile, tabChan)

	for tab := range tabChan {
		presenceSet[tab.MountPath] = true
	}

	for mount, cmd := range w.MissingMountToCmd {
		// skip any commands that are running from a previous epoch
		if _, found := w.activeCmds[mount]; found {
			continue
		}

		// run the command
		if _, found := presenceSet[mount]; !found {
			log.Printf("%s is missing. Running its command...\n", mount)
			runner := exec.Command("sh", "-c", cmd)
			runner.Start()
			w.activeCmds[mount] = runner
		}
	}
}

/*
 * Reaps all child processes. Should be called before exiting.
 */
func (w *WatchDog) reapChildren() {
	var wg sync.WaitGroup
	wg.Add(len(w.activeCmds))
	for mount, child := range w.activeCmds {
		go func(mount string, child *exec.Cmd) {
			defer wg.Done()

			log.Printf("Terminating process for mount %s...\n", mount)

			// first, try to interrupt the process
			child.Process.Signal(syscall.SIGINT)
			time.Sleep(time.Millisecond * 5)
			if pid, _ := syscall.Wait4(child.Process.Pid, nil, syscall.WNOHANG, nil); pid != child.Process.Pid {
				return
			}

			// then, try to terminate the process
			child.Process.Signal(syscall.SIGTERM)
			time.Sleep(time.Millisecond * 5)
			if pid, _ := syscall.Wait4(child.Process.Pid, nil, syscall.WNOHANG, nil); pid != child.Process.Pid {
				return
			}

			// finally, commit infanticide
			child.Process.Signal(syscall.SIGKILL)
			time.Sleep(time.Millisecond * 5)
			syscall.Wait4(child.Process.Pid, nil, syscall.WNOHANG, nil)
		}(mount, child)
	}
	wg.Wait()
}
