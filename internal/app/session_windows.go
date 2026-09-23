package app

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Create the child suspended: it cannot spawn descendants before job assignment.
// The job is unnamed, non-inheritable and disallows breakaway. Closing its last
// handle terminates only the tree launched by this session.
func runSession(cmd *exec.Cmd, signals <-chan os.Signal) error {
	job, err := windows.CreateJobObject(nil, nil)
	if err != nil {
		return fmt.Errorf("créer le Job Object : %w", err)
	}
	defer windows.CloseHandle(job)
	limits := windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION{}
	limits.BasicLimitInformation.LimitFlags = windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE
	if _, err = windows.SetInformationJobObject(job, windows.JobObjectExtendedLimitInformation, uintptr(unsafe.Pointer(&limits)), uint32(unsafe.Sizeof(limits))); err != nil {
		return fmt.Errorf("configurer le Job Object : %w", err)
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: windows.CREATE_SUSPENDED | windows.CREATE_NEW_PROCESS_GROUP}
	cmd.WaitDelay = 2 * time.Second
	if err = cmd.Start(); err != nil {
		return err
	}
	// Keep a handle until cleanup is finished, even after Wait, so that the
	// process-group ID cannot refer to a recycled, unrelated process.
	process, err := windows.OpenProcess(windows.PROCESS_SET_QUOTA|windows.PROCESS_TERMINATE|windows.PROCESS_QUERY_LIMITED_INFORMATION, false, uint32(cmd.Process.Pid))
	if err != nil {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
		return fmt.Errorf("ouvrir le processus suspendu : %w", err)
	}
	defer windows.CloseHandle(process)
	if err = windows.AssignProcessToJobObject(job, process); err != nil {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
		return fmt.Errorf("affecter le processus au Job Object : %w", err)
	}
	if err = resumeSessionProcess(uint32(cmd.Process.Pid)); err != nil {
		_ = windows.TerminateJobObject(job, 1)
		_ = cmd.Wait()
		return fmt.Errorf("reprendre le processus : %w", err)
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case err = <-done:
		stopWindowsJob(job, uint32(cmd.Process.Pid))
		return err
	case sig := <-signals:
		stopWindowsJob(job, uint32(cmd.Process.Pid))
		<-done
		return &Interrupted{Signal: sig}
	}
}

// os/exec closes the initial thread handle. Locate the still-suspended thread
// through the documented Toolhelp API instead of using undocumented Nt APIs.
func resumeSessionProcess(pid uint32) error {
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPTHREAD, 0)
	if err != nil {
		return err
	}
	defer windows.CloseHandle(snapshot)
	entry := windows.ThreadEntry32{Size: uint32(unsafe.Sizeof(windows.ThreadEntry32{}))}
	for err = windows.Thread32First(snapshot, &entry); err == nil; err = windows.Thread32Next(snapshot, &entry) {
		if entry.OwnerProcessID != pid {
			continue
		}
		thread, err := windows.OpenThread(windows.THREAD_SUSPEND_RESUME|windows.THREAD_QUERY_LIMITED_INFORMATION, false, entry.ThreadID)
		if err != nil {
			return err
		}
		// The snapshot is not an ownership handle. Check the opened thread before
		// acting in case it exited and its ID was reused concurrently.
		owner, _, ownerErr := windows.NewLazySystemDLL("kernel32.dll").NewProc("GetProcessIdOfThread").Call(uintptr(thread))
		if owner != uintptr(pid) {
			windows.CloseHandle(thread)
			return fmt.Errorf("identité du thread modifiée (owner=%d) : %v", owner, ownerErr)
		}
		_, err = windows.ResumeThread(thread)
		windows.CloseHandle(thread)
		return err
	}
	return fmt.Errorf("thread initial introuvable : %w", err)
}

type jobAccounting struct {
	TotalUserTime, TotalKernelTime, ThisPeriodTotalUserTime, ThisPeriodTotalKernelTime int64
	TotalPageFaultCount, TotalProcesses, ActiveProcesses, TotalTerminatedProcesses     uint32
}

func stopWindowsJob(job windows.Handle, group uint32) {
	// CTRL_C cannot be scoped to a group; CTRL_BREAK can. Never broadcast to 0.
	_ = windows.GenerateConsoleCtrlEvent(windows.CTRL_BREAK_EVENT, group)
	until := time.Now().Add(2 * time.Second)
	for time.Now().Before(until) {
		var info jobAccounting
		err := windows.QueryInformationJobObject(job, windows.JobObjectBasicAccountingInformation, uintptr(unsafe.Pointer(&info)), uint32(unsafe.Sizeof(info)), nil)
		if err != nil || info.ActiveProcesses == 0 {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	_ = windows.TerminateJobObject(job, 1)
}
