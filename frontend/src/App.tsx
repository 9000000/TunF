import { useState, useEffect } from 'react';
import './style.css';
import { StartProxy, StopProxy, IsProxyRunning, GetConfig, SetAutoStart } from "../wailsjs/go/main/App";

function App() {
    const [isRunning, setIsRunning] = useState(false);
    const [listenPort, setListenPort] = useState("5678");
    const [targetAddr, setTargetAddr] = useState("localhost:1234");
    const [autoOpenFirewall, setAutoOpenFirewall] = useState(false);
    const [autoStart, setAutoStart] = useState(false);
    const [history, setHistory] = useState<string[]>([]);
    const [targetHistory, setTargetHistory] = useState<string[]>([]);
    const [proxyPortHistory, setProxyPortHistory] = useState<string[]>([]);
    const [statusText, setStatusText] = useState("Ready to proxy");
    const [logs, setLogs] = useState<string[]>(["Application started"]);

    const addLog = (msg: string) => {
        setLogs(prev => [msg, ...prev].slice(0, 5));
    };

    const loadSettings = async () => {
        const config = await GetConfig();
        if (config) {
            setListenPort(config.lastListenPort);
            setTargetAddr(config.lastTargetAddr);
            setAutoOpenFirewall(config.autoOpenFirewall);
            setAutoStart(config.autoStart);
            setHistory(config.history || []);
            setTargetHistory(config.targetHistory || []);
            setProxyPortHistory(config.proxyPortHistory || []);
        }
    };

    // Check proxy status
    const checkStatus = async () => {
        const running = await IsProxyRunning();
        setIsRunning(running);
    };

    // Load settings and check status on startup
    useEffect(() => {
        loadSettings();
        checkStatus();
    }, []);

    const handleAutoStartToggle = async () => {
        const newValue = !autoStart;
        const result = await SetAutoStart(newValue);
        if (result === "Success") {
            setAutoStart(newValue);
            addLog(`Auto-start ${newValue ? 'enabled' : 'disabled'}`);
        } else {
            addLog(`Failed to set auto-start: ${result}`);
        }
    };

    const toggleProxy = async () => {
        if (isRunning) {
            await StopProxy();
            setIsRunning(false);
            setStatusText("Proxy stopped");
            addLog("Stopped proxy server");
        } else {
            setStatusText("Starting...");
            const result = await StartProxy(listenPort, targetAddr, autoOpenFirewall);
            if (result === "Success") {
                setIsRunning(true);
                setStatusText(`Running: :${listenPort} -> ${targetAddr}`);
                addLog(`Started proxy on port ${listenPort}`);
                if (autoOpenFirewall) addLog("Firewall rule added");
                // Refresh history after success
                const config = await GetConfig();
                setHistory(config.history || []);
                setTargetHistory(config.targetHistory || []);
                setProxyPortHistory(config.proxyPortHistory || []);
            } else {
                setStatusText(result);
                addLog(`Failed: ${result}`);
            }
        }
    };

    return (
        <div id="app">
            <div className="container">
                <div className="logo-container">
                    <div className="logo">TF</div>
                    <div className="brand-name">TunF</div>
                </div>

                <div className={`status-badge ${isRunning ? 'status-active' : 'status-inactive'}`}>
                    {isRunning ? '● ACTIVE' : '○ INACTIVE'}
                </div>

                <div className="input-group">
                    <label>Proxy Port (Destination)</label>
                    <input 
                        type="text" 
                        value={listenPort} 
                        onChange={(e) => setListenPort(e.target.value)}
                        placeholder="e.g. 5678"
                        disabled={isRunning}
                    />
                    {!isRunning && proxyPortHistory.length > 0 && (
                        <div className="history-chips">
                            {proxyPortHistory.map(port => (
                                <span key={port} onClick={() => setListenPort(port)}>
                                    {port}
                                </span>
                            ))}
                        </div>
                    )}
                </div>

                <div className="input-group">
                    <label>Target Address (Source)</label>
                    <input 
                        type="text" 
                        value={targetAddr} 
                        onChange={(e) => setTargetAddr(e.target.value)}
                        placeholder="e.g. localhost:1234"
                        disabled={isRunning}
                    />
                    {!isRunning && targetHistory.length > 0 && (
                        <div className="history-chips">
                            {targetHistory.map(target => (
                                <span key={target} onClick={() => setTargetAddr(target)}>
                                    {target}
                                </span>
                            ))}
                        </div>
                    )}
                </div>

                <div className="settings-row">
                    <div 
                        className="checkbox-group" 
                        onClick={() => !isRunning && setAutoOpenFirewall(!autoOpenFirewall)}
                    >
                        <input 
                            type="checkbox" 
                            checked={autoOpenFirewall} 
                            onChange={() => {}} 
                            disabled={isRunning}
                        />
                        <div className="checkbox-label">
                            Firewall
                            <span>Open Port</span>
                        </div>
                    </div>

                    <div 
                        className="checkbox-group" 
                        onClick={handleAutoStartToggle}
                    >
                        <input 
                            type="checkbox" 
                            checked={autoStart} 
                            onChange={() => {}} 
                        />
                        <div className="checkbox-label">
                            Auto-start
                            <span>With Windows</span>
                        </div>
                    </div>
                </div>

                <button 
                    className={`toggle-btn ${isRunning ? 'btn-stop' : 'btn-start'}`}
                    onClick={toggleProxy}
                >
                    {isRunning ? 'Stop Proxy' : 'Start Proxy'}
                </button>

                <div className="log-area">
                    {logs.map((log, i) => (
                        <div key={i}>[{new Date().toLocaleTimeString()}] {log}</div>
                    ))}
                </div>
            </div>
        </div>
    );
}

export default App;
