# Architecture

The launcher runs one foreground session at a time. It delegates transport, encoding,
audio and display to independently installed adb and scrcpy.

| Boundary | Responsibility |
|---|---|
| `cmd/scrcpy-tui` | Process entrypoint, streams, release version, exit status |
| `internal/cli` | Argument parsing and noninteractive command orchestration |
| `internal/app/domain.go` | Device, preset and persisted configuration data |
| `internal/app/presets.go` | Defaults, free arguments, validation |
| `internal/app/config.go` | Recovery copies and atomic restrictive writes |
| `internal/app/adb.go` | Injectable ADB runner, inventory parser, bounded discovery |
| `internal/app/plan.go` | Exact argv, preview and pre-launch target revalidation |
| `internal/app/session.go` | Owned process groups and interruption cleanup |
| `internal/tui` | Keyboard state, viewport, styling, bounded diagnostics |

Both frontends use `App.Plan` and `App.Prepare`. Neither uses a shell to launch
scrcpy. Interactive Enter and noninteractive `launch --yes` are the launch actions.
Discovery runs asynchronously; the UI locks changing actions during discovery and
preparation. Quit remains available. Device selection is retained by serial ID.

Configuration is loaded without writing. Invalid records are filtered; corrupt JSON
uses defaults. Before replacing a recovered configuration, its original bytes are
copied to a mode-0600 recovery file. Reading errors prevent overwriting. Normal
writes use a temporary file, fsync, close and rename. This is a single-user profile;
concurrent manual edits or multiple instances sharing a profile are not merged.

Session subprocesses get their own process group. Interruptions request graceful
termination, then may force only that group after two seconds. Parent exit also
cleans up descendants. WaitDelay prevents inherited output pipes from hanging the
launcher indefinitely. scrcpy diagnostics retain 32 KiB in memory and strip terminal
control sequences when rendered. No persistent log or telemetry is added.

Tests exercise three layers: pure parsers/state/rendering; CLI with injectable
streams and fake executables; real Linux PTYs with disposable profiles, real
subprocess groups and termios restoration. Automated tests cannot prove a real
Android device displays video correctly.

Windows supervision (under qualification) creates scrcpy suspended in a new
console process group, assigns it to a non-inheritable Job Object with
`KILL_ON_JOB_CLOSE`, then resumes its initial thread. Assignment failure fails
closed. No name lookup or system-wide process termination is used. Interrupts
request a group-scoped Ctrl+Break, wait up to two seconds, then terminate the job.
A retained process handle protects the group ID from reuse during cleanup. The
job also owns descendants that outlive their parent or retain an output pipe.

Graceful control events require a shared console and a cooperating child; the
Job Object remains the forced-cleanup boundary when either is unavailable.
There is a narrow startup crash boundary between creating the suspended process
and assigning it to the job: abrupt termination of the launcher in that interval
can leave a suspended child. Normal error paths kill and reap it. Runtime crash
cleanup after successful assignment is provided by closing the last job handle.
Native automated helpers validate process behavior; Windows Terminal, actual
scrcpy video/audio and MSI qualification remain separate release gates.

API references: [Job Objects](https://learn.microsoft.com/en-us/windows/win32/procthread/job-objects),
[control events](https://learn.microsoft.com/en-us/windows/console/generateconsolectrlevent).
