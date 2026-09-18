#!/usr/bin/env python3
"""Recette Linux isolée avec ADB/scrcpy factices et restauration termios."""
import atexit, os, pty, termios, tempfile, pathlib, subprocess, select, time, struct, fcntl, signal, sys
spawned=[]
def spawn(*args,**kwargs):
    proc=subprocess.Popen(*args,**kwargs);spawned.append(proc);return proc
def cleanup():
    for proc in spawned:
        if proc.poll() is None:
            proc.terminate()
            try:proc.wait(timeout=5)
            except subprocess.TimeoutExpired:proc.kill();proc.wait()
atexit.register(cleanup)
binary=str(pathlib.Path(sys.argv[1] if len(sys.argv)>1 else 'bin/scrcpy-tui').resolve())
with tempfile.TemporaryDirectory() as root:
    root=pathlib.Path(root)
    for name,body in {'adb':'#!/bin/sh\nprintf "List of devices attached\\nUSB device model:Pixel_8\\nBAD unauthorized\\n"\n','scrcpy':'#!/bin/sh\nprintf "%s\\n" "$@" > "$TEST_ROOT/launched"\nprintf \"ERROR synthetic diagnostic\\n\" >&2\necho $$ > "$TEST_ROOT/session-pid"\nsleep "${TEST_SLEEP:-0.1}" &\necho $! > "$TEST_ROOT/child-pid"\nwait\nexit "${TEST_EXIT:-0}"\n'}.items():
        path=root/name;path.write_text(body);path.chmod(0o755)
    env={**os.environ,'PATH':str(root)+':'+os.environ['PATH'],'TEST_ROOT':str(root),'TERM':'xterm-256color'}
    for finish in ('q','failure','ctrl-c','sigterm'):
        master,slave=pty.openpty();before=termios.tcgetattr(slave)
        fcntl.ioctl(slave,termios.TIOCSWINSZ,struct.pack('HHHH',24,100,0,0))
        env['TEST_EXIT']='7' if finish=='failure' else '0'
        proc=spawn([binary,'--config',str(root/'presets.json')],stdin=slave,stdout=slave,stderr=slave,env=env)
        output=b''
        def read(seconds=.5):
            global output
            deadline=time.monotonic()+seconds
            while time.monotonic()<deadline:
                ready,_,_=select.select([master],[],[],.05)
                if ready:
                    try:output+=os.read(master,65536)
                    except OSError:break
        read(1)
        assert b'USB' in output,output
        if finish=='q':
            os.write(master,b'\x1b[B\t\x1b[B\r');read(1)
            assert (root/'launched').exists(),'scrcpy not launched'
            args=(root/'launched').read_text().splitlines();assert args[:2]==['-s','USB'],args
            assert '--max-fps=25' in args,args
            os.write(master,b'n');read(.2);os.write(master,b'Perso\tDescription\t--no-audio\r');read(.4)
            assert 'Perso' in (root/'presets.json').read_text()
            os.write(master,b'dn');read(.2)
            assert 'Perso' in (root/'presets.json').read_text()
            os.write(master,b'c');read(.2);assert b'Commande exacte' in output
            os.write(master,b'\x1b');read(.2);os.write(master,b'q')
        elif finish=='failure':
            # A preset with an invalid option is represented by a failing fake process.
            os.write(master,b'\r');read(1);assert 'Échec scrcpy'.encode() in output
            output=b'';os.write(master,b'l');read(.2);assert b'Diagnostic scrcpy' in output and b'ERROR synthetic diagnostic' in output
            os.write(master,b'\x1b');read(.2);os.write(master,b'q')
        elif finish=='ctrl-c':os.write(master,b'\x03')
        else:proc.send_signal(signal.SIGTERM)
        proc.wait(timeout=5);read(.2)
        assert termios.tcgetattr(slave)==before,'terminal not restored'
        os.close(master);os.close(slave)
    preview=subprocess.run([binary,'--config',str(root/'presets.json'),'--device','USB','preview'],env=env,capture_output=True)
    assert preview.returncode==0,preview.stderr
    after=subprocess.run([binary,'preview','--config',str(root/'presets.json'),'--device','USB'],env=env,capture_output=True)
    assert after.returncode==0 and after.stdout==preview.stdout
    refused=subprocess.run([binary,'--config',str(root/'presets.json'),'--device','USB','launch'],env=env,capture_output=True)
    assert refused.returncode==4
    env['TEST_EXIT']='7'
    failed=subprocess.run([binary,'--config',str(root/'presets.json'),'--device','USB','--yes','launch'],env=env,capture_output=True)
    assert failed.returncode==1
    # Exercise the real launch bootstrap: corrupt configuration must be recoverable.
    import json
    env['TEST_EXIT']='0'
    for text in ('{broken', '{"presets":[{"name":"Keep","args":[]},{"name":"","args":[]}]}'):
        recovered=root/'recover.json';recovered.write_text(text)
        chosen='Keep' if text.startswith('{"presets"') else 'Léger Wi-Fi'
        result=subprocess.run([binary,'launch','--config',str(recovered),'--device','USB','--preset',chosen,'--yes'],env=env,capture_output=True)
        assert result.returncode==0,result.stderr
        assert any(p.read_text()==text for p in root.glob('recover.json.recovery-*')), 'original lost at launch bootstrap'
        assert json.loads(recovered.read_text())['last_device']=='USB'
    # Signals during a live session must stop both owned processes and restore termios.
    env['TEST_SLEEP']='20'
    for sig in (signal.SIGTERM,signal.SIGINT):
        for name in ('session-pid','child-pid'):(root/name).unlink(missing_ok=True)
        master,slave=pty.openpty();before=termios.tcgetattr(slave)
        fcntl.ioctl(slave,termios.TIOCSWINSZ,struct.pack('HHHH',24,100,0,0))
        proc=spawn([binary,'--config',str(root/'presets.json')],stdin=slave,stdout=slave,stderr=slave,env=env)
        output=b'';read(.5);os.write(master,b'\r');read(.5)
        assert (root/'session-pid').exists(), 'live process not started'
        live_pids=[int((root/name).read_text()) for name in ('session-pid','child-pid')]
        for pid in live_pids:assert (pathlib.Path('/proc')/str(pid)/'stat').exists(), 'session already exited'
        proc.send_signal(sig);read(3)
        if sig==signal.SIGINT:
            assert 'Session interrompue'.encode() in output
            os.write(master,b'q')
        proc.wait(timeout=5)
        for pid in live_pids:
            stat=pathlib.Path('/proc')/str(pid)/'stat'
            assert not stat.exists() or stat.read_text().split(') ')[1].startswith('Z '),'owned process still running'
        assert termios.tcgetattr(slave)==before, 'terminal not restored after live-session signal'
        os.close(master);os.close(slave)
    proc=spawn([binary,'launch','--config',str(root/'presets.json'),'--device','USB','--yes'],env=env,stdout=subprocess.DEVNULL,stderr=subprocess.DEVNULL)
    time.sleep(.5);proc.send_signal(signal.SIGTERM);assert proc.wait(timeout=6)==130
    env.pop('TEST_SLEEP')
    # Invalid selector variants must not reach either the launcher or persistence.
    before_launch=(root/'launched').read_bytes();before_config=(root/'presets.json').read_bytes()
    for arg in ('-fsOTHER','--ser=OTHER'):
        rejected=subprocess.run([binary,'launch','--config',str(root/'presets.json'),'--device','USB','--yes','--args='+arg],env=env,capture_output=True)
        assert rejected.returncode==3, rejected.stderr
        assert (root/'launched').read_bytes()==before_launch
        assert (root/'presets.json').read_bytes()==before_config
    # Charged inventory at three terminal sizes: search, cancel and keyboard paging.
    (root/'adb').write_text('#!/bin/sh\nprintf "List of devices attached\\n"\n'+''.join('printf "10.0.0.1:%d device model:Dongle_G_4K\\n"\n'%(6201+i) for i in range(42)))
    for width,height in ((40,14),(100,24),(140,48)):
        master,slave=pty.openpty();before=termios.tcgetattr(slave)
        fcntl.ioctl(slave,termios.TIOCSWINSZ,struct.pack('HHHH',height,width,0,0))
        proc=spawn([binary,'--config',str(root/'charged.json')],stdin=slave,stdout=slave,stderr=slave,env=env)
        output=b'';read(.5)
        assert b'6201' in output
        os.write(master,b'/6242');read(.3);assert b'6242' in output
        os.write(master,b'\r\x1b');read(.3)
        os.write(master,b'/absent');read(.3);assert b'Aucun appareil' in output
        os.write(master,b'\r');read(.2);os.write(master,b'\r');read(.2)
        assert not (root/'charged.json').exists(),'empty filter launched'
        os.write(master,b'\x1b');read(.2);os.write(master,b'\x1b[6~\x1b[F');read(.3)
        os.write(master,b'?');read(.2);assert b'Aide clavier' in output
        os.write(master,b'\x1b');read(.2);os.write(master,b'q');proc.wait(timeout=5)
        assert termios.tcgetattr(slave)==before
        os.close(master);os.close(slave)
print('PTY : navigation, lancement, édition, annulation, commande, q/Ctrl+C/SIGTERM, restauration et CLI OK')
